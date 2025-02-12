/*
Copyright 2024 Chia Network Inc.
*/

package chiadatalayer

import (
	"context"
	stdErrors "errors"
	"fmt"
	"strings"
	"time"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	"github.com/chia-network/chia-operator/internal/metrics"
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
)

// ChiaDataLayerReconciler reconciles a ChiaDataLayer object
type ChiaDataLayerReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chiadatalayers = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiadatalayers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiadatalayers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiadatalayers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaDataLayerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Running reconciler...")

	// Get the custom resource
	var datalayer k8schianetv1.ChiaDataLayer
	err := r.Get(ctx, req.NamespacedName, &datalayer)
	if err != nil && errors.IsNotFound(err) {
		// Remove this object from the map for tracking and subtract this CR's total metric by 1
		_, exists := chiadatalayers[req.NamespacedName.String()]
		if exists {
			delete(chiadatalayers, req.NamespacedName.String())
			metrics.ChiaDataLayers.Sub(1.0)
		}
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, "unable to fetch ChiaDataLayer resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiadatalayers[req.NamespacedName.String()]
	if !exists {
		chiadatalayers[req.NamespacedName.String()] = true
		metrics.ChiaDataLayers.Add(1.0)
	}

	// Check for ChiaNetwork, retrieve matching ConfigMap if specified
	networkData, err := kube.GetChiaNetworkData(ctx, r.Client, datalayer.Spec.ChiaConfig.CommonSpecChia, datalayer.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Assemble Daemon Service
	daemonSrv := assembleDaemonService(datalayer)
	if err := controllerutil.SetControllerReference(&datalayer, &daemonSrv, r.Scheme); err != nil {
		r.Recorder.Event(&datalayer, corev1.EventTypeWarning, "Failed", "Failed to assemble datalayer daemon Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("encountered error assembling daemon Service: %v", err)
	}
	// Reconcile Daemon Service
	res, err := kube.ReconcileService(ctx, r.Client, datalayer.Spec.ChiaConfig.DaemonService, daemonSrv, true)
	if err != nil {
		r.Recorder.Event(&datalayer, corev1.EventTypeWarning, "Failed", "Failed to reconcile datalayer daemon Service -- Check operator logs.")
		return res, err
	}

	// Assemble RPC Service
	rpcSrv := assembleRPCService(datalayer)
	if err := controllerutil.SetControllerReference(&datalayer, &rpcSrv, r.Scheme); err != nil {
		r.Recorder.Event(&datalayer, corev1.EventTypeWarning, "Failed", "Failed to assemble datalayer RPC Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("encountered error assembling RPC Service: %v", err)
	}
	// Reconcile RPC Service
	res, err = kube.ReconcileService(ctx, r.Client, datalayer.Spec.ChiaConfig.RPCService, rpcSrv, true)
	if err != nil {
		r.Recorder.Event(&datalayer, corev1.EventTypeWarning, "Failed", "Failed to reconcile datalayer RPC Service -- Check operator logs.")
		return res, err
	}

	// Assemble HTTP Service
	httpSrv := assembleDataLayerHTTPService(datalayer)
	if err := controllerutil.SetControllerReference(&datalayer, &httpSrv, r.Scheme); err != nil {
		r.Recorder.Event(&datalayer, corev1.EventTypeWarning, "Failed", "Failed to assemble datalayer HTTP Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("encountered error assembling HTTP Service: %v", err)
	}
	// Reconcile HTTP Service
	res, err = kube.ReconcileService(ctx, r.Client, datalayer.Spec.DataLayerHTTPConfig.Service, httpSrv, true)
	if err != nil {
		r.Recorder.Event(&datalayer, corev1.EventTypeWarning, "Failed", "Failed to reconcile datalayer HTTP Service -- Check operator logs.")
		return res, err
	}

	// Assemble Chia-Exporter Service
	exporterSrv := assembleChiaExporterService(datalayer)
	if err := controllerutil.SetControllerReference(&datalayer, &exporterSrv, r.Scheme); err != nil {
		r.Recorder.Event(&datalayer, corev1.EventTypeWarning, "Failed", "Failed to assemble datalayer chia-exporter Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("encountered error assembling chia-exporter Service: %v", err)
	}
	// Reconcile Chia-Exporter Service
	res, err = kube.ReconcileService(ctx, r.Client, datalayer.Spec.ChiaExporterConfig.Service, exporterSrv, true)
	if err != nil {
		r.Recorder.Event(&datalayer, corev1.EventTypeWarning, "Failed", "Failed to reconcile datalayer chia-exporter Service -- Check operator logs.")
		return res, err
	}

	// Creates a persistent volume claim if the GenerateVolumeClaims setting was set to true
	if kube.ShouldMakeChiaRootVolumeClaim(datalayer.Spec.Storage) {
		pvc, err := assembleChiaRootVolumeClaim(datalayer)
		if err != nil {
			r.Recorder.Event(&datalayer, corev1.EventTypeWarning, "Failed", "Failed to assemble datalayer PVC -- Check operator logs.")
			return reconcile.Result{}, err
		}

		if pvc != nil {
			res, err = kube.ReconcilePersistentVolumeClaim(ctx, r.Client, datalayer.Spec.Storage, *pvc)
			if err != nil {
				r.Recorder.Event(&datalayer, corev1.EventTypeWarning, "Failed", "Failed to create datalayer CHIA_ROOT PVC -- Check operator logs.")
				return res, err
			}
		} else {
			return reconcile.Result{}, stdErrors.New("CHIA_ROOT PVC could not be created")
		}
	}

	// Creates a persistent volume claim if the GenerateVolumeClaims setting was set to true
	if kube.ShouldMakeDataLayerServerFilesVolumeClaim(datalayer.Spec.Storage) {
		pvc, err := assembleDataLayerFilesVolumeClaim(datalayer)
		if err != nil {
			r.Recorder.Event(&datalayer, corev1.EventTypeWarning, "Failed", "Failed to assemble datalayer server files PVC -- Check operator logs.")
			return reconcile.Result{}, err
		}

		if pvc != nil {
			res, err = kube.ReconcilePersistentVolumeClaim(ctx, r.Client, datalayer.Spec.Storage, *pvc)
			if err != nil {
				r.Recorder.Event(&datalayer, corev1.EventTypeWarning, "Failed", "Failed to create datalayer server files PVC -- Check operator logs.")
				return res, err
			}
		} else {
			return reconcile.Result{}, stdErrors.New("server files PVC could not be created")
		}
	}

	// Assemble Deployment
	deploy, err := assembleDeployment(ctx, datalayer, networkData)
	if err != nil {
		r.Recorder.Event(&datalayer, corev1.EventTypeWarning, "Failed", "Failed to assemble datalayer Deployment -- Check operator logs.")
		return reconcile.Result{}, err
	}
	if err := controllerutil.SetControllerReference(&datalayer, &deploy, r.Scheme); err != nil {
		r.Recorder.Event(&datalayer, corev1.EventTypeWarning, "Failed", "Failed to assemble datalayer Deployment -- Check operator logs.")
		return reconcile.Result{}, err
	}
	// Reconcile Deployment
	res, err = kube.ReconcileDeployment(ctx, r.Client, deploy)
	if err != nil {
		r.Recorder.Event(&datalayer, corev1.EventTypeWarning, "Failed", "Failed to create datalayer Deployment -- Check operator logs.")
		return res, err
	}

	// Update CR status
	r.Recorder.Event(&datalayer, corev1.EventTypeNormal, "Created", "Successfully created ChiaDataLayer resources.")
	datalayer.Status.Ready = true
	err = r.Status().Update(ctx, &datalayer)
	if err != nil {
		if strings.Contains(err.Error(), kube.ObjectModifiedTryAgainError) {
			return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
		}
		log.Error(err, "unable to update ChiaDataLayer status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaDataLayerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaDataLayer{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
