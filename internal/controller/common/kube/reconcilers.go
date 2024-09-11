/*
Copyright 2023 Chia Network Inc.
*/

package kube

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

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

// ReconcileService uses the controller-runtime client to determine if the service resource needs to be created or updated
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
	} else {
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
	}

	return ctrl.Result{}, nil
}

// ReconcileDeployment uses the controller-runtime client to determine if the deployment resource needs to be created or updated
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
	} else {
		// Need to handle a case where Deployment's spec.Selector.MatchLabels changed, since the field is immutable
		if !reflect.DeepEqual(current.Spec.Selector.MatchLabels, desired.Spec.Selector.MatchLabels) {
			klog.Info("Recreating Deployment for new Selector labels -- selector labels are immutable")

			if err := c.Delete(ctx, &current); err != nil {
				if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
					return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
				}
				return ctrl.Result{}, fmt.Errorf("error deleting Deployment \"%s\": %v", current.Name, err)
			}

			// Wait for the deployment to be deleted
			for {
				var tmp appsv1.Deployment
				err = c.Get(ctx, types.NamespacedName{
					Namespace: current.Namespace,
					Name:      current.Name,
				}, &tmp)
				if err != nil {
					if client.IgnoreNotFound(err) == nil {
						break
					}
					return ctrl.Result{}, fmt.Errorf("error waiting for Deployment to be deleted \"%s\": %v", desired.Name, err)
				}
				time.Sleep(2 * time.Second)
			}

			if err := c.Create(ctx, &desired); err != nil {
				return ctrl.Result{}, fmt.Errorf("error creating Deployment \"%s\": %v", desired.Name, err)
			}

			return ctrl.Result{}, nil // Exit reconciler here because we created the desired Deployment
		}

		// Deployment exists, so we need to update it if there are any changes.
		// We'll make a copy of the current Deployment to make sure we only change mutable Deployment fields.
		// Then we will compare the current and updated Deployments, and send an Update request if there was any diff.
		updated := current

		desiredAnnotations := CombineMaps(current.Annotations, desired.Annotations)
		if !reflect.DeepEqual(current.Annotations, desiredAnnotations) {
			updated.Annotations = desiredAnnotations
		}

		if !reflect.DeepEqual(current.Labels, desired.Labels) {
			updated.Labels = desired.Labels
		}

		if !reflect.DeepEqual(current.Spec, desired.Spec) {
			updated.Spec = desired.Spec
		}

		if !reflect.DeepEqual(current, updated) {
			if err := c.Update(ctx, &updated); err != nil {
				if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
					return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
				}
				return ctrl.Result{}, fmt.Errorf("error updating Deployment \"%s\": %v", updated.Name, err)
			}
		}
	}

	return ctrl.Result{}, nil
}

