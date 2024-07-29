/*
Copyright 2023 Chia Network Inc.
*/

package kube

import (
	"fmt"
	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetCommonLabels gives some common labels for chia-operator related objects
func GetCommonLabels(kind string, meta metav1.ObjectMeta, additionalLabels ...map[string]string) map[string]string {
	var labels = make(map[string]string)
	labels = CombineMaps(additionalLabels...)
	labels["app.kubernetes.io/instance"] = meta.Name
	labels["app.kubernetes.io/name"] = meta.Name
	labels["app.kubernetes.io/managed-by"] = "chia-operator"
	labels["k8s.chia.net/provenance"] = fmt.Sprintf("%s.%s.%s", kind, meta.Namespace, meta.Name)
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

func GetChiaExporterServicePorts() []corev1.ServicePort {
	return []corev1.ServicePort{
		{
			Port:       consts.ChiaExporterPort,
			TargetPort: intstr.FromString("metrics"),
			Protocol:   "TCP",
			Name:       "metrics",
		},
	}
}

func GetChiaDaemonServicePorts() []corev1.ServicePort {
	return []corev1.ServicePort{
		{
			Port:       consts.DaemonPort,
			TargetPort: intstr.FromString("daemon"),
			Protocol:   "TCP",
			Name:       "daemon",
		},
	}
}
