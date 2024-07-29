/*
Copyright 2023 Chia Network Inc.
*/

package chiaca

import (
	"context"
	"fmt"
	"time"

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

// ChiaCAReconciler reconciles a ChiaCA object
type ChiaCAReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chiacas = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiacas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiacas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiacas/finalizers,verbs=update
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaCAReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log))
	log.Info(fmt.Sprintf("ChiaCAReconciler ChiaCA=%s running reconciler...", req.NamespacedName.String()))

	// Get the custom resource
	var ca k8schianetv1.ChiaCA
	err := r.Get(ctx, req.NamespacedName, &ca)
	if err != nil && errors.IsNotFound(err) {
		// Remove this object from the map for tracking and subtract this CR's total metric by 1
		_, exists := chiacas[req.NamespacedName.String()]
		if exists {
			delete(chiacas, req.NamespacedName.String())
			metrics.ChiaCAs.Sub(1.0)
		}
		return ctrl.Result{}, nil
	}
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		log.Error(err, fmt.Sprintf("ChiaCAReconciler ChiaCA=%s unable to fetch ChiaCA resource", req.NamespacedName))
		return ctrl.Result{}, err
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiacas[req.NamespacedName.String()]
	if !exists {
		chiacas[req.NamespacedName.String()] = true
		metrics.ChiaCAs.Add(1.0)
	}

	// Reconcile resources, creating them if they don't exist
	sa := r.assembleServiceAccount(ca)
	res, err := kube.ReconcileServiceAccount(ctx, resourceReconciler, sa)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&ca, corev1.EventTypeWarning, "Failed", "Failed to create ServiceAccount -- Check operator logs.")
		return *res, fmt.Errorf("ChiaCAReconciler ChiaCA=%s encountered error reconciling CA generator ServiceAccount: %v", req.NamespacedName, err)
	}

	role := r.assembleRole(ca)
	res, err = kube.ReconcileRole(ctx, resourceReconciler, role)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&ca, corev1.EventTypeWarning, "Failed", "Failed to create Role -- Check operator logs.")
		return *res, fmt.Errorf("ChiaCAReconciler ChiaCA=%s encountered error reconciling CA generator Role: %v", req.NamespacedName, err)
	}

	rb := r.assembleRoleBinding(ca)
	res, err = kube.ReconcileRoleBinding(ctx, resourceReconciler, rb)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&ca, corev1.EventTypeWarning, "Failed", "Failed to create RoleBinding -- Check operator logs.")
		return *res, fmt.Errorf("ChiaCAReconciler ChiaCA=%s encountered error reconciling CA generator RoleBinding: %v", req.NamespacedName, err)
	}

	// Query CA Secret
	_, notFound, err := r.getCASecret(ctx, ca)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		log.Error(err, fmt.Sprintf("ChiaCAReconciler ChiaCA=%s unable to query for ChiaCA secret", req.NamespacedName))
		return ctrl.Result{}, err
	}
	// Create CA generating Job if Secret does not already exist
	if notFound {
		job := r.assembleJob(ca)
		res, err = kube.ReconcileJob(ctx, resourceReconciler, job)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&ca, corev1.EventTypeWarning, "Failed", "Failed to create the CA generating Job -- Check operator logs.")
			return *res, fmt.Errorf("ChiaCAReconciler ChiaCA=%s encountered error reconciling CA generator Job: %v", req.NamespacedName, err)
		}

		// Loop to determine if Secret was made, set to Ready once done
		for i := 1; i <= 100; i++ {
			log.Info(fmt.Sprintf("ChiaCAReconciler ChiaCA=%s waiting for ChiaCA Job to create CA Secret, iteration %d...", req.NamespacedName.String(), i))

			_, notFound, err := r.getCASecret(ctx, ca)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaCAReconciler ChiaCA=%s unable to query for ChiaCA secret", req.NamespacedName))
				return ctrl.Result{}, err
			}

			if !notFound {
				r.Recorder.Event(&ca, corev1.EventTypeNormal, "Created",
					fmt.Sprintf("Successfully created CA Secret in %s/%s", ca.Namespace, ca.Name))

				ca.Status.Ready = true
				err = r.Status().Update(ctx, &ca)
				if err != nil {
					metrics.OperatorErrors.Add(1.0)
					log.Error(err, fmt.Sprintf("ChiaCAReconciler ChiaCA=%s unable to update ChiaCA status", req.NamespacedName))
					return ctrl.Result{}, err
				}

				break
			}

			time.Sleep(10 * time.Second)
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaCAReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaCA{}).
		Complete(r)
}
