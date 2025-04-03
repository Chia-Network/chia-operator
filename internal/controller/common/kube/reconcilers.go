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
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
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

// ReconcilePersistentVolumeClaim uses the controller-runtime client to determine if the PVC resource needs to be created or updated
func ReconcilePersistentVolumeClaim(ctx context.Context, c client.Client, storage *k8schianetv1.StorageConfig, desired corev1.PersistentVolumeClaim) (reconcile.Result, error) {
	klog := log.FromContext(ctx).WithValues("PersistentVolumeClaim.Namespace", desired.Namespace, "PersistentVolumeClaim.Name", desired.Name)
	ensurePVCExists := ShouldMakeChiaRootVolumeClaim(storage)

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

// ReconcileConfigMap uses the controller-runtime client to determine if the ConfigMap resource needs to be created or updated
func ReconcileConfigMap(ctx context.Context, c client.Client, desired corev1.ConfigMap) (reconcile.Result, error) {
	klog := log.FromContext(ctx).WithValues("ConfigMap.Namespace", desired.Namespace, "ConfigMap.Name", desired.Name)

	// Get existing ConfigMap
	var current corev1.ConfigMap
	err := c.Get(ctx, types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, &current)
	if err != nil && errors.IsNotFound(err) {
		// ConfigMap not found - create it
		klog.Info("Creating new ConfigMap")
		if err := c.Create(ctx, &desired); err != nil {
			return ctrl.Result{}, fmt.Errorf("error creating ConfigMap \"%s\": %v", desired.Name, err)
		}
	} else if err != nil {
		// Getting ConfigMap failed, but it wasn't because it doesn't exist, can't continue
		return ctrl.Result{}, fmt.Errorf("error getting existing ConfigMap \"%s\": %v", desired.Name, err)
	} else {
		updated := current

		if !reflect.DeepEqual(current.Data, desired.Data) {
			updated.Data = desired.Data

			if err := c.Update(ctx, &updated); err != nil {
				if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
					return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
				}
				return ctrl.Result{}, fmt.Errorf("error updating ConfigMap \"%s\": %v", desired.Name, err)
			}
		}
	}

	return ctrl.Result{}, nil
}

// ReconcileIngress uses the controller-runtime client to determine if the Ingress resource needs to be created or updated
func ReconcileIngress(ctx context.Context, c client.Client, ingress k8schianetv1.IngressConfig, desired networkingv1.Ingress) (reconcile.Result, error) {
	klog := log.FromContext(ctx).WithValues("Ingress.Namespace", desired.Namespace, "Ingress.Name", desired.Name)

	ensureIngressExists := false
	if ingress.Enabled != nil {
		ensureIngressExists = *ingress.Enabled
	}

	// Get existing Ingress
	var current networkingv1.Ingress
	err := c.Get(ctx, types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, &current)
	if err != nil && errors.IsNotFound(err) {
		// Ingress not found - create if it should exist, or return here if it shouldn't
		if ensureIngressExists {
			klog.Info("Creating new Ingress")
			if err := c.Create(ctx, &desired); err != nil {
				return ctrl.Result{}, fmt.Errorf("error creating Ingress \"%s\": %v", desired.Name, err)
			}
		} else {
			return ctrl.Result{}, nil
		}
	} else if err != nil {
		// Getting Ingress failed, but it wasn't because it doesn't exist, can't do anything
		return ctrl.Result{}, fmt.Errorf("error getting existing Ingress \"%s\": %v", desired.Name, err)
	} else {
		// Ingress exists, so we need to update it if there are any changes, or delete if it was disabled
		if ensureIngressExists {
			if !reflect.DeepEqual(current.Spec, desired.Spec) || !reflect.DeepEqual(current.Labels, desired.Labels) || !reflect.DeepEqual(current.Annotations, desired.Annotations) {
				current.Labels = desired.Labels
				current.Annotations = desired.Annotations
				current.Spec = desired.Spec
				if err := c.Update(ctx, &current); err != nil {
					if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
						return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
					}
					return ctrl.Result{}, fmt.Errorf("error updating Ingress \"%s\": %v", desired.Name, err)
				}
			}
		} else {
			klog.Info("Deleting Ingress because it was disabled")
			if err := c.Delete(ctx, &current); err != nil {
				return ctrl.Result{}, fmt.Errorf("error deleting Ingress \"%s\": %v", desired.Name, err)
			}
		}
	}

	return ctrl.Result{}, nil
}
