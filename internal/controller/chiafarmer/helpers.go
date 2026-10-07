/*
Copyright 2023 Chia Network Inc.
*/

package chiafarmer

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/chia-network/chia-operator/internal/controller/common/kube"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
)

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func getChiaVolumes(farmer k8schianetv1.ChiaFarmer) []corev1.Volume {
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

	// CHIA_ROOT volume
	if kube.ShouldMakeChiaRootVolumeClaim(farmer.Spec.Storage) {
		v = append(v, corev1.Volume{
			Name: "chiaroot",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: fmt.Sprintf(chiafarmerNamePattern, farmer.Name),
				},
			},
		})
	} else {
		v = append(v, kube.GetExistingChiaRootVolume(farmer.Spec.Storage))
	}

	return v
}

// getChiaVolumeMounts retrieves the requisite volume mounts from the Chia config struct
func getChiaVolumeMounts() []corev1.VolumeMount {
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
func getChiaEnv(ctx context.Context, farmer k8schianetv1.ChiaFarmer, networkData *map[string]string) ([]corev1.EnvVar, error) {
	logr := log.FromContext(ctx)
	var env []corev1.EnvVar

	// service env var
	env = append(env, corev1.EnvVar{
		Name:  "service",
		Value: "farmer-only",
	})

	// keys env var
	env = append(env, corev1.EnvVar{
		Name:  "keys",
		Value: fmt.Sprintf("/key/%s", farmer.Spec.ChiaConfig.SecretKey.Key),
	})

	// node peer env var
	if farmer.Spec.ChiaConfig.FullNodePeers != nil {
		fnp, err := kube.MarshalFullNodePeers(*farmer.Spec.ChiaConfig.FullNodePeers)
		if err != nil {
			logr.Error(err, "given full_node peers could not be marshaled to JSON, they may not appear in your chia configuration")
		} else {
			env = append(env, corev1.EnvVar{
				Name:  "chia.farmer.full_node_peers",
				Value: string(fnp),
			})
		}
	} else if farmer.Spec.ChiaConfig.FullNodePeer != nil {
		env = append(env, corev1.EnvVar{
			Name:  "full_node_peer",
			Value: *farmer.Spec.ChiaConfig.FullNodePeer,
		})
	}

	// Add common env
	commonEnv, err := kube.GetCommonChiaEnv(farmer.Spec.ChiaConfig.CommonSpecChia, networkData)
	if err != nil {
		return env, err
	}
	env = append(env, commonEnv...)

	return env, nil
}
