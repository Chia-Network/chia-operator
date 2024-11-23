/*
Copyright 2023 Chia Network Inc.
*/

package v1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func TestUnmarshalChiaTimelord(t *testing.T) {
	yamlData := []byte(`
apiVersion: k8s.chia.net/v1
kind: ChiaTimelord
metadata:
  labels:
    app.kubernetes.io/name: chiatimelord
    app.kubernetes.io/instance: chiatimelord-sample
    app.kubernetes.io/part-of: chia-operator
    app.kubernetes.io/created-by: chia-operator
  name: chiatimelord-sample
spec:
  chia:
    caSecretName: chiaca-secret
    testnet: true
    network: testnet68419
    networkPort: 8080
    introducerAddress: introducer.svc.cluster.local
    dnsIntroducerAddress: dns-introducer.svc.cluster.local
    timezone: "UTC"
    logLevel: "INFO"
    fullNodePeers: 
      - host: "node.default.svc.cluster.local"
        port: 58444
  chiaExporter:
    enabled: true
    serviceLabels:
      network: testnet
`)

	var (
		testnet                     = true
		timezone                    = "UTC"
		logLevel                    = "INFO"
		network                     = "testnet68419"
		networkPort          uint16 = 8080
		introducerAddress           = "introducer.svc.cluster.local"
		dnsIntroducerAddress        = "dns-introducer.svc.cluster.local"
	)
	expect := ChiaTimelord{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "k8s.chia.net/v1",
			Kind:       "ChiaTimelord",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "chiatimelord-sample",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "chiatimelord",
				"app.kubernetes.io/instance":   "chiatimelord-sample",
				"app.kubernetes.io/part-of":    "chia-operator",
				"app.kubernetes.io/created-by": "chia-operator",
			},
		},
		Spec: ChiaTimelordSpec{
			ChiaConfig: ChiaTimelordSpecChia{
				CommonSpecChia: CommonSpecChia{
					Testnet:              &testnet,
					Timezone:             &timezone,
					LogLevel:             &logLevel,
					Network:              &network,
					NetworkPort:          &networkPort,
					IntroducerAddress:    &introducerAddress,
					DNSIntroducerAddress: &dnsIntroducerAddress,
				},
				CASecretName: "chiaca-secret",
				FullNodePeers: &[]FullNodePeer{
					{
						Host: "node.default.svc.cluster.local",
						Port: 58444,
					},
				},
			},
			CommonSpec: CommonSpec{
				ChiaExporterConfig: SpecChiaExporter{
					Enabled: true,
				},
			},
		},
	}

	var actual ChiaTimelord
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
