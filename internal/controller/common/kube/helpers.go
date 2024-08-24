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
	labels := CombineMaps(additionalLabels...)
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

// ShouldMakeVolumeClaim returns true if the related PersistentVolumeClaim was configured to be made
func ShouldMakeVolumeClaim(storage *k8schianetv1.StorageConfig) bool {
	if storage != nil && storage.ChiaRoot != nil && storage.ChiaRoot.PersistentVolumeClaim != nil && storage.ChiaRoot.PersistentVolumeClaim.GenerateVolumeClaims {
		return storage.ChiaRoot.PersistentVolumeClaim.GenerateVolumeClaims
	}
	return false
}

// ShouldMakeService returns true if the related Service was configured to be made, otherwise returns the specified default value
func ShouldMakeService(srv k8schianetv1.Service, def bool) bool {
	if srv.Enabled != nil {
		return *srv.Enabled
	}
	return def
}

// ShouldRollIntoMainPeerService returns true if the related Service's ports were meant to be rolled into the main peer Service's ports
func ShouldRollIntoMainPeerService(srv k8schianetv1.Service) bool {
	if srv.Enabled != nil && *srv.Enabled && srv.RollIntoPeerService != nil && *srv.RollIntoPeerService {
		return true
	}
	return false
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

// GetChiaHealthcheckServicePorts returns the Service ports for chia-healthcheck Services
func GetChiaHealthcheckServicePorts() []corev1.ServicePort {
	return []corev1.ServicePort{
		{
			Port:       consts.ChiaHealthcheckPort,
			TargetPort: intstr.FromString("health"),
			Protocol:   "TCP",
			Name:       "health",
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
