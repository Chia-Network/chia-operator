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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log))
	log.Info(fmt.Sprintf("ChiaTimelordController ChiaTimelord=%s running reconciler...", req.NamespacedName.String()))

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

	// Reconcile ChiaTimelord owned objects
	if kube.ShouldMakeService(tl.Spec.ChiaConfig.PeerService, true) {
		srv := assemblePeerService(tl)
		if err := controllerutil.SetControllerReference(&tl, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&tl, corev1.EventTypeWarning, "Failed", "Failed to create timelord peer Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaFarmerReconciler ChiaFarmer=%s encountered error reconciling timelord peer Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiatimelordNamePattern, tl.Name),
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("ChiaTimelordReconciler ChiaTimelord=%s unable to GET ChiaTimelord peer Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				log.Error(err, fmt.Sprintf("ChiaTimelordReconciler ChiaTimelord=%s unable to DELETE ChiaTimelord peer Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(tl.Spec.ChiaConfig.DaemonService, true) {
		srv := assembleDaemonService(tl)
		if err := controllerutil.SetControllerReference(&tl, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&tl, corev1.EventTypeWarning, "Failed", "Failed to create timelord daemon Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaFarmerReconciler ChiaFarmer=%s encountered error reconciling timelord daemon Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiatimelordNamePattern, tl.Name) + "-daemon",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaTimelordReconciler ChiaTimelord=%s unable to GET ChiaTimelord daemon Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaTimelordReconciler ChiaTimelord=%s unable to DELETE ChiaTimelord daemon Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(tl.Spec.ChiaConfig.RPCService, true) {
		srv := assembleRPCService(tl)
		if err := controllerutil.SetControllerReference(&tl, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&tl, corev1.EventTypeWarning, "Failed", "Failed to create timelord RPC Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaFarmerReconciler ChiaFarmer=%s encountered error reconciling timelord RPC Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiatimelordNamePattern, tl.Name) + "-rpc",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaTimelordReconciler ChiaTimelord=%s unable to GET ChiaTimelord RPC Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaTimelordReconciler ChiaTimelord=%s unable to DELETE ChiaTimelord RPC Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(tl.Spec.ChiaExporterConfig.Service, true) {
		srv := assembleChiaExporterService(tl)
		if err := controllerutil.SetControllerReference(&tl, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&tl, corev1.EventTypeWarning, "Failed", "Failed to create timelord metrics Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaFarmerReconciler ChiaFarmer=%s encountered error reconciling timelord chia-exporter Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiatimelordNamePattern, tl.Name) + "-metrics",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaTimelordReconciler ChiaTimelord=%s unable to GET ChiaTimelord metrics Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaTimelordReconciler ChiaTimelord=%s unable to DELETE ChiaTimelord metrics Service resource", req.NamespacedName))
			}
		}
	}

	// Creates a persistent volume claim if the GenerateVolumeClaims setting was set to true
	if tl.Spec.Storage != nil && tl.Spec.Storage.ChiaRoot != nil && tl.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil && tl.Spec.Storage.ChiaRoot.PersistentVolumeClaim.GenerateVolumeClaims {
		pvc, err := assembleVolumeClaim(tl)
		if err != nil {
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&tl, corev1.EventTypeWarning, "Failed", "Failed to create timelord PVC -- Check operator logs.")
			return reconcile.Result{}, fmt.Errorf("ChiaTimelordReconciler ChiaTimelord=%s encountered error scaffolding a generated PersistentVolumeClaim: %v", req.NamespacedName, err)
		}

		res, err := kube.ReconcilePersistentVolumeClaim(ctx, resourceReconciler, pvc)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&tl, corev1.EventTypeWarning, "Failed", "Failed to create timelord PVC -- Check operator logs.")
			return *res, fmt.Errorf("ChiaTimelordReconciler ChiaTimelord=%s encountered error reconciling PersistentVolumeClaim: %v", req.NamespacedName, err)
		}
	}

	deploy := assembleDeployment(tl)

	if err := controllerutil.SetControllerReference(&tl, &deploy, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	res, err := kube.ReconcileDeployment(ctx, resourceReconciler, deploy)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&tl, corev1.EventTypeWarning, "Failed", "Failed to create timelord Deployment -- Check operator logs.")
		return *res, fmt.Errorf("ChiaTimelordController ChiaTimelord=%s encountered error reconciling node StatefulSet: %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&tl, corev1.EventTypeNormal, "Created", "Successfully created ChiaTimelord resources.")
	tl.Status.Ready = true
	err = r.Status().Update(ctx, &tl)
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
