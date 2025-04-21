/*
Copyright 2024 Chia Network Inc.
*/

package v1

import (
	"testing"

	"github.com/chia-network/go-chia-libs/pkg/config"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func TestUnmarshalChiaNetwork(t *testing.T) {
	yamlData := []byte(`
apiVersion: k8s.chia.net/v1
kind: ChiaNetwork
metadata:
  labels:
    app.kubernetes.io/name: chianetwork
    app.kubernetes.io/instance: chianetwork-sample
    app.kubernetes.io/part-of: chia-operator
    app.kubernetes.io/created-by: chia-operator
  name: chianetwork-sample
spec:
  constants:
    MIN_PLOT_SIZE: 18
    GENESIS_CHALLENGE: fb00c54298fc1c149afbf4c8996fb2317ae41e4649b934ca495991b7852b841
    GENESIS_PRE_FARM_POOL_PUZZLE_HASH: asdlsakldlskalskdsasdasdsadsadsadsadsdsadsas
    GENESIS_PRE_FARM_FARMER_PUZZLE_HASH: testestestestestestestesrestestestestestest
  config:
    address_prefix: txch
    default_full_node_port: 58444
  networkName: testnetz
  networkPort: 58444
  introducerAddress: intro.testnetz.example.com
  dnsIntroducerAddress: dnsintro.testnetz.example.com
`)

	var (
		networkConfig = config.NetworkConfig{
			AddressPrefix:       "txch",
			DefaultFullNodePort: 58444,
		}
		minPlotSize   uint8 = 18
		networkConsts       = NetworkConstants{
			MinPlotSize:                    &minPlotSize,
			GenesisChallenge:               "fb00c54298fc1c149afbf4c8996fb2317ae41e4649b934ca495991b7852b841",
			GenesisPreFarmPoolPuzzleHash:   "asdlsakldlskalskdsasdasdsadsadsadsadsdsadsas",
			GenesisPreFarmFarmerPuzzleHash: "testestestestestestestesrestestestestestest",
		}
		network                     = "testnetz"
		networkPort          uint16 = 58444
		introducerAddress           = "intro.testnetz.example.com"
		dnsIntroducerAddress        = "dnsintro.testnetz.example.com"
	)

	expect := ChiaNetwork{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "k8s.chia.net/v1",
			Kind:       "ChiaNetwork",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "chianetwork-sample",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "chianetwork",
				"app.kubernetes.io/instance":   "chianetwork-sample",
				"app.kubernetes.io/part-of":    "chia-operator",
				"app.kubernetes.io/created-by": "chia-operator",
			},
		},
		Spec: ChiaNetworkSpec{
			NetworkConfig:        &networkConfig,
			NetworkConstants:     &networkConsts,
			NetworkName:          &network,
			NetworkPort:          &networkPort,
			IntroducerAddress:    &introducerAddress,
			DNSIntroducerAddress: &dnsIntroducerAddress,
		},
	}

	var actual ChiaNetwork
	err := yaml.Unmarshal(yamlData, &actual)
	if err != nil {
		t.Errorf("Error unmarshaling yaml: %v", err)
		return
	}

	diff := cmp.Diff(actual, expect)
	if diff != "" {
		t.Errorf("Unmarshaled struct does not match the expected struct. Actual: %+v\nExpected: %+v\nDiff: %s", actual, expect, diff)
		return
	}
}
