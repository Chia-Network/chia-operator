/*
Copyright 2024 Chia Network Inc.
*/

package chianetwork

import (
	"strconv"
	"testing"

	"github.com/chia-network/go-chia-libs/pkg/config"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testNetwork = k8schianetv1.ChiaNetwork{
	TypeMeta: metav1.TypeMeta{
		Kind:       "ChiaNetwork",
		APIVersion: "k8s.chia.net/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "testname",
		Namespace: "testnamespace",
	},
	Spec: k8schianetv1.ChiaNetworkSpec{
		NetworkConstants:     &networkConsts,
		NetworkConfig:        &networkConfig,
		NetworkName:          &network,
		NetworkPort:          &networkPort,
		IntroducerAddress:    &introducerAddress,
		DNSIntroducerAddress: &dnsIntroducerAddress,
	},
}

var (
	networkConfig = config.NetworkConfig{
		AddressPrefix:       "txch",
		DefaultFullNodePort: 58444,
	}
	minPlotSize    uint8  = 18
	hardForkHeight uint32 = 0
	networkConsts         = k8schianetv1.NetworkConstants{
		MinPlotSize:                    &minPlotSize,
		HardForkHeight:                 &hardForkHeight,
		GenesisChallenge:               "fb00c54298fc1c149afbf4c8996fb2317ae41e4649b934ca495991b7852b841",
		GenesisPreFarmPoolPuzzleHash:   "asdlsakldlskalskdsasdasdsadsadsadsadsdsadsas",
		GenesisPreFarmFarmerPuzzleHash: "testestestestestestestesrestestestestestest",
	}
	network                     = "testnetz"
	networkPort          uint16 = 58444
	introducerAddress           = "intro.testnetz.example.com"
	dnsIntroducerAddress        = "dnsintro.testnetz.example.com"
)

func TestAssembleConfigMap(t *testing.T) {
	expected := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testname",
			Namespace: "testnamespace",
		},
		Data: map[string]string{
			"chia.network_overrides.constants": `{"testnetz":{"GENESIS_CHALLENGE":"fb00c54298fc1c149afbf4c8996fb2317ae41e4649b934ca495991b7852b841","GENESIS_PRE_FARM_POOL_PUZZLE_HASH":"asdlsakldlskalskdsasdasdsadsadsadsadsdsadsas","GENESIS_PRE_FARM_FARMER_PUZZLE_HASH":"testestestestestestestesrestestestestestest","MIN_PLOT_SIZE":18,"HARD_FORK_HEIGHT":0}}`,
			"chia.network_overrides.config":    `{"testnetz":{"address_prefix":"txch","default_full_node_port":58444}}`,
			"network":                          network,
			"network_port":                     strconv.Itoa(int(networkPort)),
			"introducer_address":               introducerAddress,
			"dns_introducer_address":           dnsIntroducerAddress,
		},
	}
	actual, err := assembleConfigMap(testNetwork)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
