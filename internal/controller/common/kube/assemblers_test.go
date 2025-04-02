package kube

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

	"github.com/chia-network/chia-operator/internal/controller/common/consts"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

var testControllerOwner = true

func TestAssembleCommonService_Minimal(t *testing.T) {
	expected := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "default",
			Labels: map[string]string{
				"app.kubernetes.io/name": "test",
			},
			Annotations: map[string]string{
				"test-annotation": "testing",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaTest",
					Name:       "test",
					UID:        "testuid",
					Controller: &testControllerOwner,
				},
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:       8444,
					TargetPort: intstr.FromString("peers"),
					Protocol:   "TCP",
					Name:       "peers",
				},
			},
			Selector: map[string]string{
				"app.kubernetes.io/name": "test",
			},
		},
	}
	actual := AssembleCommonService(AssembleCommonServiceInputs{
		Name:           expected.Name,
		Namespace:      expected.Namespace,
		Labels:         expected.Labels,
		Annotations:    expected.Annotations,
		OwnerReference: expected.OwnerReferences,
		Ports:          expected.Spec.Ports,
		SelectorLabels: expected.Spec.Selector,
	})
	require.Equal(t, expected, actual)
}

func TestAssembleCommonService_Full(t *testing.T) {
	familyPolicy := corev1.IPFamilyPolicyPreferDualStack
	expected := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "default",
			Labels: map[string]string{
				"app.kubernetes.io/name": "test",
			},
			Annotations: map[string]string{
				"test-annotation": "testing",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "k8s.chia.net/v1",
					Kind:       "ChiaTest",
					Name:       "test",
					UID:        "testuid",
					Controller: &testControllerOwner,
				},
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType("NodePort"),
			Ports: []corev1.ServicePort{
				{
					Port:       8444,
					TargetPort: intstr.FromString("peers"),
					Protocol:   "TCP",
					Name:       "peers",
				},
			},
			IPFamilyPolicy: &familyPolicy,
			IPFamilies: []corev1.IPFamily{
				"IPv4",
				"IPv6",
			},
			Selector: map[string]string{
				"app.kubernetes.io/name": "test",
			},
			SessionAffinity: corev1.ServiceAffinityClientIP,
			SessionAffinityConfig: &corev1.SessionAffinityConfig{
				ClientIP: &corev1.ClientIPConfig{
					TimeoutSeconds: ptr.To(int32(300)),
				},
			},
		},
	}
	actual := AssembleCommonService(AssembleCommonServiceInputs{
		Name:                  expected.Name,
		Namespace:             expected.Namespace,
		Labels:                expected.Labels,
		Annotations:           expected.Annotations,
		OwnerReference:        expected.OwnerReferences,
		Ports:                 expected.Spec.Ports,
		SelectorLabels:        expected.Spec.Selector,
		ServiceType:           &expected.Spec.Type,
		IPFamilyPolicy:        expected.Spec.IPFamilyPolicy,
		IPFamilies:            &expected.Spec.IPFamilies,
		SessionAffinity:       &expected.Spec.SessionAffinity,
		SessionAffinityConfig: expected.Spec.SessionAffinityConfig,
	})
	require.Equal(t, expected, actual)
}

func TestAssembleChiaContainer_Minimal(t *testing.T) {
	expected := corev1.Container{
		Name:            "chia",
		Image:           "test:latest",
		ImagePullPolicy: "Always",
		Env: []corev1.EnvVar{
			{
				Name:  "testkey",
				Value: "testvalue",
			},
		},
		Ports: []corev1.ContainerPort{
			{
				Name:          "http",
				Protocol:      "TCP",
				ContainerPort: 8080,
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "chiaroot",
				MountPath: "/chia-data",
			},
		},
	}
	actual := AssembleChiaContainer(AssembleChiaContainerInputs{
		Image:           &expected.Image,
		ImagePullPolicy: expected.ImagePullPolicy,
		Env:             expected.Env,
		Ports:           expected.Ports,
		VolumeMounts:    expected.VolumeMounts,
	})
	require.Equal(t, expected, actual)
}

func TestAssembleChiaContainer_Full(t *testing.T) {
	expected := corev1.Container{
		Name:            "chia",
		Image:           "test:latest",
		ImagePullPolicy: "Always",
		Env: []corev1.EnvVar{
			{
				Name:  "testkey",
				Value: "testvalue",
			},
		},
		Ports: []corev1.ContainerPort{
			{
				Name:          "http",
				Protocol:      "TCP",
				ContainerPort: 8080,
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "chiaroot",
				MountPath: "/chia-data",
			},
		},
		SecurityContext: &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{
					"NET_BIND_SERVICE",
				},
			},
		},
		LivenessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/liveness",
					Port: intstr.FromInt32(consts.ChiaHealthcheckPort),
				},
			},
		},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/readiness",
					Port: intstr.FromInt32(consts.ChiaHealthcheckPort),
				},
			},
		},
		StartupProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/startup",
					Port: intstr.FromInt32(consts.ChiaHealthcheckPort),
				},
			},
		},
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				"CPU":    resource.MustParse("200m"),
				"Memory": resource.MustParse("512Mi"),
			},
			Requests: corev1.ResourceList{
				"CPU":    resource.MustParse("100m"),
				"Memory": resource.MustParse("256Mi"),
			},
		},
	}
	actual := AssembleChiaContainer(AssembleChiaContainerInputs{
		Image:                &expected.Image,
		ImagePullPolicy:      expected.ImagePullPolicy,
		Env:                  expected.Env,
		Ports:                expected.Ports,
		VolumeMounts:         expected.VolumeMounts,
		SecurityContext:      expected.SecurityContext,
		ResourceRequirements: &expected.Resources,
		LivenessProbe:        expected.LivenessProbe,
		ReadinessProbe:       expected.ReadinessProbe,
		StartupProbe:         expected.StartupProbe,
	})
	require.Equal(t, expected, actual)
}

