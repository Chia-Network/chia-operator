/*
Copyright 2023 Chia Network Inc.
*/

package chiaca

import (
	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultChiaCASecretName = "chiaca"

func assembleCASecret(ca k8schianetv1.ChiaCA, publicCACrt, publicCAKey, privateCACrt, privateCAKey string) corev1.Secret {
	secretName := defaultChiaCASecretName
	if ca.Spec.Secret != "" {
		secretName = ca.Spec.Secret
	}
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: ca.Namespace,
		},
		StringData: map[string]string{
			"chia_ca.crt":    publicCACrt,
			"chia_ca.key":    publicCAKey,
			"private_ca.crt": privateCACrt,
			"private_ca.key": privateCAKey,
		},
	}
}
