/*
Copyright 2023 Chia Network Inc.
*/

package chiaharvester

import (
	"context"
	"fmt"
	"strconv"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func (r *ChiaHarvesterReconciler) getChiaVolumes(ctx context.Context, harvester k8schianetv1.ChiaHarvester) []corev1.Volume {
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
	var chiaRootAdded bool = false
	if harvester.Spec.Storage != nil && harvester.Spec.Storage.ChiaRoot != nil {
		if harvester.Spec.Storage.ChiaRoot.PersistentVolumeClaim != nil {
			v = append(v, corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: harvester.Spec.Storage.ChiaRoot.HostPathVolume.Path,
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

	return v
}

// getChiaVolumeMounts retrieves the requisite volume mounts from the Chia config struct
func (r *ChiaHarvesterReconciler) getChiaVolumeMounts(ctx context.Context, harvester k8schianetv1.ChiaHarvester) []corev1.VolumeMount {
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
func (r *ChiaHarvesterReconciler) getChiaEnv(ctx context.Context, harvester k8schianetv1.ChiaHarvester) []corev1.EnvVar {
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

// getLabels gives some common labels for ChiaHarvester related objects
func (r *ChiaHarvesterReconciler) getLabels(ctx context.Context, harvester k8schianetv1.ChiaHarvester, additionalLabels ...map[string]string) map[string]string {
	var labels = make(map[string]string)
	for _, addition := range additionalLabels {
		for k, v := range addition {
			labels[k] = v
		}
	}
	labels["app.kubernetes.io/instance"] = harvester.Name
	labels["app.kubernetes.io/name"] = harvester.Name
	labels = kube.GetCommonLabels(ctx, labels)
	return labels
}

// getOwnerReference gives the common owner reference spec for ChiaHarvester related objects
func (r *ChiaHarvesterReconciler) getOwnerReference(ctx context.Context, harvester k8schianetv1.ChiaHarvester) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: harvester.APIVersion,
			Kind:       harvester.Kind,
			Name:       harvester.Name,
			UID:        harvester.UID,
			Controller: &consts.ControllerOwner,
		},
	}
}
