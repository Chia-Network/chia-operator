/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Copyright 2023 Chia Network Inc.
*/

package chiaseeder

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

// ChiaSeederReconciler reconciles a ChiaSeeder object
type ChiaSeederReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var chiaseeders = make(map[string]bool)

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiaseeders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiaseeders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiaseeders/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaSeederReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log))
	log.Info(fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s running reconciler...", req.NamespacedName.String()))

	// Get the custom resource
	var seeder k8schianetv1.ChiaSeeder
	err := r.Get(ctx, req.NamespacedName, &seeder)
	if err != nil && errors.IsNotFound(err) {
		// Remove this object from the map for tracking and subtract this CR's total metric by 1
		_, exists := chiaseeders[req.NamespacedName.String()]
		if exists {
			delete(chiaseeders, req.NamespacedName.String())
			metrics.ChiaSeeders.Sub(1.0)
		}
		return ctrl.Result{}, nil
	}
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		log.Error(err, fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s unable to fetch ChiaSeeder resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Add this object to the tracking map and increment the gauge by 1, if it wasn't already added
	_, exists := chiaseeders[req.NamespacedName.String()]
	if !exists {
		chiaseeders[req.NamespacedName.String()] = true
		metrics.ChiaSeeders.Add(1.0)
	}

	if kube.ShouldMakeService(seeder.Spec.ChiaConfig.PeerService, true) {
		srv := assemblePeerService(seeder)
		if err := controllerutil.SetControllerReference(&seeder, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to create seeder peer Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error reconciling seeder peer Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiaseederNamePattern, seeder.Name),
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				log.Error(err, fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s unable to GET ChiaSeeder peer Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				log.Error(err, fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s unable to DELETE ChiaSeeder peer Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(seeder.Spec.ChiaConfig.DaemonService, true) {
		srv := assembleDaemonService(seeder)
		if err := controllerutil.SetControllerReference(&seeder, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to create seeder daemon Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error reconciling seeder daemon Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiaseederNamePattern, seeder.Name) + "-daemon",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s unable to GET ChiaSeeder daemon Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s unable to DELETE ChiaSeeder daemon Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(seeder.Spec.ChiaConfig.RPCService, true) {
		srv := assembleRPCService(seeder)
		if err := controllerutil.SetControllerReference(&seeder, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to create seeder RPC Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error reconciling seeder RPC Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiaseederNamePattern, seeder.Name) + "-rpc",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s unable to GET ChiaSeeder RPC Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s unable to DELETE ChiaSeeder RPC Service resource", req.NamespacedName))
			}
		}
	}

	if kube.ShouldMakeService(seeder.Spec.ChiaExporterConfig.Service, true) {
		srv := assembleChiaExporterService(seeder)
		if err := controllerutil.SetControllerReference(&seeder, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to create seeder metrics Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error reconciling seeder metrics Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiaseederNamePattern, seeder.Name) + "-metrics",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s unable to GET ChiaSeeder metrics Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s unable to DELETE ChiaSeeder metrics Service resource", req.NamespacedName))
			}
		}
	}

	// Defaults the Service to false, and adds a check for the RollIntoPeerService parameter
	if kube.ShouldMakeService(seeder.Spec.ChiaHealthcheckConfig.Service, false) && !kube.ShouldRollIntoMainPeerService(seeder.Spec.ChiaHealthcheckConfig.Service) {
		srv := assembleChiaHealthcheckService(seeder)
		if err := controllerutil.SetControllerReference(&seeder, &srv, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to create seeder healthcheck Service -- Check operator logs.")
			return *res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error reconciling seeder healthcheck Service: %v", req.NamespacedName, err)
		}
	} else {
		// Need to check if the resource exists and delete if it does
		var srv corev1.Service
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      fmt.Sprintf(chiaseederNamePattern, seeder.Name) + "-healthcheck",
		}, &srv)
		if err != nil {
			if !errors.IsNotFound(err) {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s unable to GET ChiaSeeder healthcheck Service resource", req.NamespacedName))
			}
		} else {
			err = r.Delete(ctx, &srv)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				log.Error(err, fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s unable to DELETE ChiaSeeder healthcheck Service resource", req.NamespacedName))
			}
		}
	}

	// Creates a persistent volume claim if the GenerateVolumeClaims setting was set to true
	if seeder.Spec.Storage != nil && seeder.Spec.Storage.ChiaRoot != nil && seeder.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil && seeder.Spec.Storage.ChiaRoot.PersistentVolumeClaim.GenerateVolumeClaims {
		pvc, err := assembleVolumeClaim(seeder)
		if err != nil {
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to create seeder PVC -- Check operator logs.")
			return reconcile.Result{}, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error scaffolding a generated PersistentVolumeClaim: %v", req.NamespacedName, err)
		}

		res, err := kube.ReconcilePersistentVolumeClaim(ctx, resourceReconciler, pvc)
		if err != nil {
			if res == nil {
				res = &reconcile.Result{}
			}
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to create seeder PVC -- Check operator logs.")
			return *res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error reconciling PersistentVolumeClaim: %v", req.NamespacedName, err)
		}
	}

	deploy := assembleDeployment(seeder)

	if err := controllerutil.SetControllerReference(&seeder, &deploy, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	res, err := kube.ReconcileDeployment(ctx, resourceReconciler, deploy)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to create seeder Deployment -- Check operator logs.")
		return *res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error reconciling Deployment: %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&seeder, corev1.EventTypeNormal, "Created", "Successfully created ChiaSeeder resources.")
	seeder.Status.Ready = true
	err = r.Status().Update(ctx, &seeder)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		log.Error(err, fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s unable to update ChiaSeeder status", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaSeederReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaSeeder{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
