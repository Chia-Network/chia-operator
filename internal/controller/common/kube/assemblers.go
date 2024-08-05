package kube

import (
	"fmt"
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

// AssembleChiaContainerInputs contains configuration inputs to the AssembleChiaContainer function
type AssembleChiaContainerInputs struct {
	Image                *string
	ImagePullPolicy      corev1.PullPolicy
	Env                  []corev1.EnvVar
	Ports                []corev1.ContainerPort
	VolumeMounts         []corev1.VolumeMount
	SecurityContext      *corev1.SecurityContext
	LivenessProbe        *corev1.Probe
	ReadinessProbe       *corev1.Probe
	StartupProbe         *corev1.Probe
	ResourceRequirements *corev1.ResourceRequirements
}

// AssembleChiaContainer assembles a chia container spec
func AssembleChiaContainer(input AssembleChiaContainerInputs) corev1.Container {
	container := corev1.Container{
		Name:            "chia",
		ImagePullPolicy: input.ImagePullPolicy,
		Env:             input.Env,
		Ports:           input.Ports,
		VolumeMounts:    input.VolumeMounts,
		SecurityContext: input.SecurityContext,
		LivenessProbe:   input.LivenessProbe,
		ReadinessProbe:  input.ReadinessProbe,
		StartupProbe:    input.StartupProbe,
	}

	if input.Image != nil && *input.Image != "" {
		container.Image = *input.Image
	} else {
		container.Image = fmt.Sprintf("%s:%s", consts.DefaultChiaImageName, consts.DefaultChiaImageTag)
	}

	if input.ResourceRequirements != nil {
		container.Resources = *input.ResourceRequirements
	}

	return container
}

// AssembleChiaExporterContainerInputs contains configuration inputs to the AssembleChiaExporterContainer function
type AssembleChiaExporterContainerInputs struct {
	Image                *string
	ImagePullPolicy      corev1.PullPolicy
	ResourceRequirements corev1.ResourceRequirements
	ConfigSecretName     *string
	SecurityContext      *corev1.SecurityContext
}

// AssembleChiaExporterContainer assembles a chia-exporter container spec
func AssembleChiaExporterContainer(input AssembleChiaExporterContainerInputs) corev1.Container {
	container := corev1.Container{
		Name:            "chia-exporter",
		SecurityContext: input.SecurityContext,
		ImagePullPolicy: input.ImagePullPolicy,
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

	if input.Image != nil && *input.Image != "" {
		container.Image = *input.Image
	} else {
		container.Image = fmt.Sprintf("%s:%s", consts.DefaultChiaExporterImageName, consts.DefaultChiaExporterImageTag)
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

// AssembleChiaHealthcheckContainerInputs contains configuration inputs to the AssembleChiaHealthcheckContainer function
type AssembleChiaHealthcheckContainerInputs struct {
	Image                *string
	ImagePullPolicy      corev1.PullPolicy
	ResourceRequirements corev1.ResourceRequirements
	DNSHostname          *string
	SecurityContext      *corev1.SecurityContext
}

// AssembleChiaHealthcheckContainer assembles a chia-healthcheck container spec
func AssembleChiaHealthcheckContainer(input AssembleChiaHealthcheckContainerInputs) corev1.Container {
	container := corev1.Container{
		Name:            "chia-healthcheck",
		SecurityContext: input.SecurityContext,
		ImagePullPolicy: input.ImagePullPolicy,
		Env: []corev1.EnvVar{
			{
				Name:  "CHIA_ROOT",
				Value: "/chia-data",
			},
			{
				Name:  "CHIA_HEALTHCHECK_HOSTNAME",
				Value: "127.0.0.1",
			},
		},
		Ports: []corev1.ContainerPort{
			{
				Name:          "health",
				ContainerPort: consts.ChiaHealthcheckPort,
				Protocol:      "TCP",
			},
		},
		Resources: input.ResourceRequirements,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "chiaroot",
				MountPath: "/chia-data",
			},
		},
	}

	if input.Image != nil && *input.Image != "" {
		container.Image = *input.Image
	} else {
		container.Image = fmt.Sprintf("%s:%s", consts.DefaultChiaHealthcheckImageName, consts.DefaultChiaHealthcheckImageTag)
	}

	if input.DNSHostname != nil && *input.DNSHostname != "" {
		container.Env = append(container.Env, corev1.EnvVar{
			Name:  "CHIA_HEALTHCHECK_DNS_HOSTNAME",
			Value: *input.DNSHostname,
		})
	}

	return container
}

// AssembleChiaHealthcheckProbeInputs contains configuration inputs to the AssembleChiaHealthcheckProbe function
type AssembleChiaHealthcheckProbeInputs struct {
	Kind             consts.ChiaKind
	FailureThreshold *int32
	PeriodSeconds    *int32
}

func AssembleChiaHealthcheckProbe(input AssembleChiaHealthcheckProbeInputs) *corev1.Probe {
	probe := corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/",
				Port: intstr.FromInt32(consts.ChiaHealthcheckPort),
			},
		},
	}

	if input.FailureThreshold != nil {
		probe.FailureThreshold = *input.FailureThreshold
	}

	if input.PeriodSeconds != nil {
		probe.PeriodSeconds = *input.PeriodSeconds
	}

	switch input.Kind {
	case consts.ChiaNodeKind:
		probe.ProbeHandler.HTTPGet.Path = "/full_node"
	case consts.ChiaSeederKind:
		probe.ProbeHandler.HTTPGet.Path = "/seeder"
	default:
		return nil
	}
	return &probe
}
