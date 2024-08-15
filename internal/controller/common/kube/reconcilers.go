/*
Copyright 2023 Chia Network Inc.
*/

package kube

import (
	"context"
	"fmt"
	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
	"time"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// ObjectModifiedTryAgainError contains the error text for an error that can happen when multiple reconciliation loops are called for the same object at nearly the same time.
	// When this happens, we just want to requeue the reconcile after some amount of time to ensure the latest changes were applied to the sub-resources
	ObjectModifiedTryAgainError = "the object has been modified; please apply your changes to the latest version and try again"
)

// ReconcileService uses the ResourceReconciler to determine if the service resource needs to be created or updated
func ReconcileService(ctx context.Context, c client.Client, service k8schianetv1.Service, desired corev1.Service, defaultEnabled bool) (reconcile.Result, error) {
	klog := log.FromContext(ctx).WithValues("Service.Namespace", desired.Namespace, "Service.Name", desired.Name)
	ensureServiceExists := ShouldMakeService(service, defaultEnabled)

	// Get existing Service
	var current corev1.Service
	err := c.Get(ctx, types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, &current)
	if err != nil && errors.IsNotFound(err) {
		// Service not found - create if it should exist, or return here if it shouldn't
		if ensureServiceExists {
			klog.Info("Creating new Service")
			if err := c.Create(ctx, &desired); err != nil {
				return ctrl.Result{}, fmt.Errorf("error creating Service \"%s\": %v", desired.Name, err)
			}
		} else {
			return ctrl.Result{}, nil
		}
	} else if err != nil {
		// Getting Service failed, but it wasn't because it doesn't exist, can't do anything
		return ctrl.Result{}, fmt.Errorf("error getting existing Service \"%s\": %v", desired.Name, err)
	}

	// Service exists, so we need to update it if there are any changes, or delete if it was disabled
	if ensureServiceExists {
		desiredAnnotations := CombineMaps(current.Annotations, desired.Annotations)
		if !reflect.DeepEqual(current.Spec, desired.Spec) || !reflect.DeepEqual(current.Labels, desired.Labels) || !reflect.DeepEqual(current.Annotations, desiredAnnotations) {
			current.Labels = desired.Labels
			current.Annotations = desiredAnnotations
			current.Spec = desired.Spec
			if err := c.Update(ctx, &current); err != nil {
				if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
					return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
				}
				return ctrl.Result{}, fmt.Errorf("error updating Service \"%s\": %v", desired.Name, err)
			}
		}
	} else {
		klog.Info("Deleting Service because it was disabled")
		if err := c.Delete(ctx, &current); err != nil {
			return ctrl.Result{}, fmt.Errorf("error deleting Service \"%s\": %v", desired.Name, err)
		}
	}

	return ctrl.Result{}, nil
}

// ReconcileDeployment uses the ResourceReconciler to determine if the deployment resource needs to be created or updated
func ReconcileDeployment(ctx context.Context, c client.Client, desired appsv1.Deployment) (reconcile.Result, error) {
	klog := log.FromContext(ctx).WithValues("Deployment.Namespace", desired.Namespace, "Deployment.Name", desired.Name)

	// Get existing PVC
	var current appsv1.Deployment
	err := c.Get(ctx, types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, &current)
	if err != nil && errors.IsNotFound(err) {
		klog.Info("Creating new Deployment")
		if err := c.Create(ctx, &desired); err != nil {
			return ctrl.Result{}, fmt.Errorf("error creating Deployment \"%s\": %v", desired.Name, err)
		}
	} else if err != nil {
		return ctrl.Result{}, fmt.Errorf("error getting existing Deployment \"%s\": %v", desired.Name, err)
	}

	// Deployment exists, so we need to update it if there are any changes
	desiredAnnotations := CombineMaps(current.Annotations, desired.Annotations)
	if !reflect.DeepEqual(current.Spec, desired.Spec) || !reflect.DeepEqual(current.Labels, desired.Labels) || !reflect.DeepEqual(current.Annotations, desiredAnnotations) {
		current.Labels = desired.Labels
		current.Annotations = desiredAnnotations
		current.Spec = desired.Spec
		if err := c.Update(ctx, &current); err != nil {
			if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
				return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
			}
			return ctrl.Result{}, fmt.Errorf("error updating Deployment \"%s\": %v", desired.Name, err)
		}
	}

	return ctrl.Result{}, nil
}

