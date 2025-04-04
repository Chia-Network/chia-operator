package chiawallet

import (
	"testing"

	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/stretchr/testify/assert"
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
	volumeMounts := getChiaVolumeMounts()

	assert.Len(t, volumeMounts, 3, "Expected 3 volume mounts")

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
