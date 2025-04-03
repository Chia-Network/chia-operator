/*
Copyright 2024 Chia Network Inc.
*/

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChiaDataLayerSpec defines the desired state of ChiaDataLayer
type ChiaDataLayerSpec struct {
	CommonSpec `json:",inline"`

	// ChiaConfig defines the configuration options available to Chia component containers
	ChiaConfig ChiaDataLayerSpecChia `json:"chia"`

	// FileserverConfig defines the desired state of an optional fileserver sidecar to server datalayer server files
	// +optional
	FileserverConfig FileserverConfig `json:"fileserver"`

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

// FileserverConfig defines the desired state of an optional fileserver sidecar
// data_layer_http is the default fileserver but can be configured to use nginx or any other webserver application
type FileserverConfig struct {
	// Enabled defines whether a fileserver container should run as a sidecar to the chia container.
	// Disabled by default.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Image defines the image (registry/name:tag) to use for the sidecar container.
	// Defaults to the official chia image.
	// +optional
	Image *string `json:"image,omitempty"`

	// ServerFileMountpath defines the mount path for the server files volume in the container.
	// The volume will be mounted as a read-only volume.
	// Defaults to "/datalayer/server".
	// +optional
	ServerFileMountpath *string `json:"serverFileMountpath,omitempty"`

	// ContainerPort defines the port of the http server in the container
	// Defaults to 8575.
	// NOTE: If you use a custom image for the fileserver make sure you set this to the port that the fileserver binds to in the container.
	// +optional
	ContainerPort *int `json:"containerPort,omitempty"`

	// Service defines settings for the Service optionally installed with any fileserver resource.
	// Defaults to being enabled with a ClusterIP Service type if fileserver is enabled.
	// +optional
	Service Service `json:"service,omitempty"`

	// Ingress defines settings for the Ingress optionally installed with any fileserver resource.
	// Defaults to being disabled.
	// +optional
	Ingress IngressConfig `json:"ingress,omitempty"`

	// AdditionalEnv contain a list of additional environment variables to be supplied to the chia container.
	// These variables will be placed at the end of the environment variable list in the resulting container,
	// this means they overwrite variables of the same name created by the operator in the container env.
	// +optional
	AdditionalEnv *[]corev1.EnvVar `json:"additionalEnv,omitempty"`

	// LivenessProbe used to determine if a container is running properly and will restart the container if the probe fails
	// +optional
	LivenessProbe *corev1.Probe `json:"livenessProbe,omitempty"`

	// ReadinessProbe used to indicate when a container is ready to accept traffic and prevent traffic from being sent to pods that aren't ready.
	// +optional
	ReadinessProbe *corev1.Probe `json:"readinessProbe,omitempty"`

	// StartupProbe used to give applications time to initialize fully before liveness and readiness probes begin checking, preventing premature restarts of slow-starting containers.
	// +optional
	StartupProbe *corev1.Probe `json:"startupProbe,omitempty"`

	// Resources defines the compute resources (limits/requests) for the fileserver container.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// SecurityContext defines the security context for the fileserver container.
	// +optional
	SecurityContext *corev1.SecurityContext `json:"securityContext,omitempty"`
}

// IngressConfig defines the configuration for a Kubernetes Ingress resource
type IngressConfig struct {
	AdditionalMetadata `json:",inline"`

	// Enabled defines whether an Ingress should be created for the fileserver
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// IngressClassName defines the IngressClass to use for this Ingress
	// +optional
	IngressClassName *string `json:"ingressClassName,omitempty"`

	// Host defines the hostname for the Ingress
	// +optional
	Host *string `json:"host,omitempty"`

	// TLS defines TLS configuration for the Ingress
	// +optional
	TLS *[]networkingv1.IngressTLS `json:"tls,omitempty"`

	// Rules defines the routing rules for the Ingress
	// +optional
	Rules *[]networkingv1.IngressRule `json:"rules,omitempty"`
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
