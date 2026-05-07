/*
Copyright 2023 Chia Network Inc.
*/

// Package v1 contains API Schema definitions for the k8s.chia.net v1 API group
// +kubebuilder:object:generate=true
// +groupName=k8s.chia.net
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "k8s.chia.net", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &builder{}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// builder is a thin wrapper around runtime.SchemeBuilder that preserves the
// concise `Register(&Foo{}, &FooList{})` call shape used in this api package's
// per-types init() functions. We avoid sigs.k8s.io/controller-runtime/pkg/scheme
// (deprecated) so that this api package stays cheap to import.
type builder struct {
	runtime.SchemeBuilder
}

// Register adds one or more types to the SchemeBuilder under GroupVersion and
// also ensures the standard meta/v1 types are added to the same group-version
// the first time the registered functions are applied to a Scheme.
func (b *builder) Register(objects ...runtime.Object) {
	b.SchemeBuilder.Register(func(scheme *runtime.Scheme) error {
		scheme.AddKnownTypes(GroupVersion, objects...)
		metav1.AddToGroupVersion(scheme, GroupVersion)
		return nil
	})
}
