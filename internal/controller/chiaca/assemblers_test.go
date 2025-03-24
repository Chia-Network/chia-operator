package chiaca

import (
	"testing"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testCA = k8schianetv1.ChiaCA{
	TypeMeta: metav1.TypeMeta{
		Kind:       "ChiaCA",
		APIVersion: "k8s.chia.net/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "testname",
		Namespace: "testnamespace",
	},
}

func TestAssembleCASecret_DefaultSecretName(t *testing.T) {
	expected := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testname",
			Namespace: "testnamespace",
		},
		StringData: map[string]string{
			"chia_ca.crt":    "publicCACert",
			"chia_ca.key":    "publicCAKey",
			"private_ca.crt": "privateCACert",
			"private_ca.key": "privateCAKey",
		},
	}
	actual := assembleCASecret(testCA, "publicCACert", "publicCAKey", "privateCACert", "privateCAKey")
	require.Equal(t, expected, actual)
}

func TestAssembleCASecret_CustomSecretName(t *testing.T) {
	customCA := testCA
	customCA.Spec.Secret = "chiaca-custom"
	expected := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "chiaca-custom",
			Namespace: "testnamespace",
		},
		StringData: map[string]string{
			"chia_ca.crt":    "publicCACert",
			"chia_ca.key":    "publicCAKey",
			"private_ca.crt": "privateCACert",
			"private_ca.key": "privateCAKey",
		},
	}
	actual := assembleCASecret(customCA, "publicCACert", "publicCAKey", "privateCACert", "privateCAKey")
	require.Equal(t, expected, actual)
}
