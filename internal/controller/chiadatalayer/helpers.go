/*
Copyright 2024 Chia Network Inc.
*/

package chiadatalayer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
)

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func getChiaVolumes(datalayer k8schianetv1.ChiaDataLayer) []corev1.Volume {
	var v []corev1.Volume

	// secret ca volume
	if datalayer.Spec.ChiaConfig.CASecretName != nil {
		v = append(v, corev1.Volume{
			Name: "secret-ca",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: *datalayer.Spec.ChiaConfig.CASecretName,
				},
			},
		})
	}

	// mnemonic key volume
	v = append(v, corev1.Volume{
		Name: "key",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: datalayer.Spec.ChiaConfig.SecretKey.Name,
			},
		},
	})

	// CHIA_ROOT volume
	if kube.ShouldMakeChiaRootVolumeClaim(datalayer.Spec.Storage) {
		v = append(v, corev1.Volume{
			Name: "chiaroot",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: fmt.Sprintf(chiadatalayerNamePattern, datalayer.Name),
				},
			},
		})
	} else {
		v = append(v, kube.GetExistingChiaRootVolume(datalayer.Spec.Storage))
	}

	// data_layer server files volume
	// TODO finish this
	if kube.ShouldMakeChiaRootVolumeClaim(datalayer.Spec.Storage) {
		v = append(v, corev1.Volume{
			Name: "server",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: fmt.Sprintf(chiadatalayerNamePattern, datalayer.Name),
				},
			},
		})
	} else {
		v = append(v, kube.GetExistingChiaRootVolume(datalayer.Spec.Storage))
	}

	return v
}

// getChiaVolumeMounts retrieves the requisite volume mounts from the Chia config struct
func getChiaVolumeMounts(datalayer k8schianetv1.ChiaDataLayer) []corev1.VolumeMount {
	var v []corev1.VolumeMount

	// secret ca volume
	if datalayer.Spec.ChiaConfig.CASecretName != nil {
		v = append(v, corev1.VolumeMount{
			Name:      "secret-ca",
			MountPath: "/chia-ca",
		})
	}

	// key volume
	v = append(v, corev1.VolumeMount{
		Name:      "key",
		MountPath: "/key",
	})

	// CHIA_ROOT volume
	v = append(v, corev1.VolumeMount{
		Name:      "chiaroot",
		MountPath: "/chia-data",
	})

	// data_layer server files volume
	v = append(v, corev1.VolumeMount{
		Name:      "server",
		MountPath: "/datalayer/server_files",
	})

	return v
}

// getChiaEnv retrieves the environment variables from the Chia config struct
func getChiaEnv(ctx context.Context, datalayer k8schianetv1.ChiaDataLayer, networkData *map[string]string) ([]corev1.EnvVar, error) {
	logr := log.FromContext(ctx)
	var env []corev1.EnvVar

	// service env var
	env = append(env, corev1.EnvVar{
		Name:  "service",
		Value: "data",
	})

	// keys env var
	env = append(env, corev1.EnvVar{
		Name:  "keys",
		Value: fmt.Sprintf("/key/%s", datalayer.Spec.ChiaConfig.SecretKey.Key),
	})

	env = append(env, corev1.EnvVar{
		Name:  "chia.data_layer.server_files_location",
		Value: "/datalayer/server_files",
	})

	// node peer env var
	if datalayer.Spec.ChiaConfig.FullNodePeers != nil {
		fnp, err := kube.MarshalFullNodePeers(*datalayer.Spec.ChiaConfig.FullNodePeers)
		if err != nil {
			logr.Error(err, "given full_node peers could not be marshaled to JSON, they may not appear in your chia configuration")
		} else {
			env = append(env, corev1.EnvVar{
				Name:  "chia.wallet.full_node_peers",
				Value: string(fnp),
			})
		}
	}

	// trusted_cidrs env var
	if datalayer.Spec.ChiaConfig.TrustedCIDRs != nil {
		// TODO should any special CIDR input checking happen here
		cidrs, err := json.Marshal(*datalayer.Spec.ChiaConfig.TrustedCIDRs)
		if err != nil {
			logr.Error(err, "given CIDRs could not be marshalled to json. Peer connections that you would expect to be trusted might not be trusted.")
		} else {
			env = append(env, corev1.EnvVar{
				Name:  "trusted_cidrs",
				Value: string(cidrs),
			})
		}
	}

	// Add common env
	commonEnv, err := kube.GetCommonChiaEnv(datalayer.Spec.ChiaConfig.CommonSpecChia, networkData)
	if err != nil {
		return env, err
	}
	env = append(env, commonEnv...)

	return env, nil
}

// getChiaPorts returns the ports to a chia container
func getChiaPorts() []corev1.ContainerPort {
	return []corev1.ContainerPort{
		{
			Name:          "daemon",
			ContainerPort: consts.DaemonPort,
			Protocol:      "TCP",
		},
		{
			Name:          "rpc",
			ContainerPort: consts.DataLayerRPCPort,
			Protocol:      "TCP",
		},
		{
			Name:          "wallet-rpc",
			ContainerPort: consts.WalletRPCPort,
			Protocol:      "TCP",
		},
	}
}
