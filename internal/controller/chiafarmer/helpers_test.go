package chiafarmer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
)

func stringPtr(s string) *string {
	return &s
}

func TestGetChiaVolumeMounts(t *testing.T) {
	volumeMounts := getChiaVolumeMounts()

	assert.Len(t, volumeMounts, 3, "Expected 3 volume mounts")

	// Check each volume mount
	expectedVolumeMounts := []struct {
		name      string
		mountPath string
	}{
		{"secret-ca", "/chia-ca"},
		{"key", "/key"},
		{"chiaroot", "/chia-data"},
	}

	for i, expected := range expectedVolumeMounts {
		assert.Equal(t, expected.name, volumeMounts[i].Name, "Volume mount name should match")
		assert.Equal(t, expected.mountPath, volumeMounts[i].MountPath, "Mount path should match")
	}
}

func TestGetChiaVolumes(t *testing.T) {
	testCases := []struct {
		name            string
		farmer          k8schianetv1.ChiaFarmer
		expectedVolumes []corev1.Volume
	}{
		{
			name: "With Generated ChiaRoot",
			farmer: k8schianetv1.ChiaFarmer{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: k8schianetv1.ChiaFarmerSpec{
					ChiaConfig: k8schianetv1.ChiaFarmerSpecChia{
						CASecretName: "test-ca-secret",
						SecretKey: k8schianetv1.ChiaSecretKey{
							Name: "test-key-secret",
							Key:  "test-key",
						},
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
						},
					},
				},
			},
			expectedVolumes: []corev1.Volume{
				{
					Name: "secret-ca",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-ca-secret",
						},
					},
				},
				{
					Name: "key",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-key-secret",
						},
					},
				},
				{
					Name: "chiaroot",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "test-farmer",
						},
					},
				},
			},
		},
		{
			name: "With Specified ChiaRoot Storage",
			farmer: k8schianetv1.ChiaFarmer{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: k8schianetv1.ChiaFarmerSpec{
					ChiaConfig: k8schianetv1.ChiaFarmerSpecChia{
						CASecretName: "test-ca-secret",
						SecretKey: k8schianetv1.ChiaSecretKey{
							Name: "test-key-secret",
							Key:  "test-key",
						},
					},
					CommonSpec: k8schianetv1.CommonSpec{
						Storage: &k8schianetv1.StorageConfig{
							ChiaRoot: &k8schianetv1.ChiaRootConfig{
								PersistentVolumeClaim: &k8schianetv1.PersistentVolumeClaimConfig{
									ClaimName: "specified-chiaroot",
								},
							},
						},
					},
				},
			},
			expectedVolumes: []corev1.Volume{
				{
					Name: "secret-ca",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-ca-secret",
						},
					},
				},
				{
					Name: "key",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-key-secret",
						},
					},
				},
				{
					Name: "chiaroot",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "specified-chiaroot",
						},
					},
				},
			},
		},
		{
			name: "With HostPath Storage",
			farmer: k8schianetv1.ChiaFarmer{
				Spec: k8schianetv1.ChiaFarmerSpec{
					ChiaConfig: k8schianetv1.ChiaFarmerSpecChia{
						CASecretName: "test-ca-secret",
						SecretKey: k8schianetv1.ChiaSecretKey{
							Name: "test-key-secret",
							Key:  "test-key",
						},
					},
					CommonSpec: k8schianetv1.CommonSpec{
						Storage: &k8schianetv1.StorageConfig{
							ChiaRoot: &k8schianetv1.ChiaRootConfig{
								HostPathVolume: &k8schianetv1.HostPathVolumeConfig{
									Path: "/test/path",
								},
							},
						},
					},
				},
			},
			expectedVolumes: []corev1.Volume{
				{
					Name: "secret-ca",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-ca-secret",
						},
					},
				},
				{
					Name: "key",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-key-secret",
						},
					},
				},
				{
					Name: "chiaroot",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/test/path",
						},
					},
				},
			},
		},
		{
			name: "Without Storage Config",
			farmer: k8schianetv1.ChiaFarmer{
				Spec: k8schianetv1.ChiaFarmerSpec{
					ChiaConfig: k8schianetv1.ChiaFarmerSpecChia{
						CASecretName: "test-ca-secret",
						SecretKey: k8schianetv1.ChiaSecretKey{
							Name: "test-key-secret",
							Key:  "test-key",
						},
					},
				},
			},
			expectedVolumes: []corev1.Volume{
				{
					Name: "secret-ca",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-ca-secret",
						},
					},
				},
				{
					Name: "key",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-key-secret",
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			volumes := getChiaVolumes(tc.farmer)

			assert.Equal(t, len(tc.expectedVolumes), len(volumes), "Number of volumes should match")
			for i, expectedVolume := range tc.expectedVolumes {
				assert.Equal(t, expectedVolume.Name, volumes[i].Name, "Volume name should match")
				assert.Equal(t, expectedVolume.VolumeSource, volumes[i].VolumeSource, "Volume source should match")
			}
		})
	}
}
