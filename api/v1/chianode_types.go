/*
Copyright 2023 Chia Network Inc.
*/

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChiaNodeSpec defines the desired state of ChiaNode
type ChiaNodeSpec struct {
	CommonSpec `json:",inline"`

	// ChiaConfig defines the configuration options available to Chia component containers
	ChiaConfig ChiaNodeSpecChia `json:"chia"`

	// ChiaHealthcheckConfig defines the configuration options available to an optional Chia healthcheck sidecar
	// +optional
	ChiaHealthcheckConfig SpecChiaHealthcheck `json:"chiaHealthcheck,omitempty"`

	// ChiaDBPullConfig defines the configuration options available to an optional chia-db-pull init container
	// that downloads a chia blockchain database from an S3-compatible bucket into CHIA_ROOT before the chia
	// container starts.
	// +optional
	ChiaDBPullConfig SpecChiaDBPull `json:"chiaDBPull,omitempty"`

	// Replicas is the desired number of replicas of the given Statefulset. defaults to 1.
	// +optional
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// UpdateStrategy indicates the strategy that the StatefulSet controller will use to perform updates.
	// +optional
	UpdateStrategy *appsv1.StatefulSetUpdateStrategy `json:"updateStrategy,omitempty"`
}

// SpecChiaDBPull defines the desired state of an optional chia-db-pull init container
type SpecChiaDBPull struct {
	// Enabled defines whether a chia-db-pull init container should run before the chia container.
	// Defaults to false.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Image defines the image to use for the chia-db-pull init container
	// +optional
	Image *string `json:"image,omitempty"`

	// S3Prefix is the S3 URI prefix the chia-db-pull container will download the database from.
	// Required when Enabled is true. Mapped to the S3_PREFIX env var.
	// +optional
	S3Prefix string `json:"s3Prefix,omitempty"`

	// Network is the chia network name the database belongs to. Mapped to the NETWORK env var.
	// +optional
	Network *string `json:"network,omitempty"`

	// MinHeight is the minimum block height the downloaded database should be at. Mapped to the MIN_HEIGHT env var.
	// +optional
	MinHeight *int64 `json:"minHeight,omitempty"`

	// AWSCredentialsSecret is the name of a kubernetes Secret in the same namespace whose keys will be loaded
	// into the chia-db-pull container as environment variables via envFrom.
	// Use this to inject AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY (or any other credentials) without putting them in plaintext.
	// +optional
	AWSCredentialsSecret *string `json:"awsCredentialsSecret,omitempty"`

	// AdditionalEnv contain a list of additional environment variables to be supplied to the chia-db-pull container.
	// These variables will be placed at the end of the environment variable list in the resulting container,
	// this means they overwrite variables of the same name created by the operator in the container env.
	// +optional
	AdditionalEnv *[]corev1.EnvVar `json:"additionalEnv,omitempty"`

	// Resources defines the compute resources (limits/requests) for the chia-db-pull container.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// SecurityContext defines the security context for the chia-db-pull container
	// +optional
	SecurityContext *corev1.SecurityContext `json:"securityContext,omitempty"`
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

	// FullNodePeers is a list of hostnames/IPs and port numbers to full_node peers.
	// +optional
	FullNodePeers *[]Peer `json:"fullNodePeers,omitempty"`
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
