/*
Copyright 2023 Chia Network Inc.
*/

package chiaharvester

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

// ChiaHarvesterReconciler reconciles a ChiaHarvester object
type ChiaHarvesterReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chiaharvesters = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiaharvesters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiaharvesters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiaharvesters/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaHarvesterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Running reconciler...")

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
		return ctrl.Result{}, nil
	}
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		log.Error(err, fmt.Sprintf("ChiaHarvesterReconciler ChiaHarvester=%s unable to fetch ChiaHarvester resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiaharvesters[req.NamespacedName.String()]
	if !exists {
		chiaharvesters[req.NamespacedName.String()] = true
		metrics.ChiaHarvesters.Add(1.0)
	}

	// Assemble Peer Service
	peerSrv := assemblePeerService(harvester)
	if err := controllerutil.SetControllerReference(&harvester, &peerSrv, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to assemble harvester peer Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error assembling peer Service: %v", req.NamespacedName, err)
	}
	// Reconcile Peer Service
	res, err := kube.ReconcileService(ctx, r.Client, harvester.Spec.ChiaConfig.PeerService, peerSrv, true)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		return res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
	}

	// Assemble Daemon Service
	daemonSrv := assembleDaemonService(harvester)
	if err := controllerutil.SetControllerReference(&harvester, &daemonSrv, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to assemble harvester daemon Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error assembling daemon Service: %v", req.NamespacedName, err)
	}
	// Reconcile Daemon Service
	res, err = kube.ReconcileService(ctx, r.Client, harvester.Spec.ChiaConfig.DaemonService, daemonSrv, true)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		return res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
	}

	// Assemble RPC Service
	rpcSrv := assembleRPCService(harvester)
	if err := controllerutil.SetControllerReference(&harvester, &rpcSrv, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to assemble harvester RPC Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error assembling RPC Service: %v", req.NamespacedName, err)
	}
	// Reconcile RPC Service
	res, err = kube.ReconcileService(ctx, r.Client, harvester.Spec.ChiaConfig.RPCService, rpcSrv, true)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		return res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
	}

	// Assemble Chia-Exporter Service
	exporterSrv := assembleChiaExporterService(harvester)
	if err := controllerutil.SetControllerReference(&harvester, &exporterSrv, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to assemble harvester chia-exporter Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error assembling chia-exporter Service: %v", req.NamespacedName, err)
	}
	// Reconcile Chia-Exporter Service
	res, err = kube.ReconcileService(ctx, r.Client, harvester.Spec.ChiaExporterConfig.Service, exporterSrv, true)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		return res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
	}

	// Creates a persistent volume claim if the GenerateVolumeClaims setting was set to true
	if kube.ShouldMakeVolumeClaim(harvester.Spec.Storage) {
		pvc, err := assembleVolumeClaim(harvester)
		if err != nil {
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to assemble harvester PVC -- Check operator logs.")
			return reconcile.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
		}

		if pvc != nil {
			res, err = kube.ReconcilePersistentVolumeClaim(ctx, r.Client, harvester.Spec.Storage, *pvc)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to create harvester PVC -- Check operator logs.")
				return res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
			}
		} else {
			return reconcile.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s PVC could not be created", req.NamespacedName)
		}
	}

	// Assemble Deployment
	deploy := assembleDeployment(harvester)
	if err := controllerutil.SetControllerReference(&harvester, &deploy, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to assemble harvester Deployment -- Check operator logs.")
		return reconcile.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
	}
	// Reconcile Deployment
	res, err = kube.ReconcileDeployment(ctx, r.Client, deploy)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to create harvester Deployment -- Check operator logs.")
		return res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&harvester, corev1.EventTypeNormal, "Created", "Successfully created ChiaHarvester resources.")
	harvester.Status.Ready = true
	err = r.Status().Update(ctx, &harvester)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		log.Error(err, fmt.Sprintf("ChiaHarvesterReconciler ChiaHarvester=%s unable to update ChiaHarvester status", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaHarvesterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaHarvester{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
