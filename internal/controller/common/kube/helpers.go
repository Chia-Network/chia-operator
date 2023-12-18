/*
Copyright 2023 Chia Network Inc.
*/

package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/chia-network/chia-operator/internal/controller/common/consts"
)

// GetCommonLabels gives some common labels for chia-operator related objects
func GetCommonLabels(ctx context.Context, meta metav1.ObjectMeta, additionalLabels ...map[string]string) map[string]string {
	var labels = make(map[string]string)
	for _, addition := range additionalLabels {
		for k, v := range addition {
			labels[k] = v
		}
	}
	labels["app.kubernetes.io/instance"] = meta.Name
	labels["app.kubernetes.io/name"] = meta.Name
	labels["app.kubernetes.io/managed-by"] = "chia-operator"
	return labels
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
