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

The source code is available at <https://github.com/kubewg-net/operator>
*/

package controller

import (
	"context"

	kubewgv1 "github.com/kubewg-net/operator/api/v1"
	routerObj "github.com/kubewg-net/operator/internal/controller/objects/router"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// RouterReconciler reconciles a Router object
type RouterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	instanceLabel  = "app.kubernetes.io/instance"
	componentLabel = "app.kubernetes.io/component"
	managedByLabel = "app.kubernetes.io/managed-by"
	networkLabel   = "kubewg.net/network"
)

//+kubebuilder:rbac:groups=kubewg.net,resources=routers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubewg.net,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubewg.net,resources=networks,verbs=get;list;watch
//+kubebuilder:rbac:groups=kubewg.net,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubewg.net,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubewg.net,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubewg.net,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubewg.net,resources=routers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubewg.net,resources=routers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *RouterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling Router")

	var router kubewgv1.Router
	if err := r.Get(ctx, req.NamespacedName, &router); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "unable to fetch Router")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Router", "Router", router)

	matcherLabels := client.MatchingLabels{
		instanceLabel:  "kubewg-" + router.GetName() + "-" + string(router.GetUID()),
		componentLabel: "router",
		managedByLabel: "kubewg",
		networkLabel:   router.Spec.Network.Name,
	}

	log.Info("Matcher Labels", "Matcher Labels", matcherLabels)

	finalizerName := "kubewg.net/finalizer"

	// Object not deleted
	if router.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(&router, finalizerName) {
			controllerutil.AddFinalizer(&router, finalizerName)
			if err := r.Update(ctx, &router); err != nil {
				return ctrl.Result{}, err
			}
		}

		var associatedNetwork kubewgv1.Network
		if err := r.Get(ctx, client.ObjectKey{
			Name: router.Spec.Network.Name,
		}, &associatedNetwork); err != nil {
			log.Error(err, "unable to fetch associated Network")
			return ctrl.Result{}, err
		}

		// We need to manage 4 resources:
		// A secret config for the router
		// A network policy for the router
		// A service for the router
		// A deployment for the router
		err := r.reconcileSecret(ctx, req, &matcherLabels, &router, &associatedNetwork)
		if err != nil {
			log.Error(err, "unable to reconcile secret")
			return ctrl.Result{}, err
		}

		err = r.reconcileNetworkPolicy(ctx, req, &matcherLabels, &router, &associatedNetwork)
		if err != nil {
			log.Error(err, "unable to reconcile network policy")
			return ctrl.Result{}, err
		}

		err = r.reconcileService(ctx, req, &matcherLabels, &router)
		if err != nil {
			log.Error(err, "unable to reconcile service")
			return ctrl.Result{}, err
		}

		err = r.reconcileDeployment(ctx, req, &matcherLabels, &router, &associatedNetwork)
		if err != nil {
			log.Error(err, "unable to reconcile deployment")
			return ctrl.Result{}, err
		}
	} else {
		if controllerutil.ContainsFinalizer(&router, finalizerName) {
			// Remove deployment, service, network policy, and secret
			// Delete the deployment
			var deployments appsv1.DeploymentList
			err := r.List(ctx, &deployments, client.InNamespace(req.Namespace), matcherLabels)
			if err != nil {
				log.Error(err, "unable to list deployments")
				return ctrl.Result{}, err
			}

			// Check if the deployment exists
			if len(deployments.Items) != 0 {
				// Delete the deployment
				log.Info("Deleting deployment")
				if err := r.Delete(ctx, &deployments.Items[0]); err != nil {
					log.Error(err, "unable to delete deployment")
					return ctrl.Result{}, err
				}
			}

			// Delete the service
			var services corev1.ServiceList
			err = r.List(ctx, &services, client.InNamespace(req.Namespace), matcherLabels)
			if err != nil {
				log.Error(err, "unable to list services")
				return ctrl.Result{}, err
			}

			// Check if the service exists
			if len(services.Items) != 0 {
				// Delete the service
				log.Info("Deleting service")
				if err := r.Delete(ctx, &services.Items[0]); err != nil {
					log.Error(err, "unable to delete service")
					return ctrl.Result{}, err
				}
			}

			// Delete the network policy
			var networkPolicies networkingv1.NetworkPolicyList
			err = r.List(ctx, &networkPolicies, client.InNamespace(req.Namespace), matcherLabels)
			if err != nil {
				log.Error(err, "unable to list network policies")
				return ctrl.Result{}, err
			}

			// Check if the network policy exists
			if len(networkPolicies.Items) != 0 {
				// Delete the network policy
				log.Info("Deleting network policy")
				if err := r.Delete(ctx, &networkPolicies.Items[0]); err != nil {
					log.Error(err, "unable to delete network policy")
					return ctrl.Result{}, err
				}
			}

			// Delete the secret
			var secrets corev1.SecretList
			err = r.List(ctx, &secrets, client.InNamespace(req.Namespace), matcherLabels)
			if err != nil {
				log.Error(err, "unable to list secrets")
				return ctrl.Result{}, err
			}

			// Check if the secret exists
			if len(secrets.Items) != 0 {
				// Delete the secret
				log.Info("Deleting secret")
				if err := r.Delete(ctx, &secrets.Items[0]); err != nil {
					log.Error(err, "unable to delete secret")
					return ctrl.Result{}, err
				}
			}

			// Remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(&router, finalizerName)
			if err := r.Update(ctx, &router); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RouterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubewgv1.Router{}).
		Complete(r)
}

