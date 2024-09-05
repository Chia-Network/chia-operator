/*
Copyright 2023 Chia Network Inc.
*/

package chiatimelord

import (
	"fmt"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	"strconv"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
)

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func getChiaVolumes(tl k8schianetv1.ChiaTimelord) []corev1.Volume {
	var v []corev1.Volume

	// secret ca volume
	v = append(v, corev1.Volume{
		Name: "secret-ca",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: tl.Spec.ChiaConfig.CASecretName,
			},
		},
	})

	// CHIA_ROOT volume
	if kube.ShouldMakeVolumeClaim(tl.Spec.Storage) {
		v = append(v, corev1.Volume{
			Name: "chiaroot",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: fmt.Sprintf(chiatimelordNamePattern, tl.Name),
				},
			},
		})
	} else {
		v = append(v, kube.GetExistingChiaRootVolume(tl.Spec.Storage))
	}

	// Add sidecar volumes if any exist
	if len(tl.Spec.Sidecars.Volumes) > 0 {
		v = append(v, tl.Spec.Sidecars.Volumes...)
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
			Name:      "chiaroot",
			MountPath: "/chia-data",
		},
	}
}

// getChiaEnv retrieves the environment variables from the Chia config struct
func getChiaEnv(tl k8schianetv1.ChiaTimelord) []corev1.EnvVar {
	var env []corev1.EnvVar

	// service env var
	env = append(env, corev1.EnvVar{
		Name:  "service",
		Value: "timelord-only timelord-launcher-only",
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
	if tl.Spec.ChiaConfig.Testnet != nil && *tl.Spec.ChiaConfig.Testnet {
		env = append(env, corev1.EnvVar{
			Name:  "testnet",
			Value: "true",
		})
	}

	// network env var
	if tl.Spec.ChiaConfig.Network != nil && *tl.Spec.ChiaConfig.Network != "" {
		env = append(env, corev1.EnvVar{
			Name:  "network",
			Value: *tl.Spec.ChiaConfig.Network,
		})
	}

	// network_port env var
	if tl.Spec.ChiaConfig.NetworkPort != nil && *tl.Spec.ChiaConfig.NetworkPort != 0 {
		env = append(env, corev1.EnvVar{
			Name:  "network_port",
			Value: strconv.Itoa(int(*tl.Spec.ChiaConfig.NetworkPort)),
		})
	}

	// introducer_address env var
	if tl.Spec.ChiaConfig.IntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "introducer_address",
			Value: *tl.Spec.ChiaConfig.IntroducerAddress,
		})
	}

	// dns_introducer_address env var
	if tl.Spec.ChiaConfig.DNSIntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "dns_introducer_address",
			Value: *tl.Spec.ChiaConfig.DNSIntroducerAddress,
		})
	}

	// TZ env var
	if tl.Spec.ChiaConfig.Timezone != nil {
		env = append(env, corev1.EnvVar{
			Name:  "TZ",
			Value: *tl.Spec.ChiaConfig.Timezone,
		})
	}

	// log_level env var
	if tl.Spec.ChiaConfig.LogLevel != nil {
		env = append(env, corev1.EnvVar{
			Name:  "log_level",
			Value: *tl.Spec.ChiaConfig.LogLevel,
		})
	}

	// self_hostname env var
	if tl.Spec.ChiaConfig.SelfHostname != nil {
		env = append(env, corev1.EnvVar{
			Name:  "self_hostname",
			Value: *tl.Spec.ChiaConfig.SelfHostname,
		})
	} else {
		env = append(env, corev1.EnvVar{
			Name:  "self_hostname",
			Value: "0.0.0.0",
		})
	}

	// node peer env var
	env = append(env, corev1.EnvVar{
		Name:  "full_node_peer",
		Value: tl.Spec.ChiaConfig.FullNodePeer,
	})

	return env
}
