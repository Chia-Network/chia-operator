/*
Copyright 2024 Chia Network Inc.
*/

package v1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func boolPtr(b bool) *bool {
	return &b
}

func TestUnmarshalChiaCrawler(t *testing.T) {
	yamlData := []byte(`
apiVersion: k8s.chia.net/v1
kind: ChiaCrawler
metadata:
  labels:
    app.kubernetes.io/name: chiacrawler
    app.kubernetes.io/instance: chiacrawler-sample
    app.kubernetes.io/part-of: chia-operator
    app.kubernetes.io/created-by: chia-operator
  name: chiacrawler-sample
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
		caSecret                    = "chiaca-secret"
	)
	expect := ChiaCrawler{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "k8s.chia.net/v1",
			Kind:       "ChiaCrawler",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "chiacrawler-sample",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "chiacrawler",
				"app.kubernetes.io/instance":   "chiacrawler-sample",
				"app.kubernetes.io/part-of":    "chia-operator",
				"app.kubernetes.io/created-by": "chia-operator",
			},
		},
		Spec: ChiaCrawlerSpec{
			ChiaConfig: ChiaCrawlerSpecChia{
				CommonSpecChia: CommonSpecChia{
					Testnet:              &testnet,
					Network:              &network,
					NetworkPort:          &networkPort,
					IntroducerAddress:    &introducerAddress,
					DNSIntroducerAddress: &dnsIntroducerAddress,
					Timezone:             &timezone,
					LogLevel:             &logLevel,
				},
				CASecretName: &caSecret,
			},
			CommonSpec: CommonSpec{
				ChiaExporterConfig: SpecChiaExporter{
					Enabled: boolPtr(true),
				},
			},
		},
	}

	var actual ChiaCrawler
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
