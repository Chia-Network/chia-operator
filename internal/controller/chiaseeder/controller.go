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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
)

// ChiaSeederReconciler reconciles a ChiaSeeder object
type ChiaSeederReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiaseeders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiaseeders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.chia.net,resources=chiaseeders/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *ChiaSeederReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	resourceReconciler := reconciler.NewReconcilerWith(r.Client, reconciler.WithLog(log))
	log.Info(fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s", req.NamespacedName.String()))

	// Get the custom resource
	var seeder k8schianetv1.ChiaSeeder
	err := r.Get(ctx, req.NamespacedName, &seeder)
	if err != nil && errors.IsNotFound(err) {
		// Return here, this can happen if the CR was deleted
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s unable to fetch ChiaSeeder resource", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	srv := r.assembleBaseService(ctx, seeder)
	res, err := kube.ReconcileService(ctx, resourceReconciler, srv)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		return *res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error reconciling Service: %v", req.NamespacedName, err)
	}

	srv = r.assembleChiaExporterService(ctx, seeder)
	res, err = kube.ReconcileService(ctx, resourceReconciler, srv)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		return *res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error reconciling chia-exporter Service: %v", req.NamespacedName, err)
	}

	deploy := r.assembleDeployment(ctx, seeder)
	res, err = kube.ReconcileDeployment(ctx, resourceReconciler, deploy)
	if err != nil {
		if res == nil {
			res = &reconcile.Result{}
		}
		return *res, fmt.Errorf("ChiaSeederReconciler ChiaSeeder=%s encountered error reconciling StatefulSet: %v", req.NamespacedName, err)
	}

	// Update CR status
	seeder.Status.Ready = true
	err = r.Status().Update(ctx, &seeder)
	if err != nil {
		log.Error(err, fmt.Sprintf("ChiaSeederReconciler ChiaSeeder=%s unable to update ChiaSeeder status", req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChiaSeederReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8schianetv1.ChiaSeeder{}).
		Complete(r)
}
