/*
SPDX-License-Identifier: AGPL-3.0-or-later
KubeWG - Wireguard in your Kubernetes cluster
Copyright (C) 2024 Jacob McSwain

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

The source code is available at <https://github.com/USA-RedDragon/kubewg>
*/

package controller

import (
	"context"

	kubewgv1 "github.com/USA-RedDragon/kubewg/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// RouterReconciler reconciles a Router object
type RouterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kubewg.mcswain.dev,resources=routers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubewg.mcswain.dev,resources=routers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubewg.mcswain.dev,resources=routers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *RouterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling Router")

	var router kubewgv1.Router
	if err := r.Get(ctx, req.NamespacedName, &router); err != nil {
		log.Error(err, "unable to fetch Router")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Router", "Router", router)

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RouterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubewgv1.Router{}).
		Complete(r)
}
