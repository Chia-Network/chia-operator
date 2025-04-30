package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}

// Helper function to create a bool pointer
func boolPtr(b bool) *bool {
	return &b
}

// Helper function to create a map pointer
func mapPtr(m map[string]string) *map[string]string {
	return &m
}

func TestGetCommonLabels(t *testing.T) {
	expected := map[string]string{
		"app.kubernetes.io/instance":   "testname",
		"app.kubernetes.io/name":       "testname",
		"app.kubernetes.io/managed-by": "chia-operator",
		"k8s.chia.net/kind":            "TestKind",
		"foo":                          "bar",
		"hello":                        "world",
	}
	actual := GetCommonLabels("TestKind",
		metav1.ObjectMeta{
			Name:      "testname",
			Namespace: "testnamespace",
		},
		map[string]string{
			"foo": "bar",
		},
		map[string]string{
			"hello": "world",
		})
	require.Equal(t, expected, actual)
}

func TestCombineMaps(t *testing.T) {
	expected := map[string]string{
		"foo":   "bar",
		"hello": "world",
	}
	actual := CombineMaps(
		map[string]string{
			"foo": "bar",
		},
		map[string]string{
			"hello": "world",
		})
	require.Equal(t, expected, actual)
}

func TestShouldMakeChiaRootVolumeClaim(t *testing.T) {
	// True case
	actual := ShouldMakeChiaRootVolumeClaim(&k8schianetv1.StorageConfig{
		ChiaRoot: &k8schianetv1.ChiaRootConfig{
			PersistentVolumeClaim: &k8schianetv1.PersistentVolumeClaimConfig{
				GenerateVolumeClaims: true,
			},
		},
	})
	require.Equal(t, true, actual, "expected should make volume claim")

	// False case - nil storage config
	actual = ShouldMakeChiaRootVolumeClaim(nil)
	require.Equal(t, false, actual, "expected should not make volume claim for nil storage config")

	// False case - non-nil storage config, nil ChiaRoot config
	actual = ShouldMakeChiaRootVolumeClaim(&k8schianetv1.StorageConfig{
		ChiaRoot: nil,
	})
	require.Equal(t, false, actual, "expected should not make volume claim for nil ChiaRoot config")

	// False case - non-nil storage config, nil PersistentVolumeClaim config
	actual = ShouldMakeChiaRootVolumeClaim(&k8schianetv1.StorageConfig{
		ChiaRoot: &k8schianetv1.ChiaRootConfig{
			PersistentVolumeClaim: nil,
		},
	})
	require.Equal(t, false, actual, "expected should not make volume claim for nil PersistentVolumeClaim config")

	// False case - non-nil storage config, GenerateVolumeClaims set to false
	actual = ShouldMakeChiaRootVolumeClaim(&k8schianetv1.StorageConfig{
		ChiaRoot: &k8schianetv1.ChiaRootConfig{
			PersistentVolumeClaim: &k8schianetv1.PersistentVolumeClaimConfig{
				GenerateVolumeClaims: false,
			},
		},
	})
	require.Equal(t, false, actual, "expected should not make volume claim for false GenerateVolumeClaims config")
}

func TestShouldMakeService(t *testing.T) {
	// True case - default true
	actual := ShouldMakeService(k8schianetv1.Service{
		Enabled: nil,
	}, true)
	require.Equal(t, true, actual, "expected should make Service, defaulted to true")

	// True case - default false with enabled set to true
	enabled := true
	actual = ShouldMakeService(k8schianetv1.Service{
		Enabled: &enabled,
	}, false)
	require.Equal(t, true, actual, "expected should make Service, defaulted to false with Enabled=true")

	// False case - default false with enabled nil
	actual = ShouldMakeService(k8schianetv1.Service{
		Enabled: nil,
	}, false)
	require.Equal(t, false, actual, "expected should not make Service, defaulted to false with Enabled=nil")

	// False case - default false with enabled nil
	disabled := false
	actual = ShouldMakeService(k8schianetv1.Service{
		Enabled: &disabled,
	}, false)
	require.Equal(t, false, actual, "expected should not make Service, defaulted to false with Enabled=false")
}

