/*
Copyright 2023 Chia Network Inc.
*/

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func TestUnmarshalChiaDataLayer(t *testing.T) {
	yamlData := []byte(`
apiVersion: k8s.chia.net/v1
kind: ChiaDataLayer
metadata:
  labels:
    app.kubernetes.io/name: chiadatalayer
    app.kubernetes.io/instance: chiadatalayer-sample
    app.kubernetes.io/part-of: chia-operator
    app.kubernetes.io/created-by: chia-operator
  name: chiadatalayer-sample
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
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
    trustedCIDRs:
      - "192.168.0.0/16"
      - "10.0.0.0/8"
  dataLayerHTTP:
    enabled: true
    service:
      enabled: true
      labels:
        key: value
  chiaExporter:
    enabled: true
    service:
    serviceLabels:
      network: testnet
`)

	var (
		testTrue                    = true
		timezone                    = "UTC"
		logLevel                    = "INFO"
		network                     = "testnet68419"
		networkPort          uint16 = 8080
		introducerAddress           = "introducer.svc.cluster.local"
		dnsIntroducerAddress        = "dns-introducer.svc.cluster.local"
		caSecret                    = "chiaca-secret"
		strategy                    = appsv1.DeploymentStrategy{
			Type: appsv1.RollingUpdateDeploymentStrategyType,
			RollingUpdate: &appsv1.RollingUpdateDeployment{
				MaxSurge:       &intstr.IntOrString{IntVal: 1},
				MaxUnavailable: &intstr.IntOrString{IntVal: 1},
			},
		}
	)
	expectCIDRs := []string{
		"192.168.0.0/16",
		"10.0.0.0/8",
	}
	expect := ChiaDataLayer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "k8s.chia.net/v1",
			Kind:       "ChiaDataLayer",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "chiadatalayer-sample",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "chiadatalayer",
				"app.kubernetes.io/instance":   "chiadatalayer-sample",
				"app.kubernetes.io/part-of":    "chia-operator",
				"app.kubernetes.io/created-by": "chia-operator",
			},
		},
		Spec: ChiaDataLayerSpec{
			Strategy: &strategy,
			ChiaConfig: ChiaDataLayerSpecChia{
				CommonSpecChia: CommonSpecChia{
					Testnet:              &testTrue,
					Timezone:             &timezone,
					LogLevel:             &logLevel,
					Network:              &network,
					NetworkPort:          &networkPort,
					IntroducerAddress:    &introducerAddress,
					DNSIntroducerAddress: &dnsIntroducerAddress,
				},
				CASecretName: &caSecret,
				FullNodePeers: &[]Peer{
					{
						Host: "node.default.svc.cluster.local",
						Port: 58444,
					},
				},
				SecretKey: ChiaSecretKey{
					Name: "chiakey-secret",
					Key:  "key.txt",
				},
				TrustedCIDRs: &expectCIDRs,
			},
			CommonSpec: CommonSpec{
				ChiaExporterConfig: SpecChiaExporter{
					Enabled: testTrue,
				},
			},
			DataLayerHTTPConfig: ChiaDataLayerHTTPSpecChia{
				Enabled: &testTrue,
				Service: Service{
					Enabled: &testTrue,
					AdditionalMetadata: AdditionalMetadata{
						Labels: map[string]string{
							"key": "value",
						},
					},
				},
			},
		},
	}

	var actual ChiaDataLayer
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
