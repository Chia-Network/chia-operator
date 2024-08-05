/*
Copyright 2023 Chia Network Inc.
*/

package chiawallet

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

// ChiaWalletReconciler reconciles a ChiaWallet object
type ChiaWalletReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chiawallets = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiawallets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiawallets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiawallets/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaWalletReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log))
	log.Info(fmt.Sprintf("ChiaWalletReconciler ChiaWallet=%s running reconciler...", req.NamespacedName.String()))

	// Get the custom resource
	var wallet k8schianetv1.ChiaWallet
	err := r.Get(ctx, req.NamespacedName, &wallet)
	if err != nil && errors.IsNotFound(err) {
		// Remove this object from the map for tracking and subtract this CR's total metric by 1
		_, exists := chiawallets[req.NamespacedName.String()]
		if exists {
			delete(chiawallets, req.NamespacedName.String())
			metrics.ChiaWallets.Sub(1.0)
		}
		return ctrl.Result{}, nil
	}
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		log.Error(err, fmt.Sprintf("ChiaWalletReconciler ChiaWallet=%s unable to fetch ChiaWallet resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiawallets[req.NamespacedName.String()]
	if !exists {
		chiawallets[req.NamespacedName.String()] = true
		metrics.ChiaWallets.Add(1.0)
	}

	// Reconcile ChiaWallet owned objects
	if kube.ShouldMakeService(wallet.Spec.ChiaConfig.PeerService) {
		srv := assemblePeerService(wallet)
		if err := controllerutil.SetControllerReference(&wallet, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&wallet, corev1.EventTypeWarning, "Failed", "Failed to create wallet peer Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaWalletReconciler ChiaWallet=%s encountered error reconciling wallet peer Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiawalletNamePattern, wallet.Name),
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("ChiaWalletReconciler ChiaWallet=%s unable to GET ChiaWallet peer Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				log.Error(err, fmt.Sprintf("ChiaWalletReconciler ChiaWallet=%s unable to DELETE ChiaWallet peer Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(wallet.Spec.ChiaConfig.DaemonService) {
		srv := assembleDaemonService(wallet)
		if err := controllerutil.SetControllerReference(&wallet, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&wallet, corev1.EventTypeWarning, "Failed", "Failed to create wallet daemon Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaWalletReconciler ChiaWallet=%s encountered error reconciling wallet daemon Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiawalletNamePattern, wallet.Name) + "-daemon",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaWalletReconciler ChiaWallet=%s unable to GET ChiaWallet daemon Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaWalletReconciler ChiaWallet=%s unable to DELETE ChiaWallet daemon Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(wallet.Spec.ChiaConfig.RPCService) {
		srv := assembleRPCService(wallet)
		if err := controllerutil.SetControllerReference(&wallet, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&wallet, corev1.EventTypeWarning, "Failed", "Failed to create wallet RPC Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaWalletReconciler ChiaWallet=%s encountered error reconciling wallet RPC Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiawalletNamePattern, wallet.Name) + "-rpc",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaWalletReconciler ChiaWallet=%s unable to GET ChiaWallet RPC Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaWalletReconciler ChiaWallet=%s unable to DELETE ChiaWallet RPC Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(wallet.Spec.ChiaExporterConfig.Service) {
		srv := assembleChiaExporterService(wallet)
		if err := controllerutil.SetControllerReference(&wallet, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&wallet, corev1.EventTypeWarning, "Failed", "Failed to create wallet metrics Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaWalletReconciler ChiaWallet=%s encountered error reconciling wallet chia-exporter Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiawalletNamePattern, wallet.Name) + "-metrics",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaWalletReconciler ChiaWallet=%s unable to GET ChiaWallet metrics Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaWalletReconciler ChiaWallet=%s unable to DELETE ChiaWallet metrics Service resource", req.NamespacedName))
			}
		}
	}

	// Creates a persistent volume claim if the GenerateVolumeClaims setting was set to true
	if wallet.Spec.Storage != nil && wallet.Spec.Storage.ChiaRoot != nil && wallet.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil && wallet.Spec.Storage.ChiaRoot.PersistentVolumeClaim.GenerateVolumeClaims {
		pvc, err := assembleVolumeClaim(wallet)
		if err != nil {
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&wallet, corev1.EventTypeWarning, "Failed", "Failed to create wallet PVC -- Check operator logs.")
			return reconcile.Result{}, fmt.Errorf("ChiaWalletReconciler ChiaWallet=%s encountered error scaffolding a generated PersistentVolumeClaim: %v", req.NamespacedName, err)
		}

		res, err := kube.ReconcilePersistentVolumeClaim(ctx, resourceReconciler, pvc)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&wallet, corev1.EventTypeWarning, "Failed", "Failed to create wallet PVC -- Check operator logs.")
			return *res, fmt.Errorf("ChiaWalletReconciler ChiaWallet=%s encountered error reconciling PersistentVolumeClaim: %v", req.NamespacedName, err)
		}
	}

	deploy := assembleDeployment(ctx, wallet)

	if err := controllerutil.SetControllerReference(&wallet, &deploy, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	res, err := kube.ReconcileDeployment(ctx, resourceReconciler, deploy)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&wallet, corev1.EventTypeWarning, "Failed", "Failed to create wallet Deployment -- Check operator logs.")
		return *res, fmt.Errorf("ChiaWalletReconciler ChiaWallet=%s encountered error reconciling wallet Deployment: %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&wallet, corev1.EventTypeNormal, "Created", "Successfully created ChiaWallet resources.")
	wallet.Status.Ready = true
	err = r.Status().Update(ctx, &wallet)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		log.Error(err, fmt.Sprintf("ChiaWalletReconciler ChiaWallet=%s unable to update ChiaWallet status", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaWalletReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaWallet{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
