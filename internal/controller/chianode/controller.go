/*
Copyright 2023 Chia Network Inc.
*/

package chianode

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
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaNodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log))
	log.Info(fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s running reconciler...", req.NamespacedName.String()))

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
		metrics.OperatorErrors.Add(1.0)
		log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to fetch ChiaNode resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chianodes[req.NamespacedName.String()]
	if !exists {
		chianodes[req.NamespacedName.String()] = true
		metrics.ChiaNodes.Add(1.0)
	}

	// Reconcile ChiaNode owned objects
	if kube.ShouldMakeService(node.Spec.ChiaConfig.PeerService) {
		srv := assemblePeerService(node)
		if err := controllerutil.SetControllerReference(&node, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to create node peer Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error reconciling node peer Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chianodeNamePattern, node.Name),
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to GET ChiaNode peer Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to DELETE ChiaNode peer Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(node.Spec.ChiaConfig.PeerService) {
		srv := assembleHeadlessPeerService(node)
		if err := controllerutil.SetControllerReference(&node, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to create node peer headless Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error reconciling node peer headless Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chianodeNamePattern, node.Name) + "-headless",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to GET ChiaNode peer headless Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to DELETE ChiaNode peer headless Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(node.Spec.ChiaConfig.PeerService) {
		srv := assembleLocalPeerService(node)
		if err := controllerutil.SetControllerReference(&node, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to create node peer internal Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error reconciling node peer internal Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chianodeNamePattern, node.Name) + "-internal",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to GET ChiaNode peer internal Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to DELETE ChiaNode peer internal Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(node.Spec.ChiaConfig.DaemonService) {
		srv := assembleDaemonService(node)
		if err := controllerutil.SetControllerReference(&node, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to create node daemon Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error reconciling node daemon Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chianodeNamePattern, node.Name) + "-daemon",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to GET ChiaNode daemon Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to DELETE ChiaNode daemon Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(node.Spec.ChiaConfig.RPCService) {
		srv := assembleRPCService(node)
		if err := controllerutil.SetControllerReference(&node, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to create node RPC Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error reconciling node RPC Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chianodeNamePattern, node.Name) + "-rpc",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to GET ChiaNode RPC Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to DELETE ChiaNode RPC Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(node.Spec.ChiaExporterConfig.Service) {
		srv := assembleChiaExporterService(node)
		if err := controllerutil.SetControllerReference(&node, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to create node metrics Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error reconciling node metrics Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chianodeNamePattern, node.Name) + "-metrics",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to GET ChiaNode metrics Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to DELETE ChiaNode metrics Service resource", req.NamespacedName))
			}
		}
	}

	// Adds a condition check for Service.Enabled field nilness because the default for ShouldMakeService is true for other services, but should actually be false for this one
	if kube.ShouldMakeService(node.Spec.ChiaHealthcheckConfig.Service) && node.Spec.ChiaHealthcheckConfig.Enabled && node.Spec.ChiaHealthcheckConfig.Service.Enabled != nil {
		srv := assembleChiaHealthcheckService(node)
		if err := controllerutil.SetControllerReference(&node, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to create node healthcheck Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error reconciling node healthcheck Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chianodeNamePattern, node.Name) + "-healthcheck",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to GET ChiaNode healthcheck Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s unable to DELETE ChiaNode healthcheck Service resource", req.NamespacedName))
			}
		}
	}

	stateful := assembleStatefulset(ctx, node)

	if err := controllerutil.SetControllerReference(&node, &stateful, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	res, err := kube.ReconcileStatefulset(ctx, resourceReconciler, stateful)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&node, corev1.EventTypeWarning, "Failed", "Failed to create node Statefulset -- Check operator logs.")
		return *res, fmt.Errorf("ChiaNodeReconciler ChiaNode=%s encountered error reconciling node StatefulSet: %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&node, corev1.EventTypeNormal, "Created", "Successfully created ChiaNode resources.")
	node.Status.Ready = true
	err = r.Status().Update(ctx, &node)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
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
		Complete(r)
}
