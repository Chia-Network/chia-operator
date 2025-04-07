package chiaharvester

import (
	"testing"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetChiaVolumeMounts(t *testing.T) {
	tests := []struct {
		name      string
		harvester k8schianetv1.ChiaHarvester
		want      []corev1.VolumeMount
	}{
		{
			name: "With Plot Storage",
			harvester: k8schianetv1.ChiaHarvester{
				Spec: k8schianetv1.ChiaHarvesterSpec{
					ChiaConfig: k8schianetv1.ChiaHarvesterSpecChia{
						CASecretName: "test-ca-secret",
					},
					CommonSpec: k8schianetv1.CommonSpec{
						Storage: &k8schianetv1.StorageConfig{
							Plots: &k8schianetv1.PlotsConfig{
								PersistentVolumeClaim: []*k8schianetv1.PersistentVolumeClaimConfig{
									{
										ClaimName: "plot-pvc-1",
									},
								},
								HostPathVolume: []*k8schianetv1.HostPathVolumeConfig{
									{
										Path: "/plots/hostpath-1",
									},
								},
							},
						},
					},
				},
			},
			want: []corev1.VolumeMount{
				{
					Name:      "secret-ca",
					MountPath: "/chia-ca",
				},
				{
					Name:      "chiaroot",
					MountPath: "/chia-data",
				},
				{
					Name:      "pvc-plots-0",
					ReadOnly:  true,
					MountPath: "/plots/pvc-plots-0",
				},
				{
					Name:      "hostpath-plots-0",
					ReadOnly:  true,
					MountPath: "/plots/hostpath-plots-0",
				},
			},
		},
		{
			name: "Without Plot Storage",
			harvester: k8schianetv1.ChiaHarvester{
				Spec: k8schianetv1.ChiaHarvesterSpec{
					ChiaConfig: k8schianetv1.ChiaHarvesterSpecChia{
						CASecretName: "test-ca-secret",
					},
				},
			},
			want: []corev1.VolumeMount{
				{
					Name:      "secret-ca",
					MountPath: "/chia-ca",
				},
				{
					Name:      "chiaroot",
					MountPath: "/chia-data",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getChiaVolumeMounts(tt.harvester)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetChiaVolumes(t *testing.T) {
	tests := []struct {
		name      string
		harvester k8schianetv1.ChiaHarvester
		want      []corev1.Volume
	}{
		{
			name: "With Generated ChiaRoot Storage and Specified Plots",
			harvester: k8schianetv1.ChiaHarvester{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: k8schianetv1.ChiaHarvesterSpec{
					ChiaConfig: k8schianetv1.ChiaHarvesterSpecChia{
						CASecretName: "test-ca-secret",
					},
					CommonSpec: k8schianetv1.CommonSpec{
						Storage: &k8schianetv1.StorageConfig{
							ChiaRoot: &k8schianetv1.ChiaRootConfig{
								PersistentVolumeClaim: &k8schianetv1.PersistentVolumeClaimConfig{
									GenerateVolumeClaims: true,
									StorageClass:         "standard",
									ResourceRequest:      "10Gi",
									AccessModes:          []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
								},
							},
							Plots: &k8schianetv1.PlotsConfig{
								PersistentVolumeClaim: []*k8schianetv1.PersistentVolumeClaimConfig{
									{
										ClaimName: "plot-pvc-1",
									},
								},
							},
						},
					},
				},
			},
			want: []corev1.Volume{
				{
					Name: "secret-ca",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-ca-secret",
						},
					},
				},
				{
					Name: "chiaroot",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "test-harvester",
						},
					},
				},
				{
					Name: "pvc-plots-0",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "plot-pvc-1",
						},
					},
				},
			},
		},
		{
			name: "With Specified ChiaRoot Storage and Plots",
			harvester: k8schianetv1.ChiaHarvester{
				Spec: k8schianetv1.ChiaHarvesterSpec{
					ChiaConfig: k8schianetv1.ChiaHarvesterSpecChia{
						CASecretName: "test-ca-secret",
					},
					CommonSpec: k8schianetv1.CommonSpec{
						Storage: &k8schianetv1.StorageConfig{
							ChiaRoot: &k8schianetv1.ChiaRootConfig{
								PersistentVolumeClaim: &k8schianetv1.PersistentVolumeClaimConfig{
									ClaimName: "chiaroot-pvc",
								},
							},
							Plots: &k8schianetv1.PlotsConfig{
								PersistentVolumeClaim: []*k8schianetv1.PersistentVolumeClaimConfig{
									{
										ClaimName: "plot-pvc-1",
									},
								},
							},
						},
					},
				},
			},
			want: []corev1.Volume{
				{
					Name: "secret-ca",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-ca-secret",
						},
					},
				},
				{
					Name: "chiaroot",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "chiaroot-pvc",
						},
					},
				},
				{
					Name: "pvc-plots-0",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "plot-pvc-1",
						},
					},
				},
			},
		},
		{
			name: "With HostPath Storage",
			harvester: k8schianetv1.ChiaHarvester{
				Spec: k8schianetv1.ChiaHarvesterSpec{
					ChiaConfig: k8schianetv1.ChiaHarvesterSpecChia{
						CASecretName: "test-ca-secret",
					},
					CommonSpec: k8schianetv1.CommonSpec{
						Storage: &k8schianetv1.StorageConfig{
							ChiaRoot: &k8schianetv1.ChiaRootConfig{
								HostPathVolume: &k8schianetv1.HostPathVolumeConfig{
									Path: "/chia-root",
								},
							},
							Plots: &k8schianetv1.PlotsConfig{
								HostPathVolume: []*k8schianetv1.HostPathVolumeConfig{
									{
										Path: "/plots/hostpath-1",
									},
								},
							},
						},
					},
				},
			},
			want: []corev1.Volume{
				{
					Name: "secret-ca",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-ca-secret",
						},
					},
				},
				{
					Name: "chiaroot",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/chia-root",
						},
					},
				},
				{
					Name: "hostpath-plots-0",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/plots/hostpath-1",
						},
					},
				},
			},
		},
		{
			name: "Without Storage Config",
			harvester: k8schianetv1.ChiaHarvester{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: k8schianetv1.ChiaHarvesterSpec{
					ChiaConfig: k8schianetv1.ChiaHarvesterSpecChia{
						CASecretName: "test-ca-secret",
					},
				},
			},
			want: []corev1.Volume{
				{
					Name: "secret-ca",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-ca-secret",
						},
					},
				},
				{
					Name: "chiaroot",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getChiaVolumes(tt.harvester)
			assert.Equal(t, tt.want, got)
		})
	}
}
