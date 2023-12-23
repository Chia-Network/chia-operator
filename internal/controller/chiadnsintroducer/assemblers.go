/*
Copyright 2023 Chia Network Inc.
*/

package chiadnsintroducer

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

const chiadnsIntroNamePattern = "%s-dns-introducer"

// assembleBaseService assembles the main Service resource for a ChiaDNSIntroducer CR
func (r *ChiaDNSIntroducerReconciler) assembleBaseService(ctx context.Context, dnsIntro k8schianetv1.ChiaDNSIntroducer) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiadnsIntroNamePattern, dnsIntro.Name),
			Namespace:       dnsIntro.Namespace,
			Labels:          kube.GetCommonLabels(ctx, dnsIntro.ObjectMeta, dnsIntro.Spec.AdditionalMetadata.Labels),
			Annotations:     dnsIntro.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, dnsIntro),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(dnsIntro.Spec.ServiceType),
			Ports: []corev1.ServicePort{
				{
					Port:       consts.DaemonPort,
					TargetPort: intstr.FromString("daemon"),
					Protocol:   "TCP",
					Name:       "daemon",
				},
				{
					Port:       53,
					TargetPort: intstr.FromString("dns"),
					Protocol:   "UDP",
					Name:       "dns",
				},
				{
					Port:       53,
					TargetPort: intstr.FromString("dns-tcp"),
					Protocol:   "TCP",
					Name:       "dns-tcp",
				},
				{
					Port:       r.getFullNodePort(ctx, dnsIntro),
					TargetPort: intstr.FromString("peers"),
					Protocol:   "TCP",
					Name:       "peers",
				},
				{
					Port:       consts.NodeRPCPort,
					TargetPort: intstr.FromString("rpc"),
					Protocol:   "TCP",
					Name:       "rpc",
				},
			},
			Selector: kube.GetCommonLabels(ctx, dnsIntro.ObjectMeta, dnsIntro.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaDNSIntroducer CR
func (r *ChiaDNSIntroducerReconciler) assembleChiaExporterService(ctx context.Context, dnsIntro k8schianetv1.ChiaDNSIntroducer) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiadnsIntroNamePattern, dnsIntro.Name) + "-metrics",
			Namespace:       dnsIntro.Namespace,
			Labels:          kube.GetCommonLabels(ctx, dnsIntro.ObjectMeta, dnsIntro.Spec.AdditionalMetadata.Labels, dnsIntro.Spec.ChiaExporterConfig.ServiceLabels),
			Annotations:     dnsIntro.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, dnsIntro),
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
			Selector: kube.GetCommonLabels(ctx, dnsIntro.ObjectMeta, dnsIntro.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleDeployment assembles the Deployment resource for a ChiaDNSIntroducer CR
func (r *ChiaDNSIntroducerReconciler) assembleDeployment(ctx context.Context, dnsIntro k8schianetv1.ChiaDNSIntroducer) appsv1.Deployment {
	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiadnsIntroNamePattern, dnsIntro.Name),
			Namespace:       dnsIntro.Namespace,
			Labels:          kube.GetCommonLabels(ctx, dnsIntro.ObjectMeta, dnsIntro.Spec.AdditionalMetadata.Labels),
			Annotations:     dnsIntro.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, dnsIntro),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(ctx, dnsIntro.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(ctx, dnsIntro.ObjectMeta, dnsIntro.Spec.AdditionalMetadata.Labels),
					Annotations: dnsIntro.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					// TODO add: imagePullSecret, serviceAccountName config
					Containers: []corev1.Container{
						{
							Name:            "chia",
							Image:           dnsIntro.Spec.ChiaConfig.Image,
							ImagePullPolicy: dnsIntro.Spec.ImagePullPolicy,
							Env:             r.getChiaEnv(ctx, dnsIntro),
							Ports: []corev1.ContainerPort{
								{
									Name:          "daemon",
									ContainerPort: consts.DaemonPort,
									Protocol:      "TCP",
								},
								{
									Name:          "dns",
									ContainerPort: 53,
									Protocol:      "UDP",
								},
								{
									Name:          "dns-tcp",
									ContainerPort: 53,
									Protocol:      "TCP",
								},
								{
									Name:          "peers",
									ContainerPort: r.getFullNodePort(ctx, dnsIntro),
									Protocol:      "TCP",
								},
								{
									Name:          "rpc",
									ContainerPort: consts.NodeRPCPort,
									Protocol:      "TCP",
								},
							},
							VolumeMounts: r.getChiaVolumeMounts(ctx, dnsIntro),
						},
					},
					NodeSelector: dnsIntro.Spec.NodeSelector,
					Volumes:      r.getChiaVolumes(ctx, dnsIntro),
				},
			},
		},
	}

	var containerSecurityContext *corev1.SecurityContext
	if dnsIntro.Spec.ChiaConfig.SecurityContext != nil {
		containerSecurityContext = dnsIntro.Spec.ChiaConfig.SecurityContext
		deploy.Spec.Template.Spec.Containers[0].SecurityContext = dnsIntro.Spec.ChiaConfig.SecurityContext
	}

	if dnsIntro.Spec.ChiaConfig.LivenessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].LivenessProbe = dnsIntro.Spec.ChiaConfig.LivenessProbe
	}

	if dnsIntro.Spec.ChiaConfig.ReadinessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].ReadinessProbe = dnsIntro.Spec.ChiaConfig.ReadinessProbe
	}

	if dnsIntro.Spec.ChiaConfig.StartupProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].StartupProbe = dnsIntro.Spec.ChiaConfig.StartupProbe
	}

	var containerResorces corev1.ResourceRequirements
	if dnsIntro.Spec.ChiaConfig.Resources != nil {
		containerResorces = *dnsIntro.Spec.ChiaConfig.Resources
		deploy.Spec.Template.Spec.Containers[0].Resources = *dnsIntro.Spec.ChiaConfig.Resources
	}

	if dnsIntro.Spec.ChiaExporterConfig.Enabled {
		exporterContainer := kube.GetChiaExporterContainer(ctx, dnsIntro.Spec.ChiaExporterConfig.Image, containerSecurityContext, dnsIntro.Spec.ImagePullPolicy, containerResorces)
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, exporterContainer)
	}

	if dnsIntro.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = dnsIntro.Spec.PodSecurityContext
	}

	// TODO add pod affinity, tolerations

	return deploy
}
