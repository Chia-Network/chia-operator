/*
Copyright 2023 Chia Network Inc.
*/

package kube

import (
	"context"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcileService uses the ResourceReconciler to determine if the service resource needs to be created or updated
func ReconcileService(ctx context.Context, rec reconciler.ResourceReconciler, service corev1.Service) (*reconcile.Result, error) {
	return rec.ReconcileResource(&service, reconciler.StatePresent)
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
