/*
Copyright 2025 Chia Network Inc.
*/

package chiaca

import (
	"testing"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testChiaCA = k8schianetv1.ChiaCA{
	TypeMeta: metav1.TypeMeta{
		Kind:       "ChiaCA",
		APIVersion: "k8s.chia.net/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "testname",
		Namespace: "testnamespace",
	},
}

func TestGetChiaCASecretName_DefaultName(t *testing.T) {
	// Test with default name (no custom secret name specified)
	secretName := getChiaCASecretName(testChiaCA)
	assert.Equal(t, "testname", secretName)
}

func TestGetChiaCASecretName_CustomName(t *testing.T) {
	// Test with custom secret name
	customCA := testChiaCA
	customCA.Spec.Secret = "custom-secret-name"
	secretName := getChiaCASecretName(customCA)
	assert.Equal(t, "custom-secret-name", secretName)
}

func TestGetChiaCASecretName_EmptyString(t *testing.T) {
	// Test with empty string (should use default name)
	customCA := testChiaCA
	customCA.Spec.Secret = ""
	secretName := getChiaCASecretName(customCA)
	assert.Equal(t, "testname", secretName)
}

func TestGetChiaCASecretName_WhitespaceString(t *testing.T) {
	// Test with whitespace string (should use default name)
	customCA := testChiaCA
	customCA.Spec.Secret = "   "
	secretName := getChiaCASecretName(customCA)
	assert.Equal(t, "testname", secretName)
}
