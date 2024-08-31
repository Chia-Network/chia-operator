/*
Copyright 2023 Chia Network Inc.
*/

package chiaca

import (
	"context"
	"fmt"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/metrics"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
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
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaCAReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Running reconciler...")

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
		log.Error(err, "unable to fetch ChiaCA resource")
		return ctrl.Result{}, err
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiacas[req.NamespacedName.String()]
	if !exists {
		chiacas[req.NamespacedName.String()] = true
		metrics.ChiaCAs.Add(1.0)
	}

	// Assemble ServiceAccount
	serviceaccount := assembleServiceAccount(ca)
	if err := controllerutil.SetControllerReference(&ca, &serviceaccount, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&ca, corev1.EventTypeWarning, "Failed", "Failed to assemble ChiaCA ServiceAccount -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaCAReconciler ChiaCA=%s encountered error assembling ServiceAccount: %v", req.NamespacedName, err)
	}
	// Reconcile ServiceAccount
	res, err := kube.ReconcileServiceAccount(ctx, r.Client, serviceaccount)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&ca, corev1.EventTypeWarning, "Failed", "Failed to reconcile ChiaCA ServiceAccount -- Check operator logs.")
		return res, fmt.Errorf("ChiaCAReconciler ChiaCA=%s %v", req.NamespacedName, err)
	}

	// Assemble Role
	role := assembleRole(ca)
	if err := controllerutil.SetControllerReference(&ca, &role, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&ca, corev1.EventTypeWarning, "Failed", "Failed to assemble ChiaCA Role -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaCAReconciler ChiaCA=%s encountered error assembling Role: %v", req.NamespacedName, err)
	}
	// Reconcile Role
	res, err = kube.ReconcileRole(ctx, r.Client, role)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&ca, corev1.EventTypeWarning, "Failed", "Failed to reconcile ChiaCA Role -- Check operator logs.")
		return res, fmt.Errorf("ChiaCAReconciler ChiaCA=%s %v", req.NamespacedName, err)
	}

	// Assemble RoleBinding
	rolebind := assembleRoleBinding(ca)
	if err := controllerutil.SetControllerReference(&ca, &rolebind, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&ca, corev1.EventTypeWarning, "Failed", "Failed to assemble ChiaCA RoleBinding -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaCAReconciler ChiaCA=%s encountered error assembling RoleBinding: %v", req.NamespacedName, err)
	}
	// Reconcile RoleBinding
	res, err = kube.ReconcileRoleBinding(ctx, r.Client, rolebind)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&ca, corev1.EventTypeWarning, "Failed", "Failed to reconcile ChiaCA RoleBinding -- Check operator logs.")
		return res, fmt.Errorf("ChiaCAReconciler ChiaCA=%s %v", req.NamespacedName, err)
	}

	// Query CA Secret
	caExists, err := r.caSecretExists(ctx, ca)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		return ctrl.Result{}, fmt.Errorf("ChiaCAReconciler ChiaCA=%s encountered error querying for existing CA Secret: %v", req.NamespacedName, err)
	}

	// If the CA Secret doesn't already exist, attempt to create a CA generator Job
	if !caExists {
		// Assemble Job
		job := assembleJob(ca)
		if err := controllerutil.SetControllerReference(&ca, &job, r.Scheme); err != nil {
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&ca, corev1.EventTypeWarning, "Failed", "Failed to assemble ChiaCA Job -- Check operator logs.")
			return ctrl.Result{}, fmt.Errorf("ChiaCAReconciler ChiaCA=%s encountered error assembling Job: %v", req.NamespacedName, err)
		}
		// Reconcile Job
		res, err = kube.ReconcileJob(ctx, r.Client, job)
		if err != nil {
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&ca, corev1.EventTypeWarning, "Failed", "Failed to reconcile ChiaCA Job -- Check operator logs.")
			return res, fmt.Errorf("ChiaCAReconciler ChiaCA=%s %v", req.NamespacedName, err)
		}

		// Loop to determine if Secret was made, set to Ready once done
		for i := 1; i <= 100; i++ {
			log.Info(fmt.Sprintf("ChiaCAReconciler ChiaCA=%s waiting for ChiaCA Job to create CA Secret, iteration %d...", req.NamespacedName.String(), i))
			time.Sleep(10 * time.Second)

			found, err := r.caSecretExists(ctx, ca)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, "encountered error querying for ChiaCA secret")
				continue
			}

			if found {
				r.Recorder.Event(&ca, corev1.EventTypeNormal, "Created",
					fmt.Sprintf("Successfully created CA Secret in %s/%s", ca.Namespace, ca.Name))

				ca.Status.Ready = true
				err = r.Status().Update(ctx, &ca)
				if err != nil {
					metrics.OperatorErrors.Add(1.0)
					log.Error(err, "encountered error updating ChiaCA status")
					return ctrl.Result{}, err
				}

				break
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaCAReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaCA{}).
		Owns(&batchv1.Job{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&rbacv1.Role{}).
		Owns(&rbacv1.RoleBinding{}).
		Complete(r)
}
