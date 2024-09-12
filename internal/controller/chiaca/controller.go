/*
Copyright 2023 Chia Network Inc.
*/

package chiaca

import (
	"context"
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

// ChiaCAReconciler reconciles a ChiaCA object
type ChiaCAReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chiacas = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiacas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiacas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiacas/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaCAReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Running reconciler...")

	// Get the custom resource
	var ca k8schianetv1.ChiaCA
	err := r.Get(ctx, req.NamespacedName, &ca)
	if err != nil && errors.IsNotFound(err) {
		// Remove this object from the map for tracking and subtract this CR's total metric by 1
		_, exists := chiacas[req.NamespacedName.String()]
		if exists {
			delete(chiacas, req.NamespacedName.String())
			metrics.ChiaCAs.Sub(1.0)
		}
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, "unable to fetch ChiaCA resource")
		return ctrl.Result{}, err
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiacas[req.NamespacedName.String()]
	if !exists {
		chiacas[req.NamespacedName.String()] = true
		metrics.ChiaCAs.Add(1.0)
	}

	// Check if CA Secret exists
	caExists, err := r.caSecretExists(ctx, ca)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("encountered error querying for existing CA Secret: %v", err)
	}

	// If CA Secret doesn't exist, generate a CA and create one
	if !caExists {
		// Get the public CA cert and key byte slices
		publicCACrtBytes, publicCAKeyBytes := tls.GetChiaCACertAndKey()

		// Generate a private CA cert and key
		privateCACrt, privateCAKey, err := tls.GenerateNewCA()
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("encountered error generating new private CA cert and key: %v", err)
		}

		// Encode the private CA cert and key to PEM byte slices
		privateCACrtBytes, privateCAKeyBytes, err := tls.EncodeCertAndKeyToPEM(privateCACrt, privateCAKey)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("encountered error encoding private CA cert and key to PEM: %v", err)
		}

		// Assemble CA Secret and create in cluster
		secret := assembleCASecret(ca, string(publicCACrtBytes), string(publicCAKeyBytes), string(privateCACrtBytes), string(privateCAKeyBytes))
		if err = r.Create(ctx, &secret); err != nil {
			return ctrl.Result{}, fmt.Errorf("error creating CA Secret \"%s\": %v", secret.Name, err)
		}
	}

	if !ca.Status.Ready {
		r.Recorder.Event(&ca, corev1.EventTypeNormal, "Created",
			fmt.Sprintf("Successfully created CA Secret in %s/%s", ca.Namespace, ca.Name))

		ca.Status.Ready = true
		err = r.Status().Update(ctx, &ca)
		if err != nil {
			if strings.Contains(err.Error(), kube.ObjectModifiedTryAgainError) {
				return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
			}
			log.Error(err, "encountered error updating ChiaCA status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaCAReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaCA{}).
		Complete(r)
}
