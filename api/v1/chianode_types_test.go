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

func TestUnmarshalChiaNode(t *testing.T) {
	yamlData := []byte(`
apiVersion: k8s.chia.net/v1
kind: ChiaNode
metadata:
  labels:
    app.kubernetes.io/name: chianode
    app.kubernetes.io/instance: chianode-sample
    app.kubernetes.io/part-of: chia-operator
    app.kubernetes.io/created-by: chia-operator
  name: chianode-sample
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
    trustedCIDRs:
      - "192.168.0.0/16"
      - "10.0.0.0/8"
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
	expectCIDRs := []string{
		"192.168.0.0/16",
		"10.0.0.0/8",
	}
	expect := ChiaNode{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "k8s.chia.net/v1",
			Kind:       "ChiaNode",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "chianode-sample",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "chianode",
				"app.kubernetes.io/instance":   "chianode-sample",
				"app.kubernetes.io/part-of":    "chia-operator",
				"app.kubernetes.io/created-by": "chia-operator",
			},
		},
		Spec: ChiaNodeSpec{
			ChiaConfig: ChiaNodeSpecChia{
				CommonSpecChia: CommonSpecChia{
					Testnet:              &testnet,
					Network:              &network,
					NetworkPort:          &networkPort,
					IntroducerAddress:    &introducerAddress,
					DNSIntroducerAddress: &dnsIntroducerAddress,
					Timezone:             &timezone,
					LogLevel:             &logLevel,
				},
				CASecretName: "chiaca-secret",
				TrustedCIDRs: &expectCIDRs,
				FullNodePeers: &[]Peer{
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

	var actual ChiaNode
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