// ReconcileStatefulset uses the ResourceReconciler to determine if the statefulset resource needs to be created or updated
func ReconcileStatefulset(ctx context.Context, c client.Client, desired appsv1.StatefulSet) (reconcile.Result, error) {
	klog := log.FromContext(ctx).WithValues("StatefulSet.Namespace", desired.Namespace, "StatefulSet.Name", desired.Name)

	// Get existing PVC
	var current appsv1.StatefulSet
	err := c.Get(ctx, types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, &current)
	if err != nil && errors.IsNotFound(err) {
		klog.Info("Creating new StatefulSet")
		if err := c.Create(ctx, &desired); err != nil {
			return ctrl.Result{}, fmt.Errorf("error creating StatefulSet \"%s\": %v", desired.Name, err)
		}
	} else if err != nil {
		return ctrl.Result{}, fmt.Errorf("error getting existing StatefulSet \"%s\": %v", desired.Name, err)
	}

	// StatefulSet exists, so we need to update it if there are any changes
	desiredAnnotations := CombineMaps(current.Annotations, desired.Annotations)
	if !reflect.DeepEqual(current.Spec, desired.Spec) || !reflect.DeepEqual(current.Labels, desired.Labels) || !reflect.DeepEqual(current.Annotations, desiredAnnotations) {
		current.Labels = desired.Labels
		current.Annotations = desired.Annotations
		current.Spec = desired.Spec
		if err := c.Update(ctx, &current); err != nil {
			if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
				return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
			}
			return ctrl.Result{}, fmt.Errorf("error updating StatefulSet \"%s\": %v", desired.Name, err)
		}
	}

	return ctrl.Result{}, nil
}

// ReconcileServiceAccount uses the ResourceReconciler to determine if the serviceaccount resource needs to be created or updated
func ReconcileServiceAccount(ctx context.Context, rec reconciler.ResourceReconciler, sa corev1.ServiceAccount) (*reconcile.Result, error) {
	return rec.ReconcileResource(&sa, reconciler.StatePresent)
}

// ReconcileRole uses the ResourceReconciler to determine if the role resource needs to be created or updated
func ReconcileRole(ctx context.Context, rec reconciler.ResourceReconciler, role rbacv1.Role) (*reconcile.Result, error) {
	return rec.ReconcileResource(&role, reconciler.StatePresent)
}

// ReconcileRoleBinding uses the ResourceReconciler to determine if the rolebinding resource needs to be created or updated
func ReconcileRoleBinding(ctx context.Context, rec reconciler.ResourceReconciler, rb rbacv1.RoleBinding) (*reconcile.Result, error) {
	return rec.ReconcileResource(&rb, reconciler.StatePresent)
}

// ReconcileJob uses the ResourceReconciler to determine if the job resource needs to be created or updated
func ReconcileJob(ctx context.Context, rec reconciler.ResourceReconciler, job batchv1.Job) (*reconcile.Result, error) {
	return rec.ReconcileResource(&job, reconciler.StatePresent)
}

// ReconcilePersistentVolumeClaim uses the ResourceReconciler to determine if the PVC resource needs to be created or updated
func ReconcilePersistentVolumeClaim(ctx context.Context, c client.Client, storage *k8schianetv1.StorageConfig, desired corev1.PersistentVolumeClaim) (reconcile.Result, error) {
	klog := log.FromContext(ctx).WithValues("PersistentVolumeClaim.Namespace", desired.Namespace, "PersistentVolumeClaim.Name", desired.Name)
	ensurePVCExists := ShouldMakeVolumeClaim(storage)

	// Get existing PVC
	var current corev1.PersistentVolumeClaim
	err := c.Get(ctx, types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, &current)
	if err != nil && errors.IsNotFound(err) {
		// PVC not found - create if it should exist, or return here if it shouldn't
		if ensurePVCExists {
			klog.Info("Creating new PersistentVolumeClaim")
			if err := c.Create(ctx, &desired); err != nil {
				return ctrl.Result{}, fmt.Errorf("error creating PersistentVolumeClaim \"%s\": %v", desired.Name, err)
			}
		} else {
			return ctrl.Result{}, nil
		}
	} else if err != nil {
		// Getting PVC failed, but it wasn't because it doesn't exist, can't do anything
		return ctrl.Result{}, fmt.Errorf("error getting existing PersistentVolumeClaim \"%s\": %v", desired.Name, err)
	}

	// PVC exists, so we need to update it if GeneratePersistentVolumes is enabled
	// For safety reasons we never delete PVCs, however, chia-operator users should clean up their own storage if desired
	if ensurePVCExists {
		// PVC updates are complex, many fields just cannot be changed, so we only check resource request changes from the spec
		desiredAnnotations := CombineMaps(current.Annotations, desired.Annotations)
		if !reflect.DeepEqual(current.Labels, desired.Labels) || !reflect.DeepEqual(current.Annotations, desiredAnnotations) || !reflect.DeepEqual(current.Spec.Resources.Requests, desired.Spec.Resources.Requests) {
			current.Labels = desired.Labels
			current.Annotations = desiredAnnotations
			current.Spec.Resources.Requests = desired.Spec.Resources.Requests
			if err := c.Update(ctx, &current); err != nil {
				if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
					return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
				}
				return ctrl.Result{}, fmt.Errorf("error updating PersistentVolumeClaim \"%s\": %v", desired.Name, err)
			}
		}
	}

	return ctrl.Result{}, nil
}
