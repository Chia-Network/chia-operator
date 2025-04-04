package chiadatalayer

import (
	"testing"

	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
)

func TestGetChiaPorts(t *testing.T) {
	ports := getChiaPorts()
	assert.Len(t, ports, 3, "Expected 3 ports")
	expectedPorts := []struct {
		name          string
		containerPort int32
		protocol      string
	}{
		{"daemon", consts.DaemonPort, "TCP"},
		{"rpc", consts.DataLayerRPCPort, "TCP"},
		{"wallet-rpc", consts.WalletRPCPort, "TCP"},
	}

	for i, expected := range expectedPorts {
		assert.Equal(t, expected.name, ports[i].Name, "Port name should match")
		assert.Equal(t, expected.containerPort, ports[i].ContainerPort, "Container port should match")
		assert.Equal(t, expected.protocol, string(ports[i].Protocol), "Protocol should match")
	}
}

func TestGetChiaVolumes(t *testing.T) {
	testCases := []struct {
		name            string
		datalayer       k8schianetv1.ChiaDataLayer
		expectedVolumes []struct {
			name         string
			volumeSource corev1.VolumeSource
		}
	}{
		{
			name: "With CA Secret and Generated PVC Storage",
			datalayer: k8schianetv1.ChiaDataLayer{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: k8schianetv1.ChiaDataLayerSpec{
					ChiaConfig: k8schianetv1.ChiaDataLayerSpecChia{
						CASecretName: stringPtr("test-ca-secret"),
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
							DataLayerServerFiles: &k8schianetv1.DataLayerServerFilesConfig{
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
					name: "key",
					volumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-key-secret",
						},
					},
				},
				{
					name: "chiaroot",
					volumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "test-datalayer",
						},
					},
				},
				{
					name: "server",
					volumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "test-datalayer-server",
						},
					},
				},
			},
		},
		{
			name: "With Specified PVC Storage",
			datalayer: k8schianetv1.ChiaDataLayer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-datalayer",
					Namespace: "test-namespace",
				},
				Spec: k8schianetv1.ChiaDataLayerSpec{
					ChiaConfig: k8schianetv1.ChiaDataLayerSpecChia{
						CASecretName: nil,
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
							DataLayerServerFiles: &k8schianetv1.DataLayerServerFilesConfig{
								PersistentVolumeClaim: &k8schianetv1.PersistentVolumeClaimConfig{
									ClaimName: "specified-server",
								},
							},
						},
					},
				},
			},
			expectedVolumes: []struct {
				name         string
				volumeSource corev1.VolumeSource
			}{
				{
					name: "key",
					volumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-key-secret",
						},
					},
				},
				{
					name: "chiaroot",
					volumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "specified-chiaroot",
						},
					},
				},
				{
					name: "server",
					volumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "specified-server",
						},
					},
				},
			},
		},
		{
			name: "Without CA Secret and Storage",
			datalayer: k8schianetv1.ChiaDataLayer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-datalayer",
					Namespace: "test-namespace",
				},
				Spec: k8schianetv1.ChiaDataLayerSpec{
					ChiaConfig: k8schianetv1.ChiaDataLayerSpecChia{
						CASecretName: nil,
						SecretKey: k8schianetv1.ChiaSecretKey{
							Name: "test-key-secret",
							Key:  "test-key",
						},
					},
				},
			},
			expectedVolumes: []struct {
				name         string
				volumeSource corev1.VolumeSource
			}{
				{
					name: "key",
					volumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-key-secret",
						},
					},
				},
				{
					name: "chiaroot",
					volumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
				{
					name: "server",
					volumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			volumes := getChiaVolumes(tc.datalayer)

			assert.Len(t, volumes, len(tc.expectedVolumes), "Expected %d volumes", len(tc.expectedVolumes))
			for i, expected := range tc.expectedVolumes {
				assert.Equal(t, expected.name, volumes[i].Name, "Volume name should match")
				assert.Equal(t, expected.volumeSource, volumes[i].VolumeSource, "Volume source should match")
			}
		})
	}
}

