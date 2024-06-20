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

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ChiaIntroducerSpec defines the desired state of ChiaIntroducer
type ChiaIntroducerSpec struct {
	CommonSpec `json:",inline"`

	// ChiaConfig defines the configuration options available to Chia component containers
	ChiaConfig ChiaIntroducerSpecChia `json:"chia"`

	// Strategy describes how to replace existing pods with new ones.
	// +optional
	Strategy *appsv1.DeploymentStrategy `json:"strategy,omitempty"`
}

// ChiaIntroducerSpecChia defines the desired state of Chia component configuration
type ChiaIntroducerSpecChia struct {
	CommonSpecChia `json:",inline"`

	// CASecretName is the name of the secret that contains the CA crt and key. Not required for introducers.
	// +optional
	CASecretName *string `json:"caSecretName"`
}

// ChiaIntroducerStatus defines the observed state of ChiaIntroducer
type ChiaIntroducerStatus struct {
	// Ready says whether the node is ready, this should be true when the node statefulset is in the target namespace
	// +kubebuilder:default=false
	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ChiaIntroducer is the Schema for the chiaintroducers API
type ChiaIntroducer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChiaIntroducerSpec   `json:"spec,omitempty"`
	Status ChiaIntroducerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ChiaIntroducerList contains a list of ChiaIntroducer
type ChiaIntroducerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChiaIntroducer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChiaIntroducer{}, &ChiaIntroducerList{})
}
