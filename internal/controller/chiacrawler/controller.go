/*
Copyright 2024 Chia Network Inc.
*/

package chiacrawler

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"

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

// ChiaCrawlerReconciler reconciles a ChiaCrawler object
type ChiaCrawlerReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chiacrawlers map[string]bool = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiacrawlers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiacrawlers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiacrawlers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *ChiaCrawlerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log))
	log.Info(fmt.Sprintf("ChiaCrawlerReconciler ChiaCrawler=%s", req.NamespacedName.String()))

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
		metrics.OperatorErrors.Add(1.0)
		log.Error(err, fmt.Sprintf("ChiaCrawlerReconciler ChiaCrawler=%s unable to fetch ChiaCrawler resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiacrawlers[req.NamespacedName.String()]
	if !exists {
		chiacrawlers[req.NamespacedName.String()] = true
		metrics.ChiaCrawlers.Add(1.0)
	}

	if kube.ShouldMakeService(crawler.Spec.ChiaConfig.PeerService) {
		srv := r.assemblePeerService(ctx, crawler)
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to create crawler peer Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s encountered error reconciling crawler peer Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiacrawlerNamePattern, crawler.Name),
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("ChiaCrawlerReconciler ChiaCrawler=%s unable to GET ChiaCrawler peer Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				log.Error(err, fmt.Sprintf("ChiaCrawlerReconciler ChiaCrawler=%s unable to DELETE ChiaCrawler peer Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(crawler.Spec.ChiaConfig.DaemonService) {
		srv := r.assembleDaemonService(ctx, crawler)
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to create crawler daemon Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s encountered error reconciling crawler daemon Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiacrawlerNamePattern, crawler.Name) + "-daemon",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaCrawlerReconciler ChiaCrawler=%s unable to GET ChiaCrawler daemon Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaCrawlerReconciler ChiaCrawler=%s unable to DELETE ChiaCrawler daemon Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(crawler.Spec.ChiaConfig.RPCService) {
		srv := r.assembleRPCService(ctx, crawler)
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to create crawler RPC Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s encountered error reconciling crawler RPC Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiacrawlerNamePattern, crawler.Name) + "-rpc",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaCrawlerReconciler ChiaCrawler=%s unable to GET ChiaCrawler RPC Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaCrawlerReconciler ChiaCrawler=%s unable to DELETE ChiaCrawler RPC Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(crawler.Spec.ChiaExporterConfig.Service) {
		srv := r.assembleChiaExporterService(ctx, crawler)
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to create crawler metrics Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s encountered error reconciling crawler metrics Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiacrawlerNamePattern, crawler.Name) + "-metrics",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaCrawlerReconciler ChiaCrawler=%s unable to GET ChiaCrawler metrics Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaCrawlerReconciler ChiaCrawler=%s unable to DELETE ChiaCrawler metrics Service resource", req.NamespacedName))
			}
		}
	}

	// Creates a persistent volume claim if the GenerateVolumeClaims setting was set to true
	if crawler.Spec.Storage != nil && crawler.Spec.Storage.ChiaRoot != nil && crawler.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil && crawler.Spec.Storage.ChiaRoot.PersistentVolumeClaim.GenerateVolumeClaims {
		pvc, err := r.assembleVolumeClaim(ctx, crawler)
		if err != nil {
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to create crawler PVC -- Check operator logs.")
			return reconcile.Result{}, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s encountered error scaffolding a generated PersistentVolumeClaim: %v", req.NamespacedName, err)
		}

		res, err := kube.ReconcilePersistentVolumeClaim(ctx, resourceReconciler, pvc)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to create crawler PVC -- Check operator logs.")
			return *res, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s encountered error reconciling PersistentVolumeClaim: %v", req.NamespacedName, err)
		}
	}

	deploy := r.assembleDeployment(ctx, crawler)
	res, err := kube.ReconcileDeployment(ctx, resourceReconciler, deploy)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&crawler, corev1.EventTypeWarning, "Failed", "Failed to create crawler Deployment -- Check operator logs.")
		return *res, fmt.Errorf("ChiaCrawlerReconciler ChiaCrawler=%s encountered error reconciling Deployment: %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&crawler, corev1.EventTypeNormal, "Created", "Successfully created ChiaCrawler resources.")
	crawler.Status.Ready = true
	err = r.Status().Update(ctx, &crawler)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		log.Error(err, fmt.Sprintf("ChiaCrawlerReconciler ChiaCrawler=%s unable to update ChiaCrawler status", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaCrawlerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaCrawler{}).
		Complete(r)
}
