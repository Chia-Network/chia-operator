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
	"strings"
	"time"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	"github.com/chia-network/chia-operator/internal/metrics"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is invoked on any event to a controlled Kubernetes resource
func (r *ChiaSeederReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Running reconciler...")

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

	// Assemble Peer Service
	peerSrv := assemblePeerService(seeder)
	if err := controllerutil.SetControllerReference(&seeder, &peerSrv, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to assemble seeder peer Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error assembling peer Service: %v", req.NamespacedName, err)
	}
	// Reconcile Peer Service
	res, err := kube.ReconcileService(ctx, r.Client, seeder.Spec.ChiaConfig.PeerService, peerSrv, true)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		return res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s %v", req.NamespacedName, err)
	}

	// Assemble All Service
	allSrv := assembleAllService(seeder)
	if err := controllerutil.SetControllerReference(&seeder, &allSrv, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to assemble seeder all-port Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error assembling all-port Service: %v", req.NamespacedName, err)
	}
	// Reconcile All Service
	res, err = kube.ReconcileService(ctx, r.Client, seeder.Spec.ChiaConfig.AllService, allSrv, true)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		return res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s %v", req.NamespacedName, err)
	}

	// Assemble Daemon Service
	daemonSrv := assembleDaemonService(seeder)
	if err := controllerutil.SetControllerReference(&seeder, &daemonSrv, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to assemble seeder daemon Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error assembling daemon Service: %v", req.NamespacedName, err)
	}
	// Reconcile Daemon Service
	res, err = kube.ReconcileService(ctx, r.Client, seeder.Spec.ChiaConfig.DaemonService, daemonSrv, true)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		return res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s %v", req.NamespacedName, err)
	}

	// Assemble RPC Service
	rpcSrv := assembleRPCService(seeder)
	if err := controllerutil.SetControllerReference(&seeder, &rpcSrv, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to assemble seeder RPC Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error assembling RPC Service: %v", req.NamespacedName, err)
	}
	// Reconcile RPC Service
	res, err = kube.ReconcileService(ctx, r.Client, seeder.Spec.ChiaConfig.RPCService, rpcSrv, true)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		return res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s %v", req.NamespacedName, err)
	}

	// Assemble Chia-Exporter Service
	exporterSrv := assembleChiaExporterService(seeder)
	if err := controllerutil.SetControllerReference(&seeder, &exporterSrv, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to assemble seeder chia-exporter Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error assembling chia-exporter Service: %v", req.NamespacedName, err)
	}
	// Reconcile Chia-Exporter Service
	res, err = kube.ReconcileService(ctx, r.Client, seeder.Spec.ChiaExporterConfig.Service, exporterSrv, true)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		return res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s %v", req.NamespacedName, err)
	}

	// Assemble Chia-Healthcheck Service
	healthcheckSrv := assembleChiaHealthcheckService(seeder)
	if err := controllerutil.SetControllerReference(&seeder, &healthcheckSrv, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to assemble seeder chia-healthcheck Service -- Check operator logs.")
		return ctrl.Result{}, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error assembling chia-healthcheck Service: %v", req.NamespacedName, err)
	}
	// Reconcile Chia-Healthcheck Service
	if !kube.ShouldRollIntoMainPeerService(seeder.Spec.ChiaHealthcheckConfig.Service) {
		res, err = kube.ReconcileService(ctx, r.Client, seeder.Spec.ChiaHealthcheckConfig.Service, healthcheckSrv, false)
		if err != nil {
			metrics.OperatorErrors.Add(1.0)
			return res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s %v", req.NamespacedName, err)
		}
	}

	// Creates a persistent volume claim if the GenerateVolumeClaims setting was set to true
	if kube.ShouldMakeVolumeClaim(seeder.Spec.Storage) {
		pvc, err := assembleVolumeClaim(seeder)
		if err != nil {
			metrics.OperatorErrors.Add(1.0)
			r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to assemble seeder PVC -- Check operator logs.")
			return reconcile.Result{}, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s %v", req.NamespacedName, err)
		}

		if pvc != nil {
			res, err = kube.ReconcilePersistentVolumeClaim(ctx, r.Client, seeder.Spec.Storage, *pvc)
			if err != nil {
				metrics.OperatorErrors.Add(1.0)
				r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to create seeder PVC -- Check operator logs.")
				return res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s %v", req.NamespacedName, err)
			}
		} else {
			return reconcile.Result{}, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s PVC could not be created", req.NamespacedName)
		}
	}

	// Assemble Deployment
	deploy := assembleDeployment(seeder)
	if err := controllerutil.SetControllerReference(&seeder, &deploy, r.Scheme); err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to assemble seeder Deployment -- Check operator logs.")
		return reconcile.Result{}, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s %v", req.NamespacedName, err)
	}
	// Reconcile Deployment
	res, err = kube.ReconcileDeployment(ctx, r.Client, deploy)
	if err != nil {
		metrics.OperatorErrors.Add(1.0)
		r.Recorder.Event(&seeder, corev1.EventTypeWarning, "Failed", "Failed to create seeder Deployment -- Check operator logs.")
		return res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s %v", req.NamespacedName, err)
	}

	// Update CR status
	r.Recorder.Event(&seeder, corev1.EventTypeNormal, "Created", "Successfully created ChiaSeeder resources.")
	seeder.Status.Ready = true
	err = r.Status().Update(ctx, &seeder)
	if err != nil {
		if strings.Contains(err.Error(), kube.ObjectModifiedTryAgainError) {
			return ctrl.Result{RequeueAfter: 1 * time.Second}, nil
		}
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
