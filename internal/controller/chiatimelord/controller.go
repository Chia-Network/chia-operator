/*
Copyright 2023 Chia Network Inc.
*/

package chiatimelord

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	"github.com/chia-network/chia-operator/internal/metrics"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
)

// ChiaTimelordReconciler reconciles a ChiaTimelord object
type ChiaTimelordReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chiatimelords map[string]bool = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiatimelords,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiatimelords/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiatimelords/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

func (r *ChiaTimelordReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log))
	log.Info(fmt.Sprintf("ChiaTimelordController ChiaTimelord=%s", req.NamespacedName.String()))

	// Get the custom resource
	var tl k8schianetv1.ChiaTimelord
	err := r.Get(ctx, req.NamespacedName, &tl)
	if err != nil && errors.IsNotFound(err) {
		// Remove this object from the map for tracking and subtract this CR's total metric by 1
		_, exists := chiatimelords[req.NamespacedName.String()]
		if exists {
			delete(chiatimelords, req.NamespacedName.String())
			metrics.ChiaTimelords.Sub(1.0)
		}

		// Return here, this can happen if the CR was deleted
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, fmt.Sprintf("ChiaTimelordController ChiaTimelord=%s unable to fetch ChiaTimelord resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiatimelords[req.NamespacedName.String()]
	if !exists {
		chiatimelords[req.NamespacedName.String()] = true
		metrics.ChiaTimelords.Add(1.0)
	}

	// Reconcile ChiaTimelord owned objects
	srv := r.assembleBaseService(ctx, tl)
	res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		r.Recorder.Event(&tl, corev1.EventTypeWarning, "Failed", "Failed to create timelord Service -- Check operator logs.")
		return *res, fmt.Errorf("ChiaTimelordController ChiaTimelord=%s encountered error reconciling node Service: %v", req.NamespacedName, err)
	}

	srv = r.assembleChiaExporterService(ctx, tl)
	res, err = kube.ReconcileService(ctx, resourceReconciler, srv)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		r.Recorder.Event(&tl, corev1.EventTypeWarning, "Failed", "Failed to create timelord metrics Service -- Check operator logs.")
		return *res, fmt.Errorf("ChiaTimelordController ChiaTimelord=%s encountered error reconciling node chia-exporter Service: %v", req.NamespacedName, err)
	}

	deploy := r.assembleDeployment(ctx, tl)
	res, err = kube.ReconcileDeployment(ctx, resourceReconciler, deploy)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		r.Recorder.Event(&tl, corev1.EventTypeWarning, "Failed", "Failed to create timelord Deployment -- Check operator logs.")
		return *res, fmt.Errorf("ChiaTimelordController ChiaTimelord=%s encountered error reconciling node StatefulSet: %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&tl, corev1.EventTypeNormal, "Created", "Successfully created ChiaTimelord resources.")
	tl.Status.Ready = true
	err = r.Status().Update(ctx, &tl)
	if err != nil {
		log.Error(err, fmt.Sprintf("ChiaTimelordController ChiaTimelord=%s unable to update ChiaNode status", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaTimelordReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaTimelord{}).
		Complete(r)
}
