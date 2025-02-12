/*
Copyright 2024 Chia Network Inc.
*/

package chiacrawler

import (
	"context"
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

// ChiaCrawlerReconciler reconciles a ChiaCrawler object
type ChiaCrawlerReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chiacrawlers = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiacrawlers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiacrawlers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiacrawlers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaCrawlerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Running reconciler...")

	// Get the custom resource
	var crawler k8schianetv1.ChiaCrawler
	err := r.Get(ctx, req.NamespacedName, &crawler)
	if err != nil && errors.IsNotFound(err) {
		// Remove this object from the map for tracking and subtract this CR's total metric by 1
		_, exists := chiacrawlers[req.NamespacedName.String()]
		if exists {
			delete(chiacrawlers, req.NamespacedName.String())
			metrics.ChiaCrawlers.Sub(1.0)
		}
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, fmt.Sprintf("ChiaCrawlerReconciler ChiaCrawler=%s unable to fetch ChiaCrawler resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiacrawlers[req.NamespacedName.String()]
	if !exists {
		chiacrawlers[req.NamespacedName.String()] = true
		metrics.ChiaCrawlers.Add(1.0)
	}

	// Check for ChiaNetwork, retrieve matching ConfigMap if specified
	networkData, err := kube.GetChiaNetworkData(ctx, r.Client, crawler.Spec.ChiaConfig.CommonSpecChia, crawler.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Get the full_node Port and handle the error one time instead of in every function that needs it
	fullNodePort, err := kube.GetFullNodePort(crawler.Spec.ChiaConfig.CommonSpecChia, networkData)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("encountered error retrieving the full_node Port to use: %v", err)
	}

	// Assemble Peer Service
	peerSrv := assemblePeerService(crawler, fullNodePort)
	if err := controllerutil.SetControllerReference(&crawler, &peerSrv, r.Scheme); err != nil {
		r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to assemble crawler peer Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s encountered error assembling peer Service: %v", req.NamespacedName, err)
	}
	// Reconcile Peer Service
	res, err := kube.ReconcileService(ctx, r.Client, crawler.Spec.ChiaConfig.PeerService, peerSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s %v", req.NamespacedName, err)
	}

	// Assemble All Service
	allSrv := assembleAllService(crawler, fullNodePort)
	if err := controllerutil.SetControllerReference(&crawler, &allSrv, r.Scheme); err != nil {
		r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to assemble crawler all-port Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s encountered error assembling all-port Service: %v", req.NamespacedName, err)
	}
	// Reconcile All Service
	res, err = kube.ReconcileService(ctx, r.Client, crawler.Spec.ChiaConfig.AllService, allSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s %v", req.NamespacedName, err)
	}

	// Assemble Daemon Service
	daemonSrv := assembleDaemonService(crawler)
	if err := controllerutil.SetControllerReference(&crawler, &daemonSrv, r.Scheme); err != nil {
		r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to assemble crawler daemon Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s encountered error assembling daemon Service: %v", req.NamespacedName, err)
	}
	// Reconcile Daemon Service
	res, err = kube.ReconcileService(ctx, r.Client, crawler.Spec.ChiaConfig.DaemonService, daemonSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s %v", req.NamespacedName, err)
	}

	// Assemble RPC Service
	rpcSrv := assembleRPCService(crawler)
	if err := controllerutil.SetControllerReference(&crawler, &rpcSrv, r.Scheme); err != nil {
		r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to assemble crawler RPC Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s encountered error assembling RPC Service: %v", req.NamespacedName, err)
	}
	// Reconcile RPC Service
	res, err = kube.ReconcileService(ctx, r.Client, crawler.Spec.ChiaConfig.RPCService, rpcSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s %v", req.NamespacedName, err)
	}

	// Assemble Chia-Exporter Service
	exporterSrv := assembleChiaExporterService(crawler)
	if err := controllerutil.SetControllerReference(&crawler, &exporterSrv, r.Scheme); err != nil {
		r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to assemble crawler chia-exporter Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s encountered error assembling chia-exporter Service: %v", req.NamespacedName, err)
	}
	// Reconcile Chia-Exporter Service
	res, err = kube.ReconcileService(ctx, r.Client, crawler.Spec.ChiaExporterConfig.Service, exporterSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s %v", req.NamespacedName, err)
	}

	// Creates a persistent volume claim if the GenerateVolumeClaims setting was set to true
	if kube.ShouldMakeChiaRootVolumeClaim(crawler.Spec.Storage) {
		pvc, err := assembleVolumeClaim(crawler)
		if err != nil {
			r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to assemble crawler PVC -- Check operator logs.")
			return reconcile.Result{}, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s %v", req.NamespacedName, err)
		}

		if pvc != nil {
			res, err = kube.ReconcilePersistentVolumeClaim(ctx, r.Client, crawler.Spec.Storage, *pvc)
			if err != nil {
				r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to create crawler PVC -- Check operator logs.")
				return res, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s %v", req.NamespacedName, err)
			}
		} else {
			return reconcile.Result{}, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s PVC could not be created", req.NamespacedName)
		}
	}

	// Assemble Deployment
	deploy, err := assembleDeployment(crawler, fullNodePort, networkData)
	if err != nil {
		r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to assemble crawler Deployment -- Check operator logs.")
		return reconcile.Result{}, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s %v", req.NamespacedName, err)
	}
	if err := controllerutil.SetControllerReference(&crawler, &deploy, r.Scheme); err != nil {
		r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to assemble crawler Deployment -- Check operator logs.")
		return reconcile.Result{}, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s %v", req.NamespacedName, err)
	}
	// Reconcile Deployment
	res, err = kube.ReconcileDeployment(ctx, r.Client, deploy)
	if err != nil {
		r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to create crawler Deployment -- Check operator logs.")
		return res, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&crawler, corev1.EventTypeNormal, "Created", "Successfully created ChiaCrawler resources.")
	crawler.Status.Ready = true
	err = r.Status().Update(ctx, &crawler)
	if err != nil {
		if strings.Contains(err.Error(), kube.ObjectModifiedTryAgainError) {
			return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
		}
		log.Error(err, fmt.Sprintf("ChiaCrawlerReconciler ChiaCrawler=%s unable to update ChiaCrawler status", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaCrawlerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaCrawler{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
