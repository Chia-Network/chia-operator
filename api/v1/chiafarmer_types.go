/*
Copyright 2023 Chia Network Inc.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChiaFarmerSpec defines the desired state of ChiaFarmer
type ChiaFarmerSpec struct {
	CommonSpec `json:",inline"`

	// ChiaConfig defines the configuration options available to Chia component containers
	ChiaConfig ChiaFarmerConfigSpec `json:"chia"`
}

// ChiaFarmerConfigSpec defines the desired state of Chia component configuration
type ChiaFarmerConfigSpec struct {
	CommonChiaConfigSpec `json:",inline"`

	// SecretKeySpec defines the k8s Secret name and key for a Chia mnemonic
	SecretKeySpec ChiaKeysSpec `json:"secretKey"`

	// FullNodePeer defines the farmer's full_node peer in host:port format.
	// In Kubernetes this is likely to be <node service name>.<namespace>.svc.cluster.local:8555
	FullNodePeer string `json:"fullNodePeer"`
}

// ChiaFarmerStatus defines the observed state of ChiaFarmer
type ChiaFarmerStatus struct {
	// Ready says whether the node is ready, this should be true when the node statefulset is in the target namespace
	// +kubebuilder:default=false
	Ready bool `json:"ready,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ChiaFarmer is the Schema for the chiafarmers API
type ChiaFarmer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChiaFarmerSpec   `json:"spec,omitempty"`
	Status ChiaFarmerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ChiaFarmerList contains a list of ChiaFarmer
type ChiaFarmerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChiaFarmer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChiaFarmer{}, &ChiaFarmerList{})
}
