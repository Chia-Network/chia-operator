//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AdditionalMetadata) DeepCopyInto(out *AdditionalMetadata) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdditionalMetadata.
func (in *AdditionalMetadata) DeepCopy() *AdditionalMetadata {
	if in == nil {
		return nil
	}
	out := new(AdditionalMetadata)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaCA) DeepCopyInto(out *ChiaCA) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaCA.
func (in *ChiaCA) DeepCopy() *ChiaCA {
	if in == nil {
		return nil
	}
	out := new(ChiaCA)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChiaCA) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaCAList) DeepCopyInto(out *ChiaCAList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ChiaCA, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaCAList.
func (in *ChiaCAList) DeepCopy() *ChiaCAList {
	if in == nil {
		return nil
	}
	out := new(ChiaCAList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChiaCAList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaCASpec) DeepCopyInto(out *ChiaCASpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaCASpec.
func (in *ChiaCASpec) DeepCopy() *ChiaCASpec {
	if in == nil {
		return nil
	}
	out := new(ChiaCASpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaCAStatus) DeepCopyInto(out *ChiaCAStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaCAStatus.
func (in *ChiaCAStatus) DeepCopy() *ChiaCAStatus {
	if in == nil {
		return nil
	}
	out := new(ChiaCAStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaFarmer) DeepCopyInto(out *ChiaFarmer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaFarmer.
func (in *ChiaFarmer) DeepCopy() *ChiaFarmer {
	if in == nil {
		return nil
	}
	out := new(ChiaFarmer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChiaFarmer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaFarmerList) DeepCopyInto(out *ChiaFarmerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ChiaFarmer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaFarmerList.
func (in *ChiaFarmerList) DeepCopy() *ChiaFarmerList {
	if in == nil {
		return nil
	}
	out := new(ChiaFarmerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChiaFarmerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaFarmerSpec) DeepCopyInto(out *ChiaFarmerSpec) {
	*out = *in
	in.CommonSpec.DeepCopyInto(&out.CommonSpec)
	in.ChiaConfig.DeepCopyInto(&out.ChiaConfig)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaFarmerSpec.
func (in *ChiaFarmerSpec) DeepCopy() *ChiaFarmerSpec {
	if in == nil {
		return nil
	}
	out := new(ChiaFarmerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaFarmerSpecChia) DeepCopyInto(out *ChiaFarmerSpecChia) {
	*out = *in
	in.CommonSpecChia.DeepCopyInto(&out.CommonSpecChia)
	out.SecretKey = in.SecretKey
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaFarmerSpecChia.
func (in *ChiaFarmerSpecChia) DeepCopy() *ChiaFarmerSpecChia {
	if in == nil {
		return nil
	}
	out := new(ChiaFarmerSpecChia)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaFarmerStatus) DeepCopyInto(out *ChiaFarmerStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaFarmerStatus.
func (in *ChiaFarmerStatus) DeepCopy() *ChiaFarmerStatus {
	if in == nil {
		return nil
	}
	out := new(ChiaFarmerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaHarvester) DeepCopyInto(out *ChiaHarvester) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaHarvester.
func (in *ChiaHarvester) DeepCopy() *ChiaHarvester {
	if in == nil {
		return nil
	}
	out := new(ChiaHarvester)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChiaHarvester) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaHarvesterList) DeepCopyInto(out *ChiaHarvesterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ChiaHarvester, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaHarvesterList.
func (in *ChiaHarvesterList) DeepCopy() *ChiaHarvesterList {
	if in == nil {
		return nil
	}
	out := new(ChiaHarvesterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChiaHarvesterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaHarvesterSpec) DeepCopyInto(out *ChiaHarvesterSpec) {
	*out = *in
	in.CommonSpec.DeepCopyInto(&out.CommonSpec)
	in.ChiaConfig.DeepCopyInto(&out.ChiaConfig)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaHarvesterSpec.
func (in *ChiaHarvesterSpec) DeepCopy() *ChiaHarvesterSpec {
	if in == nil {
		return nil
	}
	out := new(ChiaHarvesterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaHarvesterSpecChia) DeepCopyInto(out *ChiaHarvesterSpecChia) {
	*out = *in
	in.CommonSpecChia.DeepCopyInto(&out.CommonSpecChia)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaHarvesterSpecChia.
func (in *ChiaHarvesterSpecChia) DeepCopy() *ChiaHarvesterSpecChia {
	if in == nil {
		return nil
	}
	out := new(ChiaHarvesterSpecChia)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaHarvesterStatus) DeepCopyInto(out *ChiaHarvesterStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaHarvesterStatus.
func (in *ChiaHarvesterStatus) DeepCopy() *ChiaHarvesterStatus {
	if in == nil {
		return nil
	}
	out := new(ChiaHarvesterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaNode) DeepCopyInto(out *ChiaNode) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaNode.
func (in *ChiaNode) DeepCopy() *ChiaNode {
	if in == nil {
		return nil
	}
	out := new(ChiaNode)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChiaNode) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaNodeList) DeepCopyInto(out *ChiaNodeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ChiaNode, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaNodeList.
func (in *ChiaNodeList) DeepCopy() *ChiaNodeList {
	if in == nil {
		return nil
	}
	out := new(ChiaNodeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChiaNodeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaNodeSpec) DeepCopyInto(out *ChiaNodeSpec) {
	*out = *in
	in.CommonSpec.DeepCopyInto(&out.CommonSpec)
	in.ChiaConfig.DeepCopyInto(&out.ChiaConfig)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaNodeSpec.
func (in *ChiaNodeSpec) DeepCopy() *ChiaNodeSpec {
	if in == nil {
		return nil
	}
	out := new(ChiaNodeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaNodeSpecChia) DeepCopyInto(out *ChiaNodeSpecChia) {
	*out = *in
	in.CommonSpecChia.DeepCopyInto(&out.CommonSpecChia)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaNodeSpecChia.
func (in *ChiaNodeSpecChia) DeepCopy() *ChiaNodeSpecChia {
	if in == nil {
		return nil
	}
	out := new(ChiaNodeSpecChia)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaNodeStatus) DeepCopyInto(out *ChiaNodeStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaNodeStatus.
func (in *ChiaNodeStatus) DeepCopy() *ChiaNodeStatus {
	if in == nil {
		return nil
	}
	out := new(ChiaNodeStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaRootConfig) DeepCopyInto(out *ChiaRootConfig) {
	*out = *in
	if in.PersistentVolumeClaim != nil {
		in, out := &in.PersistentVolumeClaim, &out.PersistentVolumeClaim
		*out = new(PersistentVolumeClaimConfig)
		**out = **in
	}
	if in.HostPathVolume != nil {
		in, out := &in.HostPathVolume, &out.HostPathVolume
		*out = new(HostPathVolumeConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaRootConfig.
func (in *ChiaRootConfig) DeepCopy() *ChiaRootConfig {
	if in == nil {
		return nil
	}
	out := new(ChiaRootConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaSecretKey) DeepCopyInto(out *ChiaSecretKey) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaSecretKey.
func (in *ChiaSecretKey) DeepCopy() *ChiaSecretKey {
	if in == nil {
		return nil
	}
	out := new(ChiaSecretKey)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaTimelord) DeepCopyInto(out *ChiaTimelord) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaTimelord.
func (in *ChiaTimelord) DeepCopy() *ChiaTimelord {
	if in == nil {
		return nil
	}
	out := new(ChiaTimelord)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChiaTimelord) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaTimelordList) DeepCopyInto(out *ChiaTimelordList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ChiaTimelord, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaTimelordList.
func (in *ChiaTimelordList) DeepCopy() *ChiaTimelordList {
	if in == nil {
		return nil
	}
	out := new(ChiaTimelordList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChiaTimelordList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaTimelordSpec) DeepCopyInto(out *ChiaTimelordSpec) {
	*out = *in
	in.CommonSpec.DeepCopyInto(&out.CommonSpec)
	in.ChiaConfig.DeepCopyInto(&out.ChiaConfig)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaTimelordSpec.
func (in *ChiaTimelordSpec) DeepCopy() *ChiaTimelordSpec {
	if in == nil {
		return nil
	}
	out := new(ChiaTimelordSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaTimelordSpecChia) DeepCopyInto(out *ChiaTimelordSpecChia) {
	*out = *in
	in.CommonSpecChia.DeepCopyInto(&out.CommonSpecChia)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaTimelordSpecChia.
func (in *ChiaTimelordSpecChia) DeepCopy() *ChiaTimelordSpecChia {
	if in == nil {
		return nil
	}
	out := new(ChiaTimelordSpecChia)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaTimelordStatus) DeepCopyInto(out *ChiaTimelordStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaTimelordStatus.
func (in *ChiaTimelordStatus) DeepCopy() *ChiaTimelordStatus {
	if in == nil {
		return nil
	}
	out := new(ChiaTimelordStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaWallet) DeepCopyInto(out *ChiaWallet) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaWallet.
func (in *ChiaWallet) DeepCopy() *ChiaWallet {
	if in == nil {
		return nil
	}
	out := new(ChiaWallet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChiaWallet) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaWalletList) DeepCopyInto(out *ChiaWalletList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ChiaWallet, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaWalletList.
func (in *ChiaWalletList) DeepCopy() *ChiaWalletList {
	if in == nil {
		return nil
	}
	out := new(ChiaWalletList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ChiaWalletList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaWalletSpec) DeepCopyInto(out *ChiaWalletSpec) {
	*out = *in
	in.CommonSpec.DeepCopyInto(&out.CommonSpec)
	in.ChiaConfig.DeepCopyInto(&out.ChiaConfig)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaWalletSpec.
func (in *ChiaWalletSpec) DeepCopy() *ChiaWalletSpec {
	if in == nil {
		return nil
	}
	out := new(ChiaWalletSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaWalletSpecChia) DeepCopyInto(out *ChiaWalletSpecChia) {
	*out = *in
	in.CommonSpecChia.DeepCopyInto(&out.CommonSpecChia)
	out.SecretKey = in.SecretKey
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaWalletSpecChia.
func (in *ChiaWalletSpecChia) DeepCopy() *ChiaWalletSpecChia {
	if in == nil {
		return nil
	}
	out := new(ChiaWalletSpecChia)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ChiaWalletStatus) DeepCopyInto(out *ChiaWalletStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ChiaWalletStatus.
func (in *ChiaWalletStatus) DeepCopy() *ChiaWalletStatus {
	if in == nil {
		return nil
	}
	out := new(ChiaWalletStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonSpec) DeepCopyInto(out *CommonSpec) {
	*out = *in
	in.AdditionalMetadata.DeepCopyInto(&out.AdditionalMetadata)
	in.ChiaExporterConfig.DeepCopyInto(&out.ChiaExporterConfig)
	if in.Storage != nil {
		in, out := &in.Storage, &out.Storage
		*out = new(StorageConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.PodSecurityContext != nil {
		in, out := &in.PodSecurityContext, &out.PodSecurityContext
		*out = new(corev1.PodSecurityContext)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonSpec.
func (in *CommonSpec) DeepCopy() *CommonSpec {
	if in == nil {
		return nil
	}
	out := new(CommonSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonSpecChia) DeepCopyInto(out *CommonSpecChia) {
	*out = *in
	if in.Testnet != nil {
		in, out := &in.Testnet, &out.Testnet
		*out = new(bool)
		**out = **in
	}
	if in.Network != nil {
		in, out := &in.Network, &out.Network
		*out = new(string)
		**out = **in
	}
	if in.NetworkPort != nil {
		in, out := &in.NetworkPort, &out.NetworkPort
		*out = new(uint16)
		**out = **in
	}
	if in.IntroducerAddress != nil {
		in, out := &in.IntroducerAddress, &out.IntroducerAddress
		*out = new(string)
		**out = **in
	}
	if in.DNSIntroducerAddress != nil {
		in, out := &in.DNSIntroducerAddress, &out.DNSIntroducerAddress
		*out = new(string)
		**out = **in
	}
	if in.Timezone != nil {
		in, out := &in.Timezone, &out.Timezone
		*out = new(string)
		**out = **in
	}
	if in.LogLevel != nil {
		in, out := &in.LogLevel, &out.LogLevel
		*out = new(string)
		**out = **in
	}
	if in.LivenessProbe != nil {
		in, out := &in.LivenessProbe, &out.LivenessProbe
		*out = new(corev1.Probe)
		(*in).DeepCopyInto(*out)
	}
	if in.ReadinessProbe != nil {
		in, out := &in.ReadinessProbe, &out.ReadinessProbe
		*out = new(corev1.Probe)
		(*in).DeepCopyInto(*out)
	}
	if in.StartupProbe != nil {
		in, out := &in.StartupProbe, &out.StartupProbe
		*out = new(corev1.Probe)
		(*in).DeepCopyInto(*out)
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(corev1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.SecurityContext != nil {
		in, out := &in.SecurityContext, &out.SecurityContext
		*out = new(corev1.SecurityContext)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonSpecChia.
func (in *CommonSpecChia) DeepCopy() *CommonSpecChia {
	if in == nil {
		return nil
	}
	out := new(CommonSpecChia)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostPathVolumeConfig) DeepCopyInto(out *HostPathVolumeConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostPathVolumeConfig.
func (in *HostPathVolumeConfig) DeepCopy() *HostPathVolumeConfig {
	if in == nil {
		return nil
	}
	out := new(HostPathVolumeConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PersistentVolumeClaimConfig) DeepCopyInto(out *PersistentVolumeClaimConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PersistentVolumeClaimConfig.
func (in *PersistentVolumeClaimConfig) DeepCopy() *PersistentVolumeClaimConfig {
	if in == nil {
		return nil
	}
	out := new(PersistentVolumeClaimConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlotsConfig) DeepCopyInto(out *PlotsConfig) {
	*out = *in
	if in.PersistentVolumeClaim != nil {
		in, out := &in.PersistentVolumeClaim, &out.PersistentVolumeClaim
		*out = make([]*PersistentVolumeClaimConfig, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(PersistentVolumeClaimConfig)
				**out = **in
			}
		}
	}
	if in.HostPathVolume != nil {
		in, out := &in.HostPathVolume, &out.HostPathVolume
		*out = make([]*HostPathVolumeConfig, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(HostPathVolumeConfig)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlotsConfig.
func (in *PlotsConfig) DeepCopy() *PlotsConfig {
	if in == nil {
		return nil
	}
	out := new(PlotsConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SpecChiaExporter) DeepCopyInto(out *SpecChiaExporter) {
	*out = *in
	if in.ServiceLabels != nil {
		in, out := &in.ServiceLabels, &out.ServiceLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SpecChiaExporter.
func (in *SpecChiaExporter) DeepCopy() *SpecChiaExporter {
	if in == nil {
		return nil
	}
	out := new(SpecChiaExporter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageConfig) DeepCopyInto(out *StorageConfig) {
	*out = *in
	if in.ChiaRoot != nil {
		in, out := &in.ChiaRoot, &out.ChiaRoot
		*out = new(ChiaRootConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.Plots != nil {
		in, out := &in.Plots, &out.Plots
		*out = new(PlotsConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageConfig.
func (in *StorageConfig) DeepCopy() *StorageConfig {
	if in == nil {
		return nil
	}
	out := new(StorageConfig)
	in.DeepCopyInto(out)
	return out
}