func TestShouldRollIntoMainPeerService(t *testing.T) {
	enabled := true
	disabled := false

	// True case
	actual := ShouldRollIntoMainPeerService(k8schianetv1.Service{
		Enabled:             &enabled,
		RollIntoPeerService: &enabled,
	})
	require.Equal(t, true, actual, "expected should roll into peer Service, enabled service and enabled roll-into")

	// False cases
	actual = ShouldRollIntoMainPeerService(k8schianetv1.Service{
		Enabled:             &disabled,
		RollIntoPeerService: &enabled,
	})
	require.Equal(t, false, actual, "expected should not roll into peer Service, disabled service and enabled roll-into")

	actual = ShouldRollIntoMainPeerService(k8schianetv1.Service{
		Enabled:             &disabled,
		RollIntoPeerService: &disabled,
	})
	require.Equal(t, false, actual, "expected should not roll into peer Service, disabled service and disabled roll-into")

	actual = ShouldRollIntoMainPeerService(k8schianetv1.Service{
		Enabled:             &enabled,
		RollIntoPeerService: nil,
	})
	require.Equal(t, false, actual, "expected should not roll into peer Service, enabled service and nil roll-into")

	actual = ShouldRollIntoMainPeerService(k8schianetv1.Service{
		Enabled:             nil,
		RollIntoPeerService: &enabled,
	})
	require.Equal(t, false, actual, "expected should not roll into peer Service, nil service and enabled roll-into")
}

func TestGetFullNodePort(t *testing.T) {
	// Get Mainnet port
	actual, err := GetFullNodePort(k8schianetv1.CommonSpecChia{}, nil)
	require.NoError(t, err)
	require.Equal(t, int32(consts.MainnetNodePort), actual, "expected mainnet full_node port")

	// Get default testnet port
	testTrue := true
	actual, err = GetFullNodePort(k8schianetv1.CommonSpecChia{
		Testnet: &testTrue,
	}, nil)
	require.NoError(t, err)
	require.Equal(t, int32(consts.TestnetNodePort), actual, "expected testnet full_node port")

	// Get custom full_node port
	var customPort uint16 = 58441
	actual, err = GetFullNodePort(k8schianetv1.CommonSpecChia{
		NetworkPort: &customPort,
	}, nil)
	require.NoError(t, err)
	require.Equal(t, int32(customPort), actual, "expected custom full_node port")

	// Get custom full_node port, defined in a ChiaNetwork
	networkDataPort := 58442
	networkData := map[string]string{
		"network_port": "58442",
	}
	actual, err = GetFullNodePort(k8schianetv1.CommonSpecChia{
		NetworkPort: &customPort,
	}, &networkData)
	require.NoError(t, err)
	require.Equal(t, int32(networkDataPort), actual, "expected custom full_node port from network data")
}

