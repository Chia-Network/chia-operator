/*
Copyright 2023 Chia Network Inc.
*/

package kube

import (
	"strconv"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GetCommonLabels gives some common labels for chia-operator related objects
func GetCommonLabels(kind string, meta metav1.ObjectMeta, additionalLabels ...map[string]string) map[string]string {
	labels := CombineMaps(additionalLabels...)
	labels["app.kubernetes.io/instance"] = meta.Name
	labels["app.kubernetes.io/name"] = meta.Name
	labels["app.kubernetes.io/managed-by"] = "chia-operator"
	labels["k8s.chia.net/kind"] = kind
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

// GetFullNodePort determines the correct full_node port to use
func GetFullNodePort(chia k8schianetv1.CommonSpecChia) int32 {
	if chia.NetworkPort != nil {
		return int32(*chia.NetworkPort)
	}
	if chia.Testnet != nil && *chia.Testnet {
		return consts.TestnetNodePort
	}
	return consts.MainnetNodePort
}

// GetExistingChiaRootVolume returns a corev1 API Volume specification for CHIA_ROOT.
// If both a PV and hostPath volume are specified for CHIA_ROOT, the PV will take precedence.
// If both configs are empty, this will fall back to emptyDir so sidecars can mount CHIA_ROOT.
// NOTE: This function does not handle the mode where the controller generates a CHIA_ROOT PVC, itself.
// Therefore, if ShouldMakeVolumeClaim is true, specifying the PVC's name should be handled in the controller.
func GetExistingChiaRootVolume(storage *k8schianetv1.StorageConfig) corev1.Volume {
	volumeName := "chiaroot"
	if storage != nil && storage.ChiaRoot != nil {
		if storage.ChiaRoot.PersistentVolumeClaim != nil && storage.ChiaRoot.PersistentVolumeClaim.ClaimName != "" {
			return corev1.Volume{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: storage.ChiaRoot.PersistentVolumeClaim.ClaimName,
					},
				},
			}
		} else if storage.ChiaRoot.HostPathVolume != nil && storage.ChiaRoot.HostPathVolume.Path != "" {
			return corev1.Volume{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: storage.ChiaRoot.HostPathVolume.Path,
					},
				},
			}
		}
	}

	return corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
}

// GetCommonChiaEnv retrieves the environment variables from the CommonSpecChia config struct
func GetCommonChiaEnv(commonSpecChia k8schianetv1.CommonSpecChia) []corev1.EnvVar {
	var env []corev1.EnvVar

	// CHIA_ROOT env var
	env = append(env, corev1.EnvVar{
		Name:  "CHIA_ROOT",
		Value: "/chia-data",
	})

	// ca env var
	env = append(env, corev1.EnvVar{
		Name:  "ca",
		Value: "/chia-ca",
	})

	// testnet env var
	if commonSpecChia.Testnet != nil && *commonSpecChia.Testnet {
		env = append(env, corev1.EnvVar{
			Name:  "testnet",
			Value: "true",
		})
	}

	// network env var
	if commonSpecChia.Network != nil && *commonSpecChia.Network != "" {
		env = append(env, corev1.EnvVar{
			Name:  "network",
			Value: *commonSpecChia.Network,
		})
	}

	// network_port env var
	env = append(env, corev1.EnvVar{
		Name:  "network_port",
		Value: strconv.Itoa(int(GetFullNodePort(commonSpecChia))),
	})

	// introducer_address env var
	if commonSpecChia.IntroducerAddress != nil && *commonSpecChia.IntroducerAddress != "" {
		env = append(env, corev1.EnvVar{
			Name:  "introducer_address",
			Value: *commonSpecChia.IntroducerAddress,
		})
	}

	// dns_introducer_address env var
	if commonSpecChia.DNSIntroducerAddress != nil && *commonSpecChia.DNSIntroducerAddress != "" {
		env = append(env, corev1.EnvVar{
			Name:  "dns_introducer_address",
			Value: *commonSpecChia.DNSIntroducerAddress,
		})
	}

	// TZ env var
	if commonSpecChia.Timezone != nil && *commonSpecChia.Timezone != "" {
		env = append(env, corev1.EnvVar{
			Name:  "TZ",
			Value: *commonSpecChia.Timezone,
		})
	}

	// log_level env var
	if commonSpecChia.LogLevel != nil && *commonSpecChia.LogLevel != "" {
		env = append(env, corev1.EnvVar{
			Name:  "log_level",
			Value: *commonSpecChia.LogLevel,
		})
	}

	// source_ref env var
	if commonSpecChia.SourceRef != nil && *commonSpecChia.SourceRef != "" {
		env = append(env, corev1.EnvVar{
			Name:  "source_ref",
			Value: *commonSpecChia.SourceRef,
		})
	}

	// self_hostname env var
	if commonSpecChia.SelfHostname != nil && *commonSpecChia.SelfHostname != "" {
		env = append(env, corev1.EnvVar{
			Name:  "self_hostname",
			Value: *commonSpecChia.SelfHostname,
		})
	} else {
		env = append(env, corev1.EnvVar{
			Name:  "self_hostname",
			Value: "0.0.0.0",
		})
	}

	return env
}
