/*
Copyright 2024 Chia Network Inc.
*/

package chiaintroducer

import (
	"fmt"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	corev1 "k8s.io/api/core/v1"
	"strconv"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
)

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func getChiaVolumes(introducer k8schianetv1.ChiaIntroducer) []corev1.Volume {
	var v []corev1.Volume

	// secret ca volume
	if introducer.Spec.ChiaConfig.CASecretName != nil {
		v = append(v, corev1.Volume{
			Name: "secret-ca",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: *introducer.Spec.ChiaConfig.CASecretName,
				},
			},
		})
	}

	// CHIA_ROOT volume -- PVC is respected first if both it and hostpath are specified, falls back to hostPath if specified
	// If both are empty, fall back to emptyDir so chia-exporter can mount CHIA_ROOT
	var chiaRootAdded = false
	if introducer.Spec.Storage != nil && introducer.Spec.Storage.ChiaRoot != nil {
		if introducer.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil {
			var pvcName string
			if introducer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.GenerateVolumeClaims {
				pvcName = fmt.Sprintf(chiaintroducerNamePattern, introducer.Name)
			} else {
				pvcName = introducer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ClaimName
			}

			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: pvcName,
					},
				},
			})
			chiaRootAdded = true
		} else if introducer.Spec.Storage.ChiaRoot.HostPathVolume != nil {
			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: introducer.Spec.Storage.ChiaRoot.HostPathVolume.Path,
					},
				},
			})
			chiaRootAdded = true
		}
	}
	if !chiaRootAdded {
		v = append(v, corev1.Volume{
			Name: "chiaroot",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
	}

	// Add sidecar volumes if any exist
	if len(introducer.Spec.Sidecars.Volumes) > 0 {
		v = append(v, introducer.Spec.Sidecars.Volumes...)
	}

	return v
}

// getChiaVolumeMounts retrieves the requisite volume mounts from the Chia config struct
func getChiaVolumeMounts(introducer k8schianetv1.ChiaIntroducer) []corev1.VolumeMount {
	var v []corev1.VolumeMount

	// secret ca volume
	if introducer.Spec.ChiaConfig.CASecretName != nil {
		v = append(v, corev1.VolumeMount{
			Name:      "secret-ca",
			MountPath: "/chia-ca",
		})
	}

	// CHIA_ROOT volume
	v = append(v, corev1.VolumeMount{
		Name:      "chiaroot",
		MountPath: "/chia-data",
	})

	return v
}

// getChiaEnv retrieves the environment variables from the Chia config struct
func getChiaEnv(introducer k8schianetv1.ChiaIntroducer) []corev1.EnvVar {
	var env []corev1.EnvVar

	// service env var
	env = append(env, corev1.EnvVar{
		Name:  "service",
		Value: "introducer",
	})

	// CHIA_ROOT env var
	env = append(env, corev1.EnvVar{
		Name:  "CHIA_ROOT",
		Value: "/chia-data",
	})

	// keys env var -- no keys required for a introducer
	env = append(env, corev1.EnvVar{
		Name:  "keys",
		Value: "none",
	})

	// ca env var
	env = append(env, corev1.EnvVar{
		Name:  "ca",
		Value: "/chia-ca",
	})

	// testnet env var
	if introducer.Spec.ChiaConfig.Testnet != nil && *introducer.Spec.ChiaConfig.Testnet {
		env = append(env, corev1.EnvVar{
			Name:  "testnet",
			Value: "true",
		})
	}

	// network env var
	if introducer.Spec.ChiaConfig.Network != nil && *introducer.Spec.ChiaConfig.Network != "" {
		env = append(env, corev1.EnvVar{
			Name:  "network",
			Value: *introducer.Spec.ChiaConfig.Network,
		})
	}

	// network_port env var
	// network_port env var is required for introducers because it sets the introducer's full_node port in the config
	// The default full_node port in the initial config is 8445, which will often need overwriting
	env = append(env, corev1.EnvVar{
		Name:  "network_port",
		Value: strconv.Itoa(int(kube.GetFullNodePort(introducer.Spec.ChiaConfig.CommonSpecChia))),
	})

	// introducer_address env var
	if introducer.Spec.ChiaConfig.IntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "introducer_address",
			Value: *introducer.Spec.ChiaConfig.IntroducerAddress,
		})
	}

	// dns_introducer_address env var
	if introducer.Spec.ChiaConfig.DNSIntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "dns_introducer_address",
			Value: *introducer.Spec.ChiaConfig.DNSIntroducerAddress,
		})
	}

	// TZ env var
	if introducer.Spec.ChiaConfig.Timezone != nil {
		env = append(env, corev1.EnvVar{
			Name:  "TZ",
			Value: *introducer.Spec.ChiaConfig.Timezone,
		})
	}

	// log_level env var
	if introducer.Spec.ChiaConfig.LogLevel != nil {
		env = append(env, corev1.EnvVar{
			Name:  "log_level",
			Value: *introducer.Spec.ChiaConfig.LogLevel,
		})
	}

	// self_hostname env var
	if introducer.Spec.ChiaConfig.SelfHostname != nil {
		env = append(env, corev1.EnvVar{
			Name:  "self_hostname",
			Value: *introducer.Spec.ChiaConfig.SelfHostname,
		})
	} else {
		env = append(env, corev1.EnvVar{
			Name:  "self_hostname",
			Value: "0.0.0.0",
		})
	}

	return env
}

// getFullNodePort determines the correct full_node port to use
func getFullNodePort(introducer k8schianetv1.ChiaIntroducer) int32 {
	if introducer.Spec.ChiaConfig.Testnet != nil && *introducer.Spec.ChiaConfig.Testnet {
		return consts.TestnetNodePort
	}
	return consts.MainnetNodePort
}
