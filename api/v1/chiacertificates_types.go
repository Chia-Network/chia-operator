/*
Copyright 2025 Chia Network Inc.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChiaCertificatesSpec defines the desired state of ChiaCertificates.
type ChiaCertificatesSpec struct {
	// CASecretName is the name of a Secret in the same namespace that contains the private Chia CA
	CASecretName string `json:"caSecretName"`
}

// ChiaCertificatesStatus defines the observed state of ChiaCertificates.
type ChiaCertificatesStatus struct {
	// Ready says whether the ChiaCertificates is ready, this should be true when the SSL secret is in the target namespace
	// +kubebuilder:default=false
	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ChiaCertificates is the Schema for the chiacertificates API.
type ChiaCertificates struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChiaCertificatesSpec   `json:"spec,omitempty"`
	Status ChiaCertificatesStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ChiaCertificatesList contains a list of ChiaCertificates.
type ChiaCertificatesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChiaCertificates `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChiaCertificates{}, &ChiaCertificatesList{})
}
