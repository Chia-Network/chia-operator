/*
Copyright 2023 Chia Network Inc.
*/

package chiafarmer

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

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

// ChiaFarmerReconciler reconciles a ChiaFarmer object
type ChiaFarmerReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chiafarmers = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiafarmers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiafarmers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiafarmers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaFarmerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log))
	log.Info(fmt.Sprintf("ChiaFarmerReconciler ChiaFarmer=%s running reconciler...", req.NamespacedName.String()))

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
		return ctrl.Result{}, nil
	}
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
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
	if kube.ShouldMakeService(farmer.Spec.ChiaConfig.PeerService) {
		srv := r.assemblePeerService(ctx, farmer)
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&farmer, corev1.EventTypeWarning, "Failed", "Failed to create farmer peer Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaFarmerReconciler ChiaFarmer=%s encountered error reconciling farmer peer Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiafarmerNamePattern, farmer.Name),
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("ChiaFarmerReconciler ChiaFarmer=%s unable to GET ChiaFarmer peer Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				log.Error(err, fmt.Sprintf("ChiaFarmerReconciler ChiaFarmer=%s unable to DELETE ChiaFarmer peer Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(farmer.Spec.ChiaConfig.DaemonService) {
		srv := r.assembleDaemonService(ctx, farmer)
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&farmer, corev1.EventTypeWarning, "Failed", "Failed to create farmer daemon Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaFarmerReconciler ChiaFarmer=%s encountered error reconciling farmer daemon Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiafarmerNamePattern, farmer.Name) + "-daemon",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaFarmerReconciler ChiaFarmer=%s unable to GET ChiaFarmer daemon Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaFarmerReconciler ChiaFarmer=%s unable to DELETE ChiaFarmer daemon Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(farmer.Spec.ChiaConfig.RPCService) {
		srv := r.assembleRPCService(ctx, farmer)
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&farmer, corev1.EventTypeWarning, "Failed", "Failed to create farmer RPC Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaFarmerReconciler ChiaFarmer=%s encountered error reconciling farmer RPC Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiafarmerNamePattern, farmer.Name) + "-rpc",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaFarmerReconciler ChiaFarmer=%s unable to GET ChiaFarmer RPC Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaFarmerReconciler ChiaFarmer=%s unable to DELETE ChiaFarmer RPC Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(farmer.Spec.ChiaExporterConfig.Service) {
		srv := r.assembleChiaExporterService(ctx, farmer)
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&farmer, corev1.EventTypeWarning, "Failed", "Failed to create farmer metrics Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaFarmerReconciler ChiaFarmer=%s encountered error reconciling farmer chia-exporter Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiafarmerNamePattern, farmer.Name) + "-metrics",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaFarmerReconciler ChiaFarmer=%s unable to GET ChiaFarmer metrics Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaFarmerReconciler ChiaFarmer=%s unable to DELETE ChiaFarmer metrics Service resource", req.NamespacedName))
			}
		}
	}

	// Creates a persistent volume claim if the GenerateVolumeClaims setting was set to true
	if farmer.Spec.Storage != nil && farmer.Spec.Storage.ChiaRoot != nil && farmer.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil && farmer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.GenerateVolumeClaims {
		pvc, err := r.assembleVolumeClaim(ctx, farmer)
		if err != nil {
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&farmer, corev1.EventTypeWarning, "Failed", "Failed to create farmer PVC -- Check operator logs.")
			return reconcile.Result{}, fmt.Errorf("ChiaFarmerReconciler ChiaFarmer=%s encountered error scaffolding a generated PersistentVolumeClaim: %v", req.NamespacedName, err)
		}

		res, err := kube.ReconcilePersistentVolumeClaim(ctx, resourceReconciler, pvc)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&farmer, corev1.EventTypeWarning, "Failed", "Failed to create farmer PVC -- Check operator logs.")
			return *res, fmt.Errorf("ChiaFarmerReconciler ChiaFarmer=%s encountered error reconciling PersistentVolumeClaim: %v", req.NamespacedName, err)
		}
	}

	deploy := r.assembleDeployment(ctx, farmer)

	if err := controllerutil.SetControllerReference(&farmer, &deploy, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	res, err := kube.ReconcileDeployment(ctx, resourceReconciler, deploy)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&farmer, corev1.EventTypeWarning, "Failed", "Failed to create farmer Deployment -- Check operator logs.")
		return *res, fmt.Errorf("ChiaFarmerReconciler ChiaFarmer=%s encountered error reconciling farmer Deployment: %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&farmer, corev1.EventTypeNormal, "Created", "Successfully created ChiaFarmer resources.")
	farmer.Status.Ready = true
	err = r.Status().Update(ctx, &farmer)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		log.Error(err, fmt.Sprintf("ChiaFarmerReconciler ChiaFarmer=%s unable to update ChiaFarmer status", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaFarmerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaFarmer{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
