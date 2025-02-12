/*
Copyright 2024 Chia Network Inc.
*/

package chiaintroducer

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

// ChiaIntroducerReconciler reconciles a ChiaIntroducer object
type ChiaIntroducerReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chiaintroducers = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiaintroducers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiaintroducers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiaintroducers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaIntroducerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Running reconciler...")

	// Get the custom resource
	var introducer k8schianetv1.ChiaIntroducer
	err := r.Get(ctx, req.NamespacedName, &introducer)
	if err != nil && errors.IsNotFound(err) {
		// Remove this object from the map for tracking and subtract this CR's total metric by 1
		_, exists := chiaintroducers[req.NamespacedName.String()]
		if exists {
			delete(chiaintroducers, req.NamespacedName.String())
			metrics.ChiaIntroducers.Sub(1.0)
		}
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, fmt.Sprintf("ChiaIntroducerReconciler ChiaIntroducer=%s unable to fetch ChiaIntroducer resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiaintroducers[req.NamespacedName.String()]
	if !exists {
		chiaintroducers[req.NamespacedName.String()] = true
		metrics.ChiaIntroducers.Add(1.0)
	}

	// Check for ChiaNetwork, retrieve matching ConfigMap if specified
	networkData, err := kube.GetChiaNetworkData(ctx, r.Client, introducer.Spec.ChiaConfig.CommonSpecChia, introducer.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Get the full_node Port and handle the error one time instead of in every function that needs it
	fullNodePort, err := kube.GetFullNodePort(introducer.Spec.ChiaConfig.CommonSpecChia, networkData)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("encountered error retrieving the full_node Port to use: %v", err)
	}

	// Assemble Peer Service
	peerSrv := assemblePeerService(introducer, fullNodePort)
	if err := controllerutil.SetControllerReference(&introducer, &peerSrv, r.Scheme); err != nil {
		r.Recorder.Event(&introducer, corev1.EventTypeWarning, "Failed", "Failed to assemble introducer peer Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s encountered error assembling peer Service: %v", req.NamespacedName, err)
	}
	// Reconcile Peer Service
	res, err := kube.ReconcileService(ctx, r.Client, introducer.Spec.ChiaConfig.PeerService, peerSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s %v", req.NamespacedName, err)
	}

	// Assemble All Service
	allSrv := assembleAllService(introducer, fullNodePort)
	if err := controllerutil.SetControllerReference(&introducer, &allSrv, r.Scheme); err != nil {
		r.Recorder.Event(&introducer, corev1.EventTypeWarning, "Failed", "Failed to assemble introducer all-port Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s encountered error assembling all-port Service: %v", req.NamespacedName, err)
	}
	// Reconcile All Service
	res, err = kube.ReconcileService(ctx, r.Client, introducer.Spec.ChiaConfig.AllService, allSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s %v", req.NamespacedName, err)
	}

	// Assemble Daemon Service
	daemonSrv := assembleDaemonService(introducer)
	if err := controllerutil.SetControllerReference(&introducer, &daemonSrv, r.Scheme); err != nil {
		r.Recorder.Event(&introducer, corev1.EventTypeWarning, "Failed", "Failed to assemble introducer daemon Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s encountered error assembling daemon Service: %v", req.NamespacedName, err)
	}
	// Reconcile Daemon Service
	res, err = kube.ReconcileService(ctx, r.Client, introducer.Spec.ChiaConfig.DaemonService, daemonSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s %v", req.NamespacedName, err)
	}

	// Assemble Chia-Exporter Service
	exporterSrv := assembleChiaExporterService(introducer)
	if err := controllerutil.SetControllerReference(&introducer, &exporterSrv, r.Scheme); err != nil {
		r.Recorder.Event(&introducer, corev1.EventTypeWarning, "Failed", "Failed to assemble introducer chia-exporter Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s encountered error assembling chia-exporter Service: %v", req.NamespacedName, err)
	}
	// Reconcile Chia-Exporter Service
	res, err = kube.ReconcileService(ctx, r.Client, introducer.Spec.ChiaExporterConfig.Service, exporterSrv, true)
	if err != nil {
		return res, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s %v", req.NamespacedName, err)
	}

	// Creates a persistent volume claim if the GenerateVolumeClaims setting was set to true
	if kube.ShouldMakeChiaRootVolumeClaim(introducer.Spec.Storage) {
		pvc, err := assembleVolumeClaim(introducer)
		if err != nil {
			r.Recorder.Event(&introducer, corev1.EventTypeWarning, "Failed", "Failed to assemble introducer PVC -- Check operator logs.")
			return reconcile.Result{}, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s %v", req.NamespacedName, err)
		}

		if pvc != nil {
			res, err = kube.ReconcilePersistentVolumeClaim(ctx, r.Client, introducer.Spec.Storage, *pvc)
			if err != nil {
				r.Recorder.Event(&introducer, corev1.EventTypeWarning, "Failed", "Failed to create introducer PVC -- Check operator logs.")
				return res, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s %v", req.NamespacedName, err)
			}
		} else {
			return reconcile.Result{}, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s PVC could not be created", req.NamespacedName)
		}
	}

	// Assemble Deployment
	deploy, err := assembleDeployment(introducer, fullNodePort, networkData)
	if err != nil {
		r.Recorder.Event(&introducer, corev1.EventTypeWarning, "Failed", "Failed to assemble introducer Deployment -- Check operator logs.")
		return reconcile.Result{}, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s %v", req.NamespacedName, err)
	}
	if err := controllerutil.SetControllerReference(&introducer, &deploy, r.Scheme); err != nil {
		r.Recorder.Event(&introducer, corev1.EventTypeWarning, "Failed", "Failed to assemble introducer Deployment -- Check operator logs.")
		return reconcile.Result{}, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s %v", req.NamespacedName, err)
	}
	// Reconcile Deployment
	res, err = kube.ReconcileDeployment(ctx, r.Client, deploy)
	if err != nil {
		r.Recorder.Event(&introducer, corev1.EventTypeWarning, "Failed", "Failed to create introducer Deployment -- Check operator logs.")
		return res, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&introducer, corev1.EventTypeNormal, "Created", "Successfully created ChiaIntroducer resources.")
	introducer.Status.Ready = true
	err = r.Status().Update(ctx, &introducer)
	if err != nil {
		if strings.Contains(err.Error(), kube.ObjectModifiedTryAgainError) {
			return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
		}
		log.Error(err, fmt.Sprintf("ChiaIntroducerReconciler ChiaIntroducer=%s unable to update ChiaIntroducer status", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaIntroducerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaIntroducer{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Watches(
			&corev1.ConfigMap{},
			handler.EnqueueRequestsFromMapFunc(r.handleChiaNetworks),
		).
		Complete(r)
}

func (r *ChiaIntroducerReconciler) handleChiaNetworks(ctx context.Context, obj client.Object) []reconcile.Request {
	listOps := &client.ListOptions{
		Namespace: obj.GetNamespace(),
	}
	list := &k8schianetv1.ChiaIntroducerList{}
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
