/*
Copyright 2023 Chia Network Inc.
*/

package chiaca

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
)

// caSecretExists fetches the k8s Secret that matches this ChiaCA deployment. Returns true if the Secret exists.
func (r *ChiaCAReconciler) caSecretExists(ctx context.Context, ca k8schianetv1.ChiaCA) (bool, error) {
	var secret corev1.Secret
	err := r.Get(ctx, types.NamespacedName{
		Namespace: ca.Namespace,
		Name:      getChiaCASecretName(ca),
	}, &secret)
	if err != nil && errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// getChiaCASecretName gets the name of the Secret to check if it exists
func getChiaCASecretName(ca k8schianetv1.ChiaCA) string {
	secretName := ca.Name
	if strings.TrimSpace(ca.Spec.Secret) != "" {
		secretName = ca.Spec.Secret
	}
	return secretName
}
