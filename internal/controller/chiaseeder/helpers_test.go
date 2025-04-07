package chiaseeder

import (
	"testing"

	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
)

func TestGetChiaPorts(t *testing.T) {
	testCases := []struct {
		name          string
		fullNodePort  int32
		expectedPorts []struct {
			name          string
			containerPort int32
			protocol      string
		}
	}{
		{
			name:         "Mainnet Port",
			fullNodePort: consts.MainnetNodePort,
			expectedPorts: []struct {
				name          string
				containerPort int32
				protocol      string
			}{
				{"daemon", consts.DaemonPort, "TCP"},
				{"dns", 53, "UDP"},
				{"dns-tcp", 53, "TCP"},
				{"peers", consts.MainnetNodePort, "TCP"},
				{"rpc", consts.CrawlerRPCPort, "TCP"},
			},
		},
		{
			name:         "Testnet Port",
			fullNodePort: consts.TestnetNodePort,
			expectedPorts: []struct {
				name          string
				containerPort int32
				protocol      string
			}{
				{"daemon", consts.DaemonPort, "TCP"},
				{"dns", 53, "UDP"},
				{"dns-tcp", 53, "TCP"},
				{"peers", consts.TestnetNodePort, "TCP"},
				{"rpc", consts.CrawlerRPCPort, "TCP"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ports := getChiaPorts(tc.fullNodePort)

			assert.Len(t, ports, len(tc.expectedPorts), "Expected %d ports", len(tc.expectedPorts))
			for i, expected := range tc.expectedPorts {
				assert.Equal(t, expected.name, ports[i].Name, "Port name should match")
				assert.Equal(t, expected.containerPort, ports[i].ContainerPort, "Container port should match")
				assert.Equal(t, expected.protocol, string(ports[i].Protocol), "Protocol should match")
			}
		})
	}
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}

func TestGetChiaVolumeMounts(t *testing.T) {
	testCases := []struct {
		name           string
		seeder         k8schianetv1.ChiaSeeder
		expectedMounts []corev1.VolumeMount
	}{
		{
			name: "With CA Secret",
			seeder: k8schianetv1.ChiaSeeder{
				Spec: k8schianetv1.ChiaSeederSpec{
					ChiaConfig: k8schianetv1.ChiaSeederSpecChia{
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
					Name:      "chiaroot",
					MountPath: "/chia-data",
				},
			},
		},
		{
			name: "Without CA Secret",
			seeder: k8schianetv1.ChiaSeeder{
				Spec: k8schianetv1.ChiaSeederSpec{
					ChiaConfig: k8schianetv1.ChiaSeederSpecChia{
						CASecretName: nil,
					},
				},
			},
			expectedMounts: []corev1.VolumeMount{
				{
					Name:      "chiaroot",
					MountPath: "/chia-data",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			volumeMounts := getChiaVolumeMounts(tc.seeder)

			assert.Equal(t, len(tc.expectedMounts), len(volumeMounts), "Number of volume mounts should match")
			for i, expectedMount := range tc.expectedMounts {
				assert.Equal(t, expectedMount.Name, volumeMounts[i].Name, "Volume mount name should match")
				assert.Equal(t, expectedMount.MountPath, volumeMounts[i].MountPath, "Mount path should match")
			}
		})
	}
}

func TestGetChiaVolumes(t *testing.T) {
	testCases := []struct {
		name            string
		seeder          k8schianetv1.ChiaSeeder
		expectedVolumes []corev1.Volume
	}{
		{
			name: "With CA Secret and Generated ChiaRoot",
			seeder: k8schianetv1.ChiaSeeder{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: k8schianetv1.ChiaSeederSpec{
					ChiaConfig: k8schianetv1.ChiaSeederSpecChia{
						CASecretName: stringPtr("test-ca-secret"),
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
					Name: "chiaroot",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "test-seeder",
						},
					},
				},
			},
		},
		{
			name: "With CA Secret and Specified ChiaRoot",
			seeder: k8schianetv1.ChiaSeeder{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: k8schianetv1.ChiaSeederSpec{
					ChiaConfig: k8schianetv1.ChiaSeederSpecChia{
						CASecretName: stringPtr("test-ca-secret"),
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
			name: "With CA Secret and HostPath Storage",
			seeder: k8schianetv1.ChiaSeeder{
				Spec: k8schianetv1.ChiaSeederSpec{
					ChiaConfig: k8schianetv1.ChiaSeederSpecChia{
						CASecretName: stringPtr("test-ca-secret"),
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
		},
		{
			name: "Without CA Secret and Storage Config",
			seeder: k8schianetv1.ChiaSeeder{
				Spec: k8schianetv1.ChiaSeederSpec{
					ChiaConfig: k8schianetv1.ChiaSeederSpecChia{
						CASecretName: nil,
					},
				},
			},
			expectedVolumes: []corev1.Volume{
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
			volumes := getChiaVolumes(tc.seeder)

			assert.Equal(t, len(tc.expectedVolumes), len(volumes), "Number of volumes should match")
			for i, expectedVolume := range tc.expectedVolumes {
				assert.Equal(t, expectedVolume.Name, volumes[i].Name, "Volume name should match")
				assert.Equal(t, expectedVolume.VolumeSource, volumes[i].VolumeSource, "Volume source should match")
			}
		})
	}
}
