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

func TestUnmarshalChiaFarmer(t *testing.T) {
	yamlData := []byte(`
apiVersion: k8s.chia.net/v1
kind: ChiaFarmer
metadata:
  labels:
    app.kubernetes.io/name: chiafarmer
    app.kubernetes.io/instance: chiafarmer-sample
    app.kubernetes.io/part-of: chia-operator
    app.kubernetes.io/created-by: chia-operator
  name: chiafarmer-sample
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
    secretKey:
      name: "chiakey-secret"
      key: "key.txt"
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
	expect := ChiaFarmer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "k8s.chia.net/v1",
			Kind:       "ChiaFarmer",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "chiafarmer-sample",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "chiafarmer",
				"app.kubernetes.io/instance":   "chiafarmer-sample",
				"app.kubernetes.io/part-of":    "chia-operator",
				"app.kubernetes.io/created-by": "chia-operator",
			},
		},
		Spec: ChiaFarmerSpec{
			ChiaConfig: ChiaFarmerSpecChia{
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
				SecretKey: ChiaSecretKey{
					Name: "chiakey-secret",
					Key:  "key.txt",
				},
			},
			CommonSpec: CommonSpec{
				ChiaExporterConfig: SpecChiaExporter{
					Enabled: true,
				},
			},
		},
	}

	var actual ChiaFarmer
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
