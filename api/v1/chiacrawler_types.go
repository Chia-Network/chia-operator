/*
Copyright 2024 Chia Network Inc.
*/

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChiaCrawlerSpec defines the desired state of ChiaCrawler
type ChiaCrawlerSpec struct {
	CommonSpec `json:",inline"`

	// ChiaConfig defines the configuration options available to Chia component containers
	ChiaConfig ChiaCrawlerSpecChia `json:"chia"`

	// Strategy describes how to replace existing pods with new ones.
	// +optional
	Strategy *appsv1.DeploymentStrategy `json:"strategy,omitempty"`
}

// ChiaCrawlerSpecChia defines the desired state of Chia component configuration
type ChiaCrawlerSpecChia struct {
	CommonSpecChia `json:",inline"`

	// CASecretName is the name of the secret that contains the CA crt and key. Not required for seeders.
	// +optional
	CASecretName *string `json:"caSecretName"`
}

// ChiaCrawlerStatus defines the observed state of ChiaCrawler
type ChiaCrawlerStatus struct {
	// Ready says whether the chia component is ready deployed
	// +kubebuilder:default=false
	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ChiaCrawler is the Schema for the chiacrawlers API
type ChiaCrawler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChiaCrawlerSpec   `json:"spec,omitempty"`
	Status ChiaCrawlerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ChiaCrawlerList contains a list of ChiaCrawler
type ChiaCrawlerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChiaCrawler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChiaCrawler{}, &ChiaCrawlerList{})
}