func TestGetChiaVolumeMounts(t *testing.T) {
	testCases := []struct {
		name           string
		datalayer      k8schianetv1.ChiaDataLayer
		expectedMounts []corev1.VolumeMount
	}{
		{
			name: "With CA Secret",
			datalayer: k8schianetv1.ChiaDataLayer{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-datalayer",
				},
				Spec: k8schianetv1.ChiaDataLayerSpec{
					ChiaConfig: k8schianetv1.ChiaDataLayerSpecChia{
						CASecretName: stringPtr("test-ca-secret"),
					},
				},
			},
			expectedMounts: []corev1.VolumeMount{
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
				{
					Name:      "server",
					MountPath: "/datalayer/server",
				},
			},
		},
		{
			name: "Without CA Secret",
			datalayer: k8schianetv1.ChiaDataLayer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-datalayer",
					Namespace: "test-namespace",
				},
				Spec: k8schianetv1.ChiaDataLayerSpec{
					ChiaConfig: k8schianetv1.ChiaDataLayerSpecChia{
						CASecretName: nil,
					},
				},
			},
			expectedMounts: []corev1.VolumeMount{
				{
					Name:      "key",
					MountPath: "/key",
				},
				{
					Name:      "chiaroot",
					MountPath: "/chia-data",
				},
				{
					Name:      "server",
					MountPath: "/datalayer/server",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			volumeMounts := getChiaVolumeMounts(tc.datalayer)

			assert.Len(t, volumeMounts, len(tc.expectedMounts), "Expected %d volume mounts", len(tc.expectedMounts))

			// Check each volume mount
			for i, expected := range tc.expectedMounts {
				assert.Equal(t, expected.Name, volumeMounts[i].Name, "Volume mount name should match")
				assert.Equal(t, expected.MountPath, volumeMounts[i].MountPath, "Mount path should match")
			}
		})
	}
}

func TestGetExistingChiaDatalayerServerVolume(t *testing.T) {

	testCases := []struct {
		name           string
		storage        *k8schianetv1.StorageConfig
		expectedVolume struct {
			name         string
			volumeSource corev1.VolumeSource
		}
	}{
		{
			name: "With PVC",
			storage: &k8schianetv1.StorageConfig{
				DataLayerServerFiles: &k8schianetv1.DataLayerServerFilesConfig{
					PersistentVolumeClaim: &k8schianetv1.PersistentVolumeClaimConfig{
						ClaimName: "test-pvc",
					},
				},
			},
			expectedVolume: struct {
				name         string
				volumeSource corev1.VolumeSource
			}{
				name: "server",
				volumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: "test-pvc",
					},
				},
			},
		},
		{
			name: "With HostPath",
			storage: &k8schianetv1.StorageConfig{
				DataLayerServerFiles: &k8schianetv1.DataLayerServerFilesConfig{
					HostPathVolume: &k8schianetv1.HostPathVolumeConfig{
						Path: "/test/path",
					},
				},
			},
			expectedVolume: struct {
				name         string
				volumeSource corev1.VolumeSource
			}{
				name: "server",
				volumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/test/path",
					},
				},
			},
		},
		{
			name:    "Without Storage Config",
			storage: nil,
			expectedVolume: struct {
				name         string
				volumeSource corev1.VolumeSource
			}{
				name: "server",
				volumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		},
		{
			name: "With Empty Storage Config",
			storage: &k8schianetv1.StorageConfig{
				DataLayerServerFiles: nil,
			},
			expectedVolume: struct {
				name         string
				volumeSource corev1.VolumeSource
			}{
				name: "server",
				volumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			volume := getExistingChiaDatalayerServerVolume(tc.storage)

			assert.Equal(t, tc.expectedVolume.name, volume.Name, "Volume name should match")
			assert.Equal(t, tc.expectedVolume.volumeSource, volume.VolumeSource, "Volume source should match")
		})
	}
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
