/*
Copyright 2023 Chia Network Inc.
*/

package chiaharvester

import (
	"fmt"
	"strconv"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	corev1 "k8s.io/api/core/v1"
)

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func getChiaVolumes(harvester k8schianetv1.ChiaHarvester) []corev1.Volume {
	var v []corev1.Volume

	// secret ca volume
	v = append(v, corev1.Volume{
		Name: "secret-ca",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: harvester.Spec.ChiaConfig.CASecretName,
			},
		},
	})

	// CHIA_ROOT volume -- PVC is respected first if both it and hostpath are specified, falls back to hostPath if specified
	// If both are empty, fall back to emptyDir so chia-exporter can mount CHIA_ROOT
	var chiaRootAdded = false
	if harvester.Spec.Storage != nil && harvester.Spec.Storage.ChiaRoot != nil {
		if harvester.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil {
			var pvcName string
			if harvester.Spec.Storage.ChiaRoot.PersistentVolumeClaim.GenerateVolumeClaims {
				pvcName = fmt.Sprintf(chiaharvesterNamePattern, harvester.Name)
			} else {
				pvcName = harvester.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ClaimName
			}

			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: pvcName,
					},
				},
			})
			chiaRootAdded = true
		} else if harvester.Spec.Storage.ChiaRoot.HostPathVolume != nil {
			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: harvester.Spec.Storage.ChiaRoot.HostPathVolume.Path,
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

	// hostPath and PVC plot volumes
	if harvester.Spec.Storage != nil {
		if harvester.Spec.Storage.Plots != nil {
			// PVC plot volumes
			if harvester.Spec.Storage.Plots.PersistentVolumeClaim != nil {
				for i, vol := range harvester.Spec.Storage.Plots.PersistentVolumeClaim {
					if vol != nil {
						v = append(v, corev1.Volume{
							Name: fmt.Sprintf("pvc-plots-%d", i),
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: vol.ClaimName,
								},
							},
						})
					}
				}
			}

			// hostPath plot volumes
			if harvester.Spec.Storage.Plots.HostPathVolume != nil {
				for i, vol := range harvester.Spec.Storage.Plots.HostPathVolume {
					if vol != nil {
						v = append(v, corev1.Volume{
							Name: fmt.Sprintf("hostpath-plots-%d", i),
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: vol.Path,
								},
							},
						})
					}
				}
			}
		}
	}

	// Add sidecar volumes if any exist
	if len(harvester.Spec.Sidecars.Volumes) > 0 {
		v = append(v, harvester.Spec.Sidecars.Volumes...)
	}

	return v
}

// getChiaVolumeMounts retrieves the requisite volume mounts from the Chia config struct
func getChiaVolumeMounts(harvester k8schianetv1.ChiaHarvester) []corev1.VolumeMount {
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

	// hostPath and PVC plot volumemounts
	if harvester.Spec.Storage != nil {
		if harvester.Spec.Storage.Plots != nil {
			// PVC plot volume mounts
			if harvester.Spec.Storage.Plots.PersistentVolumeClaim != nil {
				for i, vol := range harvester.Spec.Storage.Plots.PersistentVolumeClaim {
					if vol != nil {
						v = append(v, corev1.VolumeMount{
							Name:      fmt.Sprintf("pvc-plots-%d", i),
							ReadOnly:  true,
							MountPath: fmt.Sprintf("/plots/pvc-plots-%d", i),
						})
					}
				}
			}

			// hostPath plot volume mounts
			if harvester.Spec.Storage.Plots.HostPathVolume != nil {
				for i, vol := range harvester.Spec.Storage.Plots.HostPathVolume {
					if vol != nil {
						v = append(v, corev1.VolumeMount{
							Name:      fmt.Sprintf("hostpath-plots-%d", i),
							ReadOnly:  true,
							MountPath: fmt.Sprintf("/plots/hostpath-plots-%d", i),
						})
					}
				}
			}
		}
	}

	return v
}

// getChiaEnv retrieves the environment variables from the Chia config struct
func getChiaEnv(harvester k8schianetv1.ChiaHarvester) []corev1.EnvVar {
	var env []corev1.EnvVar

	// service env var
	env = append(env, corev1.EnvVar{
		Name:  "service",
		Value: "harvester",
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
	if harvester.Spec.ChiaConfig.Testnet != nil && *harvester.Spec.ChiaConfig.Testnet {
		env = append(env, corev1.EnvVar{
			Name:  "testnet",
			Value: "true",
		})
	}

	// network env var
	if harvester.Spec.ChiaConfig.Network != nil && *harvester.Spec.ChiaConfig.Network != "" {
		env = append(env, corev1.EnvVar{
			Name:  "network",
			Value: *harvester.Spec.ChiaConfig.Network,
		})
	}

	// network_port env var
	if harvester.Spec.ChiaConfig.NetworkPort != nil && *harvester.Spec.ChiaConfig.NetworkPort != 0 {
		env = append(env, corev1.EnvVar{
			Name:  "network_port",
			Value: strconv.Itoa(int(*harvester.Spec.ChiaConfig.NetworkPort)),
		})
	}

	// introducer_address env var
	if harvester.Spec.ChiaConfig.IntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "introducer_address",
			Value: *harvester.Spec.ChiaConfig.IntroducerAddress,
		})
	}

	// dns_introducer_address env var
	if harvester.Spec.ChiaConfig.DNSIntroducerAddress != nil {
		env = append(env, corev1.EnvVar{
			Name:  "dns_introducer_address",
			Value: *harvester.Spec.ChiaConfig.DNSIntroducerAddress,
		})
	}

	// TZ env var
	if harvester.Spec.ChiaConfig.Timezone != nil {
		env = append(env, corev1.EnvVar{
			Name:  "TZ",
			Value: *harvester.Spec.ChiaConfig.Timezone,
		})
	}

	// log_level env var
	if harvester.Spec.ChiaConfig.LogLevel != nil {
		env = append(env, corev1.EnvVar{
			Name:  "log_level",
			Value: *harvester.Spec.ChiaConfig.LogLevel,
		})
	}

	// self_hostname env var
	if harvester.Spec.ChiaConfig.SelfHostname != nil {
		env = append(env, corev1.EnvVar{
			Name:  "self_hostname",
			Value: *harvester.Spec.ChiaConfig.SelfHostname,
		})
	} else {
		env = append(env, corev1.EnvVar{
			Name:  "self_hostname",
			Value: "0.0.0.0",
		})
	}

	// recursive_plot_scan env var -- needed because all plot drives are just mounted as subdirs under `/plots`.
	// TODO make plot mount paths configurable -- make this var optional
	env = append(env, corev1.EnvVar{
		Name:  "recursive_plot_scan",
		Value: "true",
	})

	// farmer peer env vars
	env = append(env, corev1.EnvVar{
		Name:  "farmer_address",
		Value: harvester.Spec.ChiaConfig.FarmerAddress,
	})
	env = append(env, corev1.EnvVar{
		Name:  "farmer_port",
		Value: strconv.Itoa(consts.FarmerPort),
	})

	return env
}
