/*
Copyright 2023 Chia Network Inc.
*/

package chianode

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

// ChiaNodeReconciler reconciles a ChiaNode object
type ChiaNodeReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chianodes map[string]bool = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chianodes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chianodes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chianodes/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *ChiaNodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log))
	log.Info(fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s", req.NamespacedName.String()))

	// Get the custom resource
	var node k8schianetv1.ChiaNode
	err := r.Get(ctx, req.NamespacedName, &node)
	if err != nil && errors.IsNotFound(err) {
		// Remove this object from the map for tracking and subtract this CR's total metric by 1
		_, exists := chianodes[req.NamespacedName.String()]
		if exists {
			delete(chianodes, req.NamespacedName.String())
			metrics.ChiaNodes.Sub(1.0)
		}

		// Return here, this can happen if the CR was deleted
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to fetch ChiaNode resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chianodes[req.NamespacedName.String()]
	if !exists {
		chianodes[req.NamespacedName.String()] = true
		metrics.ChiaNodes.Add(1.0)
	}

	// Reconcile ChiaNode owned objects
	srv := r.assembleBaseService(ctx, node)
	res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		return *res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error reconciling node Service: %v", req.NamespacedName, err)
	}

	srv = r.assembleInternalService(ctx, node)
	res, err = kube.ReconcileService(ctx, resourceReconciler, srv)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		return *res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error reconciling node Local Service: %v", req.NamespacedName, err)
	}

	srv = r.assembleHeadlessService(ctx, node)
	res, err = kube.ReconcileService(ctx, resourceReconciler, srv)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		return *res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error reconciling node headless Service: %v", req.NamespacedName, err)
	}

	srv = r.assembleChiaExporterService(ctx, node)
	res, err = kube.ReconcileService(ctx, resourceReconciler, srv)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		return *res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error reconciling node chia-exporter Service: %v", req.NamespacedName, err)
	}

	stateful := r.assembleStatefulset(ctx, node)
	res, err = kube.ReconcileStatefulset(ctx, resourceReconciler, stateful)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		return *res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error reconciling node StatefulSet: %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&node, "Info", "Created", "Successfully created ChiaNode resources.")
	node.Status.Ready = true
	err = r.Status().Update(ctx, &node)
	if err != nil {
		log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to update ChiaNode status", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaNodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaNode{}).
		Complete(r)
}
