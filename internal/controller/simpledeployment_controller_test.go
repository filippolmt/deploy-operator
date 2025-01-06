/*
Copyright 2024.

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

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1alpha1 "github.com/filippolmt/deploy-operator/api/v1alpha1"
)

func int32Ptr(i int32) *int32 {
	return &i
}

var _ = Describe("SimpleDeployment Controller", func() {
	Context("When reconciling a resource", func() {
		const (
			resourceName      = "test-resource"
			resourceNamespace = "default"
		)

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: resourceNamespace,
		}
		simpledeployment := &appsv1alpha1.SimpleDeployment{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind SimpleDeployment")
			err := k8sClient.Get(ctx, typeNamespacedName, simpledeployment)
			if err != nil && errors.IsNotFound(err) {
				resource := &appsv1alpha1.SimpleDeployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: resourceNamespace,
					},
					Spec: appsv1alpha1.SimpleDeploymentSpec{
						Replicas: int32Ptr(1),
						Image:    "nginx:latest",
						Port:     80,
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			By("deleting the custom resource for the Kind SimpleDeployment")
			resource := &appsv1alpha1.SimpleDeployment{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			if err == nil {
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}
		})

		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &SimpleDeploymentReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Verifying the Deployment has been created")
			deployment := &appsv1.Deployment{}
			err = k8sClient.Get(ctx, typeNamespacedName, deployment)
			Expect(err).NotTo(HaveOccurred())

			// Verifica campi del Deployment
			Expect(*deployment.Spec.Replicas).To(Equal(int32(1)))
			Expect(deployment.Spec.Template.Spec.Containers).To(HaveLen(1))
			container := deployment.Spec.Template.Spec.Containers[0]
			Expect(container.Image).To(Equal("nginx:latest"))
			Expect(container.Ports).To(HaveLen(1))
			Expect(container.Ports[0].ContainerPort).To(Equal(int32(80)))

			By("Verifying the Deployment is owned by SimpleDeployment")
			ownerReferences := deployment.OwnerReferences
			Expect(ownerReferences).To(HaveLen(1))
			Expect(ownerReferences[0].Kind).To(Equal("SimpleDeployment"))
			Expect(ownerReferences[0].Name).To(Equal(resourceName))
		})
	})
})
