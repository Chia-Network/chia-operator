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
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcileService uses the ResourceReconciler to determine if the service resource needs to be created or updated
func ReconcileService(ctx context.Context, c client.Client, service k8schianetv1.Service, desired corev1.Service, defaultEnabled bool) error {
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
				return fmt.Errorf("error creating Service \"%s\": %v", desired.Name, err)
			}
		} else {
			return nil
		}
	} else if err != nil {
		// Getting Service failed, but it wasn't because it doesn't exist, can't do anything
		return fmt.Errorf("error getting existing Service \"%s\": %v", desired.Name, err)
	}

	// Service exists, so we need to update it if there are any changes, or delete if it was disabled
	if ensureServiceExists {
		if !reflect.DeepEqual(current.Spec, desired.Spec) || !reflect.DeepEqual(current.Labels, desired.Labels) || !reflect.DeepEqual(current.Annotations, desired.Annotations) {
			klog.Info("Updating Service with new spec or metadata")
			current.Labels = desired.Labels
			current.Annotations = desired.Annotations
			current.Spec = desired.Spec
			if err := c.Update(ctx, &current); err != nil {
				return fmt.Errorf("error updating Service \"%s\": %v", desired.Name, err)
			}
		}
	} else {
		klog.Info("Deleting Service because it was disabled")
		if err := c.Delete(ctx, &current); err != nil {
			return fmt.Errorf("error deleting Service \"%s\": %v", desired.Name, err)
		}
	}

	return nil
}

// ReconcileDeployment uses the ResourceReconciler to determine if the deployment resource needs to be created or updated
func ReconcileDeployment(ctx context.Context, rec reconciler.ResourceReconciler, deploy appsv1.Deployment) (*reconcile.Result, error) {
	return rec.ReconcileResource(&deploy, reconciler.StatePresent)
}

// ReconcileStatefulset uses the ResourceReconciler to determine if the statefulset resource needs to be created or updated
func ReconcileStatefulset(ctx context.Context, rec reconciler.ResourceReconciler, stateful appsv1.StatefulSet) (*reconcile.Result, error) {
	return rec.ReconcileResource(&stateful, reconciler.StatePresent)
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
func ReconcilePersistentVolumeClaim(ctx context.Context, rec reconciler.ResourceReconciler, pvc corev1.PersistentVolumeClaim) (*reconcile.Result, error) {
	return rec.ReconcileResource(&pvc, reconciler.StatePresent)
}
