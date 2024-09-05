/*
Copyright 2023 Chia Network Inc.
*/

package chiawallet

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	corev1 "k8s.io/api/core/v1"
)

// getChiaPorts returns the ports to a chia container
func getChiaPorts() []corev1.ContainerPort {
	return []corev1.ContainerPort{
		{
			Name:          "daemon",
			ContainerPort: consts.DaemonPort,
			Protocol:      "TCP",
		},
		{
			Name:          "peers",
			ContainerPort: consts.WalletPort,
			Protocol:      "TCP",
		},
		{
			Name:          "rpc",
			ContainerPort: consts.WalletRPCPort,
			Protocol:      "TCP",
		},
	}
}

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func getChiaVolumes(wallet k8schianetv1.ChiaWallet) []corev1.Volume {
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

	// CHIA_ROOT volume
	if kube.ShouldMakeVolumeClaim(wallet.Spec.Storage) {
		v = append(v, corev1.Volume{
			Name: "chiaroot",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: fmt.Sprintf(chiawalletNamePattern, wallet.Name),
				},
			},
		})
	} else {
		v = append(v, kube.GetChiaRootVolume(wallet.Spec.Storage))
	}

	// Add sidecar volumes if any exist
	if len(wallet.Spec.Sidecars.Volumes) > 0 {
		v = append(v, wallet.Spec.Sidecars.Volumes...)
	}

	return v
}

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
func getChiaEnv(ctx context.Context, wallet k8schianetv1.ChiaWallet) []corev1.EnvVar {
	logr := log.FromContext(ctx)
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

	// network env var
	if wallet.Spec.ChiaConfig.Network != nil && *wallet.Spec.ChiaConfig.Network != "" {
		env = append(env, corev1.EnvVar{
			Name:  "network",
			Value: *wallet.Spec.ChiaConfig.Network,
		})
	}

	// network_port env var
	if wallet.Spec.ChiaConfig.NetworkPort != nil && *wallet.Spec.ChiaConfig.NetworkPort != 0 {
		env = append(env, corev1.EnvVar{
			Name:  "network_port",
			Value: strconv.Itoa(int(*wallet.Spec.ChiaConfig.NetworkPort)),
		})
	}

	// introducer_address env var
	if wallet.Spec.ChiaConfig.IntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "introducer_address",
			Value: *wallet.Spec.ChiaConfig.IntroducerAddress,
		})
	}

	// dns_introducer_address env var
	if wallet.Spec.ChiaConfig.DNSIntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "dns_introducer_address",
			Value: *wallet.Spec.ChiaConfig.DNSIntroducerAddress,
		})
	}

	// trusted_cidrs env var
	if wallet.Spec.ChiaConfig.TrustedCIDRs != nil {
		// TODO should any special CIDR input checking happen here
		cidrs, err := json.Marshal(*wallet.Spec.ChiaConfig.TrustedCIDRs)
		if err != nil {
			logr.Error(err, fmt.Sprintf("ChiaWalletReconciler ChiaWallet=%s given CIDRs could not be marshalled to json. Peer connections that you would expect to be trusted might not be trusted.", wallet.Name))
		} else {
			env = append(env, corev1.EnvVar{
				Name:  "trusted_cidrs",
				Value: string(cidrs),
			})
		}
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

	// self_hostname env var
	if wallet.Spec.ChiaConfig.SelfHostname != nil {
		env = append(env, corev1.EnvVar{
			Name:  "self_hostname",
			Value: *wallet.Spec.ChiaConfig.SelfHostname,
		})
	} else {
		env = append(env, corev1.EnvVar{
			Name:  "self_hostname",
			Value: "0.0.0.0",
		})
	}

	// keys env var
	env = append(env, corev1.EnvVar{
		Name:  "keys",
		Value: fmt.Sprintf("/key/%s", wallet.Spec.ChiaConfig.SecretKey.Key),
	})

	// node peer env var
	if wallet.Spec.ChiaConfig.FullNodePeer != "" {
		env = append(env, corev1.EnvVar{
			Name:  "full_node_peer",
			Value: wallet.Spec.ChiaConfig.FullNodePeer,
		})
	}

	return env
}
