/*
Copyright 2023 Chia Network Inc.
*/

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChiaWalletSpec defines the desired state of ChiaWallet
type ChiaWalletSpec struct {
	CommonSpec `json:",inline"`

	// ChiaConfig defines the configuration options available to Chia component containers
	ChiaConfig ChiaWalletSpecChia `json:"chia"`

	// Strategy describes how to replace existing pods with new ones.
	// +optional
	Strategy *appsv1.DeploymentStrategy `json:"strategy,omitempty"`
}

// ChiaWalletSpecChia defines the desired state of Chia component configuration
type ChiaWalletSpecChia struct {
	CommonSpecChia `json:",inline"`

	// CASecretName is the name of the secret that contains the CA crt and key. Not required for seeders.
	CASecretName string `json:"caSecretName"`

	// SecretKey defines the k8s Secret name and key for a Chia mnemonic
	SecretKey ChiaSecretKey `json:"secretKey"`

	// FullNodePeer defines the farmer's full_node peer in host:port format.
	// In Kubernetes this is likely to be <node service name>.<namespace>.svc.cluster.local:8555
	// Either fullNodePeer or fullNodePeers should be specified. fullNodePeers takes precedence.
	// Deprecated in favor of fullNodePeers.
	// +optional
	FullNodePeer *string `json:"fullNodePeer,omitempty"`

	// FullNodePeers is a list of hostnames/IPs and port numbers to full_node peers.
	// Either fullNodePeer or fullNodePeers should be specified. fullNodePeers takes precedence.
	// +optional
	FullNodePeers *[]FullNodePeer `json:"fullNodePeers,omitempty"`

	// TrustedCIDRs is a list of CIDRs that this chia component should trust peers from
	// See: https://docs.chia.net/faq/?_highlight=trust#what-are-trusted-peers-and-how-do-i-add-them
	// +optional
	TrustedCIDRs *[]string `json:"trustedCIDRs,omitempty"`
}

// ChiaWalletStatus defines the observed state of ChiaWallet
type ChiaWalletStatus struct {
	// Ready says whether the node is ready, this should be true when the node statefulset is in the target namespace
	// +kubebuilder:default=false
	Ready bool `json:"ready,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ChiaWallet is the Schema for the chiawallets API
type ChiaWallet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChiaWalletSpec   `json:"spec,omitempty"`
	Status ChiaWalletStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ChiaWalletList contains a list of ChiaWallet
type ChiaWalletList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChiaWallet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChiaWallet{}, &ChiaWalletList{})
}
