package chiawallet

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

	// Check each port
	expectedPorts := []struct {
		name          string
		containerPort int32
		protocol      string
	}{
		{"daemon", consts.DaemonPort, "TCP"},
		{"peers", consts.WalletPort, "TCP"},
		{"rpc", consts.WalletRPCPort, "TCP"},
	}

	for i, expected := range expectedPorts {
		assert.Equal(t, expected.name, ports[i].Name, "Port name should match")
		assert.Equal(t, expected.containerPort, ports[i].ContainerPort, "Container port should match")
		assert.Equal(t, expected.protocol, string(ports[i].Protocol), "Protocol should match")
	}
}

func TestGetChiaVolumeMounts(t *testing.T) {
	testCases := []struct {
		name           string
		wallet         k8schianetv1.ChiaWallet
		expectedMounts []corev1.VolumeMount
	}{
		{
			name: "With CA Secret",
			wallet: k8schianetv1.ChiaWallet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-wallet",
				},
				Spec: k8schianetv1.ChiaWalletSpec{
					ChiaConfig: k8schianetv1.ChiaWalletSpecChia{
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
			},
		},
		{
			name: "Without CA Secret",
			wallet: k8schianetv1.ChiaWallet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-wallet",
					Namespace: "test-namespace",
				},
				Spec: k8schianetv1.ChiaWalletSpec{
					ChiaConfig: k8schianetv1.ChiaWalletSpecChia{
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
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			volumeMounts := getChiaVolumeMounts(tc.wallet)

			assert.Len(t, volumeMounts, len(tc.expectedMounts), "Expected %d volume mounts", len(tc.expectedMounts))

			// Check each volume mount
			for i, expected := range tc.expectedMounts {
				assert.Equal(t, expected.Name, volumeMounts[i].Name, "Volume mount name should match")
				assert.Equal(t, expected.MountPath, volumeMounts[i].MountPath, "Mount path should match")
			}
		})
	}
}

func TestGetChiaVolumes(t *testing.T) {
	testCases := []struct {
		name            string
		wallet          k8schianetv1.ChiaWallet
		expectedVolumes []corev1.Volume
	}{
		{
			name: "With Generated ChiaRoot",
			wallet: k8schianetv1.ChiaWallet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: k8schianetv1.ChiaWalletSpec{
					ChiaConfig: k8schianetv1.ChiaWalletSpecChia{
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
							ClaimName: "test-wallet",
						},
					},
				},
			},
		},
		{
			name: "With Specified ChiaRoot Storage",
			wallet: k8schianetv1.ChiaWallet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: k8schianetv1.ChiaWalletSpec{
					ChiaConfig: k8schianetv1.ChiaWalletSpecChia{
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
			wallet: k8schianetv1.ChiaWallet{
				Spec: k8schianetv1.ChiaWalletSpec{
					ChiaConfig: k8schianetv1.ChiaWalletSpecChia{
						CASecretName: stringPtr("test-ca-secret"),
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
			wallet: k8schianetv1.ChiaWallet{
				Spec: k8schianetv1.ChiaWalletSpec{
					ChiaConfig: k8schianetv1.ChiaWalletSpecChia{
						CASecretName: stringPtr("test-ca-secret"),
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
			volumes := getChiaVolumes(tc.wallet)

			assert.Equal(t, len(tc.expectedVolumes), len(volumes), "Number of volumes should match")
			for i, expectedVolume := range tc.expectedVolumes {
				assert.Equal(t, expectedVolume.Name, volumes[i].Name, "Volume name should match")
				assert.Equal(t, expectedVolume.VolumeSource, volumes[i].VolumeSource, "Volume source should match")
			}
		})
	}
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}
