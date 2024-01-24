/*
Copyright 2023 Chia Network Inc.
*/

package chiaseeder

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
func (r *ChiaSeederReconciler) getChiaVolumes(ctx context.Context, seeder k8schianetv1.ChiaSeeder) []corev1.Volume {
	var v []corev1.Volume

	// secret ca volume
	v = append(v, corev1.Volume{
		Name: "secret-ca",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: seeder.Spec.ChiaConfig.CASecretName,
			},
		},
	})

	// CHIA_ROOT volume -- PVC is respected first if both it and hostpath are specified, falls back to hostPath if specified
	// If both are empty, fall back to emptyDir so chia-exporter can mount CHIA_ROOT
	var chiaRootAdded bool = false
	if seeder.Spec.Storage != nil && seeder.Spec.Storage.ChiaRoot != nil {
		if seeder.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil {
			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: seeder.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ClaimName,
					},
				},
			})
			chiaRootAdded = true
		} else if seeder.Spec.Storage.ChiaRoot.HostPathVolume != nil {
			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: seeder.Spec.Storage.ChiaRoot.HostPathVolume.Path,
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
func (r *ChiaSeederReconciler) getChiaVolumeMounts(ctx context.Context, seeder k8schianetv1.ChiaSeeder) []corev1.VolumeMount {
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
func (r *ChiaSeederReconciler) getChiaEnv(ctx context.Context, seeder k8schianetv1.ChiaSeeder) []corev1.EnvVar {
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

	// keys env var -- no keys required for a seeder
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
	if seeder.Spec.ChiaConfig.Testnet != nil && *seeder.Spec.ChiaConfig.Testnet {
		env = append(env, corev1.EnvVar{
			Name:  "testnet",
			Value: "true",
		})
	}

	// network env var
	if seeder.Spec.ChiaConfig.Network != nil && *seeder.Spec.ChiaConfig.Network != "" {
		env = append(env, corev1.EnvVar{
			Name:  "network",
			Value: *seeder.Spec.ChiaConfig.Network,
		})
	}

	// network_port env var
	if seeder.Spec.ChiaConfig.NetworkPort != nil && *seeder.Spec.ChiaConfig.NetworkPort != 0 {
		env = append(env, corev1.EnvVar{
			Name:  "network_port",
			Value: strconv.Itoa(int(*seeder.Spec.ChiaConfig.NetworkPort)),
		})
	}

	// introducer_address env var
	if seeder.Spec.ChiaConfig.IntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "introducer_address",
			Value: *seeder.Spec.ChiaConfig.IntroducerAddress,
		})
	}

	// dns_introducer_address env var
	if seeder.Spec.ChiaConfig.DNSIntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "dns_introducer_address",
			Value: *seeder.Spec.ChiaConfig.DNSIntroducerAddress,
		})
	}

	// TZ env var
	if seeder.Spec.ChiaConfig.Timezone != nil {
		env = append(env, corev1.EnvVar{
			Name:  "TZ",
			Value: *seeder.Spec.ChiaConfig.Timezone,
		})
	}

	// log_level env var
	if seeder.Spec.ChiaConfig.LogLevel != nil {
		env = append(env, corev1.EnvVar{
			Name:  "log_level",
			Value: *seeder.Spec.ChiaConfig.LogLevel,
		})
	}

	// seeder_bootstrap_peers env var
	if seeder.Spec.ChiaConfig.BootstrapPeer != nil {
		env = append(env, corev1.EnvVar{
			Name:  "seeder_bootstrap_peers",
			Value: *seeder.Spec.ChiaConfig.BootstrapPeer,
		})
	}

	// seeder_minimum_height env var
	if seeder.Spec.ChiaConfig.MinimumHeight != nil {
		env = append(env, corev1.EnvVar{
			Name:  "seeder_minimum_height",
			Value: fmt.Sprintf("%d", *seeder.Spec.ChiaConfig.MinimumHeight),
		})
	}

	// seeder_domain_name env var
	env = append(env, corev1.EnvVar{
		Name:  "seeder_domain_name",
		Value: seeder.Spec.ChiaConfig.DomainName,
	})

	// seeder_nameserver env var
	env = append(env, corev1.EnvVar{
		Name:  "seeder_nameserver",
		Value: seeder.Spec.ChiaConfig.Nameserver,
	})

	// seeder_soa_rname env var
	env = append(env, corev1.EnvVar{
		Name:  "seeder_soa_rname",
		Value: seeder.Spec.ChiaConfig.Rname,
	})

	// seeder_ttl env var
	if seeder.Spec.ChiaConfig.TTL != nil {
		env = append(env, corev1.EnvVar{
			Name:  "seeder_ttl",
			Value: fmt.Sprintf("%d", *seeder.Spec.ChiaConfig.TTL),
		})
	}

	return env
}

// getOwnerReference gives the common owner reference spec for ChiaSeeder related objects
func (r *ChiaSeederReconciler) getOwnerReference(ctx context.Context, seeder k8schianetv1.ChiaSeeder) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: seeder.APIVersion,
			Kind:       seeder.Kind,
			Name:       seeder.Name,
			UID:        seeder.UID,
			Controller: &consts.ControllerOwner,
		},
	}
}

// getFullNodePort determines the correct full_node port to use
func (r *ChiaSeederReconciler) getFullNodePort(ctx context.Context, seeder k8schianetv1.ChiaSeeder) int32 {
	if seeder.Spec.ChiaConfig.Testnet != nil && *seeder.Spec.ChiaConfig.Testnet {
		return consts.TestnetNodePort
	}
	return consts.MainnetNodePort
}