func TestAssembleChiaExporterContainer_Minimal(t *testing.T) {
	expected := corev1.Container{
		Name:            "chia-exporter",
		Image:           "test:latest",
		ImagePullPolicy: "Always",
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
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "chiaroot",
				MountPath: "/chia-data",
			},
		},
	}
	actual := AssembleChiaExporterContainer(AssembleChiaExporterContainerInputs{
		Image:           &expected.Image,
		ImagePullPolicy: expected.ImagePullPolicy,
	})
	require.Equal(t, expected, actual)
}

func TestAssembleChiaExporterContainer_Full(t *testing.T) {
	secretName := "configsecret"
	expected := corev1.Container{
		Name:            "chia-exporter",
		Image:           "test:latest",
		ImagePullPolicy: "Always",
		EnvFrom: []corev1.EnvFromSource{
			{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secretName,
					},
				},
			},
		},
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
		SecurityContext: &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{
					"NET_BIND_SERVICE",
				},
			},
		},
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				"CPU":    resource.MustParse("200m"),
				"Memory": resource.MustParse("512Mi"),
			},
			Requests: corev1.ResourceList{
				"CPU":    resource.MustParse("100m"),
				"Memory": resource.MustParse("256Mi"),
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "chiaroot",
				MountPath: "/chia-data",
			},
		},
	}
	actual := AssembleChiaExporterContainer(AssembleChiaExporterContainerInputs{
		Image:                &expected.Image,
		ImagePullPolicy:      expected.ImagePullPolicy,
		ResourceRequirements: expected.Resources,
		ConfigSecretName:     &secretName,
		SecurityContext:      expected.SecurityContext,
	})
	require.Equal(t, expected, actual)
}

func TestAssembleChiaHealthcheckContainer_Minimal(t *testing.T) {
	expected := corev1.Container{
		Name:            "chia-healthcheck",
		Image:           "test:latest",
		ImagePullPolicy: "Always",
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
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "chiaroot",
				MountPath: "/chia-data",
			},
		},
	}
	actual := AssembleChiaHealthcheckContainer(AssembleChiaHealthcheckContainerInputs{
		Image:           &expected.Image,
		ImagePullPolicy: expected.ImagePullPolicy,
	})
	require.Equal(t, expected, actual)
}

func TestAssembleChiaHealthcheckContainer_Full(t *testing.T) {
	dnsHostname := "test.testing.test"
	expected := corev1.Container{
		Name:            "chia-healthcheck",
		Image:           "test:latest",
		ImagePullPolicy: "Always",
		Env: []corev1.EnvVar{
			{
				Name:  "CHIA_ROOT",
				Value: "/chia-data",
			},
			{
				Name:  "CHIA_HEALTHCHECK_HOSTNAME",
				Value: "127.0.0.1",
			},
			{
				Name:  "CHIA_HEALTHCHECK_DNS_HOSTNAME",
				Value: dnsHostname,
			},
		},
		Ports: []corev1.ContainerPort{
			{
				Name:          "health",
				ContainerPort: consts.ChiaHealthcheckPort,
				Protocol:      "TCP",
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "chiaroot",
				MountPath: "/chia-data",
			},
		},
		SecurityContext: &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{
					"NET_BIND_SERVICE",
				},
			},
		},
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				"CPU":    resource.MustParse("200m"),
				"Memory": resource.MustParse("512Mi"),
			},
			Requests: corev1.ResourceList{
				"CPU":    resource.MustParse("100m"),
				"Memory": resource.MustParse("256Mi"),
			},
		},
	}
	actual := AssembleChiaHealthcheckContainer(AssembleChiaHealthcheckContainerInputs{
		Image:                &expected.Image,
		ImagePullPolicy:      expected.ImagePullPolicy,
		ResourceRequirements: expected.Resources,
		SecurityContext:      expected.SecurityContext,
		DNSHostname:          &dnsHostname,
	})
	require.Equal(t, expected, actual)
}

func TestAssembleChiaHealthcheckProbe_Minimal(t *testing.T) {
	expected := corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/seeder",
				Port: intstr.FromInt32(consts.ChiaHealthcheckPort),
			},
		},
	}
	actual := AssembleChiaHealthcheckProbe(AssembleChiaHealthcheckProbeInputs{
		Path: "/seeder",
	})
	require.Equal(t, expected, *actual)
}

func TestAssembleChiaHealthcheckProbe_Full(t *testing.T) {
	failThresh := int32(11)
	periodSec := int32(17)
	expected := corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/seeder/readiness",
				Port: intstr.FromInt32(consts.ChiaHealthcheckPort),
			},
		},
		FailureThreshold: failThresh,
		PeriodSeconds:    periodSec,
	}
	actual := AssembleChiaHealthcheckProbe(AssembleChiaHealthcheckProbeInputs{
		Path:             "/seeder/readiness",
		FailureThreshold: &failThresh,
		PeriodSeconds:    &periodSec,
	})
	require.Equal(t, expected, *actual)
}
