/*
Copyright 2023 Chia Network Inc.
*/

package chianode

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
)

// getChiaVolumes retrieves the requisite volumes from the Chia config struct
func getChiaVolumesAndTemplates(node k8schianetv1.ChiaNode) ([]corev1.Volume, []corev1.PersistentVolumeClaim) {
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

	// CHIA_ROOT volume
	rootVol, rootVolumeClaimTempl := getChiaRootVolume(node.Spec.Storage)
	if rootVolumeClaimTempl != nil {
		vcts = append(vcts, *rootVolumeClaimTempl)
	} else if rootVol != nil {
		v = append(v, *rootVol)
	}

	// Add sidecar volumes if any exist
	if len(node.Spec.Sidecars.Volumes) > 0 {
		v = append(v, node.Spec.Sidecars.Volumes...)
	}

	return v, vcts
}

// getChiaRootVolume gets the CHIA_ROOT volume for a Chia full_node.
// This function is unique to ChiaNodes because it's the only Kind that deploys a StatefulSet that can use PersistentVolumeClaimTemplates.
func getChiaRootVolume(storage *k8schianetv1.StorageConfig) (*corev1.Volume, *corev1.PersistentVolumeClaim) {
	volumeName := "chiaroot"
	if storage != nil && storage.ChiaRoot != nil {
		if storage.ChiaRoot.PersistentVolumeClaim != nil {
			// Get AccessModes, default to RWO
			accessModes := []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
			if len(storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
				accessModes = storage.ChiaRoot.PersistentVolumeClaim.AccessModes
			}

			// Parses the resource requests, and if there's an error this will fall through to hostPath config or emptyDir
			resourceReq, err := resource.ParseQuantity(storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
			if err == nil {
				return nil, &corev1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{
						Name: volumeName,
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes:      accessModes,
						StorageClassName: &storage.ChiaRoot.PersistentVolumeClaim.StorageClass,
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resourceReq,
							},
						},
					},
				}
			}
		} else if storage.ChiaRoot.HostPathVolume != nil && storage.ChiaRoot.HostPathVolume.Path != "" {
			return &corev1.Volume{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: storage.ChiaRoot.HostPathVolume.Path,
					},
				},
			}, nil
		}
	}

	return &corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}, nil
}

// getChiaVolumeMounts retrieves the requisite volume mounts from the Chia config struct
func getChiaVolumeMounts() []corev1.VolumeMount {
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
func getChiaEnv(ctx context.Context, node k8schianetv1.ChiaNode, networkData *map[string]string) ([]corev1.EnvVar, error) {
	logr := log.FromContext(ctx)
	var env []corev1.EnvVar

	// service env var
	env = append(env, corev1.EnvVar{
		Name:  "service",
		Value: "node",
	})

	// keys env var -- no keys required for a node
	env = append(env, corev1.EnvVar{
		Name:  "keys",
		Value: "none",
	})

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

	// Add common env
	commonEnv, err := kube.GetCommonChiaEnv(node.Spec.ChiaConfig.CommonSpecChia, networkData)
	if err != nil {
		return env, err
	}
	env = append(env, commonEnv...)

	return env, nil
}
