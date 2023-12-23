/*
Copyright 2023 Chia Network Inc.
*/

package chiadnsintroducer

import (
	"context"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
)

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func (r *ChiaDNSIntroducerReconciler) getChiaVolumes(ctx context.Context, dnsintro k8schianetv1.ChiaDNSIntroducer) []corev1.Volume {
	var v []corev1.Volume

	// secret ca volume
	v = append(v, corev1.Volume{
		Name: "secret-ca",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: dnsintro.Spec.ChiaConfig.CASecretName,
			},
		},
	})

	// CHIA_ROOT volume -- PVC is respected first if both it and hostpath are specified, falls back to hostPath if specified
	// If both are empty, fall back to emptyDir so chia-exporter can mount CHIA_ROOT
	var chiaRootAdded bool = false
	if dnsintro.Spec.Storage != nil && dnsintro.Spec.Storage.ChiaRoot != nil {
		if dnsintro.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil {
			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: dnsintro.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ClaimName,
					},
				},
			})
			chiaRootAdded = true
		} else if dnsintro.Spec.Storage.ChiaRoot.HostPathVolume != nil {
			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: dnsintro.Spec.Storage.ChiaRoot.HostPathVolume.Path,
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

	return v
}

// getChiaVolumeMounts retrieves the requisite volume mounts from the Chia config struct
func (r *ChiaDNSIntroducerReconciler) getChiaVolumeMounts(ctx context.Context, dnsIntro k8schianetv1.ChiaDNSIntroducer) []corev1.VolumeMount {
	var v []corev1.VolumeMount

	// secret ca volume
	v = append(v, corev1.VolumeMount{
		Name:      "secret-ca",
		MountPath: "/chia-ca",
	})

	// CHIA_ROOT volume
	v = append(v, corev1.VolumeMount{
		Name:      "chiaroot",
		MountPath: "/chia-data",
	})

	return v
}

// getChiaEnv retrieves the environment variables from the Chia config struct
func (r *ChiaDNSIntroducerReconciler) getChiaEnv(ctx context.Context, dnsIntro k8schianetv1.ChiaDNSIntroducer) []corev1.EnvVar {
	var env []corev1.EnvVar

	// service env var
	env = append(env, corev1.EnvVar{
		Name:  "service",
		Value: "seeder",
	})

	// CHIA_ROOT env var
	env = append(env, corev1.EnvVar{
		Name:  "CHIA_ROOT",
		Value: "/chia-data",
	})

	// keys env var -- no keys required for a dnsIntro
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
	if dnsIntro.Spec.ChiaConfig.Testnet != nil && *dnsIntro.Spec.ChiaConfig.Testnet {
		env = append(env, corev1.EnvVar{
			Name:  "testnet",
			Value: "true",
		})
	}

	// network env var
	if dnsIntro.Spec.ChiaConfig.Network != nil && *dnsIntro.Spec.ChiaConfig.Network != "" {
		env = append(env, corev1.EnvVar{
			Name:  "network",
			Value: *dnsIntro.Spec.ChiaConfig.Network,
		})
	}

	// network_port env var
	if dnsIntro.Spec.ChiaConfig.NetworkPort != nil && *dnsIntro.Spec.ChiaConfig.NetworkPort != 0 {
		env = append(env, corev1.EnvVar{
			Name:  "network_port",
			Value: strconv.Itoa(int(*dnsIntro.Spec.ChiaConfig.NetworkPort)),
		})
	}

	// introducer_address env var
	if dnsIntro.Spec.ChiaConfig.IntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "introducer_address",
			Value: *dnsIntro.Spec.ChiaConfig.IntroducerAddress,
		})
	}

	// dns_introducer_address env var
	if dnsIntro.Spec.ChiaConfig.DNSIntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "dns_introducer_address",
			Value: *dnsIntro.Spec.ChiaConfig.DNSIntroducerAddress,
		})
	}

	// TZ env var
	if dnsIntro.Spec.ChiaConfig.Timezone != nil {
		env = append(env, corev1.EnvVar{
			Name:  "TZ",
			Value: *dnsIntro.Spec.ChiaConfig.Timezone,
		})
	}

	// log_level env var
	if dnsIntro.Spec.ChiaConfig.LogLevel != nil {
		env = append(env, corev1.EnvVar{
			Name:  "log_level",
			Value: *dnsIntro.Spec.ChiaConfig.LogLevel,
		})
	}

	return env
}

// getOwnerReference gives the common owner reference spec for ChiaDNSIntroducer related objects
func (r *ChiaDNSIntroducerReconciler) getOwnerReference(ctx context.Context, dnsIntro k8schianetv1.ChiaDNSIntroducer) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: dnsIntro.APIVersion,
			Kind:       dnsIntro.Kind,
			Name:       dnsIntro.Name,
			UID:        dnsIntro.UID,
			Controller: &consts.ControllerOwner,
		},
	}
}

// getFullNodePort determines the correct full_node port to use
func (r *ChiaDNSIntroducerReconciler) getFullNodePort(ctx context.Context, dnsIntro k8schianetv1.ChiaDNSIntroducer) int32 {
	if dnsIntro.Spec.ChiaConfig.Testnet != nil && *dnsIntro.Spec.ChiaConfig.Testnet {
		return consts.TestnetNodePort
	}
	return consts.MainnetNodePort
}