func TestGetChiaRootVolume(t *testing.T) {
	// emptyDir cases
	expected := corev1.Volume{
		Name: "chiaroot",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}

	// emptyDir case - nil storage config
	actual := GetExistingChiaRootVolume(nil)
	require.Equal(t, expected, actual, "expected emptyDir volume - nil storage config")

	// emptyDir case - nil ChiaRoot config
	actual = GetExistingChiaRootVolume(&k8schianetv1.StorageConfig{
		ChiaRoot: nil,
	})
	require.Equal(t, expected, actual, "expected emptyDir volume - nil ChiaRoot config")

	// emptyDir case - nil pvc and hpv configs
	actual = GetExistingChiaRootVolume(&k8schianetv1.StorageConfig{
		ChiaRoot: &k8schianetv1.ChiaRootConfig{
			PersistentVolumeClaim: nil,
			HostPathVolume:        nil,
		},
	})
	require.Equal(t, expected, actual, "expected emptyDir volume - nil PVC and HostPathVolume configs")

	// emptyDir case - empty claim name
	actual = GetExistingChiaRootVolume(&k8schianetv1.StorageConfig{
		ChiaRoot: &k8schianetv1.ChiaRootConfig{
			PersistentVolumeClaim: &k8schianetv1.PersistentVolumeClaimConfig{
				ClaimName: "",
			},
		},
	})
	require.Equal(t, expected, actual, "expected emptyDir volume - empty claim name")

	// emptyDir case - empty host path
	actual = GetExistingChiaRootVolume(&k8schianetv1.StorageConfig{
		ChiaRoot: &k8schianetv1.ChiaRootConfig{
			HostPathVolume: &k8schianetv1.HostPathVolumeConfig{
				Path: "",
			},
		},
	})
	require.Equal(t, expected, actual, "expected emptyDir volume - empty host path")

	// PVC case
	expected = corev1.Volume{
		Name: "chiaroot",
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: "testname",
			},
		},
	}
	actual = GetExistingChiaRootVolume(&k8schianetv1.StorageConfig{
		ChiaRoot: &k8schianetv1.ChiaRootConfig{
			PersistentVolumeClaim: &k8schianetv1.PersistentVolumeClaimConfig{
				ClaimName: "testname",
			},
		},
	})
	require.Equal(t, expected, actual, "expected PVC volume")

	// HostPath case
	expected = corev1.Volume{
		Name: "chiaroot",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "testpath",
			},
		},
	}
	actual = GetExistingChiaRootVolume(&k8schianetv1.StorageConfig{
		ChiaRoot: &k8schianetv1.ChiaRootConfig{
			HostPathVolume: &k8schianetv1.HostPathVolumeConfig{
				Path: "testpath",
			},
		},
	})
	require.Equal(t, expected, actual, "expected hostPath volume")
}

