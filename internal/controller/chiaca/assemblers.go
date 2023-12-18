/*
Copyright 2023 Chia Network Inc.
*/

package chiaca

import (
	"context"
	"fmt"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const chiacaNamePattern = "%s-chiaca-generator"

// assembleJob assembles the Job resource for a ChiaCA CR
func (r *ChiaCAReconciler) assembleJob(ctx context.Context, ca k8schianetv1.ChiaCA) batchv1.Job {
	var job batchv1.Job = batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiacaNamePattern, ca.Name),
			Namespace:       ca.Namespace,
			Labels:          r.getLabels(ctx, ca),
			OwnerReferences: r.getOwnerReference(ctx, ca),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:      "Never",
					ServiceAccountName: fmt.Sprintf(chiacaNamePattern, ca.Name),
					Containers: []corev1.Container{
						{
							Name:  "chiaca-generator",
							Image: ca.Spec.Image,
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
	if ca.Spec.ImagePullSecret != "" {
		job.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
			{
				Name: ca.Spec.ImagePullSecret,
			},
		}
	}

	return job
}

// assembleServiceAccount assembles the ServiceAccount resource for a ChiaCA CR
func (r *ChiaCAReconciler) assembleServiceAccount(ctx context.Context, ca k8schianetv1.ChiaCA) corev1.ServiceAccount {
	return corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiacaNamePattern, ca.Name),
			Namespace:       ca.Namespace,
			Labels:          r.getLabels(ctx, ca),
			OwnerReferences: r.getOwnerReference(ctx, ca),
		},
	}
}

// assembleRole assembles the Role resource for a ChiaCA CR
func (r *ChiaCAReconciler) assembleRole(ctx context.Context, ca k8schianetv1.ChiaCA) rbacv1.Role {
	return rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiacaNamePattern, ca.Name),
			Namespace:       ca.Namespace,
			Labels:          r.getLabels(ctx, ca),
			OwnerReferences: r.getOwnerReference(ctx, ca),
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
func (r *ChiaCAReconciler) assembleRoleBinding(ctx context.Context, ca k8schianetv1.ChiaCA) rbacv1.RoleBinding {
	return rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiacaNamePattern, ca.Name),
			Namespace:       ca.Namespace,
			Labels:          r.getLabels(ctx, ca),
			OwnerReferences: r.getOwnerReference(ctx, ca),
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
