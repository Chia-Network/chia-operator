/*
Copyright 2023 Chia Network Inc.
*/

package chiaharvester

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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
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
		_, exists := chiaharvesters[req.String()]
		if exists {
			delete(chiaharvesters, req.String())
			metrics.ChiaHarvesters.Sub(1.0)
		}
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, fmt.Sprintf("ChiaHarvesterReconciler ChiaHarvester=%s unable to fetch ChiaHarvester resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiaharvesters[req.String()]
	if !exists {
		chiaharvesters[req.String()] = true
		metrics.ChiaHarvesters.Add(1.0)
	}

	// Check for ChiaNetwork, retrieve matching ConfigMap if specified
	networkData, err := kube.GetChiaNetworkData(ctx, r.Client, harvester.Spec.ChiaConfig.CommonSpecChia, harvester.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Assemble Peer Service
	peerSrv := assemblePeerService(harvester)
	if err := controllerutil.SetControllerReference(&harvester, &peerSrv, r.Scheme); err != nil {
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to assemble harvester peer Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error assembling peer Service: %v", req.NamespacedName, err)
	}
	// Reconcile Peer Service
	res, err := kube.ReconcileService(ctx, r.Client, harvester.Spec.ChiaConfig.PeerService, peerSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
	}

	// Assemble All Service
	allSrv := assembleAllService(harvester)
	if err := controllerutil.SetControllerReference(&harvester, &allSrv, r.Scheme); err != nil {
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to assemble harvester all-port Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error assembling all-port Service: %v", req.NamespacedName, err)
	}
	// Reconcile All Service
	res, err = kube.ReconcileService(ctx, r.Client, harvester.Spec.ChiaConfig.AllService, allSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
	}

	// Assemble Daemon Service
	daemonSrv := assembleDaemonService(harvester)
	if err := controllerutil.SetControllerReference(&harvester, &daemonSrv, r.Scheme); err != nil {
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to assemble harvester daemon Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error assembling daemon Service: %v", req.NamespacedName, err)
	}
	// Reconcile Daemon Service
	res, err = kube.ReconcileService(ctx, r.Client, harvester.Spec.ChiaConfig.DaemonService, daemonSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
	}

	// Assemble RPC Service
	rpcSrv := assembleRPCService(harvester)
	if err := controllerutil.SetControllerReference(&harvester, &rpcSrv, r.Scheme); err != nil {
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to assemble harvester RPC Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error assembling RPC Service: %v", req.NamespacedName, err)
	}
	// Reconcile RPC Service
	res, err = kube.ReconcileService(ctx, r.Client, harvester.Spec.ChiaConfig.RPCService, rpcSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
	}

	// Assemble Chia-Exporter Service
	exporterSrv := assembleChiaExporterService(harvester)
	if err := controllerutil.SetControllerReference(&harvester, &exporterSrv, r.Scheme); err != nil {
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to assemble harvester chia-exporter Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error assembling chia-exporter Service: %v", req.NamespacedName, err)
	}
	// Reconcile Chia-Exporter Service
	res, err = kube.ReconcileService(ctx, r.Client, harvester.Spec.ChiaExporterConfig.Service, exporterSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
	}

	// Assemble Chia-Healthcheck Service
	healthcheckSrv := assembleChiaHealthcheckService(harvester)
	if err := controllerutil.SetControllerReference(&harvester, &healthcheckSrv, r.Scheme); err != nil {
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to assemble harvester chia-healthcheck Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s encountered error assembling chia-healthcheck Service: %v", req.NamespacedName, err)
	}
	// Reconcile Chia-Healthcheck Service
	if !kube.ShouldRollIntoMainPeerService(harvester.Spec.ChiaHealthcheckConfig.Service) {
		res, err = kube.ReconcileService(ctx, r.Client, harvester.Spec.ChiaHealthcheckConfig.Service, healthcheckSrv, false)
		if err != nil {
			return res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
		}
	}

	// Creates a persistent volume claim if the GenerateVolumeClaims setting was set to true
	if kube.ShouldMakeChiaRootVolumeClaim(harvester.Spec.Storage) {
		pvc, err := assembleVolumeClaim(harvester)
		if err != nil {
			r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to assemble harvester PVC -- Check operator logs.")
			return reconcile.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
		}

		if pvc != nil {
			res, err = kube.ReconcilePersistentVolumeClaim(ctx, r.Client, harvester.Spec.Storage, *pvc)
			if err != nil {
				r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to create harvester PVC -- Check operator logs.")
				return res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
			}
		} else {
			return reconcile.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s PVC could not be created", req.NamespacedName)
		}
	}

	// Assemble Deployment
	deploy, err := assembleDeployment(harvester, networkData)
	if err != nil {
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to assemble harvester Deployment -- Check operator logs.")
		return reconcile.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
	}
	if err := controllerutil.SetControllerReference(&harvester, &deploy, r.Scheme); err != nil {
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to assemble harvester Deployment -- Check operator logs.")
		return reconcile.Result{}, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
	}
	// Reconcile Deployment
	res, err = kube.ReconcileDeployment(ctx, r.Client, deploy)
	if err != nil {
		r.Recorder.Event(&harvester, corev1.EventTypeWarning, "Failed", "Failed to create harvester Deployment -- Check operator logs.")
		return res, fmt.Errorf("ChiaHarvesterReconciler ChiaHarvester=%s %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&harvester, corev1.EventTypeNormal, "Created", "Successfully created ChiaHarvester resources.")
	harvester.Status.Ready = true
	err = r.Status().Update(ctx, &harvester)
	if err != nil {
		if strings.Contains(err.Error(), kube.ObjectModifiedTryAgainError) {
			return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
		}
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
		Watches(
			&corev1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(r.handleChiaNetworks),
		).
		Complete(r)
}

func (r *ChiaHarvesterReconciler) handleChiaNetworks(ctx context.Context, obj client.Object) []reconcile.Request {
	listOps := &client.ListOptions{
		Namespace: obj.GetNamespace(),
	}
	list := &k8schianetv1.ChiaHarvesterList{}
	err := r.List(ctx, list, listOps)
	if err != nil {
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, len(list.Items))
	for i, item := range list.Items {
		chiaNetwork := item.Spec.ChiaConfig.ChiaNetwork
		if chiaNetwork != nil && *chiaNetwork == obj.GetName() {
			requests[i] = reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      item.GetName(),
					Namespace: item.GetNamespace(),
				},
			}
		}
	}
	return requests
}
