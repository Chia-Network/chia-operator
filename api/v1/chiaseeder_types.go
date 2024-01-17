/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChiaSeederSpec defines the desired state of ChiaSeeder
type ChiaSeederSpec struct {
	CommonSpec `json:",inline"`

	// ChiaConfig defines the configuration options available to Chia component containers
	ChiaConfig ChiaSeederSpecChia `json:"chia"`
}

// ChiaSeederSpecChia defines the desired state of Chia component configuration
type ChiaSeederSpecChia struct {
	CommonSpecChia `json:",inline"`

	// BootstrapPeer a peer to bootstrap the seeder's peer database
	// +optional
	BootstrapPeer *string `json:"bootstrapPeer"`

	// MinimumHeight only consider nodes synced at least to this height
	// +optional
	MinimumHeight *uint64 `json:"minimumHeight"`

	// DomainName the name of the NS record for your server with a trailing period. (ex. "seeder.example.com.")
	DomainName string `json:"domainName"`

	// Nameserver the name of the A record for your server with a trailing period. (ex. "seeder-us-west-2.example.com.")
	Nameserver string `json:"nameserver"`

	// TTL ttl setting in the seeder configuration
	// +optional
	TTL *uint32 `json:"ttl"`

	// Rname an administrator's email address with '@' replaced with '.'
	// +optional
	Rname *string `json:"rname"`
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
