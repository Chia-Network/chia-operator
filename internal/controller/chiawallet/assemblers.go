/*
Copyright 2023 Chia Network Inc.
*/

package chiawallet

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
)

const chiawalletNamePattern = "%s-wallet"

// assembleBaseService reconciles the main Service resource for a ChiaWallet CR
func (r *ChiaWalletReconciler) assembleBaseService(ctx context.Context, wallet k8schianetv1.ChiaWallet) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiawalletNamePattern, wallet.Name),
			Namespace:       wallet.Namespace,
			Labels:          kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels),
			Annotations:     wallet.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, wallet),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(wallet.Spec.ServiceType),
			Ports: []corev1.ServicePort{
				{
					Port:       consts.DaemonPort,
					TargetPort: intstr.FromString("daemon"),
					Protocol:   "TCP",
					Name:       "daemon",
				},
				{
					Port:       consts.WalletPort,
					TargetPort: intstr.FromString("peers"),
					Protocol:   "TCP",
					Name:       "peers",
				},
				{
					Port:       consts.WalletRPCPort,
					TargetPort: intstr.FromString("rpc"),
					Protocol:   "TCP",
					Name:       "rpc",
				},
			},
			Selector: kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaWallet CR
func (r *ChiaWalletReconciler) assembleChiaExporterService(ctx context.Context, wallet k8schianetv1.ChiaWallet) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiawalletNamePattern, wallet.Name) + "-metrics",
			Namespace:       wallet.Namespace,
			Labels:          kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels, wallet.Spec.ChiaExporterConfig.ServiceLabels),
			Annotations:     wallet.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, wallet),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType("ClusterIP"),
			Ports: []corev1.ServicePort{
				{
					Port:       consts.ChiaExporterPort,
					TargetPort: intstr.FromString("metrics"),
					Protocol:   "TCP",
					Name:       "metrics",
				},
			},
			Selector: kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleDeployment reconciles the wallet Deployment resource for a ChiaWallet CR
func (r *ChiaWalletReconciler) assembleDeployment(ctx context.Context, wallet k8schianetv1.ChiaWallet) appsv1.Deployment {
	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiawalletNamePattern, wallet.Name),
			Namespace:       wallet.Namespace,
			Labels:          kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels),
			Annotations:     wallet.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, wallet),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels),
					Annotations: wallet.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					// TODO add: imagePullSecret, serviceAccountName config
					Containers: []corev1.Container{
						{
							Name:            "chia",
							Image:           wallet.Spec.ChiaConfig.Image,
							ImagePullPolicy: wallet.Spec.ImagePullPolicy,
							Env:             r.getChiaEnv(ctx, wallet),
							Ports: []corev1.ContainerPort{
								{
									Name:          "daemon",
									ContainerPort: consts.DaemonPort,
									Protocol:      "TCP",
								},
								{
									Name:          "peers",
									ContainerPort: consts.WalletPort,
									Protocol:      "TCP",
								},
								{
									Name:          "rpc",
									ContainerPort: consts.WalletRPCPort,
									Protocol:      "TCP",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "secret-ca",
									MountPath: "/chia-ca",
								},
								{
									Name:      "key",
									MountPath: "/key",
								},
								{
									Name:      "chiaroot",
									MountPath: "/chia-data",
								},
							},
						},
					},
					NodeSelector: wallet.Spec.NodeSelector,
					Volumes:      r.getChiaVolumes(ctx, wallet),
				},
			},
		},
	}

	var containerSecurityContext *corev1.SecurityContext
	if wallet.Spec.ChiaConfig.SecurityContext != nil {
		containerSecurityContext = wallet.Spec.ChiaConfig.SecurityContext
		deploy.Spec.Template.Spec.Containers[0].SecurityContext = wallet.Spec.ChiaConfig.SecurityContext
	}

	if wallet.Spec.ChiaConfig.LivenessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].LivenessProbe = wallet.Spec.ChiaConfig.LivenessProbe
	}

	if wallet.Spec.ChiaConfig.ReadinessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].ReadinessProbe = wallet.Spec.ChiaConfig.ReadinessProbe
	}

	if wallet.Spec.ChiaConfig.StartupProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].StartupProbe = wallet.Spec.ChiaConfig.StartupProbe
	}

	var containerResorces corev1.ResourceRequirements
	if wallet.Spec.ChiaConfig.Resources != nil {
		containerResorces = *wallet.Spec.ChiaConfig.Resources
		deploy.Spec.Template.Spec.Containers[0].Resources = *wallet.Spec.ChiaConfig.Resources
	}

	if wallet.Spec.ChiaExporterConfig.Enabled {
		exporterContainer := kube.GetChiaExporterContainer(ctx, wallet.Spec.ChiaExporterConfig.Image, containerSecurityContext, wallet.Spec.ImagePullPolicy, containerResorces)
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, exporterContainer)
	}

	if wallet.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = wallet.Spec.PodSecurityContext
	}

	if len(wallet.Spec.SidecarContainers) > 0 {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, wallet.Spec.SidecarContainers...)
	}

	// TODO add pod affinity, tolerations

	return deploy
}
