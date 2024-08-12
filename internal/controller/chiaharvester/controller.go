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
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log))
	log.Info(fmt.Sprintf("ChiaHarvesterReconciler ChiaHarvester=%s running reconciler...", req.NamespacedName.String()))

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

	// Reconcile ChiaHarvester owned objects
	if kube.ShouldMakeService(harvester.Spec.ChiaConfig.PeerService, true) {
		srv := assemblePeerService(harvester)
		if err := controllerutil.SetControllerReference(&harvester, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to create harvester peer Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error reconciling harvester peer Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiaharvesterNamePattern, harvester.Name),
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("ChiaHarvesterReconciler ChiaHarvester=%s unable to GET ChiaHarvester peer Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				log.Error(err, fmt.Sprintf("ChiaHarvesterReconciler ChiaHarvester=%s unable to DELETE ChiaHarvester peer Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(harvester.Spec.ChiaConfig.DaemonService, true) {
		srv := assembleDaemonService(harvester)
		if err := controllerutil.SetControllerReference(&harvester, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to create harvester daemon Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error reconciling harvester daemon Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiaharvesterNamePattern, harvester.Name) + "-daemon",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaHarvesterReconciler ChiaHarvester=%s unable to GET ChiaHarvester daemon Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaHarvesterReconciler ChiaHarvester=%s unable to DELETE ChiaHarvester daemon Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(harvester.Spec.ChiaConfig.RPCService, true) {
		srv := assembleRPCService(harvester)
		if err := controllerutil.SetControllerReference(&harvester, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to create harvester RPC Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error reconciling harvester RPC Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiaharvesterNamePattern, harvester.Name) + "-rpc",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaHarvesterReconciler ChiaHarvester=%s unable to GET ChiaHarvester RPC Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaHarvesterReconciler ChiaHarvester=%s unable to DELETE ChiaHarvester RPC Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(harvester.Spec.ChiaExporterConfig.Service, true) {
		srv := assembleChiaExporterService(harvester)
		if err := controllerutil.SetControllerReference(&harvester, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to create harvester metrics Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error reconciling harvester metrics Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiaharvesterNamePattern, harvester.Name) + "-metrics",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaHarvesterReconciler ChiaHarvester=%s unable to GET ChiaHarvester metrics Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaHarvesterReconciler ChiaHarvester=%s unable to DELETE ChiaHarvester metrics Service resource", req.NamespacedName))
			}
		}
	}

	// Creates a persistent volume claim if the GenerateVolumeClaims setting was set to true
	if harvester.Spec.Storage != nil && harvester.Spec.Storage.ChiaRoot != nil && harvester.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil && harvester.Spec.Storage.ChiaRoot.PersistentVolumeClaim.GenerateVolumeClaims {
		pvc, err := assembleVolumeClaim(harvester)
		if err != nil {
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to create harvester PVC -- Check operator logs.")
			return reconcile.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error scaffolding a generated PersistentVolumeClaim: %v", req.NamespacedName, err)
		}

		res, err := kube.ReconcilePersistentVolumeClaim(ctx, resourceReconciler, pvc)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to create harvester PVC -- Check operator logs.")
			return *res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error reconciling PersistentVolumeClaim: %v", req.NamespacedName, err)
		}
	}

	deploy := assembleDeployment(harvester)

	if err := controllerutil.SetControllerReference(&harvester, &deploy, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	res, err := kube.ReconcileDeployment(ctx, resourceReconciler, deploy)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to create harvester Deployment -- Check operator logs.")
		return *res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error reconciling harvester Deployment: %v", req.NamespacedName, err)
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
