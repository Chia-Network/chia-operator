/*
Copyright 2023 Chia Network Inc.
*/

package chiatimelord

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	"github.com/chia-network/chia-operator/internal/metrics"
)

// ChiaTimelordReconciler reconciles a ChiaTimelord object
type ChiaTimelordReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chiatimelords = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiatimelords,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiatimelords/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiatimelords/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaTimelordReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info(fmt.Sprintf("ChiaTimelordController ChiaTimelord=%s running reconciler...", req.NamespacedName.String()))

	// Get the custom resource
	var timelord k8schianetv1.ChiaTimelord
	err := r.Get(ctx, req.NamespacedName, &timelord)
	if err != nil && errors.IsNotFound(err) {
		// Remove this object from the map for tracking and subtract this CR's total metric by 1
		_, exists := chiatimelords[req.NamespacedName.String()]
		if exists {
			delete(chiatimelords, req.NamespacedName.String())
			metrics.ChiaTimelords.Sub(1.0)
		}
		return ctrl.Result{}, nil
	}
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		log.Error(err, fmt.Sprintf("ChiaTimelordController ChiaTimelord=%s unable to fetch ChiaTimelord resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiatimelords[req.NamespacedName.String()]
	if !exists {
		chiatimelords[req.NamespacedName.String()] = true
		metrics.ChiaTimelords.Add(1.0)
	}

	// Assemble Peer Service
	peerSrv := assemblePeerService(timelord)
	if err := controllerutil.SetControllerReference(&timelord, &peerSrv, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&timelord, corev1.EventTypeWarning, "Failed", "Failed to assemble timelord peer Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaTimelordReconciler ChiaTimelord=%s encountered error assembling peer Service: %v", req.NamespacedName, err)
	}
	// Reconcile Peer Service
	err = kube.ReconcileService(ctx, r.Client, timelord.Spec.ChiaConfig.PeerService, peerSrv, true)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		return ctrl.Result{}, fmt.Errorf("ChiaTimelordReconciler ChiaTimelord=%s %v", req.NamespacedName, err)
	}

	// Assemble Daemon Service
	daemonSrv := assembleDaemonService(timelord)
	if err := controllerutil.SetControllerReference(&timelord, &daemonSrv, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&timelord, corev1.EventTypeWarning, "Failed", "Failed to assemble timelord daemon Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaTimelordReconciler ChiaTimelord=%s encountered error assembling daemon Service: %v", req.NamespacedName, err)
	}
	// Reconcile Daemon Service
	err = kube.ReconcileService(ctx, r.Client, timelord.Spec.ChiaConfig.DaemonService, daemonSrv, true)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		return ctrl.Result{}, fmt.Errorf("ChiaTimelordReconciler ChiaTimelord=%s %v", req.NamespacedName, err)
	}

	// Assemble RPC Service
	rpcSrv := assembleRPCService(timelord)
	if err := controllerutil.SetControllerReference(&timelord, &rpcSrv, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&timelord, corev1.EventTypeWarning, "Failed", "Failed to assemble timelord RPC Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaTimelordReconciler ChiaTimelord=%s encountered error assembling RPC Service: %v", req.NamespacedName, err)
	}
	// Reconcile RPC Service
	err = kube.ReconcileService(ctx, r.Client, timelord.Spec.ChiaConfig.RPCService, rpcSrv, true)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		return ctrl.Result{}, fmt.Errorf("ChiaTimelordReconciler ChiaTimelord=%s %v", req.NamespacedName, err)
	}

	// Assemble Chia-Exporter Service
	exporterSrv := assembleChiaExporterService(timelord)
	if err := controllerutil.SetControllerReference(&timelord, &exporterSrv, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&timelord, corev1.EventTypeWarning, "Failed", "Failed to assemble timelord chia-exporter Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaTimelordReconciler ChiaTimelord=%s encountered error assembling chia-exporter Service: %v", req.NamespacedName, err)
	}
	// Reconcile Chia-Exporter Service
	err = kube.ReconcileService(ctx, r.Client, timelord.Spec.ChiaExporterConfig.Service, exporterSrv, true)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		return ctrl.Result{}, fmt.Errorf("ChiaTimelordReconciler ChiaTimelord=%s %v", req.NamespacedName, err)
	}

	// Creates a persistent volume claim if the GenerateVolumeClaims setting was set to true
	if kube.ShouldMakeVolumeClaim(timelord.Spec.Storage) {
		pvc, err := assembleVolumeClaim(timelord)
		if err != nil {
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&timelord, corev1.EventTypeWarning, "Failed", "Failed to assemble timelord PVC -- Check operator logs.")
			return reconcile.Result{}, fmt.Errorf("ChiaTimelordReconciler ChiaTimelord=%s %v", req.NamespacedName, err)
		}

		if pvc != nil {
			err = kube.ReconcilePersistentVolumeClaim(ctx, r.Client, timelord.Spec.Storage, *pvc)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				r.Recorder.Event(&timelord, corev1.EventTypeWarning, "Failed", "Failed to create timelord PVC -- Check operator logs.")
				return reconcile.Result{}, fmt.Errorf("ChiaTimelordReconciler ChiaTimelord=%s %v", req.NamespacedName, err)
			}
		}
	}

	// Assemble Deployment
	deploy := assembleDeployment(timelord)
	if err := controllerutil.SetControllerReference(&timelord, &deploy, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&timelord, corev1.EventTypeWarning, "Failed", "Failed to assemble timelord Deployment -- Check operator logs.")
		return reconcile.Result{}, fmt.Errorf("ChiaTimelordReconciler ChiaTimelord=%s %v", req.NamespacedName, err)
	}
	// Reconcile Deployment
	err = kube.ReconcileDeployment(ctx, r.Client, deploy)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&timelord, corev1.EventTypeWarning, "Failed", "Failed to create timelord Deployment -- Check operator logs.")
		return reconcile.Result{}, fmt.Errorf("ChiaTimelordReconciler ChiaTimelord=%s %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&timelord, corev1.EventTypeNormal, "Created", "Successfully created ChiaTimelord resources.")
	timelord.Status.Ready = true
	err = r.Status().Update(ctx, &timelord)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		log.Error(err, fmt.Sprintf("ChiaTimelordController ChiaTimelord=%s unable to update ChiaNode status", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaTimelordReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaTimelord{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
