/*
Copyright 2024 Chia Network Inc.
*/

package chianetwork

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	"github.com/chia-network/chia-operator/internal/metrics"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
)

// ChiaNetworkReconciler reconciles a ChiaNetwork object
type ChiaNetworkReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chianetworks = make(map[string]bool)

// +kubebuilder:rbac:groups=k8s.chia.net,resources=chianetworks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=k8s.chia.net,resources=chianetworks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=k8s.chia.net,resources=chianetworks/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaNetworkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog := log.FromContext(ctx)
	klog.Info("Running reconciler...")

	// Get the custom resource
	var network k8schianetv1.ChiaNetwork
	err := r.Get(ctx, req.NamespacedName, &network)
	if err != nil && errors.IsNotFound(err) {
		// Remove this object from the map for tracking and subtract this CR's total metric by 1
		_, exists := chianetworks[req.NamespacedName.String()]
		if exists {
			delete(chianetworks, req.NamespacedName.String())
			metrics.ChiaNetworks.Sub(1.0)
		}
		return ctrl.Result{}, nil
	}
	if err != nil {
		klog.Error(err, "unable to fetch ChiaNetwork resource")
		return ctrl.Result{}, err
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chianetworks[req.NamespacedName.String()]
	if !exists {
		chianetworks[req.NamespacedName.String()] = true
		metrics.ChiaNetworks.Add(1.0)
	}

	// Assemble configmap
	configmap, err := assembleConfigMap(network)
	if err != nil {
		r.Recorder.Event(&network, corev1.EventTypeWarning, "Failed", "Failed to assemble network ConfigMap -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("encountered error assembling network ConfigMap: %v", err)
	}
	if err := controllerutil.SetControllerReference(&network, &configmap, r.Scheme); err != nil {
		r.Recorder.Event(&network, corev1.EventTypeWarning, "Failed", "Failed to set controller reference on network ConfigMap -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("encountered error setting controller reference on network ConfigMap: %v", err)
	}

	// Reconcile configmap
	res, err := kube.ReconcileConfigMap(ctx, r.Client, configmap)
	if err != nil {
		r.Recorder.Event(&network, corev1.EventTypeWarning, "Failed", "Failed to reconcile network ConfigMap -- Check operator logs.")
		return res, fmt.Errorf("encountered error reconciling network ConfigMap: %v", err)
	}

	if !network.Status.Ready {
		r.Recorder.Event(&network, corev1.EventTypeNormal, "Created",
			fmt.Sprintf("Successfully created network ConfigMap in %s/%s", network.Namespace, network.Name))

		network.Status.Ready = true
		err = r.Status().Update(ctx, &network)
		if err != nil {
			if strings.Contains(err.Error(), kube.ObjectModifiedTryAgainError) {
				return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
			}
			klog.Error(err, "encountered error updating ChiaNetwork status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaNetworkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaNetwork{}).
		Complete(r)
}
