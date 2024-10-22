/*
Copyright 2023 Chia Network Inc.
*/

package chiaseeder

import (
	"fmt"

	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	corev1 "k8s.io/api/core/v1"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
)

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func getChiaVolumes(seeder k8schianetv1.ChiaSeeder) []corev1.Volume {
	var v []corev1.Volume

	// secret ca volume
	if seeder.Spec.ChiaConfig.CASecretName != nil {
		v = append(v, corev1.Volume{
			Name: "secret-ca",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: *seeder.Spec.ChiaConfig.CASecretName,
				},
			},
		})
	}

	// CHIA_ROOT volume
	if kube.ShouldMakeVolumeClaim(seeder.Spec.Storage) {
		v = append(v, corev1.Volume{
			Name: "chiaroot",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: fmt.Sprintf(chiaseederNamePattern, seeder.Name),
				},
			},
		})
	} else {
		v = append(v, kube.GetExistingChiaRootVolume(seeder.Spec.Storage))
	}

	return v
}

// getChiaVolumeMounts retrieves the requisite volume mounts from the Chia config struct
func getChiaVolumeMounts(seeder k8schianetv1.ChiaSeeder) []corev1.VolumeMount {
	var v []corev1.VolumeMount

	// secret ca volume
	if seeder.Spec.ChiaConfig.CASecretName != nil {
		v = append(v, corev1.VolumeMount{
			Name:      "secret-ca",
			MountPath: "/chia-ca",
		})
	}

	// CHIA_ROOT volume
	v = append(v, corev1.VolumeMount{
		Name:      "chiaroot",
		MountPath: "/chia-data",
	})

	return v
}

// getChiaEnv retrieves the environment variables from the Chia config struct
func getChiaEnv(seeder k8schianetv1.ChiaSeeder, networkData *map[string]string) ([]corev1.EnvVar, error) {
	var env []corev1.EnvVar

	// service env var
	env = append(env, corev1.EnvVar{
		Name:  "service",
		Value: "seeder",
	})

	// keys env var -- no keys required for a seeder
	env = append(env, corev1.EnvVar{
		Name:  "keys",
		Value: "none",
	})

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

	// Add common env
	commonEnv, err := kube.GetCommonChiaEnv(seeder.Spec.ChiaConfig.CommonSpecChia, networkData)
	if err != nil {
		return env, err
	}
	env = append(env, commonEnv...)

	return env, nil
}

// getChiaPorts returns the ports to a chia container
func getChiaPorts(fullNodePort int32) []corev1.ContainerPort {
	return []corev1.ContainerPort{
		{
			Name:          "daemon",
			ContainerPort: consts.DaemonPort,
			Protocol:      "TCP",
		},
		{
			Name:          "dns",
			ContainerPort: 53,
			Protocol:      "UDP",
		},
		{
			Name:          "dns-tcp",
			ContainerPort: 53,
			Protocol:      "TCP",
		},
		{
			Name:          "peers",
			ContainerPort: fullNodePort,
			Protocol:      "TCP",
		},
		{
			Name:          "rpc",
			ContainerPort: consts.CrawlerRPCPort,
			Protocol:      "TCP",
		},
	}
}
