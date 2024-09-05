package kube

import (
	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGetCommonLabels(t *testing.T) {
	expected := map[string]string{
		"app.kubernetes.io/instance":   "testname",
		"app.kubernetes.io/name":       "testname",
		"app.kubernetes.io/managed-by": "chia-operator",
		"k8s.chia.net/provenance":      "TestKind.testnamespace.testname",
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

func TestShouldMakeVolumeClaim(t *testing.T) {
	// True case
	actual := ShouldMakeVolumeClaim(&k8schianetv1.StorageConfig{
		ChiaRoot: &k8schianetv1.ChiaRootConfig{
			PersistentVolumeClaim: &k8schianetv1.PersistentVolumeClaimConfig{
				GenerateVolumeClaims: true,
			},
		},
	})
	require.Equal(t, true, actual, "expected should make volume claim")

	// False case - nil storage config
	actual = ShouldMakeVolumeClaim(nil)
	require.Equal(t, false, actual, "expected should not make volume claim for nil storage config")

	// False case - non-nil storage config, nil ChiaRoot config
	actual = ShouldMakeVolumeClaim(&k8schianetv1.StorageConfig{
		ChiaRoot: nil,
	})
	require.Equal(t, false, actual, "expected should not make volume claim for nil ChiaRoot config")

	// False case - non-nil storage config, nil PersistentVolumeClaim config
	actual = ShouldMakeVolumeClaim(&k8schianetv1.StorageConfig{
		ChiaRoot: &k8schianetv1.ChiaRootConfig{
			PersistentVolumeClaim: nil,
		},
	})
	require.Equal(t, false, actual, "expected should not make volume claim for nil PersistentVolumeClaim config")

	// False case - non-nil storage config, GenerateVolumeClaims set to false
	actual = ShouldMakeVolumeClaim(&k8schianetv1.StorageConfig{
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
	actual := GetFullNodePort(k8schianetv1.CommonSpecChia{})
	require.Equal(t, int32(consts.MainnetNodePort), actual, "expected mainnet full_node port")

	// Get default testnet port
	testTrue := true
	actual = GetFullNodePort(k8schianetv1.CommonSpecChia{
		Testnet: &testTrue,
	})
	require.Equal(t, int32(consts.TestnetNodePort), actual, "expected testnet full_node port")

	// Get custom full_node port
	var port uint16 = 58441
	actual = GetFullNodePort(k8schianetv1.CommonSpecChia{
		NetworkPort: &port,
	})
	require.Equal(t, int32(port), actual, "expected custom full_node port")
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
