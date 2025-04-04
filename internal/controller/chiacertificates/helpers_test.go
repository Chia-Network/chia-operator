/*
Copyright 2025 Chia Network Inc.
*/

package chiacertificates

import (
	"testing"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/go-chia-libs/pkg/tls"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testChiaCertificates = k8schianetv1.ChiaCertificates{
	TypeMeta: metav1.TypeMeta{
		Kind:       "ChiaCertificates",
		APIVersion: "k8s.chia.net/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "testname",
		Namespace: "testnamespace",
	},
	Spec: k8schianetv1.ChiaCertificatesSpec{
		CASecretName: "ca-secret",
	},
}

func TestAssembleSecret_DefaultName(t *testing.T) {
	// Test with default name (no custom secret name specified)
	certMap := map[string]string{
		"test.crt": "test-cert-data",
		"test.key": "test-key-data",
	}

	secret := assembleSecret(testChiaCertificates, certMap)

	assert.Equal(t, "testname", secret.Name)
	assert.Equal(t, "testnamespace", secret.Namespace)
	assert.Equal(t, "ca-secret", secret.Labels["k8s.chia.net/chiaca.secret"])
	assert.Equal(t, "test-cert-data", secret.StringData["test.crt"])
	assert.Equal(t, "test-key-data", secret.StringData["test.key"])
}

func TestAssembleSecret_CustomName(t *testing.T) {
	// Test with custom secret name
	customCertificates := testChiaCertificates
	customCertificates.Spec.Secret = "custom-secret-name"

	certMap := map[string]string{
		"test.crt": "test-cert-data",
		"test.key": "test-key-data",
	}

	secret := assembleSecret(customCertificates, certMap)

	assert.Equal(t, "custom-secret-name", secret.Name)
	assert.Equal(t, "testnamespace", secret.Namespace)
	assert.Equal(t, "ca-secret", secret.Labels["k8s.chia.net/chiaca.secret"])
	assert.Equal(t, "test-cert-data", secret.StringData["test.crt"])
	assert.Equal(t, "test-key-data", secret.StringData["test.key"])
}

func TestGetChiaCertificatesSecretName_DefaultName(t *testing.T) {
	// Test with default name (no custom secret name specified)
	secretName := getChiaCertificatesSecretName(testChiaCertificates)
	assert.Equal(t, "testname", secretName)
}

func TestGetChiaCertificatesSecretName_CustomName(t *testing.T) {
	// Test with custom secret name
	customCertificates := testChiaCertificates
	customCertificates.Spec.Secret = "custom-secret-name"
	secretName := getChiaCertificatesSecretName(customCertificates)
	assert.Equal(t, "custom-secret-name", secretName)
}

func TestGetChiaCertificatesSecretName_EmptyString(t *testing.T) {
	// Test with empty string (should use default name)
	customCertificates := testChiaCertificates
	customCertificates.Spec.Secret = ""
	secretName := getChiaCertificatesSecretName(customCertificates)
	assert.Equal(t, "testname", secretName)
}

func TestGetChiaCertificatesSecretName_WhitespaceString(t *testing.T) {
	// Test with whitespace string (should use default name)
	customCertificates := testChiaCertificates
	customCertificates.Spec.Secret = "   "
	secretName := getChiaCertificatesSecretName(customCertificates)
	assert.Equal(t, "testname", secretName)
}

func TestConstructCertMap_NilCertificate(t *testing.T) {
	// Create a test ChiaCertificates object with a nil certificate pair
	allCerts := &tls.ChiaCertificates{
		PrivateCrawler: nil, // This should cause an error
	}

	// Call the function
	certMap, err := constructCertMap(allCerts)

	// Assert the results
	assert.Error(t, err)
	assert.Nil(t, certMap)
	assert.Contains(t, err.Error(), "key pair nil")
}
