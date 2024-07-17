package kube

import (
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// AssembleCommonServiceInputs contains configuration inputs to the AssembleCommonService function
type AssembleCommonServiceInputs struct {
	Name           string
	Namespace      string
	Labels         map[string]string
	Annotations    map[string]string
	OwnerReference []metav1.OwnerReference
	IPFamilyPolicy *corev1.IPFamilyPolicy
	IPFamilies     *[]corev1.IPFamily
	ServiceType    *corev1.ServiceType
	Ports          []corev1.ServicePort
	SelectorLabels map[string]string
}

// AssembleCommonService accepts some values and outputs a kubernetes Service definition in a standard way
func AssembleCommonService(input AssembleCommonServiceInputs) corev1.Service {
	srv := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            input.Name,
			Namespace:       input.Namespace,
			Labels:          input.Labels,
			Annotations:     input.Annotations,
			OwnerReferences: input.OwnerReference,
		},
		Spec: corev1.ServiceSpec{
			IPFamilyPolicy: input.IPFamilyPolicy,
			Ports:          input.Ports,
			Selector:       input.SelectorLabels,
		},
	}

	if input.ServiceType != nil {
		srv.Spec.Type = *input.ServiceType
	}

	if input.IPFamilies != nil {
		srv.Spec.IPFamilies = *input.IPFamilies
	}

	return srv
}

// AssembleChiaExporterContainerInputs contains configuration inputs to the AssembleChiaExporterContainer function
type AssembleChiaExporterContainerInputs struct {
	Image                string
	ConfigSecretName     *string
	SecurityContext      *corev1.SecurityContext
	PullPolicy           corev1.PullPolicy
	ResourceRequirements corev1.ResourceRequirements
}

// AssembleChiaExporterContainer assembles a chia-exporter container spec
func AssembleChiaExporterContainer(input AssembleChiaExporterContainerInputs) corev1.Container {
	container := corev1.Container{
		Name:            "chia-exporter",
		SecurityContext: input.SecurityContext,
		Image:           input.Image,
		ImagePullPolicy: input.PullPolicy,
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
					Port: intstr.FromInt32(consts.ChiaExporterPort),
				},
			},
		},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/healthz",
					Port: intstr.FromInt32(consts.ChiaExporterPort),
				},
			},
		},
		StartupProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/healthz",
					Port: intstr.FromInt32(consts.ChiaExporterPort),
				},
			},
			FailureThreshold: 30,
			PeriodSeconds:    10,
		},
		Resources: input.ResourceRequirements,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "chiaroot",
				MountPath: "/chia-data",
			},
		},
	}

	if input.ConfigSecretName != nil && *input.ConfigSecretName != "" {
		container.EnvFrom = append(container.EnvFrom, corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: *input.ConfigSecretName,
				},
			},
		})
	}

	return container
}
