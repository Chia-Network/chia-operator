package chiaintroducer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
)

func TestGetChiaVolumeMounts(t *testing.T) {
	testCases := []struct {
		name           string
		introducer     k8schianetv1.ChiaIntroducer
		expectedMounts []struct {
			name      string
			mountPath string
		}
	}{
		{
			name: "With CA Secret",
			introducer: k8schianetv1.ChiaIntroducer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-introducer",
					Namespace: "test-namespace",
				},
				Spec: k8schianetv1.ChiaIntroducerSpec{
					ChiaConfig: k8schianetv1.ChiaIntroducerSpecChia{
						CASecretName: stringPtr("test-ca-secret"),
					},
				},
			},
			expectedMounts: []struct {
				name      string
				mountPath string
			}{
				{"secret-ca", "/chia-ca"},
				{"chiaroot", "/chia-data"},
			},
		},
		{
			name: "Without CA Secret",
			introducer: k8schianetv1.ChiaIntroducer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-introducer",
					Namespace: "test-namespace",
				},
				Spec: k8schianetv1.ChiaIntroducerSpec{
					ChiaConfig: k8schianetv1.ChiaIntroducerSpecChia{
						CASecretName: nil,
					},
				},
			},
			expectedMounts: []struct {
				name      string
				mountPath string
			}{
				{"chiaroot", "/chia-data"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			volumeMounts := getChiaVolumeMounts(tc.introducer)

			assert.Len(t, volumeMounts, len(tc.expectedMounts), "Expected %d volume mounts", len(tc.expectedMounts))

			// Check each volume mount
			for i, expected := range tc.expectedMounts {
				assert.Equal(t, expected.name, volumeMounts[i].Name, "Volume mount name should match")
				assert.Equal(t, expected.mountPath, volumeMounts[i].MountPath, "Mount path should match")
			}
		})
	}
}

func TestGetChiaVolumes(t *testing.T) {
	testCases := []struct {
		name            string
		introducer      k8schianetv1.ChiaIntroducer
		expectedVolumes []struct {
			name         string
			volumeSource corev1.VolumeSource
		}
	}{
		{
			name: "With CA Secret and Generated PVC",
			introducer: k8schianetv1.ChiaIntroducer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "test",
				},
				Spec: k8schianetv1.ChiaIntroducerSpec{
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
					ChiaConfig: k8schianetv1.ChiaIntroducerSpecChia{
						CASecretName: stringPtr("test-ca-secret"),
					},
				},
			},
			expectedVolumes: []struct {
				name         string
				volumeSource corev1.VolumeSource
			}{
				{
					name: "secret-ca",
					volumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-ca-secret",
						},
					},
				},
				{
					name: "chiaroot",
					volumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "test-introducer",
						},
					},
				},
			},
		},
		{
			name: "With CA Secret and Specified PVC",
			introducer: k8schianetv1.ChiaIntroducer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-introducer",
					Namespace: "test-namespace",
				},
				Spec: k8schianetv1.ChiaIntroducerSpec{
					CommonSpec: k8schianetv1.CommonSpec{
						Storage: &k8schianetv1.StorageConfig{
							ChiaRoot: &k8schianetv1.ChiaRootConfig{
								PersistentVolumeClaim: &k8schianetv1.PersistentVolumeClaimConfig{
									ClaimName: "test-pvc",
								},
							},
						},
					},
					ChiaConfig: k8schianetv1.ChiaIntroducerSpecChia{
						CASecretName: stringPtr("test-ca-secret"),
					},
				},
			},
			expectedVolumes: []struct {
				name         string
				volumeSource corev1.VolumeSource
			}{
				{
					name: "secret-ca",
					volumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-ca-secret",
						},
					},
				},
				{
					name: "chiaroot",
					volumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "test-pvc",
						},
					},
				},
			},
		},
		{
			name: "Without CA Secret and Storage",
			introducer: k8schianetv1.ChiaIntroducer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-introducer",
					Namespace: "test-namespace",
				},
				Spec: k8schianetv1.ChiaIntroducerSpec{
					ChiaConfig: k8schianetv1.ChiaIntroducerSpecChia{
						CASecretName: nil,
					},
				},
			},
			expectedVolumes: []struct {
				name         string
				volumeSource corev1.VolumeSource
			}{
				{
					name: "chiaroot",
					volumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			volumes := getChiaVolumes(tc.introducer)

			assert.Len(t, volumes, len(tc.expectedVolumes), "Expected %d volumes", len(tc.expectedVolumes))

			for i, expected := range tc.expectedVolumes {
				assert.Equal(t, expected.name, volumes[i].Name, "Volume name should match")
				assert.Equal(t, expected.volumeSource, volumes[i].VolumeSource, "Volume source should match")
			}
		})
	}
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
