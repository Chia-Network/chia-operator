/*
Copyright 2023 Chia Network Inc.
*/

package kube

import (
	"context"
	"fmt"
	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/chia-network/chia-operator/internal/controller/common/consts"
)

// GetCommonLabels gives some common labels for chia-operator related objects
func GetCommonLabels(ctx context.Context, kind string, meta metav1.ObjectMeta, additionalLabels ...map[string]string) map[string]string {
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

// GetChiaExporterContainer assembles a chia-exporter container spec
func GetChiaExporterContainer(ctx context.Context, image string, secContext *corev1.SecurityContext, pullPolicy corev1.PullPolicy, resReq corev1.ResourceRequirements) corev1.Container {
	return corev1.Container{
		Name:            "chia-exporter",
		SecurityContext: secContext,
		Image:           image,
		ImagePullPolicy: pullPolicy,
		Env: []corev1.EnvVar{
			{
				Name:  "CHIA_ROOT",
				Value: "/chia-data",
			},
		},
		Ports: []corev1.ContainerPort{
			{
				Name:          "metrics",
				ContainerPort: consts.ChiaExporterPort,
				Protocol:      "TCP",
			},
		},
		LivenessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/healthz",
					Port: intstr.FromInt(consts.ChiaExporterPort),
				},
			},
		},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/healthz",
					Port: intstr.FromInt(consts.ChiaExporterPort),
				},
			},
		},
		StartupProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/healthz",
					Port: intstr.FromInt(consts.ChiaExporterPort),
				},
			},
			FailureThreshold: 30,
			PeriodSeconds:    10,
		},
		Resources: resReq,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "chiaroot",
				MountPath: "/chia-data",
			},
		},
	}
}

// ShouldMakeService returns true if the related Service was configured to be made
func ShouldMakeService(srv *k8schianetv1.Service) bool {
	if srv != nil && srv.Enabled != nil {
		return *srv.Enabled
	}
	return true // default to true if the Service wasn't declared
}
