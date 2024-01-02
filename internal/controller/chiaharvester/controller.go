/*
Copyright 2023 Chia Network Inc.
*/

package chiaharvester

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	"github.com/chia-network/chia-operator/internal/metrics"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
)

// ChiaHarvesterReconciler reconciles a ChiaHarvester object
type ChiaHarvesterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var chiaharvesters map[string]bool = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiaharvesters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiaharvesters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiaharvesters/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *ChiaHarvesterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log))
	log.Info(fmt.Sprintf("ChiaHarvesterReconciler ChiaHarvester=%s", req.NamespacedName.String()))

	// Get the custom resource
	var harvester k8schianetv1.ChiaHarvester
	err := r.Get(ctx, req.NamespacedName, &harvester)
	if err != nil && errors.IsNotFound(err) {
		// Remove this object from the map for tracking and subtract this CR's total metric by 1
		_, exists := chiaharvesters[req.NamespacedName.String()]
		if exists {
			delete(chiaharvesters, req.NamespacedName.String())
			metrics.ChiaHarvesters.Sub(1.0)
		}

		// Return here, this can happen if the CR was deleted
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, fmt.Sprintf("ChiaHarvesterReconciler ChiaHarvester=%s unable to fetch ChiaHarvester resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiaharvesters[req.NamespacedName.String()]
	if !exists {
		chiaharvesters[req.NamespacedName.String()] = true
		metrics.ChiaHarvesters.Add(1.0)
	}

	// Reconcile ChiaHarvester owned objects
	srv := r.assembleBaseService(ctx, harvester)
	res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		return *res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error reconciling harvester Service: %v", req.NamespacedName, err)
	}

	srv = r.assembleChiaExporterService(ctx, harvester)
	res, err = kube.ReconcileService(ctx, resourceReconciler, srv)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		return *res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error reconciling harvester chia-exporter Service: %v", req.NamespacedName, err)
	}

	deploy := r.assembleDeployment(ctx, harvester)
	res, err = kube.ReconcileDeployment(ctx, resourceReconciler, deploy)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		return *res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error reconciling harvester Deployment: %v", req.NamespacedName, err)
	}

	// Update CR status
	harvester.Status.Ready = true
	err = r.Status().Update(ctx, &harvester)
	if err != nil {
		log.Error(err, fmt.Sprintf("ChiaHarvesterReconciler ChiaHarvester=%s unable to update ChiaHarvester status", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaHarvesterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaHarvester{}).
		Complete(r)
}
