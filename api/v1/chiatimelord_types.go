/*
Copyright 2023 Chia Network Inc.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChiaTimelordSpec defines the desired state of ChiaTimelord
type ChiaTimelordSpec struct {
	CommonSpec `json:",inline"`

	// ChiaConfig defines the configuration options available to Chia component containers
	ChiaConfig ChiaTimelordConfigSpec `json:"chia"`
}

// ChiaTimelordConfigSpec defines the desired state of Chia component configuration
type ChiaTimelordConfigSpec struct {
	CommonChiaConfigSpec `json:",inline"`

	// FullNodePeer defines the timelord's full_node peer in host:port format.
	// In Kubernetes this is likely to be <node service name>.<namespace>.svc.cluster.local:8555
	FullNodePeer string `json:"fullNodePeer"`
}

// ChiaTimelordStatus defines the observed state of ChiaTimelord
type ChiaTimelordStatus struct {
	// Ready says whether the CA is ready, this should be true when the SSL secret is in the target namespace
	// +kubebuilder:default=false
	Ready bool `json:"ready,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ChiaTimelord is the Schema for the chiatimelords API
type ChiaTimelord struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChiaTimelordSpec   `json:"spec,omitempty"`
	Status ChiaTimelordStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ChiaTimelordList contains a list of ChiaTimelord
type ChiaTimelordList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChiaTimelord `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChiaTimelord{}, &ChiaTimelordList{})
}