// ReconcileStatefulset uses the controller-runtime client to determine if the statefulset resource needs to be created or updated
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
	} else {
		// Need to handle a case where StatefulSet's spec.Selector.MatchLabels changed, since the field is immutable
		if !reflect.DeepEqual(current.Spec.Selector.MatchLabels, desired.Spec.Selector.MatchLabels) {
			klog.Info("Recreating StatefulSet for new Selector labels -- selector labels are immutable")

			if err := c.Delete(ctx, &current); err != nil {
				if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
					return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
				}
				return ctrl.Result{}, fmt.Errorf("error deleting StatefulSet \"%s\": %v", current.Name, err)
			}

			// Wait for the statefulset to be deleted
			for {
				var tmp appsv1.StatefulSet
				err = c.Get(ctx, types.NamespacedName{
					Namespace: current.Namespace,
					Name:      current.Name,
				}, &tmp)
				if err != nil {
					if client.IgnoreNotFound(err) == nil {
						break
					}
					return ctrl.Result{}, fmt.Errorf("error waiting for StatefulSet to be deleted \"%s\": %v", desired.Name, err)
				}
				time.Sleep(2 * time.Second)
			}

			if err := c.Create(ctx, &desired); err != nil {
				return ctrl.Result{}, fmt.Errorf("error creating StatefulSet \"%s\": %v", desired.Name, err)
			}

			return ctrl.Result{}, nil // Exit reconciler here because we created the desired StatefulSet
		}

		// StatefulSet exists, so we need to update it if there are any changes.
		// We'll make a copy of the current StatefulSet to make sure we only change mutable StatefulSet fields.
		// Then we will compare the current and updated StatefulSets, and send an Update request if there was any diff.
		updated := current

		desiredAnnotations := CombineMaps(current.Annotations, desired.Annotations)
		if !reflect.DeepEqual(current.Annotations, desiredAnnotations) {
			updated.Annotations = desiredAnnotations
		}

		if !reflect.DeepEqual(current.Labels, desired.Labels) {
			updated.Labels = desired.Labels
		}

		if !reflect.DeepEqual(current.Spec.UpdateStrategy, desired.Spec.UpdateStrategy) {
			updated.Spec.UpdateStrategy = desired.Spec.UpdateStrategy
		}

		if !reflect.DeepEqual(current.Spec.Replicas, desired.Spec.Replicas) {
			updated.Spec.Replicas = desired.Spec.Replicas
		}

		if !reflect.DeepEqual(current.Spec.Template, desired.Spec.Template) {
			updated.Spec.Template = desired.Spec.Template
		}

		if !reflect.DeepEqual(current, updated) {
			if err := c.Update(ctx, &updated); err != nil {
				if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
					return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
				}
				return ctrl.Result{}, fmt.Errorf("error updating StatefulSet \"%s\": %v", desired.Name, err)
			}
		}
	}

	return ctrl.Result{}, nil
}

// ReconcileServiceAccount uses the controller-runtime client to determine if the serviceaccount resource needs to be created or updated
func ReconcileServiceAccount(ctx context.Context, c client.Client, desired corev1.ServiceAccount) (reconcile.Result, error) {
	klog := log.FromContext(ctx).WithValues("ServiceAccount.Namespace", desired.Namespace, "ServiceAccount.Name", desired.Name)

	var current corev1.ServiceAccount
	err := c.Get(ctx, types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, &current)
	if err != nil && errors.IsNotFound(err) {
		klog.Info("Creating new ServiceAccount")
		if err := c.Create(ctx, &desired); err != nil {
			return ctrl.Result{}, fmt.Errorf("error creating ServiceAccount \"%s\": %v", desired.Name, err)
		}
	} else if err != nil {
		return ctrl.Result{}, fmt.Errorf("error getting existing ServiceAccount \"%s\": %v", desired.Name, err)
	}

	// No need to update ServiceAccounts if they already exist because we don't manage any of the fields on a ServiceAccount

	return ctrl.Result{}, nil
}

// ReconcileRole uses the controller-runtime client to determine if the role resource needs to be created or updated
func ReconcileRole(ctx context.Context, c client.Client, desired rbacv1.Role) (reconcile.Result, error) {
	klog := log.FromContext(ctx).WithValues("Role.Namespace", desired.Namespace, "Role.Name", desired.Name)

	var current rbacv1.Role
	err := c.Get(ctx, types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, &current)
	if err != nil && errors.IsNotFound(err) {
		klog.Info("Creating new Role")
		if err := c.Create(ctx, &desired); err != nil {
			return ctrl.Result{}, fmt.Errorf("error creating Role \"%s\": %v", desired.Name, err)
		}
	} else if err != nil {
		return ctrl.Result{}, fmt.Errorf("error getting existing Role \"%s\": %v", desired.Name, err)
	} else {
		// Role exists, so we need to update it if it's changed
		if !reflect.DeepEqual(current.Rules, desired.Rules) {
			current.Rules = desired.Rules
			if err := c.Update(ctx, &current); err != nil {
				if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
					return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
				}
				return ctrl.Result{}, fmt.Errorf("error updating Role \"%s\": %v", desired.Name, err)
			}
		}
	}

	return ctrl.Result{}, nil
}

