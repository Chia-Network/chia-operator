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

	// CHIA_ROOT volume
	if kube.ShouldMakeVolumeClaim(introducer.Spec.Storage) {
		v = append(v, corev1.Volume{
			Name: "chiaroot",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: fmt.Sprintf(chiaintroducerNamePattern, introducer.Name),
				},
			},
		})
	} else {
		v = append(v, kube.GetExistingChiaRootVolume(introducer.Spec.Storage))
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
