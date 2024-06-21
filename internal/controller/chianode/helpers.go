/*
Copyright 2023 Chia Network Inc.
*/

package chianode

import (
	"context"
	"encoding/json"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
)

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func (r *ChiaNodeReconciler) getChiaVolumesAndTemplates(ctx context.Context, node k8schianetv1.ChiaNode) ([]corev1.Volume, []corev1.PersistentVolumeClaim) {
	var v []corev1.Volume
	var vcts []corev1.PersistentVolumeClaim

	// secret ca volume
	v = append(v, corev1.Volume{
		Name: "secret-ca",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: node.Spec.ChiaConfig.CASecretName,
			},
		},
	})

	// CHIA_ROOT volume -- PVC is respected first if both it and hostpath are specified, falls back to hostPath if specified
	// If both are empty, fall back to emptyDir so chia-exporter can mount CHIA_ROOT
	var chiaRootAdded bool = false
	if node.Spec.Storage != nil && node.Spec.Storage.ChiaRoot != nil {
		if node.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil {
			vcts = append(vcts, corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: "chiaroot",
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes:      []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
					StorageClassName: &node.Spec.Storage.ChiaRoot.PersistentVolumeClaim.StorageClass,
					Resources: corev1.VolumeResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceStorage: resource.MustParse(node.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest),
						},
					},
				},
			})
			chiaRootAdded = true
		} else if node.Spec.Storage.ChiaRoot.HostPathVolume != nil {
			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: node.Spec.Storage.ChiaRoot.HostPathVolume.Path,
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
	if len(node.Spec.Sidecars.Volumes) > 0 {
		v = append(v, node.Spec.Sidecars.Volumes...)
	}

	return v, vcts
}

// getChiaVolumeMounts retrieves the requisite volume mounts from the Chia config struct
func (r *ChiaNodeReconciler) getChiaVolumeMounts(ctx context.Context, node k8schianetv1.ChiaNode) []corev1.VolumeMount {
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
func (r *ChiaNodeReconciler) getChiaEnv(ctx context.Context, node k8schianetv1.ChiaNode) []corev1.EnvVar {
	logr := log.FromContext(ctx)
	var env []corev1.EnvVar

	// service env var
	env = append(env, corev1.EnvVar{
		Name:  "service",
		Value: "node",
	})

	// CHIA_ROOT env var
	env = append(env, corev1.EnvVar{
		Name:  "CHIA_ROOT",
		Value: "/chia-data",
	})

	// keys env var -- no keys required for a node
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
	if node.Spec.ChiaConfig.Testnet != nil && *node.Spec.ChiaConfig.Testnet {
		env = append(env, corev1.EnvVar{
			Name:  "testnet",
			Value: "true",
		})
	}

	// network env var
	if node.Spec.ChiaConfig.Network != nil && *node.Spec.ChiaConfig.Network != "" {
		env = append(env, corev1.EnvVar{
			Name:  "network",
			Value: *node.Spec.ChiaConfig.Network,
		})
	}

	// network_port env var
	if node.Spec.ChiaConfig.NetworkPort != nil && *node.Spec.ChiaConfig.NetworkPort != 0 {
		env = append(env, corev1.EnvVar{
			Name:  "network_port",
			Value: strconv.Itoa(int(*node.Spec.ChiaConfig.NetworkPort)),
		})
	}

	// introducer_address env var
	if node.Spec.ChiaConfig.IntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "introducer_address",
			Value: *node.Spec.ChiaConfig.IntroducerAddress,
		})
	}

	// dns_introducer_address env var
	if node.Spec.ChiaConfig.DNSIntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "dns_introducer_address",
			Value: *node.Spec.ChiaConfig.DNSIntroducerAddress,
		})
	}

	// trusted_cidrs env var
	if node.Spec.ChiaConfig.TrustedCIDRs != nil {
		// TODO should any special CIDR input checking happen here
		cidrs, err := json.Marshal(*node.Spec.ChiaConfig.TrustedCIDRs)
		if err != nil {
			logr.Error(err, fmt.Sprintf("ChiaNodeReconciler ChiaNode=%s given CIDRs could not be marshalled to json. Peer connections that you would expect to be trusted might not be trusted.", node.Name))
		} else {
			env = append(env, corev1.EnvVar{
				Name:  "trusted_cidrs",
				Value: string(cidrs),
			})
		}
	}

	// TZ env var
	if node.Spec.ChiaConfig.Timezone != nil {
		env = append(env, corev1.EnvVar{
			Name:  "TZ",
			Value: *node.Spec.ChiaConfig.Timezone,
		})
	}

	// log_level env var
	if node.Spec.ChiaConfig.LogLevel != nil {
		env = append(env, corev1.EnvVar{
			Name:  "log_level",
			Value: *node.Spec.ChiaConfig.LogLevel,
		})
	}

	return env
}

// getOwnerReference gives the common owner reference spec for ChiaNode related objects
func (r *ChiaNodeReconciler) getOwnerReference(ctx context.Context, node k8schianetv1.ChiaNode) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: node.APIVersion,
			Kind:       node.Kind,
			Name:       node.Name,
			UID:        node.UID,
			Controller: &consts.ControllerOwner,
		},
	}
}

// getFullNodePort determines the correct full node port to use
func (r *ChiaNodeReconciler) getFullNodePort(ctx context.Context, node k8schianetv1.ChiaNode) int32 {
	if node.Spec.ChiaConfig.Testnet != nil && *node.Spec.ChiaConfig.Testnet {
		return consts.TestnetNodePort
	}
	return consts.MainnetNodePort
}
