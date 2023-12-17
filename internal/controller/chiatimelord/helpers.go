/*
Copyright 2023 Chia Network Inc.
*/

package chiatimelord

import (
	"context"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
						ClaimName: tl.Spec.Storage.ChiaRoot.HostPathVolume.Path,
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

	return v
}

// getChiaEnv retrieves the environment variables from the Chia config struct
func (r *ChiaTimelordReconciler) getChiaEnv(ctx context.Context, tl k8schianetv1.ChiaTimelord) []corev1.EnvVar {
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

// getLabels gives some common labels for ChiaTimelord related objects
func (r *ChiaTimelordReconciler) getLabels(ctx context.Context, tl k8schianetv1.ChiaTimelord, additionalLabels ...map[string]string) map[string]string {
	var labels = make(map[string]string)
	for _, addition := range additionalLabels {
		for k, v := range addition {
			labels[k] = v
		}
	}
	labels["app.kubernetes.io/instance"] = tl.Name
	labels["app.kubernetes.io/name"] = tl.Name
	labels = kube.GetCommonLabels(ctx, labels)
	return labels
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
