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

// ChiaDNSIntroducerSpec defines the desired state of ChiaDNSIntroducer
type ChiaDNSIntroducerSpec struct {
	CommonSpec `json:",inline"`

	// ChiaConfig defines the configuration options available to Chia component containers
	ChiaConfig ChiaDNSIntroducerSpecChia `json:"chia"`
}

// ChiaDNSIntroducerSpecChia defines the desired state of Chia component configuration
type ChiaDNSIntroducerSpecChia struct {
	CommonSpecChia `json:",inline"`

	// BootstrapPeer a peer to bootstrap the seeder's peer database
	BootstrapPeer string `json:"bootstrapPeer"`

	// MinimumHeight only consider nodes synced at least to this height
	MinimumHeight uint64 `json:"minimumHeight"`

	// DomainName the domain name of the server
	DomainName string `json:"domainName"`

	// Nameserver the adddress the dns server is running on
	Nameserver string `json:"nameserver"`

	// Rname an administrator's email address with '@' replaced with '.'
	Rname string `json:"rname"`
}

// ChiaDNSIntroducerStatus defines the observed state of ChiaDNSIntroducer
type ChiaDNSIntroducerStatus struct {
	// Ready says whether the chia component is ready deployed
	// +kubebuilder:default=false
	Ready bool `json:"ready,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ChiaDNSIntroducer is the Schema for the chiadnsintroducers API
type ChiaDNSIntroducer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChiaDNSIntroducerSpec   `json:"spec,omitempty"`
	Status ChiaDNSIntroducerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ChiaDNSIntroducerList contains a list of ChiaDNSIntroducer
type ChiaDNSIntroducerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChiaDNSIntroducer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChiaDNSIntroducer{}, &ChiaDNSIntroducerList{})
}
