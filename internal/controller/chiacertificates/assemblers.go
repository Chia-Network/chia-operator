/*
Copyright 2025 Chia Network Inc.
*/

package chiacertificates

import (
	"strings"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultChiaCertificatesSecretName = "chiacertificates"

func assembleSecret(cr k8schianetv1.ChiaCertificates, certMap map[string]string) corev1.Secret {
	secretName := defaultChiaCertificatesSecretName
	if strings.TrimSpace(cr.Spec.Secret) != "" {
		secretName = cr.Spec.Secret
	}
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"k8s.chia.net/chiaca.secret": cr.Spec.CASecretName,
			},
		},
		StringData: certMap,
	}
}
