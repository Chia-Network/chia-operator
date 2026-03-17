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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// ObjectModifiedTryAgainError contains the error text for an error that can happen when multiple reconciliation loops are called for the same object at nearly the same time.
	// When this happens, we just want to requeue the reconcile after some amount of time to ensure the latest changes were applied to the sub-resources
	ObjectModifiedTryAgainError = "the object has been modified; please apply your changes to the latest version and try again"
)

// serverSideApply attempts to apply the desired object server-side.
func serverSideApply(ctx context.Context, c client.Client, desired runtime.Object, kind, apiVersion string) error {
	u := &unstructured.Unstructured{}
	objMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(desired)
	if err != nil {
		return fmt.Errorf("error converting object: %w", err)
	}
	u.Object = objMap
	u.SetKind(kind)
	u.SetAPIVersion(apiVersion)
	u.SetManagedFields(nil)

	err = c.Patch(ctx, u, client.Apply, client.ForceOwnership, client.FieldOwner("chia-operator"))
	if err != nil {
		klog.V(1).Info("object that failed to apply", "applyError", err, "object", objMap)
		return fmt.Errorf("error applying object: %w", err)
	}
	return nil
}

// filterStaleContainers returns a container list with only the containers
// whose name appears in the desired list, and reports whether any were removed.
func filterStaleContainers(current, desired []corev1.Container) ([]corev1.Container, bool) {
	desiredNames := make(map[string]struct{}, len(desired))
	for _, ctr := range desired {
		desiredNames[ctr.Name] = struct{}{}
	}

	filtered := make([]corev1.Container, 0, len(current))
	for _, ctr := range current {
		if _, ok := desiredNames[ctr.Name]; ok {
			filtered = append(filtered, ctr)
		}
	}

	return filtered, len(filtered) != len(current)
}

// filterStaleContainerFields nils out optional fields on current containers
// that are set in the live object but absent in the corresponding desired
// container. This is necessary because the merge-patch used to remove stale
// containers replaces the entire containers array, which can re-assert field
// values under a different field manager and prevent SSA from removing them.
// Returns true if any fields were cleared.
func filterStaleContainerFields(current, desired []corev1.Container) bool {
	desiredMap := make(map[string]corev1.Container, len(desired))
	for _, ctr := range desired {
		desiredMap[ctr.Name] = ctr
	}

	changed := false
	for i, ctr := range current {
		desiredCtr, ok := desiredMap[ctr.Name]
		if !ok {
			continue
		}
		if ctr.LivenessProbe != nil && desiredCtr.LivenessProbe == nil {
			current[i].LivenessProbe = nil
			changed = true
		}
		if ctr.ReadinessProbe != nil && desiredCtr.ReadinessProbe == nil {
			current[i].ReadinessProbe = nil
			changed = true
		}
		if ctr.StartupProbe != nil && desiredCtr.StartupProbe == nil {
			current[i].StartupProbe = nil
			changed = true
		}
		if ctr.SecurityContext != nil && desiredCtr.SecurityContext == nil {
			current[i].SecurityContext = nil
			changed = true
		}
		if ctr.Resources.Limits != nil && desiredCtr.Resources.Limits == nil {
			current[i].Resources.Limits = nil
			changed = true
		}
		if ctr.Resources.Requests != nil && desiredCtr.Resources.Requests == nil {
			current[i].Resources.Requests = nil
			changed = true
		}
	}
	return changed
}

// removeStaleWorkloadContainers patches out containers and init containers from a live
// workload object that are not present in the desired pod spec, and clears optional
// fields from remaining containers when the desired spec no longer includes them. This
// is necessary because Kubernetes SSA treats containers as an associative list keyed by
// name and does not remove entries that are simply omitted from the applied configuration.
func removeStaleWorkloadContainers(ctx context.Context, c client.Client, obj client.Object, currentPodSpec, desiredPodSpec *corev1.PodSpec) error {
	filteredContainers, staleContainers := filterStaleContainers(currentPodSpec.Containers, desiredPodSpec.Containers)
	filteredInitContainers, staleInitContainers := filterStaleContainers(currentPodSpec.InitContainers, desiredPodSpec.InitContainers)
	staleFields := filterStaleContainerFields(filteredContainers, desiredPodSpec.Containers)
	if !staleContainers && !staleInitContainers && !staleFields {
		return nil
	}

	original := obj.DeepCopyObject().(client.Object)
	if staleContainers || staleFields {
		currentPodSpec.Containers = filteredContainers
	}
	if staleInitContainers {
		currentPodSpec.InitContainers = filteredInitContainers
	}
	return c.Patch(ctx, obj, client.MergeFrom(original))
}

// ReconcileService uses the controller-runtime client to determine if the service resource needs to be created or updated
func ReconcileService(ctx context.Context, c client.Client, service k8schianetv1.Service, desired corev1.Service, defaultEnabled bool) (reconcile.Result, error) {
	klog := log.FromContext(ctx).WithValues("Service.Namespace", desired.Namespace, "Service.Name", desired.Name)

	if ShouldMakeService(service, defaultEnabled) {
		if err := serverSideApply(ctx, c, &desired, "Service", "v1"); err != nil {
			if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
				return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
			}
			return ctrl.Result{}, fmt.Errorf("error applying Service \"%s\": %v", desired.Name, err)
		}
		return ctrl.Result{}, nil
	}

	var current corev1.Service
	err := c.Get(ctx, types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, &current)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("error getting Service \"%s\": %v", desired.Name, err)
	}

	klog.Info("Deleting Service because it was disabled")
	if err := c.Delete(ctx, &current); err != nil {
		return ctrl.Result{}, fmt.Errorf("error deleting Service \"%s\": %v", desired.Name, err)
	}

	return ctrl.Result{}, nil
}

