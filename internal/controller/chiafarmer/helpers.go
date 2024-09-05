/*
Copyright 2023 Chia Network Inc.
*/

package chiafarmer

import (
	"fmt"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	"strconv"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
)

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func getChiaVolumes(farmer k8schianetv1.ChiaFarmer) []corev1.Volume {
	var v []corev1.Volume

	// secret ca volume
	v = append(v, corev1.Volume{
		Name: "secret-ca",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: farmer.Spec.ChiaConfig.CASecretName,
			},
		},
	})

	// mnemonic key volume
	v = append(v, corev1.Volume{
		Name: "key",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: farmer.Spec.ChiaConfig.SecretKey.Name,
			},
		},
	})

	// CHIA_ROOT volume
	if kube.ShouldMakeVolumeClaim(farmer.Spec.Storage) {
		v = append(v, corev1.Volume{
			Name: "chiaroot",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: fmt.Sprintf(chiafarmerNamePattern, farmer.Name),
				},
			},
		})
	} else {
		v = append(v, kube.GetExistingChiaRootVolume(farmer.Spec.Storage))
	}

	// Add sidecar volumes if any exist
	if len(farmer.Spec.Sidecars.Volumes) > 0 {
		v = append(v, farmer.Spec.Sidecars.Volumes...)
	}

	return v
}

// getChiaVolumeMounts retrieves the requisite volume mounts from the Chia config struct
func getChiaVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "secret-ca",
			MountPath: "/chia-ca",
		},
		{
			Name:      "key",
			MountPath: "/key",
		},
		{
			Name:      "chiaroot",
			MountPath: "/chia-data",
		},
	}
}

// getChiaEnv retrieves the environment variables from the Chia config struct
func getChiaEnv(farmer k8schianetv1.ChiaFarmer) []corev1.EnvVar {
	var env []corev1.EnvVar

	// service env var
	env = append(env, corev1.EnvVar{
		Name:  "service",
		Value: "farmer-only",
	})

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
	if farmer.Spec.ChiaConfig.Testnet != nil && *farmer.Spec.ChiaConfig.Testnet {
		env = append(env, corev1.EnvVar{
			Name:  "testnet",
			Value: "true",
		})
	}

	// network env var
	if farmer.Spec.ChiaConfig.Network != nil && *farmer.Spec.ChiaConfig.Network != "" {
		env = append(env, corev1.EnvVar{
			Name:  "network",
			Value: *farmer.Spec.ChiaConfig.Network,
		})
	}

	// network_port env var
	if farmer.Spec.ChiaConfig.NetworkPort != nil && *farmer.Spec.ChiaConfig.NetworkPort != 0 {
		env = append(env, corev1.EnvVar{
			Name:  "network_port",
			Value: strconv.Itoa(int(*farmer.Spec.ChiaConfig.NetworkPort)),
		})
	}

	// introducer_address env var
	if farmer.Spec.ChiaConfig.IntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "introducer_address",
			Value: *farmer.Spec.ChiaConfig.IntroducerAddress,
		})
	}

	// dns_introducer_address env var
	if farmer.Spec.ChiaConfig.DNSIntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "dns_introducer_address",
			Value: *farmer.Spec.ChiaConfig.DNSIntroducerAddress,
		})
	}

	// TZ env var
	if farmer.Spec.ChiaConfig.Timezone != nil {
		env = append(env, corev1.EnvVar{
			Name:  "TZ",
			Value: *farmer.Spec.ChiaConfig.Timezone,
		})
	}

	// log_level env var
	if farmer.Spec.ChiaConfig.LogLevel != nil {
		env = append(env, corev1.EnvVar{
			Name:  "log_level",
			Value: *farmer.Spec.ChiaConfig.LogLevel,
		})
	}

	// self_hostname env var
	if farmer.Spec.ChiaConfig.SelfHostname != nil {
		env = append(env, corev1.EnvVar{
			Name:  "self_hostname",
			Value: *farmer.Spec.ChiaConfig.SelfHostname,
		})
	} else {
		env = append(env, corev1.EnvVar{
			Name:  "self_hostname",
			Value: "0.0.0.0",
		})
	}

	// keys env var
	env = append(env, corev1.EnvVar{
		Name:  "keys",
		Value: fmt.Sprintf("/key/%s", farmer.Spec.ChiaConfig.SecretKey.Key),
	})

	// node peer env var
	env = append(env, corev1.EnvVar{
		Name:  "full_node_peer",
		Value: farmer.Spec.ChiaConfig.FullNodePeer,
	})

	return env
}
