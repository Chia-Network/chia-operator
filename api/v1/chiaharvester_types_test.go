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

func TestUnmarshalChiaHarvester(t *testing.T) {
	yamlData := []byte(`
apiVersion: k8s.chia.net/v1
kind: ChiaHarvester
metadata:
  labels:
    app.kubernetes.io/name: chiaharvester
    app.kubernetes.io/instance: chiaharvester-sample
    app.kubernetes.io/part-of: chia-operator
    app.kubernetes.io/created-by: chia-operator
  name: chiaharvester-sample
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
    farmerAddress: "farmer.default.svc.cluster.local:58444"
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
	expect := ChiaHarvester{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "k8s.chia.net/v1",
			Kind:       "ChiaHarvester",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "chiaharvester-sample",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "chiaharvester",
				"app.kubernetes.io/instance":   "chiaharvester-sample",
				"app.kubernetes.io/part-of":    "chia-operator",
				"app.kubernetes.io/created-by": "chia-operator",
			},
		},
		Spec: ChiaHarvesterSpec{
			ChiaConfig: ChiaHarvesterSpecChia{
				CommonSpecChia: CommonSpecChia{
					CASecretName:         "chiaca-secret",
					Testnet:              &testnet,
					Network:              &network,
					NetworkPort:          &networkPort,
					IntroducerAddress:    &introducerAddress,
					DNSIntroducerAddress: &dnsIntroducerAddress,
					Timezone:             &timezone,
					LogLevel:             &logLevel,
				},
				FarmerAddress: "farmer.default.svc.cluster.local:58444",
			},
			CommonSpec: CommonSpec{
				ChiaExporterConfig: SpecChiaExporter{
					Enabled: true,
					ServiceLabels: map[string]string{
						"network": "testnet",
					},
				},
			},
		},
	}

	var actual ChiaHarvester
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
