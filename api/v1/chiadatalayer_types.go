/*
Copyright 2024 Chia Network Inc.
*/

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChiaDataLayerSpec defines the desired state of ChiaDataLayer
type ChiaDataLayerSpec struct {
	CommonSpec `json:",inline"`

	// ChiaConfig defines the configuration options available to Chia component containers
	ChiaConfig ChiaDataLayerSpecChia `json:"chia"`

	// DataLayerHTTPConfig defines the desired state of an optional data_layer_http sidecar
	DataLayerHTTPConfig ChiaDataLayerHTTPSpecChia `json:"dataLayerHTTP"`

	// Strategy describes how to replace existing pods with new ones.
	// +optional
	Strategy *appsv1.DeploymentStrategy `json:"strategy,omitempty"`
}

// ChiaDataLayerSpecChia defines the desired state of Chia component configuration
type ChiaDataLayerSpecChia struct {
	CommonSpecChia `json:",inline"`

	// CASecretName is the name of the secret that contains the CA crt and key.
	// +optional
	CASecretName *string `json:"caSecretName"`

	// SecretKey defines the k8s Secret name and key for a Chia mnemonic
	SecretKey ChiaSecretKey `json:"secretKey"`

	// FullNodePeers is a list of hostnames/IPs and port numbers to full_node peers.
	// Either fullNodePeer or fullNodePeers should be specified. fullNodePeers takes precedence.
	// +optional
	FullNodePeers *[]Peer `json:"fullNodePeers,omitempty"`

	// TrustedCIDRs is a list of CIDRs that this chia component should trust peers from
	// See: https://docs.chia.net/faq/?_highlight=trust#what-are-trusted-peers-and-how-do-i-add-them
	// +optional
	TrustedCIDRs *[]string `json:"trustedCIDRs,omitempty"`
}

// ChiaDataLayerHTTPSpecChia defines the desired state of an optional data_layer_http sidecar
// data_layer_http is a chia component, and therefore inherits most of the generic configuration options for any chia component
type ChiaDataLayerHTTPSpecChia struct {
	CommonSpecChia `json:",inline"`

	// Enabled defines whether a data_layer_http sidecar container should run as a sidecar to the chia container
	// +kubebuilder:default=true
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// CASecretName is the name of the secret that contains the CA crt and key.
	// +optional
	CASecretName *string `json:"caSecretName"`

	// Service defines settings for the Service optionally installed with any data_layer_http resource.
	// This Service will default to being enabled with a ClusterIP Service type if data_layer_http is enabled.
	// +optional
	Service Service `json:"service,omitempty"`
}

// ChiaDataLayerStatus defines the observed state of ChiaDataLayer
type ChiaDataLayerStatus struct {
	// Ready says whether the chia component is ready, this should be true when the data_layer resource is in the target namespace
	// +kubebuilder:default=false
	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ChiaDataLayer is the Schema for the chiadatalayers API
type ChiaDataLayer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChiaDataLayerSpec   `json:"spec,omitempty"`
	Status ChiaDataLayerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ChiaDataLayerList contains a list of ChiaDataLayer
type ChiaDataLayerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChiaDataLayer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChiaDataLayer{}, &ChiaDataLayerList{})
}
