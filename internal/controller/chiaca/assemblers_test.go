package chiaca

import (
	"fmt"
	"testing"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
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
	Spec: k8schianetv1.ChiaCASpec{
		Secret: "test-ca-secret",
	},
}

var testObjMeta = metav1.ObjectMeta{
	Name:      "testname-chiaca-generator",
	Namespace: "testnamespace",
	Labels: map[string]string{
		"app.kubernetes.io/instance":   "testname",
		"app.kubernetes.io/name":       "testname",
		"app.kubernetes.io/managed-by": "chia-operator",
		"k8s.chia.net/provenance":      "ChiaCA.testnamespace.testname",
	},
}

func TestAssembleJob(t *testing.T) {
	var backoffLimit int32 = 3
	expected := batchv1.Job{
		ObjectMeta: testObjMeta,
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:      "Never",
					ServiceAccountName: "testname-chiaca-generator",
					Containers: []corev1.Container{
						{
							Name:  "chiaca-generator",
							Image: fmt.Sprintf("%s:%s", consts.DefaultChiaCAImageName, consts.DefaultChiaCAImageTag),
							Env: []corev1.EnvVar{
								{
									Name:  "NAMESPACE",
									Value: "testnamespace",
								},
								{
									Name:  "SECRET_NAME",
									Value: "test-ca-secret",
								},
							},
						},
					},
				},
			},
		},
	}
	actual := assembleJob(testCA)
	require.Equal(t, expected, actual)
}

func TestAssembleServiceAccount(t *testing.T) {
	expected := corev1.ServiceAccount{
		ObjectMeta: testObjMeta,
	}
	actual := assembleServiceAccount(testCA)
	require.Equal(t, expected, actual)
}

func TestAssembleRole(t *testing.T) {
	expected := rbacv1.Role{
		ObjectMeta: testObjMeta,
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"secrets",
				},
				Verbs: []string{
					"create",
				},
			},
		},
	}
	actual := assembleRole(testCA)
	require.Equal(t, expected, actual)
}

func TestAssembleRoleBinding(t *testing.T) {
	expected := rbacv1.RoleBinding{
		ObjectMeta: testObjMeta,
		Subjects: []rbacv1.Subject{
			{
				Kind: "ServiceAccount",
				Name: "testname-chiaca-generator",
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind: "Role",
			Name: "testname-chiaca-generator",
		},
	}
	actual := assembleRoleBinding(testCA)
	require.Equal(t, expected, actual)
}
