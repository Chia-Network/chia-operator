/*
Copyright 2023 Chia Network Inc.
*/

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChiaSeederSpec defines the desired state of ChiaSeeder
type ChiaSeederSpec struct {
	CommonSpec `json:",inline"`

	// ChiaConfig defines the configuration options available to Chia component containers
	ChiaConfig ChiaSeederSpecChia `json:"chia"`

	// ChiaHealthcheckConfig defines the configuration options available to an optional Chia healthcheck sidecar
	// +optional
	ChiaHealthcheckConfig SpecChiaHealthcheck `json:"chiaHealthcheck,omitempty"`

	// Strategy describes how to replace existing pods with new ones.
	// +optional
	Strategy *appsv1.DeploymentStrategy `json:"strategy,omitempty"`
}

// ChiaSeederSpecChia defines the desired state of Chia component configuration
type ChiaSeederSpecChia struct {
	CommonSpecChia `json:",inline"`

	// BootstrapPeer a peer to bootstrap the seeder's peer database.
	// DEPRECATED: Use BootstrapPeers instead.
	// +optional
	BootstrapPeer *string `json:"bootstrapPeer,omitempty"`

	// BootstrapPeers a list of peers to bootstrap the seeder's peer database
	// +optional
	BootstrapPeers *[]string `json:"bootstrapPeers,omitempty"`

	// MinimumHeight only consider nodes synced at least to this height
	// +optional
	MinimumHeight *uint64 `json:"minimumHeight,omitempty"`

	// DomainName the name of the NS record for your server with a trailing period. (ex. "seeder.example.com.")
	DomainName string `json:"domainName"`

	// Nameserver the name of the A record for your server with a trailing period. (ex. "seeder-us-west-2.example.com.")
	Nameserver string `json:"nameserver"`

	// Rname an administrator's email address with '@' replaced with '.'
	Rname string `json:"rname"`

	// CASecretName is the name of the secret that contains the CA crt and key. Not required for seeders.
	// +optional
	CASecretName *string `json:"caSecretName,omitempty"`

	// TTL field on DNS records that controls the length of time that a record is considered valid
	// +optional
	TTL *uint32 `json:"ttl,omitempty"`
}

// ChiaSeederStatus defines the observed state of ChiaSeeder
type ChiaSeederStatus struct {
	// Ready says whether the chia component is ready deployed
	// +kubebuilder:default=false
	Ready bool `json:"ready,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ChiaSeeder is the Schema for the chiaseeders API
type ChiaSeeder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChiaSeederSpec   `json:"spec,omitempty"`
	Status ChiaSeederStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ChiaSeederList contains a list of ChiaSeeder
type ChiaSeederList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChiaSeeder `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChiaSeeder{}, &ChiaSeederList{})
}
