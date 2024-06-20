/*
Copyright 2023 Chia Network Inc.
*/

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChiaNodeSpec defines the desired state of ChiaNode
type ChiaNodeSpec struct {
	CommonSpec `json:",inline"`

	// ChiaConfig defines the configuration options available to Chia component containers
	ChiaConfig ChiaNodeSpecChia `json:"chia"`

	// Replicas is the desired number of replicas of the given Statefulset. defaults to 1.
	// +optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// UpdateStrategy indicates the strategy that the StatefulSet controller will use to perform updates.
	// +optional
	UpdateStrategy *appsv1.StatefulSetUpdateStrategy `json:"updateStrategy,omitempty"`
}

// ChiaNodeSpecChia defines the desired state of Chia component configuration
type ChiaNodeSpecChia struct {
	CommonSpecChia `json:",inline"`

	// CASecretName is the name of the secret that contains the CA crt and key. Not required for seeders.
	CASecretName string `json:"caSecretName"`

	// TrustedCIDRs is a list of CIDRs that this chia component should trust peers from
	// See: https://docs.chia.net/faq/?_highlight=trust#what-are-trusted-peers-and-how-do-i-add-them
	// +optional
	TrustedCIDRs *[]string `json:"trustedCIDRs,omitempty"`
}

// ChiaNodeStatus defines the observed state of ChiaNode
type ChiaNodeStatus struct {
	// Ready says whether the node is ready, this should be true when the node statefulset is in the target namespace
	// +kubebuilder:default=false
	Ready bool `json:"ready,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ChiaNode is the Schema for the chianodes API
type ChiaNode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChiaNodeSpec   `json:"spec,omitempty"`
	Status ChiaNodeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ChiaNodeList contains a list of ChiaNode
type ChiaNodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChiaNode `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChiaNode{}, &ChiaNodeList{})
}