// ReconcileDeployment uses the controller-runtime client to determine if the deployment resource needs to be created or updated
func ReconcileDeployment(ctx context.Context, c client.Client, desired appsv1.Deployment) (reconcile.Result, error) {
	klog := log.FromContext(ctx).WithValues("Deployment.Namespace", desired.Namespace, "Deployment.Name", desired.Name)

	var current appsv1.Deployment
	err := c.Get(ctx, types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, &current)
	if err != nil && !errors.IsNotFound(err) {
		return ctrl.Result{}, fmt.Errorf("error getting Deployment \"%s\": %v", desired.Name, err)
	}

	if err == nil {
		if !reflect.DeepEqual(current.Spec.Selector.MatchLabels, desired.Spec.Selector.MatchLabels) {
			klog.Info("Recreating Deployment for new Selector labels -- selector labels are immutable")

			if err := c.Delete(ctx, &current); err != nil {
				if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
					return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
				}
				return ctrl.Result{}, fmt.Errorf("error deleting Deployment \"%s\": %v", current.Name, err)
			}

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
		} else {
			if err := removeStaleWorkloadContainers(ctx, c, &current, &current.Spec.Template.Spec, &desired.Spec.Template.Spec); err != nil {
				if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
					return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
				}
				return ctrl.Result{}, fmt.Errorf("error removing stale containers from Deployment \"%s\": %v", current.Name, err)
			}
		}
	}

	if err := serverSideApply(ctx, c, &desired, "Deployment", "apps/v1"); err != nil {
		if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
			return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
		}
		return ctrl.Result{}, fmt.Errorf("error applying Deployment \"%s\": %v", desired.Name, err)
	}

	return ctrl.Result{}, nil
}

// ReconcileStatefulset uses the controller-runtime client to determine if the statefulset resource needs to be created or updated
func ReconcileStatefulset(ctx context.Context, c client.Client, desired appsv1.StatefulSet) (reconcile.Result, error) {
	klog := log.FromContext(ctx).WithValues("StatefulSet.Namespace", desired.Namespace, "StatefulSet.Name", desired.Name)

	var current appsv1.StatefulSet
	err := c.Get(ctx, types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, &current)
	if err != nil && !errors.IsNotFound(err) {
		return ctrl.Result{}, fmt.Errorf("error getting StatefulSet \"%s\": %v", desired.Name, err)
	}

	if err == nil {
		if !reflect.DeepEqual(current.Spec.Selector.MatchLabels, desired.Spec.Selector.MatchLabels) {
			klog.Info("Recreating StatefulSet for new Selector labels -- selector labels are immutable")

			if err := c.Delete(ctx, &current); err != nil {
				if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
					return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
				}
				return ctrl.Result{}, fmt.Errorf("error deleting StatefulSet \"%s\": %v", current.Name, err)
			}

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
		} else {
			if err := removeStaleWorkloadContainers(ctx, c, &current, &current.Spec.Template.Spec, &desired.Spec.Template.Spec); err != nil {
				if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
					return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
				}
				return ctrl.Result{}, fmt.Errorf("error removing stale containers from StatefulSet \"%s\": %v", current.Name, err)
			}
		}
	}

	if err := serverSideApply(ctx, c, &desired, "StatefulSet", "apps/v1"); err != nil {
		if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
			return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
		}
		return ctrl.Result{}, fmt.Errorf("error applying StatefulSet \"%s\": %v", desired.Name, err)
	}

	return ctrl.Result{}, nil
}

// ReconcilePersistentVolumeClaim uses the controller-runtime client to determine if the PVC resource needs to be created or updated
func ReconcilePersistentVolumeClaim(ctx context.Context, c client.Client, storage *k8schianetv1.StorageConfig, desired corev1.PersistentVolumeClaim) (reconcile.Result, error) {
	if !ShouldMakeChiaRootVolumeClaim(storage) {
		return ctrl.Result{}, nil
	}

	if err := serverSideApply(ctx, c, &desired, "PersistentVolumeClaim", "v1"); err != nil {
		if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
			return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
		}
		return ctrl.Result{}, fmt.Errorf("error applying PersistentVolumeClaim \"%s\": %v", desired.Name, err)
	}

	return ctrl.Result{}, nil
}

// ReconcileConfigMap uses the controller-runtime client to determine if the ConfigMap resource needs to be created or updated
func ReconcileConfigMap(ctx context.Context, c client.Client, desired corev1.ConfigMap) (reconcile.Result, error) {
	if err := serverSideApply(ctx, c, &desired, "ConfigMap", "v1"); err != nil {
		if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
			return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
		}
		return ctrl.Result{}, fmt.Errorf("error applying ConfigMap \"%s\": %v", desired.Name, err)
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

	if ensureIngressExists {
		if err := serverSideApply(ctx, c, &desired, "Ingress", "networking.k8s.io/v1"); err != nil {
			if strings.Contains(err.Error(), ObjectModifiedTryAgainError) {
				return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
			}
			return ctrl.Result{}, fmt.Errorf("error applying Ingress \"%s\": %v", desired.Name, err)
		}
		return ctrl.Result{}, nil
	}

	var current networkingv1.Ingress
	err := c.Get(ctx, types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, &current)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("error getting Ingress \"%s\": %v", desired.Name, err)
	}

	klog.Info("Deleting Ingress because it was disabled")
	if err := c.Delete(ctx, &current); err != nil {
		return ctrl.Result{}, fmt.Errorf("error deleting Ingress \"%s\": %v", desired.Name, err)
	}

	return ctrl.Result{}, nil
}
