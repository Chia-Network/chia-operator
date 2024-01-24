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

func TestUnmarshalChiaSeeder(t *testing.T) {
	yamlData := []byte(`
apiVersion: k8s.chia.net/v1
kind: ChiaSeeder
metadata:
  labels:
    app.kubernetes.io/name: chiaseeder
    app.kubernetes.io/instance: chiaseeder-sample
    app.kubernetes.io/part-of: chia-operator
    app.kubernetes.io/created-by: chia-operator
  name: chiaseeder-sample
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
    bootstrapPeer: "node.default.svc.cluster.local"
    minimumHeight: 100
    domainName: seeder.example.com.
    nameserver: example.com.
    rname: admin.example.com.
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
		bootstrapPeer               = "node.default.svc.cluster.local"
		minimumHeight        uint64 = 100
		domainName                  = "seeder.example.com."
		nameserver                  = "example.com."
		rname                       = "admin.example.com."
	)
	expect := ChiaSeeder{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "k8s.chia.net/v1",
			Kind:       "ChiaSeeder",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "chiaseeder-sample",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "chiaseeder",
				"app.kubernetes.io/instance":   "chiaseeder-sample",
				"app.kubernetes.io/part-of":    "chia-operator",
				"app.kubernetes.io/created-by": "chia-operator",
			},
		},
		Spec: ChiaSeederSpec{
			ChiaConfig: ChiaSeederSpecChia{
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
				BootstrapPeer: &bootstrapPeer,
				MinimumHeight: &minimumHeight,
				DomainName:    domainName,
				Nameserver:    nameserver,
				Rname:         &rname,
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

	var actual ChiaSeeder
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
