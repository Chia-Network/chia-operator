/*
Copyright 2023 Chia Network Inc.
*/

package chiafarmer

import (
	"context"
	"fmt"

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

// ChiaFarmerReconciler reconciles a ChiaFarmer object
type ChiaFarmerReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chiafarmers map[string]bool = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiafarmers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiafarmers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiafarmers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *ChiaFarmerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log))
	log.Info(fmt.Sprintf("ChiaFarmerReconciler ChiaFarmer=%s", req.NamespacedName.String()))

	// Get the custom resource
	var farmer k8schianetv1.ChiaFarmer
	err := r.Get(ctx, req.NamespacedName, &farmer)
	if err != nil && errors.IsNotFound(err) {
		// Remove this object from the map for tracking and subtract this CR's total metric by 1
		_, exists := chiafarmers[req.NamespacedName.String()]
		if exists {
			delete(chiafarmers, req.NamespacedName.String())
			metrics.ChiaFarmers.Sub(1.0)
		}

		// Return here, this can happen if the CR was deleted
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, fmt.Sprintf("ChiaFarmerReconciler ChiaFarmer=%s unable to fetch ChiaFarmer resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiafarmers[req.NamespacedName.String()]
	if !exists {
		chiafarmers[req.NamespacedName.String()] = true
		metrics.ChiaFarmers.Add(1.0)
	}

	// Reconcile ChiaFarmer owned objects
	srv := r.assembleBaseService(ctx, farmer)
	res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		r.Recorder.Event(&farmer, "Warning", "Failed", "Failed to create farmer Service -- Check operator logs.")
		return *res, fmt.Errorf("ChiaFarmerReconciler ChiaFarmer=%s encountered error reconciling farmer Service: %v", req.NamespacedName, err)
	}

	srv = r.assembleChiaExporterService(ctx, farmer)
	res, err = kube.ReconcileService(ctx, resourceReconciler, srv)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		r.Recorder.Event(&farmer, "Warning", "Failed", "Failed to create farmer metrics Service -- Check operator logs.")
		return *res, fmt.Errorf("ChiaFarmerReconciler ChiaFarmer=%s encountered error reconciling farmer chia-exporter Service: %v", req.NamespacedName, err)
	}

	deploy := r.assembleDeployment(ctx, farmer)
	res, err = kube.ReconcileDeployment(ctx, resourceReconciler, deploy)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		r.Recorder.Event(&farmer, "Warning", "Failed", "Failed to create farmer Deployment -- Check operator logs.")
		return *res, fmt.Errorf("ChiaFarmerReconciler ChiaFarmer=%s encountered error reconciling farmer Deployment: %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&farmer, "Normal", "Created", "Successfully created ChiaFarmer resources.")
	farmer.Status.Ready = true
	err = r.Status().Update(ctx, &farmer)
	if err != nil {
		log.Error(err, fmt.Sprintf("ChiaFarmerReconciler ChiaFarmer=%s unable to update ChiaFarmer status", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaFarmerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaFarmer{}).
		Complete(r)
}
