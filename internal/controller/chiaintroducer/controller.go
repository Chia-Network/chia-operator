/*
Copyright 2024 Chia Network Inc.
*/

package chiaintroducer

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"

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
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaIntroducerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log))
	log.Info(fmt.Sprintf("ChiaIntroducerReconciler ChiaIntroducer=%s running reconciler...", req.NamespacedName.String()))

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
		metrics.OperatorErrors.Add(1.0)
		log.Error(err, fmt.Sprintf("ChiaIntroducerReconciler ChiaIntroducer=%s unable to fetch ChiaIntroducer resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiaintroducers[req.NamespacedName.String()]
	if !exists {
		chiaintroducers[req.NamespacedName.String()] = true
		metrics.ChiaIntroducers.Add(1.0)
	}

	if kube.ShouldMakeService(introducer.Spec.ChiaConfig.PeerService) {
		srv := assemblePeerService(introducer)
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&introducer, corev1.EventTypeWarning, "Failed", "Failed to create introducer peer Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s encountered error reconciling introducer peer Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiaintroducerNamePattern, introducer.Name),
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("ChiaIntroducerReconciler ChiaIntroducer=%s unable to GET ChiaIntroducer peer Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				log.Error(err, fmt.Sprintf("ChiaIntroducerReconciler ChiaIntroducer=%s unable to DELETE ChiaIntroducer peer Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(introducer.Spec.ChiaConfig.DaemonService) {
		srv := assembleDaemonService(introducer)
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&introducer, corev1.EventTypeWarning, "Failed", "Failed to create introducer daemon Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s encountered error reconciling introducer daemon Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiaintroducerNamePattern, introducer.Name) + "-daemon",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaIntroducerReconciler ChiaIntroducer=%s unable to GET ChiaIntroducer daemon Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaIntroducerReconciler ChiaIntroducer=%s unable to DELETE ChiaIntroducer daemon Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(introducer.Spec.ChiaExporterConfig.Service) {
		srv := assembleChiaExporterService(introducer)
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&introducer, corev1.EventTypeWarning, "Failed", "Failed to create introducer metrics Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s encountered error reconciling introducer metrics Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiaintroducerNamePattern, introducer.Name) + "-metrics",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaIntroducerReconciler ChiaIntroducer=%s unable to GET ChiaIntroducer metrics Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaIntroducerReconciler ChiaIntroducer=%s unable to DELETE ChiaIntroducer metrics Service resource", req.NamespacedName))
			}
		}
	}

	// Creates a persistent volume claim if the GenerateVolumeClaims setting was set to true
	if introducer.Spec.Storage != nil && introducer.Spec.Storage.ChiaRoot != nil && introducer.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil && introducer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.GenerateVolumeClaims {
		pvc, err := assembleVolumeClaim(introducer)
		if err != nil {
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&introducer, corev1.EventTypeWarning, "Failed", "Failed to create introducer PVC -- Check operator logs.")
			return reconcile.Result{}, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s encountered error scaffolding a generated PersistentVolumeClaim: %v", req.NamespacedName, err)
		}

		res, err := kube.ReconcilePersistentVolumeClaim(ctx, resourceReconciler, pvc)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&introducer, corev1.EventTypeWarning, "Failed", "Failed to create introducer PVC -- Check operator logs.")
			return *res, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s encountered error reconciling PersistentVolumeClaim: %v", req.NamespacedName, err)
		}
	}

	deploy := assembleDeployment(introducer)

	if err := controllerutil.SetControllerReference(&introducer, &deploy, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	res, err := kube.ReconcileDeployment(ctx, resourceReconciler, deploy)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&introducer, corev1.EventTypeWarning, "Failed", "Failed to create introducer Deployment -- Check operator logs.")
		return *res, fmt.Errorf("ChiaIntroducerReconciler ChiaIntroducer=%s encountered error reconciling Deployment: %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&introducer, corev1.EventTypeNormal, "Created", "Successfully created ChiaIntroducer resources.")
	introducer.Status.Ready = true
	err = r.Status().Update(ctx, &introducer)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
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
		Watches(
			&corev1.Service{},
			handler.EnqueueRequestsFromMapFunc(r.findObjectsForService),
		).
		Complete(r)
}

// findObjectsForService since we get an event for any changes to every Service in the cluster,
// this lists the Chia custom resource this controller manages in the same namespace as the Service, and then
// checks if this Service has an OwnerReference to any of the Chia custom resources returned in the list,
// and sends a reconcile request for the resource this Service was owned by, if any
func (r *ChiaIntroducerReconciler) findObjectsForService(ctx context.Context, obj client.Object) []reconcile.Request {
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
		for _, ref := range obj.GetOwnerReferences() {
			if ref.Kind == item.Kind && ref.Name == item.GetName() {
				requests[i] = reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      item.GetName(),
						Namespace: item.GetNamespace(),
					},
				}
			}
		}
	}
	return requests
}
