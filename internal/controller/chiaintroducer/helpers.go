/*
Copyright 2024 Chia Network Inc.
*/

package chiaintroducer

import (
	"fmt"

	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	corev1 "k8s.io/api/core/v1"

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
	if kube.ShouldMakeChiaRootVolumeClaim(introducer.Spec.Storage) {
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
func getChiaEnv(introducer k8schianetv1.ChiaIntroducer, networkData *map[string]string) ([]corev1.EnvVar, error) {
	var env []corev1.EnvVar

	// service env var
	env = append(env, corev1.EnvVar{
		Name:  "service",
		Value: "introducer",
	})

	// keys env var -- no keys required for a introducer
	env = append(env, corev1.EnvVar{
		Name:  "keys",
		Value: "none",
	})

	// Add common env
	commonEnv, err := kube.GetCommonChiaEnv(introducer.Spec.ChiaConfig.CommonSpecChia, networkData)
	if err != nil {
		return env, err
	}
	env = append(env, commonEnv...)

	return env, nil
}
