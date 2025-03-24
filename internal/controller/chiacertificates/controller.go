/*
Copyright 2025 Chia Network Inc.
*/

package chiacertificates

import (
	"context"
	stdlibErrors "errors"
	"fmt"
	"strings"
	"time"

	"github.com/chia-network/go-chia-libs/pkg/tls"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	"github.com/chia-network/chia-operator/internal/metrics"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ChiaCertificatesReconciler reconciles a ChiaCertificates object
type ChiaCertificatesReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chiacertificates = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiacertificates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiacertificates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiacertificates/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaCertificatesReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Running reconciler...")

	// Get the custom resource
	var cr k8schianetv1.ChiaCertificates
	err := r.Get(ctx, req.NamespacedName, &cr)
	if err != nil && errors.IsNotFound(err) {
		// Remove this object from the map for tracking and subtract this CR's total metric by 1
		_, exists := chiacertificates[req.NamespacedName.String()]
		if exists {
			delete(chiacertificates, req.NamespacedName.String())
			metrics.ChiaCertificates.Sub(1.0)
		}
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, "unable to fetch ChiaCertificates resource")
		return ctrl.Result{}, err
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiacertificates[req.NamespacedName.String()]
	if !exists {
		chiacertificates[req.NamespacedName.String()] = true
		metrics.ChiaCertificates.Add(1.0)
	}

	// Verify that certificate Secret name does not match the CA Secret name
	certSecretName := getChiaCertificatesSecretName(cr)
	caSecretName := cr.Spec.CASecretName
	if certSecretName == caSecretName {
		log.Error(stdlibErrors.New("certificate Secret cannot be the same name as the CA Secret"),
			"Invalid certificate Secret name", "certificate Secret name", certSecretName, "CA Secret name", caSecretName)
		return ctrl.Result{}, nil
	}

	// Check if certificate Secret exists
	_, certSecretExists, err := r.getSecret(ctx, cr.Namespace, certSecretName)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("encountered error querying for existing Certificates Secret: %v", err)
	}

	// If Certificates Secret doesn't exist, generate certificates and create one
	if !certSecretExists {
		caSecret, caSecretExists, err := r.getSecret(ctx, cr.Namespace, caSecretName)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("encountered error querying for existing CA Secret: %v", err)
		}
		if !caSecretExists {
			log.Info("CA Secret not found, cancelling reconciliation and retrying in 10 seconds")
			return ctrl.Result{
				RequeueAfter: 10 * time.Second,
			}, nil
		}

		privateCACertData, ok := caSecret.Data["private_ca.crt"]
		if !ok {
			return ctrl.Result{}, fmt.Errorf("private CA certificate not present in CA Secret")
		}
		privateCACert, err := tls.ParsePemCertificate(privateCACertData)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("error parsing private CA certificate from Secret: %v", err)
		}

		privateCAKeyData, ok := caSecret.Data["private_ca.key"]
		if !ok {
			return ctrl.Result{}, fmt.Errorf("private CA key not present in CA Secret")
		}
		privateCAKey, err := tls.ParsePemKey(privateCAKeyData)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("error parsing private CA key from Secret: %v", err)
		}

		allCerts, err := tls.GenerateAllCerts(privateCACert, privateCAKey)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("error generating new certificates: %v", err)
		}

		certMap, err := constructCertMap(allCerts)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("error converting certificates to map: %v", err)
		}

		secret := assembleSecret(cr, certMap)
		if err = r.Create(ctx, &secret); err != nil {
			return ctrl.Result{}, fmt.Errorf("error creating certificate Secret \"%s\": %v", secret.Name, err)
		}
	}

	if !cr.Status.Ready {
		r.Recorder.Event(&cr, corev1.EventTypeNormal, "Created",
			fmt.Sprintf("Successfully created Certificates Secret in %s/%s", cr.Namespace, cr.Name))

		cr.Status.Ready = true
		err = r.Status().Update(ctx, &cr)
		if err != nil {
			if strings.Contains(err.Error(), kube.ObjectModifiedTryAgainError) {
				return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
			}
			log.Error(err, "encountered error updating ChiaCertificates status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaCertificatesReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaCertificates{}).
		Complete(r)
}