func (r *RouterReconciler) reconcileSecret(
	ctx context.Context,
	req ctrl.Request,
	labels *client.MatchingLabels,
	router *kubewgv1.Router,
	network *kubewgv1.Network) error {

	log := log.FromContext(ctx)

	var secrets corev1.SecretList
	err := r.List(ctx, &secrets, client.InNamespace(req.Namespace), labels)
	if err != nil {
		log.Error(err, "unable to list secrets")
		return err
	}

	if len(secrets.Items) > 1 {
		log.Info("Found more than one secret, deleting all")
		for _, cm := range secrets.Items {
			if err := r.Delete(ctx, &cm); err != nil {
				log.Error(err, "unable to delete secret")
				return err
			}
		}
		secrets.Items = []corev1.Secret{}
	}

	secret := routerObj.SecretFromRouterCRD(router, network, labels)
	// Check if the secret exists
	if len(secrets.Items) == 0 {
		// Create the secret
		log.Info("Creating secret")
		if err := r.Create(ctx, secret); err != nil {
			log.Error(err, "unable to create secret")
			return err
		}
	} else {
		// Update the secret
		log.Info("Updating secret")
		if err := r.Update(ctx, secret); err != nil {
			log.Error(err, "unable to update secret")
			return err
		}
	}

	return nil
}

func (r *RouterReconciler) reconcileNetworkPolicy(
	ctx context.Context,
	req ctrl.Request,
	labels *client.MatchingLabels,
	router *kubewgv1.Router,
	network *kubewgv1.Network) error {

	log := log.FromContext(ctx)

	var networkPolicies networkingv1.NetworkPolicyList
	err := r.List(ctx, &networkPolicies, client.InNamespace(req.Namespace), labels)
	if err != nil {
		log.Error(err, "unable to list network policies")
		return err
	}

	if len(networkPolicies.Items) > 1 {
		log.Info("Found more than one network policy, deleting all")
		for _, np := range networkPolicies.Items {
			if err := r.Delete(ctx, &np); err != nil {
				log.Error(err, "unable to delete network policy")
				return err
			}
		}
		networkPolicies.Items = []networkingv1.NetworkPolicy{}
	}

	networkPolicy := routerObj.NetworkPolicyFromRouterCRD(router, network, labels)
	// Check if the network policy exists
	if len(networkPolicies.Items) == 0 {
		// Create the network policy
		log.Info("Creating network policy")
		if err := r.Create(ctx, networkPolicy); err != nil {
			log.Error(err, "unable to create network policy")
			return err
		}
	} else {
		// Update the network policy
		log.Info("Updating network policy")
		if err := r.Update(ctx, networkPolicy); err != nil {
			log.Error(err, "unable to update network policy")
			return err
		}
	}

	return nil
}

func (r *RouterReconciler) reconcileService(
	ctx context.Context,
	req ctrl.Request,
	labels *client.MatchingLabels,
	router *kubewgv1.Router) error {

	log := log.FromContext(ctx)

	var services corev1.ServiceList
	err := r.List(ctx, &services, client.InNamespace(req.Namespace), labels)
	if err != nil {
		log.Error(err, "unable to list services")
		return err
	}

	if len(services.Items) > 1 {
		log.Info("Found more than one service, deleting all")
		for _, svc := range services.Items {
			if err := r.Delete(ctx, &svc); err != nil {
				log.Error(err, "unable to delete service")
				return err
			}
		}
		services.Items = []corev1.Service{}
	}

	service := routerObj.ServiceFromRouterCRD(router, labels)
	// Check if the service exists
	if len(services.Items) == 0 {
		// Create the service
		log.Info("Creating service")
		if err := r.Create(ctx, service); err != nil {
			log.Error(err, "unable to create service")
			return err
		}
	} else {
		// Update the service
		log.Info("Updating service")
		if err := r.Update(ctx, service); err != nil {
			log.Error(err, "unable to update service")
			return err
		}
	}

	return nil
}

func (r *RouterReconciler) reconcileDeployment(
	ctx context.Context,
	req ctrl.Request,
	labels *client.MatchingLabels,
	router *kubewgv1.Router,
	network *kubewgv1.Network) error {

	log := log.FromContext(ctx)

	var deployments appsv1.DeploymentList
	err := r.List(ctx, &deployments, client.InNamespace(req.Namespace), labels)
	if err != nil {
		log.Error(err, "unable to list deployments")
		return err
	}

	if len(deployments.Items) > 1 {
		log.Info("Found more than one deployment, deleting all")
		for _, dep := range deployments.Items {
			if err := r.Delete(ctx, &dep); err != nil {
				log.Error(err, "unable to delete deployment")
				return err
			}
		}
		deployments.Items = []appsv1.Deployment{}
	}

	deployment := routerObj.DeploymentFromRouterCRD(router, network, labels)
	// Check if the deployment exists
	if len(deployments.Items) == 0 {
		// Create the deployment
		log.Info("Creating deployment")
		if err := r.Create(ctx, deployment); err != nil {
			log.Error(err, "unable to create deployment")
			return err
		}
	} else {
		// Update the deployment
		log.Info("Updating deployment")
		if err := r.Update(ctx, deployment); err != nil {
			log.Error(err, "unable to update deployment")
			return err
		}
	}

	return nil
}
