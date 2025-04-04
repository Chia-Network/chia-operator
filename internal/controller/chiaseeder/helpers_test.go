package chiaseeder

import (
	"testing"

	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/stretchr/testify/assert"
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
