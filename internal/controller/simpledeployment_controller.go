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

	appsv1alpha1 "github.com/filippolmt/deploy-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// SimpleDeploymentReconciler manages the reconciliation logic for SimpleDeployment.
type SimpleDeploymentReconciler struct {
	Client client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile implements the reconciliation logic for SimpleDeployment.
func (r *SimpleDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("simpledeployment", req.NamespacedName)
	log.Info("Reconciling SimpleDeployment")

	// Get the SimpleDeployment resource with the name specified in req.
	var simpleDeployment appsv1alpha1.SimpleDeployment
	if err := r.Client.Get(ctx, req.NamespacedName, &simpleDeployment); err != nil {
		if errors.IsNotFound(err) {
			// SimpleDeployment resource not found. Ignoring since object must be deleted.
			log.Info("SimpleDeployment resource not found. Ignoring since it must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get SimpleDeployment")
		return ctrl.Result{}, err
	}

	// Define a new Deployment object.
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      simpleDeployment.Name,
			Namespace: simpleDeployment.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: simpleDeployment.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": simpleDeployment.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": simpleDeployment.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "app",
							Image: simpleDeployment.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: simpleDeployment.Spec.Port,
								},
							},
						},
					},
				},
			},
		},
	}

	// Set the controller as the owner of the Deployment.
	if err := controllerutil.SetControllerReference(&simpleDeployment, deployment, r.Scheme); err != nil {
		log.Error(err, "Failed to set owner reference on Deployment")
		return ctrl.Result{}, err
	}

	// Check if the Deployment already exists.
	var existingDeployment appsv1.Deployment
	err := r.Client.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, &existingDeployment)
	if err != nil && errors.IsNotFound(err) {
		// Create a new Deployment
		log.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		if err := r.Client.Create(ctx, deployment); err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	} else {
		// Use Patch instead of overwriting the entire spec
		log.Info("Patching existing Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)

		// Prepare patch to update only necessary fields
		patch := client.MergeFrom(existingDeployment.DeepCopy())
		existingDeployment.Spec.Replicas = deployment.Spec.Replicas
		existingDeployment.Spec.Template = deployment.Spec.Template

		if err := r.Client.Patch(ctx, &existingDeployment, patch); err != nil {
			log.Error(err, "Failed to patch Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return ctrl.Result{}, err
		}
	}

	// Update SimpleDeployment status
	simpleDeployment.Status.AvailableReplicas = existingDeployment.Status.AvailableReplicas
	if err := r.Client.Status().Update(ctx, &simpleDeployment); err != nil {
		log.Error(err, "Failed to update SimpleDeployment status")
		return ctrl.Result{}, err
	}

	log.Info("Successfully reconciled SimpleDeployment")
	return ctrl.Result{}, nil
}

// SetupWithManager registers the controller with the manager.
func (r *SimpleDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.SimpleDeployment{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
