/*
Copyright 2025 Chia Network Inc.
*/

package chiacertificates

import (
	"context"
	"fmt"
	"github.com/chia-network/go-chia-libs/pkg/tls"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

// getSecret fetches the k8s Secret that matches this ChiaCertificates deployment. Returns true if the Secret exists.
func (r *ChiaCertificatesReconciler) getSecret(ctx context.Context, namespace, name string) (corev1.Secret, bool, error) {
	var secret corev1.Secret
	err := r.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, &secret)
	if err != nil && errors.IsNotFound(err) {
		return corev1.Secret{}, false, nil
	}
	if err != nil {
		return corev1.Secret{}, false, err
	}
	return secret, true, nil
}

type fetchCertKeyPair func(certs *tls.ChiaCertificates) *tls.CertificateKeyPair

var certNodes = map[string]fetchCertKeyPair{
	"private_crawler":    func(c *tls.ChiaCertificates) *tls.CertificateKeyPair { return c.PrivateCrawler },
	"private_daemon":     func(c *tls.ChiaCertificates) *tls.CertificateKeyPair { return c.PrivateDaemon },
	"private_data_layer": func(c *tls.ChiaCertificates) *tls.CertificateKeyPair { return c.PrivateDatalayer },
	"public_data_layer":  func(c *tls.ChiaCertificates) *tls.CertificateKeyPair { return c.PublicDatalayer },
	"private_farmer":     func(c *tls.ChiaCertificates) *tls.CertificateKeyPair { return c.PrivateFarmer },
	"public_farmer":      func(c *tls.ChiaCertificates) *tls.CertificateKeyPair { return c.PublicFarmer },
	"private_full_node":  func(c *tls.ChiaCertificates) *tls.CertificateKeyPair { return c.PrivateFullNode },
	"public_full_node":   func(c *tls.ChiaCertificates) *tls.CertificateKeyPair { return c.PublicFullNode },
	"private_harvester":  func(c *tls.ChiaCertificates) *tls.CertificateKeyPair { return c.PrivateHarvester },
	"public_introducer":  func(c *tls.ChiaCertificates) *tls.CertificateKeyPair { return c.PublicIntroducer },
	"private_timelord":   func(c *tls.ChiaCertificates) *tls.CertificateKeyPair { return c.PrivateTimelord },
	"public_timelord":    func(c *tls.ChiaCertificates) *tls.CertificateKeyPair { return c.PublicTimelord },
	"private_wallet":     func(c *tls.ChiaCertificates) *tls.CertificateKeyPair { return c.PrivateWallet },
	"public_wallet":      func(c *tls.ChiaCertificates) *tls.CertificateKeyPair { return c.PublicWallet },
}

func constructCertMap(allCerts *tls.ChiaCertificates) (map[string]string, error) {
	certMap := make(map[string]string)

	for filenameBase, fetchCertKeyPairFunc := range certNodes {
		crtKey := fetchCertKeyPairFunc(allCerts)
		if crtKey == nil {
			return nil, fmt.Errorf("key pair nil, but expected data for %s", filenameBase)
		}

		cert, key, err := tls.EncodeCertAndKeyToPEM(crtKey.CertificateDER, crtKey.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("encoding cert and key to PEM for %s: %w", filenameBase, err)
		}

		certFilename := filenameBase + ".crt"
		certMap[certFilename] = string(cert)

		keyFilename := filenameBase + ".key"
		certMap[keyFilename] = string(key)
	}

	return certMap, nil
}
