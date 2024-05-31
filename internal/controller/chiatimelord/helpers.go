/*
Copyright 2023 Chia Network Inc.
*/

package chiatimelord

import (
	"context"
	"encoding/json"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
)

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func (r *ChiaTimelordReconciler) getChiaVolumes(ctx context.Context, tl k8schianetv1.ChiaTimelord) []corev1.Volume {
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

	// CHIA_ROOT volume -- PVC is respected first if both it and hostpath are specified, falls back to hostPath if specified
	// If both are empty, fall back to emptyDir so chia-exporter can mount CHIA_ROOT
	var chiaRootAdded bool = false
	if tl.Spec.Storage != nil && tl.Spec.Storage.ChiaRoot != nil {
		if tl.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil {
			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: tl.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ClaimName,
					},
				},
			})
			chiaRootAdded = true
		} else if tl.Spec.Storage.ChiaRoot.HostPathVolume != nil {
			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: tl.Spec.Storage.ChiaRoot.HostPathVolume.Path,
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
	if len(tl.Spec.Sidecars.Volumes) > 0 {
		v = append(v, tl.Spec.Sidecars.Volumes...)
	}

	return v
}

// getChiaEnv retrieves the environment variables from the Chia config struct
func (r *ChiaTimelordReconciler) getChiaEnv(ctx context.Context, tl k8schianetv1.ChiaTimelord) []corev1.EnvVar {
	logr := log.FromContext(ctx)
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

	// trusted_cidrs env var
	if tl.Spec.ChiaConfig.TrustedCIDRs != nil {
		// TODO should any special CIDR input checking happen here? Chia might also do that for me.
		cidrs, err := json.Marshal(*tl.Spec.ChiaConfig.TrustedCIDRs)
		if err != nil {
			logr.Error(err, fmt.Sprintf("ChiaTimelordReconciler ChiaTimelord=%s given CIDRs could not be marshalled to json. Peer connections that you would expect to be trusted might not be trusted.", tl.Name))
		} else {
			env = append(env, corev1.EnvVar{
				Name:  "trusted_cidrs",
				Value: string(cidrs),
			})
		}
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

	// node peer env var
	env = append(env, corev1.EnvVar{
		Name:  "full_node_peer",
		Value: tl.Spec.ChiaConfig.FullNodePeer,
	})

	return env
}

// getOwnerReference gives the common owner reference spec for ChiaTimelord related objects
func (r *ChiaTimelordReconciler) getOwnerReference(ctx context.Context, tl k8schianetv1.ChiaTimelord) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: tl.APIVersion,
			Kind:       tl.Kind,
			Name:       tl.Name,
			UID:        tl.UID,
			Controller: &consts.ControllerOwner,
		},
	}
}
