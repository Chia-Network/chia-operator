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

func TestUnmarshalChiaCA(t *testing.T) {
	yamlData := []byte(`
apiVersion: k8s.chia.net/v1
kind: ChiaCA
metadata:
  labels:
    app.kubernetes.io/name: chiaca
    app.kubernetes.io/instance: chiaca-sample
    app.kubernetes.io/part-of: chia-operator
    app.kubernetes.io/created-by: chia-operator
  name: chiaca-sample
spec:
  image: "ca-gen-image:latest"
  imagePullSecret: registrypullsecret
  secret: chiaca-secret
`)

	expect := ChiaCA{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "k8s.chia.net/v1",
			Kind:       "ChiaCA",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "chiaca-sample",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "chiaca",
				"app.kubernetes.io/instance":   "chiaca-sample",
				"app.kubernetes.io/part-of":    "chia-operator",
				"app.kubernetes.io/created-by": "chia-operator",
			},
		},
		Spec: ChiaCASpec{
			Image:           "ca-gen-image:latest",
			ImagePullSecret: "registrypullsecret",
			Secret:          "chiaca-secret",
		},
	}

	var actual ChiaCA
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
