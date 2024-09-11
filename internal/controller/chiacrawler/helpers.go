/*
Copyright 2024 Chia Network Inc.
*/

package chiacrawler

import (
	"fmt"
	"strconv"

	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	corev1 "k8s.io/api/core/v1"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
)

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func getChiaVolumes(crawler k8schianetv1.ChiaCrawler) []corev1.Volume {
	var v []corev1.Volume

	// secret ca volume
	if crawler.Spec.ChiaConfig.CASecretName != nil {
		v = append(v, corev1.Volume{
			Name: "secret-ca",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: *crawler.Spec.ChiaConfig.CASecretName,
				},
			},
		})
	}

	// CHIA_ROOT volume
	if kube.ShouldMakeVolumeClaim(crawler.Spec.Storage) {
		v = append(v, corev1.Volume{
			Name: "chiaroot",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: fmt.Sprintf(chiacrawlerNamePattern, crawler.Name),
				},
			},
		})
	} else {
		v = append(v, kube.GetExistingChiaRootVolume(crawler.Spec.Storage))
	}

	// Add sidecar volumes if any exist
	if len(crawler.Spec.Sidecars.Volumes) > 0 {
		v = append(v, crawler.Spec.Sidecars.Volumes...)
	}

	return v
}

// getChiaVolumeMounts retrieves the requisite volume mounts from the Chia config struct
func getChiaVolumeMounts(crawler k8schianetv1.ChiaCrawler) []corev1.VolumeMount {
	var v []corev1.VolumeMount

	// secret ca volume
	if crawler.Spec.ChiaConfig.CASecretName != nil {
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
func getChiaEnv(crawler k8schianetv1.ChiaCrawler) []corev1.EnvVar {
	var env []corev1.EnvVar

	// service env var
	env = append(env, corev1.EnvVar{
		Name:  "service",
		Value: "crawler",
	})

	// CHIA_ROOT env var
	env = append(env, corev1.EnvVar{
		Name:  "CHIA_ROOT",
		Value: "/chia-data",
	})

	// keys env var -- no keys required for a crawler
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
	if crawler.Spec.ChiaConfig.Testnet != nil && *crawler.Spec.ChiaConfig.Testnet {
		env = append(env, corev1.EnvVar{
			Name:  "testnet",
			Value: "true",
		})
	}

	// network env var
	if crawler.Spec.ChiaConfig.Network != nil && *crawler.Spec.ChiaConfig.Network != "" {
		env = append(env, corev1.EnvVar{
			Name:  "network",
			Value: *crawler.Spec.ChiaConfig.Network,
		})
	}

	// network_port env var
	if crawler.Spec.ChiaConfig.NetworkPort != nil && *crawler.Spec.ChiaConfig.NetworkPort != 0 {
		env = append(env, corev1.EnvVar{
			Name:  "network_port",
			Value: strconv.Itoa(int(*crawler.Spec.ChiaConfig.NetworkPort)),
		})
	}

	// introducer_address env var
	if crawler.Spec.ChiaConfig.IntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "introducer_address",
			Value: *crawler.Spec.ChiaConfig.IntroducerAddress,
		})
	}

	// dns_introducer_address env var
	if crawler.Spec.ChiaConfig.DNSIntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "dns_introducer_address",
			Value: *crawler.Spec.ChiaConfig.DNSIntroducerAddress,
		})
	}

	// TZ env var
	if crawler.Spec.ChiaConfig.Timezone != nil {
		env = append(env, corev1.EnvVar{
			Name:  "TZ",
			Value: *crawler.Spec.ChiaConfig.Timezone,
		})
	}

	// log_level env var
	if crawler.Spec.ChiaConfig.LogLevel != nil {
		env = append(env, corev1.EnvVar{
			Name:  "log_level",
			Value: *crawler.Spec.ChiaConfig.LogLevel,
		})
	}

	// source_ref env var
	if crawler.Spec.ChiaConfig.SourceRef != nil && *crawler.Spec.ChiaConfig.SourceRef != "" {
		env = append(env, corev1.EnvVar{
			Name:  "source_ref",
			Value: *crawler.Spec.ChiaConfig.SourceRef,
		})
	}

	// self_hostname env var
	if crawler.Spec.ChiaConfig.SelfHostname != nil {
		env = append(env, corev1.EnvVar{
			Name:  "self_hostname",
			Value: *crawler.Spec.ChiaConfig.SelfHostname,
		})
	} else {
		env = append(env, corev1.EnvVar{
			Name:  "self_hostname",
			Value: "0.0.0.0",
		})
	}

	return env
}