func TestGetExtraContainers(t *testing.T) {
	expected := []corev1.Container{
		{
			Name:  "testcar",
			Image: "image:tag",
			Env: []corev1.EnvVar{
				{
					Name:  "foo",
					Value: "bar",
				},
				{
					Name:  "CHIA_ROOT",
					Value: "/chia-data",
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "foo-volume",
					MountPath: "/bar/path",
				},
				{
					Name:      "chia-data",
					MountPath: "/chia-data",
				},
			},
		},
	}

	actual := GetExtraContainers([]k8schianetv1.ExtraContainer{
		{
			Container: corev1.Container{
				Name:  "testcar",
				Image: "image:tag",
				Env: []corev1.EnvVar{
					{
						Name:  "foo",
						Value: "bar",
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "foo-volume",
						MountPath: "/bar/path",
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "foo-volume",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
			ShareVolumeMounts: true,
			ShareEnv:          true,
		},
	}, corev1.Container{
		Name:  "chia",
		Image: "chia:tag",
		Env: []corev1.EnvVar{
			{
				Name:  "CHIA_ROOT",
				Value: "/chia-data",
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "chia-data",
				MountPath: "/chia-data",
			},
		},
	})
	require.Equal(t, expected, actual)
}

func TestChiaHealthcheckEnabled(t *testing.T) {
	// True case - default true
	actual := ChiaHealthcheckEnabled(k8schianetv1.SpecChiaHealthcheck{
		Enabled: nil,
	})
	require.Equal(t, true, actual, "expected healthcheck enabled, defaulted to true")

	// True case - set to true
	enabled := true
	actual = ChiaHealthcheckEnabled(k8schianetv1.SpecChiaHealthcheck{
		Enabled: &enabled,
	})
	require.Equal(t, true, actual, "expected healthcheck enabled, set to true")

	// False case - set to false
	disabled := false
	actual = ChiaHealthcheckEnabled(k8schianetv1.SpecChiaHealthcheck{
		Enabled: &disabled,
	})
	require.Equal(t, false, actual, "expected healthcheck disabled, set to false")
}

func TestChiaExporterEnabled(t *testing.T) {
	// True case - default true
	actual := ChiaExporterEnabled(k8schianetv1.SpecChiaExporter{
		Enabled: nil,
	})
	require.Equal(t, true, actual, "expected exporter enabled, defaulted to true")

	// True case - set to true
	enabled := true
	actual = ChiaExporterEnabled(k8schianetv1.SpecChiaExporter{
		Enabled: &enabled,
	})
	require.Equal(t, true, actual, "expected exporter enabled, set to true")

	// False case - set to false
	disabled := false
	actual = ChiaExporterEnabled(k8schianetv1.SpecChiaExporter{
		Enabled: &disabled,
	})
	require.Equal(t, false, actual, "expected exporter disabled, set to false")
}

func TestGetCommonChiaEnv(t *testing.T) {
	testCases := []struct {
		name        string
		spec        k8schianetv1.CommonSpecChia
		networkData *map[string]string
		expectedEnv []corev1.EnvVar
	}{
		{
			name:        "Basic Common Spec",
			spec:        k8schianetv1.CommonSpecChia{},
			networkData: nil,
			expectedEnv: []corev1.EnvVar{
				{
					Name:  "CHIA_ROOT",
					Value: "/chia-data",
				},
				{
					Name:  "ca",
					Value: "/chia-ca",
				},
				{
					Name:  "network_port",
					Value: "8444",
				},
				{
					Name:  "self_hostname",
					Value: "0.0.0.0",
				},
			},
		},
		{
			name: "Full Common Spec",
			spec: k8schianetv1.CommonSpecChia{
				DNSIntroducerAddress: stringPtr("test-dns.address"),
				IntroducerAddress:    stringPtr("test.address"),
				LogLevel:             stringPtr("DEBUG"),
				Network:              stringPtr("testnet11"),
				SourceRef:            stringPtr("main"),
				SelfHostname:         stringPtr("127.0.0.1"),
				Testnet:              boolPtr(true),
				Timezone:             stringPtr("America/Los_Angeles"),
			},
			networkData: nil,
			expectedEnv: []corev1.EnvVar{
				{
					Name:  "CHIA_ROOT",
					Value: "/chia-data",
				},
				{
					Name:  "TZ",
					Value: "America/Los_Angeles",
				},
				{
					Name:  "ca",
					Value: "/chia-ca",
				},
				{
					Name:  "dns_introducer_address",
					Value: "test-dns.address",
				},
				{
					Name:  "introducer_address",
					Value: "test.address",
				},
				{
					Name:  "log_level",
					Value: "DEBUG",
				},
				{
					Name:  "network",
					Value: "testnet11",
				},
				{
					Name:  "network_port",
					Value: "58444",
				},
				{
					Name:  "self_hostname",
					Value: "127.0.0.1",
				},
				{
					Name:  "source_ref",
					Value: "main",
				},
				{
					Name:  "testnet",
					Value: "true",
				},
			},
		},
		{
			name: "Full Common Spec and Network Data",
			spec: k8schianetv1.CommonSpecChia{
				DNSIntroducerAddress: stringPtr("test-dns.address"),
				IntroducerAddress:    stringPtr("test.address"),
				LogLevel:             stringPtr("DEBUG"),
				Network:              stringPtr("testnet11"),
				SourceRef:            stringPtr("main"),
				SelfHostname:         stringPtr("127.0.0.1"),
				Testnet:              boolPtr(true),
				Timezone:             stringPtr("America/Los_Angeles"),
			},
			networkData: mapPtr(map[string]string{
				"chia.network_overrides.config":    `{"testnetwork":{"address_prefix":"txch","default_full_node_port":58445}}`,
				"chia.network_overrides.constants": `{"testnetwork":{"GENESIS_CHALLENGE":"fcb55f73488f2959f45823cf795b7567061eba768bc985dfaef70aa3af0448cc","GENESIS_PRE_FARM_POOL_PUZZLE_HASH":"66c86d91a50d56d74f8f2884b42d65211d7384f0c48a6701d845b8681d93f2c6","GENESIS_PRE_FARM_FARMER_PUZZLE_HASH":"66c86d91a50d56d74f8f2884b42d65211d7384f0c48a6701d845b8681d93f2c6","AGG_SIG_ME_ADDITIONAL_DATA":"fcb55f73488f2959f45823cf795b7567061eba768bc985dfaef70aa3af0448cc","DIFFICULTY_CONSTANT_FACTOR":10052721566054,"DIFFICULTY_STARTING":30,"EPOCH_BLOCKS":768,"MEMPOOL_BLOCK_BUFFER":10,"MIN_PLOT_SIZE":18,"NETWORK_TYPE":1,"SUB_SLOT_ITERS_STARTING":67108864,"HARD_FORK_HEIGHT":0}}`,
				"dns_introducer_address":           "dns-introducer-testnetwork.address",
				"introducer_address":               "introducer-testnetwork.address",
				"network":                          "testnetwork",
				"network_port":                     "58445",
			}),
			expectedEnv: []corev1.EnvVar{
				{
					Name:  "CHIA_ROOT",
					Value: "/chia-data",
				},
				{
					Name:  "TZ",
					Value: "America/Los_Angeles",
				},
				{
					Name:  "ca",
					Value: "/chia-ca",
				},
				{
					Name:  "chia.network_overrides.config",
					Value: `{"testnetwork":{"address_prefix":"txch","default_full_node_port":58445}}`,
				},
				{
					Name:  "chia.network_overrides.constants",
					Value: `{"testnetwork":{"GENESIS_CHALLENGE":"fcb55f73488f2959f45823cf795b7567061eba768bc985dfaef70aa3af0448cc","GENESIS_PRE_FARM_POOL_PUZZLE_HASH":"66c86d91a50d56d74f8f2884b42d65211d7384f0c48a6701d845b8681d93f2c6","GENESIS_PRE_FARM_FARMER_PUZZLE_HASH":"66c86d91a50d56d74f8f2884b42d65211d7384f0c48a6701d845b8681d93f2c6","AGG_SIG_ME_ADDITIONAL_DATA":"fcb55f73488f2959f45823cf795b7567061eba768bc985dfaef70aa3af0448cc","DIFFICULTY_CONSTANT_FACTOR":10052721566054,"DIFFICULTY_STARTING":30,"EPOCH_BLOCKS":768,"MEMPOOL_BLOCK_BUFFER":10,"MIN_PLOT_SIZE":18,"NETWORK_TYPE":1,"SUB_SLOT_ITERS_STARTING":67108864,"HARD_FORK_HEIGHT":0}}`,
				},
				{
					Name:  "dns_introducer_address",
					Value: "dns-introducer-testnetwork.address",
				},
				{
					Name:  "introducer_address",
					Value: "introducer-testnetwork.address",
				},
				{
					Name:  "log_level",
					Value: "DEBUG",
				},
				{
					Name:  "network",
					Value: "testnetwork",
				},
				{
					Name:  "network_port",
					Value: "58445",
				},
				{
					Name:  "self_hostname",
					Value: "127.0.0.1",
				},
				{
					Name:  "source_ref",
					Value: "main",
				},
				{
					Name:  "testnet",
					Value: "true",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			env, err := GetCommonChiaEnv(tc.spec, tc.networkData)
			assert.NoError(t, err)
			assert.Equal(t, len(tc.expectedEnv), len(env), "Number of environment variables should match")
			for i, expectedEnv := range tc.expectedEnv {
				assert.Equal(t, expectedEnv.Name, env[i].Name, "Environment variable name should match")
				assert.Equal(t, expectedEnv.Value, env[i].Value, "Environment variable value should match")
			}
		})
	}
}
