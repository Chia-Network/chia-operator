package chianode

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
)

func TestGetChiaVolumeMounts(t *testing.T) {
	volumeMounts := getChiaVolumeMounts()

	assert.Len(t, volumeMounts, 2, "Expected 2 volume mounts")

	expectedVolumeMounts := []struct {
		name      string
		mountPath string
	}{
		{"secret-ca", "/chia-ca"},
		{"chiaroot", "/chia-data"},
	}

	for i, expected := range expectedVolumeMounts {
		assert.Equal(t, expected.name, volumeMounts[i].Name, "Volume mount name should match")
		assert.Equal(t, expected.mountPath, volumeMounts[i].MountPath, "Mount path should match")
	}
}

func TestGetChiaRootVolume(t *testing.T) {
	testCases := []struct {
		name           string
		storage        *k8schianetv1.StorageConfig
		expectedVolume *corev1.Volume
		expectedPVC    *corev1.PersistentVolumeClaim
	}{
		{
			name: "With PVC",
			storage: &k8schianetv1.StorageConfig{
				ChiaRoot: &k8schianetv1.ChiaRootConfig{
					PersistentVolumeClaim: &k8schianetv1.PersistentVolumeClaimConfig{
						StorageClass:    "standard",
						ResourceRequest: "10Gi",
						AccessModes:     []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					},
				},
			},
			expectedVolume: nil,
			expectedPVC: &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: "chiaroot",
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					StorageClassName: stringPtr("standard"),
					Resources: corev1.VolumeResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceStorage: resource.MustParse("10Gi"),
						},
					},
				},
			},
		},
		{
			name: "With HostPath",
			storage: &k8schianetv1.StorageConfig{
				ChiaRoot: &k8schianetv1.ChiaRootConfig{
					HostPathVolume: &k8schianetv1.HostPathVolumeConfig{
						Path: "/test/path",
					},
				},
			},
			expectedVolume: &corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/test/path",
					},
				},
			},
			expectedPVC: nil,
		},
		{
			name:    "Without Storage Config",
			storage: nil,
			expectedVolume: &corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
			expectedPVC: nil,
		},
		{
			name: "With Empty Storage Config",
			storage: &k8schianetv1.StorageConfig{
				ChiaRoot: nil,
			},
			expectedVolume: &corev1.Volume{
				Name: "chiaroot",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
			expectedPVC: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			volume, pvc := getChiaRootVolume(tc.storage)

			if tc.expectedVolume != nil {
				assert.NotNil(t, volume, "Volume should not be nil")
				assert.Equal(t, tc.expectedVolume.Name, volume.Name, "Volume name should match")
				assert.Equal(t, tc.expectedVolume.VolumeSource, volume.VolumeSource, "Volume source should match")
			} else {
				assert.Nil(t, volume, "Volume should be nil")
			}

			if tc.expectedPVC != nil {
				assert.NotNil(t, pvc, "PVC should not be nil")
				assert.Equal(t, tc.expectedPVC.Name, pvc.Name, "PVC name should match")
				assert.Equal(t, tc.expectedPVC.Spec, pvc.Spec, "PVC spec should match")
			} else {
				assert.Nil(t, pvc, "PVC should be nil")
			}
		})
	}
}

func TestGetChiaVolumesAndTemplates(t *testing.T) {
	testCases := []struct {
		name                    string
		node                    k8schianetv1.ChiaNode
		expectedVolumes         []corev1.Volume
		expectedVolumeTemplates []corev1.PersistentVolumeClaim
	}{
		{
			name: "With PVC Storage",
			node: k8schianetv1.ChiaNode{
				Spec: k8schianetv1.ChiaNodeSpec{
					ChiaConfig: k8schianetv1.ChiaNodeSpecChia{
						CASecretName: "test-ca-secret",
					},
					CommonSpec: k8schianetv1.CommonSpec{
						Storage: &k8schianetv1.StorageConfig{
							ChiaRoot: &k8schianetv1.ChiaRootConfig{
								PersistentVolumeClaim: &k8schianetv1.PersistentVolumeClaimConfig{
									StorageClass:    "standard",
									ResourceRequest: "10Gi",
									AccessModes:     []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
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
			},
			expectedVolumeTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "chiaroot",
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
						StorageClassName: stringPtr("standard"),
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse("10Gi"),
							},
						},
					},
				},
			},
		},
		{
			name: "With HostPath Storage",
			node: k8schianetv1.ChiaNode{
				Spec: k8schianetv1.ChiaNodeSpec{
					ChiaConfig: k8schianetv1.ChiaNodeSpecChia{
						CASecretName: "test-ca-secret",
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
					Name: "chiaroot",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/test/path",
						},
					},
				},
			},
			expectedVolumeTemplates: nil,
		},
		{
			name: "Without Storage Config",
			node: k8schianetv1.ChiaNode{
				Spec: k8schianetv1.ChiaNodeSpec{
					ChiaConfig: k8schianetv1.ChiaNodeSpecChia{
						CASecretName: "test-ca-secret",
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
					Name: "chiaroot",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
			expectedVolumeTemplates: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			volumes, volumeTemplates := getChiaVolumesAndTemplates(tc.node)

			assert.Equal(t, len(tc.expectedVolumes), len(volumes), "Number of volumes should match")
			for i, expectedVolume := range tc.expectedVolumes {
				assert.Equal(t, expectedVolume.Name, volumes[i].Name, "Volume name should match")
				assert.Equal(t, expectedVolume.VolumeSource, volumes[i].VolumeSource, "Volume source should match")
			}

			if tc.expectedVolumeTemplates == nil {
				assert.Empty(t, volumeTemplates, "Volume templates should be empty")
			} else {
				assert.Equal(t, len(tc.expectedVolumeTemplates), len(volumeTemplates), "Number of volume templates should match")
				for i, expectedTemplate := range tc.expectedVolumeTemplates {
					assert.Equal(t, expectedTemplate.Name, volumeTemplates[i].Name, "Template name should match")
					assert.Equal(t, expectedTemplate.Spec, volumeTemplates[i].Spec, "Template spec should match")
				}
			}
		})
	}
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
