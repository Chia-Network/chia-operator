/*
Copyright 2023 Chia Network Inc.
*/

package chianode

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

// ChiaNodeReconciler reconciles a ChiaNode object
type ChiaNodeReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chianodes = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chianodes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chianodes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chianodes/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaNodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Running reconciler...")

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

	// Check for ChiaNetwork, retrieve matching ConfigMap if specified
	networkData, err := kube.GetChiaNetworkData(ctx, r.Client, node.Spec.ChiaConfig.CommonSpecChia, node.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Get the full_node Port and handle the error one time instead of in every function that needs it
	fullNodePort, err := kube.GetFullNodePort(node.Spec.ChiaConfig.CommonSpecChia, networkData)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("encountered error retrieving the full_node Port to use: %v", err)
	}

	// Assemble Peer Service
	peerSrv := assemblePeerService(node, fullNodePort)
	if err := controllerutil.SetControllerReference(&node, &peerSrv, r.Scheme); err != nil {
		r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to assemble node peer Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error assembling peer Service: %v", req.NamespacedName, err)
	}
	// Reconcile Peer Service
	res, err := kube.ReconcileService(ctx, r.Client, node.Spec.ChiaConfig.PeerService, peerSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s %v", req.NamespacedName, err)
	}

	// Assemble All Service
	allSrv := assembleAllService(node, fullNodePort)
	if err := controllerutil.SetControllerReference(&node, &allSrv, r.Scheme); err != nil {
		r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to assemble node all-port Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error assembling all-port Service: %v", req.NamespacedName, err)
	}
	// Reconcile All Service
	res, err = kube.ReconcileService(ctx, r.Client, node.Spec.ChiaConfig.AllService, allSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s %v", req.NamespacedName, err)
	}

	// Assemble Headless Peer Service
	headlessPeerSrv := assembleHeadlessPeerService(node, fullNodePort)
	if err := controllerutil.SetControllerReference(&node, &headlessPeerSrv, r.Scheme); err != nil {
		r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to assemble node headless peer Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error assembling headless peer Service: %v", req.NamespacedName, err)
	}
	// Reconcile Headless Peer Service
	res, err = kube.ReconcileService(ctx, r.Client, node.Spec.ChiaConfig.PeerService, headlessPeerSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s %v", req.NamespacedName, err)
	}

	// Assemble Local Peer Service
	localPeerSrv := assembleLocalPeerService(node, fullNodePort)
	if err := controllerutil.SetControllerReference(&node, &localPeerSrv, r.Scheme); err != nil {
		r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to assemble node local peer Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error assembling local peer Service: %v", req.NamespacedName, err)
	}
	// Reconcile Local Peer Service
	res, err = kube.ReconcileService(ctx, r.Client, node.Spec.ChiaConfig.PeerService, localPeerSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s %v", req.NamespacedName, err)
	}

	// Assemble Daemon Service
	daemonSrv := assembleDaemonService(node)
	if err := controllerutil.SetControllerReference(&node, &daemonSrv, r.Scheme); err != nil {
		r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to assemble node daemon Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error assembling daemon Service: %v", req.NamespacedName, err)
	}
	// Reconcile Daemon Service
	res, err = kube.ReconcileService(ctx, r.Client, node.Spec.ChiaConfig.DaemonService, daemonSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s %v", req.NamespacedName, err)
	}

	// Assemble RPC Service
	rpcSrv := assembleRPCService(node)
	if err := controllerutil.SetControllerReference(&node, &rpcSrv, r.Scheme); err != nil {
		r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to assemble node RPC Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error assembling RPC Service: %v", req.NamespacedName, err)
	}
	// Reconcile RPC Service
	res, err = kube.ReconcileService(ctx, r.Client, node.Spec.ChiaConfig.RPCService, rpcSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s %v", req.NamespacedName, err)
	}

	// Assemble Chia-Exporter Service
	exporterSrv := assembleChiaExporterService(node)
	if err := controllerutil.SetControllerReference(&node, &exporterSrv, r.Scheme); err != nil {
		r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to assemble node chia-exporter Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error assembling chia-exporter Service: %v", req.NamespacedName, err)
	}
	// Reconcile Chia-Exporter Service
	res, err = kube.ReconcileService(ctx, r.Client, node.Spec.ChiaExporterConfig.Service, exporterSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s %v", req.NamespacedName, err)
	}

	// Assemble Chia-Healthcheck Service
	healthcheckSrv := assembleChiaHealthcheckService(node)
	if err := controllerutil.SetControllerReference(&node, &healthcheckSrv, r.Scheme); err != nil {
		r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to assemble node chia-healthcheck Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error assembling chia-healthcheck Service: %v", req.NamespacedName, err)
	}
	// Reconcile Chia-Healthcheck Service
	if !kube.ShouldRollIntoMainPeerService(node.Spec.ChiaHealthcheckConfig.Service) {
		res, err = kube.ReconcileService(ctx, r.Client, node.Spec.ChiaHealthcheckConfig.Service, healthcheckSrv, false)
		if err != nil {
			return res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s %v", req.NamespacedName, err)
		}
	}

	// Assemble StatefulSet
	stateful, err := assembleStatefulset(ctx, node, fullNodePort, networkData)
	if err != nil {
		r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to assemble node Deployment -- Check operator logs.")
		return reconcile.Result{}, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s %v", req.NamespacedName, err)
	}
	if err := controllerutil.SetControllerReference(&node, &stateful, r.Scheme); err != nil {
		r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to assemble seeder Statefulset -- Check operator logs.")
		return reconcile.Result{}, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s %v", req.NamespacedName, err)
	}
	// Reconcile StatefulSet
	res, err = kube.ReconcileStatefulset(ctx, r.Client, stateful)
	if err != nil {
		r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to create seeder Statefulset -- Check operator logs.")
		return res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&node, corev1.EventTypeNormal, "Created", "Successfully created ChiaNode resources.")
	node.Status.Ready = true
	err = r.Status().Update(ctx, &node)
	if err != nil {
		if strings.Contains(err.Error(), kube.ObjectModifiedTryAgainError) {
			return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
		}
		log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to update ChiaNode status", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaNodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaNode{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Watches(
			&corev1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(r.handleChiaNetworks),
		).
		Complete(r)
}

func (r *ChiaNodeReconciler) handleChiaNetworks(ctx context.Context, obj client.Object) []reconcile.Request {
	listOps := &client.ListOptions{
		Namespace: obj.GetNamespace(),
	}
	list := &k8schianetv1.ChiaNodeList{}
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
