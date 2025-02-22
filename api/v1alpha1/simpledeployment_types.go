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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SimpleDeploymentSpec defines the desired state of SimpleDeployment.
type SimpleDeploymentSpec struct {
	Replicas *int32 `json:"replicas,omitempty"`
	Image    string `json:"image"`
	Port     int32  `json:"port,omitempty"`
}

// SimpleDeploymentStatus defines the observed state of SimpleDeployment.
type SimpleDeploymentStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SimpleDeployment is the Schema for the simpledeployments API.
type SimpleDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SimpleDeploymentSpec   `json:"spec,omitempty"`
	Status SimpleDeploymentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SimpleDeploymentList contains a list of SimpleDeployment.
type SimpleDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SimpleDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SimpleDeployment{}, &SimpleDeploymentList{})
}
