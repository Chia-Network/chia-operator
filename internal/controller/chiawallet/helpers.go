/*
Copyright 2023 Chia Network Inc.
*/

package chiawallet

import (
	"context"
	"fmt"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func (r *ChiaWalletReconciler) getChiaVolumes(ctx context.Context, wallet k8schianetv1.ChiaWallet) []corev1.Volume {
	var v []corev1.Volume

	// secret ca volume
	v = append(v, corev1.Volume{
		Name: "secret-ca",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: wallet.Spec.ChiaConfig.CASecretName,
			},
		},
	})

	// mnemonic key volume
	v = append(v, corev1.Volume{
		Name: "key",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: wallet.Spec.ChiaConfig.SecretKey.Name,
			},
		},
	})

	// CHIA_ROOT volume -- PVC is respected first if both it and hostpath are specified, falls back to hostPath if specified
	// If both are empty, fall back to emptyDir so chia-exporter can mount CHIA_ROOT
	var chiaRootAdded bool = false
	if wallet.Spec.Storage != nil && wallet.Spec.Storage.ChiaRoot != nil {
		if wallet.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil {
			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: wallet.Spec.Storage.ChiaRoot.HostPathVolume.Path,
					},
				},
			})
			chiaRootAdded = true
		} else if wallet.Spec.Storage.ChiaRoot.HostPathVolume != nil {
			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: wallet.Spec.Storage.ChiaRoot.HostPathVolume.Path,
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
func (r *ChiaWalletReconciler) getChiaEnv(ctx context.Context, wallet k8schianetv1.ChiaWallet) []corev1.EnvVar {
	var env []corev1.EnvVar

	// service env var
	env = append(env, corev1.EnvVar{
		Name:  "service",
		Value: "wallet",
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
	if wallet.Spec.ChiaConfig.Testnet != nil && *wallet.Spec.ChiaConfig.Testnet {
		env = append(env, corev1.EnvVar{
			Name:  "testnet",
			Value: "true",
		})
	}

	// TZ env var
	if wallet.Spec.ChiaConfig.Timezone != nil {
		env = append(env, corev1.EnvVar{
			Name:  "TZ",
			Value: *wallet.Spec.ChiaConfig.Timezone,
		})
	}

	// log_level env var
	if wallet.Spec.ChiaConfig.LogLevel != nil {
		env = append(env, corev1.EnvVar{
			Name:  "log_level",
			Value: *wallet.Spec.ChiaConfig.LogLevel,
		})
	}

	// keys env var
	env = append(env, corev1.EnvVar{
		Name:  "keys",
		Value: fmt.Sprintf("/key/%s", wallet.Spec.ChiaConfig.SecretKey.Key),
	})

	// node peer env var
	env = append(env, corev1.EnvVar{
		Name:  "full_node_peer",
		Value: wallet.Spec.ChiaConfig.FullNodePeer,
	})

	return env
}

// getLabels gives some common labels for ChiaWallet related objects
func (r *ChiaWalletReconciler) getLabels(ctx context.Context, wallet k8schianetv1.ChiaWallet, additionalLabels ...map[string]string) map[string]string {
	var labels = make(map[string]string)
	for _, addition := range additionalLabels {
		for k, v := range addition {
			labels[k] = v
		}
	}
	labels["app.kubernetes.io/instance"] = wallet.Name
	labels["app.kubernetes.io/name"] = wallet.Name
	labels = kube.GetCommonLabels(ctx, labels)
	return labels
}

// getOwnerReference gives the common owner reference spec for ChiaWallet related objects
func (r *ChiaWalletReconciler) getOwnerReference(ctx context.Context, wallet k8schianetv1.ChiaWallet) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: wallet.APIVersion,
			Kind:       wallet.Kind,
			Name:       wallet.Name,
			UID:        wallet.UID,
			Controller: &consts.ControllerOwner,
		},
	}
}
