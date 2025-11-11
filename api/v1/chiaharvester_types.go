/*
Copyright 2023 Chia Network Inc.
*/

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChiaHarvesterSpec defines the desired state of ChiaHarvester
type ChiaHarvesterSpec struct {
	CommonSpec `json:",inline"`

	// ChiaConfig defines the configuration options available to Chia component containers
	ChiaConfig ChiaHarvesterSpecChia `json:"chia"`

	// ChiaHealthcheckConfig defines the configuration options available to an optional Chia healthcheck sidecar
	// +optional
	ChiaHealthcheckConfig SpecChiaHealthcheck `json:"chiaHealthcheck,omitempty"`

	// Strategy describes how to replace existing pods with new ones.
	// +optional
	Strategy *appsv1.DeploymentStrategy `json:"strategy,omitempty"`
}

// ChiaHarvesterSpecChia defines the desired state of Chia component configuration
type ChiaHarvesterSpecChia struct {
	CommonSpecChia `json:",inline"`

	// CASecretName is the name of the secret that contains the CA crt and key. Not required for seeders.
	CASecretName string `json:"caSecretName"`

	// FarmerAddress defines the harvester's farmer peer's hostname. The farmer's port is inferred.
	// In Kubernetes this is likely to be <farmer service name>.<namespace>.svc.cluster.local
	FarmerAddress string `json:"farmerAddress"`
}

// ChiaHarvesterStatus defines the observed state of ChiaHarvester
type ChiaHarvesterStatus struct {
	// Ready says whether the node is ready, this should be true when the node statefulset is in the target namespace
	// +kubebuilder:default=false
	Ready bool `json:"ready,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ChiaHarvester is the Schema for the chiaharvesters API
type ChiaHarvester struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChiaHarvesterSpec   `json:"spec,omitempty"`
	Status ChiaHarvesterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ChiaHarvesterList contains a list of ChiaHarvester
type ChiaHarvesterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChiaHarvester `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChiaHarvester{}, &ChiaHarvesterList{})
}
