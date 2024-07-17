/*
Copyright 2023 Chia Network Inc.
*/

package kube

import (
	"context"
	"fmt"
	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProvenanceLabelKey the key to the provenance label automatically added to operator managed resources
const ProvenanceLabelKey = "k8s.chia.net/provenance"

// GetCommonLabels gives some common labels for chia-operator related objects
func GetCommonLabels(ctx context.Context, kind string, meta metav1.ObjectMeta, additionalLabels ...map[string]string) map[string]string {
	var labels = make(map[string]string)
	labels = CombineMaps(additionalLabels...)
	labels["app.kubernetes.io/instance"] = meta.Name
	labels["app.kubernetes.io/name"] = meta.Name
	labels["app.kubernetes.io/managed-by"] = "chia-operator"
	labels[ProvenanceLabelKey] = fmt.Sprintf("%s.%s.%s", kind, meta.Namespace, meta.Name)
	return labels
}

// CombineMaps takes an arbitrary number of maps and combines them to one map[string]string
func CombineMaps(maps ...map[string]string) map[string]string {
	var keyvalues = make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			keyvalues[k] = v
		}
	}

	return keyvalues
}

// ShouldMakeService returns true if the related Service was configured to be made
func ShouldMakeService(srv *k8schianetv1.Service) bool {
	if srv != nil && srv.Enabled != nil {
		return *srv.Enabled
	}
	return true // default to true if the Service wasn't declared
}
