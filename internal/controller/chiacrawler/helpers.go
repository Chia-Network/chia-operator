/*
Copyright 2024 Chia Network Inc.
*/

package chiacrawler

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
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

	// CHIA_ROOT volume -- PVC is respected first if both it and hostpath are specified, falls back to hostPath if specified
	// If both are empty, fall back to emptyDir so chia-exporter can mount CHIA_ROOT
	var chiaRootAdded = false
	if crawler.Spec.Storage != nil && crawler.Spec.Storage.ChiaRoot != nil {
		if crawler.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil {
			var pvcName string
			if crawler.Spec.Storage.ChiaRoot.PersistentVolumeClaim.GenerateVolumeClaims {
				pvcName = fmt.Sprintf(chiacrawlerNamePattern, crawler.Name)
			} else {
				pvcName = crawler.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ClaimName
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
		} else if crawler.Spec.Storage.ChiaRoot.HostPathVolume != nil {
			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: crawler.Spec.Storage.ChiaRoot.HostPathVolume.Path,
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

	return env
}

// getOwnerReference gives the common owner reference spec for ChiaCrawler related objects
func getOwnerReference(crawler k8schianetv1.ChiaCrawler) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: crawler.APIVersion,
			Kind:       crawler.Kind,
			Name:       crawler.Name,
			UID:        crawler.UID,
			Controller: &consts.ControllerOwner,
		},
	}
}

// getFullNodePort determines the correct full_node port to use
func getFullNodePort(crawler k8schianetv1.ChiaCrawler) int32 {
	if crawler.Spec.ChiaConfig.Testnet != nil && *crawler.Spec.ChiaConfig.Testnet {
		return consts.TestnetNodePort
	}
	return consts.MainnetNodePort
}
