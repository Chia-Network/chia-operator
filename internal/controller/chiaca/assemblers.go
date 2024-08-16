/*
Copyright 2023 Chia Network Inc.
*/

package chiaca

import (
	"fmt"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const chiacaNamePattern = "%s-chiaca-generator"

// assembleJob assembles the Job resource for a ChiaCA CR
func assembleJob(ca k8schianetv1.ChiaCA) batchv1.Job {
	var job = batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(chiacaNamePattern, ca.Name),
			Namespace: ca.Namespace,
			Labels:    kube.GetCommonLabels(ca.Kind, ca.ObjectMeta),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:      "Never",
					ServiceAccountName: fmt.Sprintf(chiacaNamePattern, ca.Name),
					Containers: []corev1.Container{
						{
							Name: "chiaca-generator",
							Env: []corev1.EnvVar{
								{
									Name:  "NAMESPACE",
									Value: ca.Namespace,
								},
								{
									Name:  "SECRET_NAME",
									Value: ca.Spec.Secret,
								},
							},
						},
					},
				},
			},
		},
	}

	if ca.Spec.Image != nil && *ca.Spec.Image != "" {
		job.Spec.Template.Spec.Containers[0].Image = *ca.Spec.Image
	} else {
		job.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("%s:%s", consts.DefaultChiaCAImageName, consts.DefaultChiaCAImageTag)
	}

	if ca.Spec.ImagePullSecret != "" {
		job.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
			{
				Name: ca.Spec.ImagePullSecret,
			},
		}
	}

	// Set backoff limit to 3, which is he maximum number of retries for ChiaCA Jobs
	backoff := int32(3)
	job.Spec.BackoffLimit = &backoff

	return job
}

// assembleServiceAccount assembles the ServiceAccount resource for a ChiaCA CR
func assembleServiceAccount(ca k8schianetv1.ChiaCA) corev1.ServiceAccount {
	return corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(chiacaNamePattern, ca.Name),
			Namespace: ca.Namespace,
			Labels:    kube.GetCommonLabels(ca.Kind, ca.ObjectMeta),
		},
	}
}

// assembleRole assembles the Role resource for a ChiaCA CR
func assembleRole(ca k8schianetv1.ChiaCA) rbacv1.Role {
	return rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(chiacaNamePattern, ca.Name),
			Namespace: ca.Namespace,
			Labels:    kube.GetCommonLabels(ca.Kind, ca.ObjectMeta),
		},
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
}

// assembleRoleBinding assembles the RoleBinding resource for a ChiaCA CR
func assembleRoleBinding(ca k8schianetv1.ChiaCA) rbacv1.RoleBinding {
	return rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(chiacaNamePattern, ca.Name),
			Namespace: ca.Namespace,
			Labels:    kube.GetCommonLabels(ca.Kind, ca.ObjectMeta),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "ServiceAccount",
				Name: fmt.Sprintf(chiacaNamePattern, ca.Name),
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind: "Role",
			Name: fmt.Sprintf(chiacaNamePattern, ca.Name),
		},
	}
}
