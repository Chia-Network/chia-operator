/*
Copyright 2024 Chia Network Inc.
*/

package chiacrawler

import (
	"fmt"

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
func getChiaEnv(crawler k8schianetv1.ChiaCrawler, networkData *map[string]string) ([]corev1.EnvVar, error) {
	var env []corev1.EnvVar

	// service env var
	env = append(env, corev1.EnvVar{
		Name:  "service",
		Value: "crawler",
	})

	// keys env var -- no keys required for a crawler
	env = append(env, corev1.EnvVar{
		Name:  "keys",
		Value: "none",
	})

	// Add common env
	commonEnv, err := kube.GetCommonChiaEnv(crawler.Spec.ChiaConfig.CommonSpecChia, networkData)
	if err != nil {
		return env, err
	}
	env = append(env, commonEnv...)

	return env, nil
}
