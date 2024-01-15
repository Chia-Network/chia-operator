/*
Copyright 2023 Chia Network Inc.
*/

package chiawallet

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

// ChiaWalletReconciler reconciles a ChiaWallet object
type ChiaWalletReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chiawallets map[string]bool = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiawallets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiawallets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiawallets/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *ChiaWalletReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log))
	log.Info(fmt.Sprintf("ChiaWalletReconciler ChiaWallet=%s", req.NamespacedName.String()))

	// Get the custom resource
	var wallet k8schianetv1.ChiaWallet
	err := r.Get(ctx, req.NamespacedName, &wallet)
	if err != nil && errors.IsNotFound(err) {
		// Remove this object from the map for tracking and subtract this CR's total metric by 1
		_, exists := chiawallets[req.NamespacedName.String()]
		if exists {
			delete(chiawallets, req.NamespacedName.String())
			metrics.ChiaWallets.Sub(1.0)
		}

		// Return here, this can happen if the CR was deleted
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, fmt.Sprintf("ChiaWalletReconciler ChiaWallet=%s unable to fetch ChiaWallet resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiawallets[req.NamespacedName.String()]
	if !exists {
		chiawallets[req.NamespacedName.String()] = true
		metrics.ChiaWallets.Add(1.0)
	}

	// Reconcile ChiaWallet owned objects
	service := r.assembleBaseService(ctx, wallet)
	res, err := kube.ReconcileService(ctx, resourceReconciler, service)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		r.Recorder.Event(&wallet, "Warning", "Failed", "Failed to create harvester Service -- Check operator logs.")
		return *res, fmt.Errorf("ChiaWalletReconciler ChiaWallet=%s encountered error reconciling wallet Service: %v", req.NamespacedName, err)
	}

	service = r.assembleChiaExporterService(ctx, wallet)
	res, err = kube.ReconcileService(ctx, resourceReconciler, service)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		r.Recorder.Event(&wallet, "Warning", "Failed", "Failed to create harvester metrics Service -- Check operator logs.")
		return *res, fmt.Errorf("ChiaWalletReconciler ChiaWallet=%s encountered error reconciling wallet chia-exporter Service: %v", req.NamespacedName, err)
	}

	deploy := r.assembleDeployment(ctx, wallet)
	res, err = kube.ReconcileDeployment(ctx, resourceReconciler, deploy)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		r.Recorder.Event(&wallet, "Warning", "Failed", "Failed to create harvester Deployment -- Check operator logs.")
		return *res, fmt.Errorf("ChiaWalletReconciler ChiaWallet=%s encountered error reconciling wallet Deployment: %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&wallet, "Normal", "Created", "Successfully created ChiaWallet resources.")
	wallet.Status.Ready = true
	err = r.Status().Update(ctx, &wallet)
	if err != nil {
		log.Error(err, fmt.Sprintf("ChiaWalletReconciler ChiaWallet=%s unable to update ChiaWallet status", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaWalletReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaWallet{}).
		Complete(r)
}