// ReconcileRoleBinding uses the controller-runtime client to determine if the rolebinding resource needs to be created or updated
func ReconcileRoleBinding(ctx context.Context, c client.Client, desired rbacv1.RoleBinding) (reconcile.Result, error) {
	klog := log.FromContext(ctx).WithValues("RoleBinding.Namespace", desired.Namespace, "RoleBinding.Name", desired.Name)

	var current rbacv1.RoleBinding
	err := c.Get(ctx, types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, &current)
	if err != nil && errors.IsNotFound(err) {
		klog.Info("Creating new RoleBinding")
		if err := c.Create(ctx, &desired); err != nil {
			return ctrl.Result{}, fmt.Errorf("error creating RoleBinding \"%s\": %v", desired.Name, err)
		}
	} else if err != nil {
		return ctrl.Result{}, fmt.Errorf("error getting existing RoleBinding \"%s\": %v", desired.Name, err)
	} else {
		// RoleBinding exists, so we need to update it if it's changed
		if !reflect.DeepEqual(current.Subjects, desired.Subjects) || !reflect.DeepEqual(current.RoleRef, desired.RoleRef) {
			current.RoleRef = desired.RoleRef
			current.Subjects = desired.Subjects
			if err := c.Update(ctx, &current); err != nil {
				if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
					return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
				}
				return ctrl.Result{}, fmt.Errorf("error updating RoleBinding \"%s\": %v", desired.Name, err)
			}
		}
	}

	return ctrl.Result{}, nil
}

// ReconcileJob uses the controller-runtime client to determine if the job resource needs to be created or updated
func ReconcileJob(ctx context.Context, c client.Client, desired batchv1.Job) (reconcile.Result, error) {
	klog := log.FromContext(ctx).WithValues("Job.Namespace", desired.Namespace, "Job.Name", desired.Name)

	var current batchv1.Job
	err := c.Get(ctx, types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, &current)
	if err != nil && errors.IsNotFound(err) {
		klog.Info("Creating new Job")
		if err := c.Create(ctx, &desired); err != nil {
			return ctrl.Result{}, fmt.Errorf("error creating Job \"%s\": %v", desired.Name, err)
		}
	} else if err != nil {
		return ctrl.Result{}, fmt.Errorf("error getting existing Job \"%s\": %v", desired.Name, err)
	} else {
		// Job exists, so we need to update it if there are any Pod env changes, which might specify a new CA Secret to make -- since the pod template in a Job is immutable, we're just going to delete and create a new Job
		currentContainerEnv := current.Spec.Template.Spec.Containers[0].Env
		desiredContainerEnv := desired.Spec.Template.Spec.Containers[0].Env
		if !reflect.DeepEqual(currentContainerEnv, desiredContainerEnv) {
			// Delete policy cleans up the old Job's Pods
			deletePolicy := v1.DeletePropagationBackground
			deleteOptions := &client.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			}
			if err := c.Delete(ctx, &current, deleteOptions); err != nil {
				return ctrl.Result{}, fmt.Errorf("error deleting old Job \"%s\": %v", desired.Name, err)
			}

			// TODO could probably have a fancier mechanism of waiting for the old Job to be deleted
			time.Sleep(5 * time.Second)

			if err := c.Create(ctx, &desired); err != nil {
				if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
					return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
				}
				return ctrl.Result{}, fmt.Errorf("error creating Job \"%s\": %v", desired.Name, err)
			}
		}
	}

	return ctrl.Result{}, nil
}

// ReconcilePersistentVolumeClaim uses the controller-runtime client to determine if the PVC resource needs to be created or updated
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
		// Getting PVC failed, but it wasn't because it doesn't exist, can't continue
		return ctrl.Result{}, fmt.Errorf("error getting existing PersistentVolumeClaim \"%s\": %v", desired.Name, err)
	} else {
		// PVC exists, so we need to update it if GeneratePersistentVolumes is enabled
		// For safety reasons we never delete PVCs, however, chia-operator users should clean up their own storage if desired
		if ensurePVCExists {
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
	}

	return ctrl.Result{}, nil
}
