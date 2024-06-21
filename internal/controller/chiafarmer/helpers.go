/*
Copyright 2023 Chia Network Inc.
*/

package chiafarmer

import (
	"context"
	"fmt"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
)

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func (r *ChiaFarmerReconciler) getChiaVolumes(ctx context.Context, farmer k8schianetv1.ChiaFarmer) []corev1.Volume {
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

	// CHIA_ROOT volume -- PVC is respected first if both it and hostpath are specified, falls back to hostPath if specified
	// If both are empty, fall back to emptyDir so chia-exporter can mount CHIA_ROOT
	var chiaRootAdded bool = false
	if farmer.Spec.Storage != nil && farmer.Spec.Storage.ChiaRoot != nil {
		if farmer.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil {
			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: farmer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ClaimName,
					},
				},
			})
			chiaRootAdded = true
		} else if farmer.Spec.Storage.ChiaRoot.HostPathVolume != nil {
			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: farmer.Spec.Storage.ChiaRoot.HostPathVolume.Path,
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
	if len(farmer.Spec.Sidecars.Volumes) > 0 {
		v = append(v, farmer.Spec.Sidecars.Volumes...)
	}

	return v
}

// getChiaVolumeMounts retrieves the requisite volume mounts from the Chia config struct
func (r *ChiaFarmerReconciler) getChiaVolumeMounts(ctx context.Context, farmer k8schianetv1.ChiaFarmer) []corev1.VolumeMount {
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
func (r *ChiaFarmerReconciler) getChiaEnv(ctx context.Context, farmer k8schianetv1.ChiaFarmer) []corev1.EnvVar {
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

// getOwnerReference gives the common owner reference spec for ChiaFarmer related objects
func (r *ChiaFarmerReconciler) getOwnerReference(ctx context.Context, farmer k8schianetv1.ChiaFarmer) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: farmer.APIVersion,
			Kind:       farmer.Kind,
			Name:       farmer.Name,
			UID:        farmer.UID,
			Controller: &consts.ControllerOwner,
		},
	}
}
