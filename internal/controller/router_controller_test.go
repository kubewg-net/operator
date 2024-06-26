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

package controller_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubewgv1 "github.com/kubewg-net/operator/api/v1"
	"github.com/kubewg-net/operator/internal/controller"
)

var _ = Describe("Router Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default",
		}
		router := &kubewgv1.Router{}
		network := &kubewgv1.Network{}

		BeforeEach(func() {
			By("creating the test network for the Router")
			err := k8sClient.Get(ctx, types.NamespacedName{
				Name:      resourceName,
				Namespace: "default",
			}, network)
			if err != nil && errors.IsNotFound(err) {
				resource := &kubewgv1.Network{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}

			By("creating the custom resource for the Kind Router")
			err = k8sClient.Get(ctx, typeNamespacedName, router)
			if err != nil && errors.IsNotFound(err) {
				resource := &kubewgv1.Router{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: kubewgv1.RouterSpec{
						Network: kubewgv1.NameSelectorSpec{
							Name: resourceName,
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			router := &kubewgv1.Router{}
			err := k8sClient.Get(ctx, typeNamespacedName, router)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance Router")
			Expect(k8sClient.Delete(ctx, router)).To(Succeed())

			network = &kubewgv1.Network{}
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      resourceName,
				Namespace: "default",
			}, network)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the test network")
			Expect(k8sClient.Delete(ctx, network)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &controller.RouterReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})
