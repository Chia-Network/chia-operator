/*
Copyright 2023 Chia Network Inc.
*/

package chiatimelord

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

const chiatimelordNamePattern = "%s-timelord"

// assembleBaseService assembles the main Service resource for a Chiatl CR
func (r *ChiaTimelordReconciler) assembleBaseService(ctx context.Context, tl k8schianetv1.ChiaTimelord) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiatimelordNamePattern, tl.Name),
			Namespace:       tl.Namespace,
			Labels:          kube.GetCommonLabels(ctx, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels),
			Annotations:     tl.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, tl),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(tl.Spec.ServiceType),
			Ports: []corev1.ServicePort{
				{
					Port:       consts.DaemonPort,
					TargetPort: intstr.FromString("daemon"),
					Protocol:   "TCP",
					Name:       "daemon",
				},
				{
					Port:       consts.TimelordPort,
					TargetPort: intstr.FromString("peers"),
					Protocol:   "TCP",
					Name:       "peers",
				},
				{
					Port:       consts.TimelordRPCPort,
					TargetPort: intstr.FromString("rpc"),
					Protocol:   "TCP",
					Name:       "rpc",
				},
			},
			Selector: kube.GetCommonLabels(ctx, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaTimelord CR
func (r *ChiaTimelordReconciler) assembleChiaExporterService(ctx context.Context, tl k8schianetv1.ChiaTimelord) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiatimelordNamePattern, tl.Name) + "metrics",
			Namespace:       tl.Namespace,
			Labels:          kube.GetCommonLabels(ctx, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels, tl.Spec.ChiaExporterConfig.ServiceLabels),
			Annotations:     tl.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, tl),
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
			Selector: kube.GetCommonLabels(ctx, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleDeployment assembles the tl Deployment resource for a ChiaTimelord CR
func (r *ChiaTimelordReconciler) assembleDeployment(ctx context.Context, tl k8schianetv1.ChiaTimelord) appsv1.Deployment {
	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiatimelordNamePattern, tl.Name),
			Namespace:       tl.Namespace,
			Labels:          kube.GetCommonLabels(ctx, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels),
			Annotations:     tl.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, tl),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(ctx, tl.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(ctx, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels),
					Annotations: tl.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					// TODO add: imagePullSecret, serviceAccountName config
					Containers: []corev1.Container{
						{
							Name:            "chia",
							Image:           tl.Spec.ChiaConfig.Image,
							ImagePullPolicy: tl.Spec.ImagePullPolicy,
							Env:             r.getChiaEnv(ctx, tl),
							Ports: []corev1.ContainerPort{
								{
									Name:          "daemon",
									ContainerPort: consts.DaemonPort,
									Protocol:      "TCP",
								},
								{
									Name:          "peers",
									ContainerPort: consts.TimelordPort,
									Protocol:      "TCP",
								},
								{
									Name:          "rpc",
									ContainerPort: consts.TimelordRPCPort,
									Protocol:      "TCP",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "secret-ca",
									MountPath: "/chia-ca",
								},
								{
									Name:      "chiaroot",
									MountPath: "/chia-data",
								},
							},
						},
					},
					NodeSelector: tl.Spec.NodeSelector,
					Volumes:      r.getChiaVolumes(ctx, tl),
				},
			},
		},
	}

	var containerSecurityContext *corev1.SecurityContext
	if tl.Spec.ChiaConfig.SecurityContext != nil {
		containerSecurityContext = tl.Spec.ChiaConfig.SecurityContext
		deploy.Spec.Template.Spec.Containers[0].SecurityContext = tl.Spec.ChiaConfig.SecurityContext
	}

	if tl.Spec.ChiaConfig.LivenessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].LivenessProbe = tl.Spec.ChiaConfig.LivenessProbe
	}

	if tl.Spec.ChiaConfig.ReadinessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].ReadinessProbe = tl.Spec.ChiaConfig.ReadinessProbe
	}

	if tl.Spec.ChiaConfig.StartupProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].StartupProbe = tl.Spec.ChiaConfig.StartupProbe
	}

	var containerResorces corev1.ResourceRequirements
	if tl.Spec.ChiaConfig.Resources != nil {
		containerResorces = *tl.Spec.ChiaConfig.Resources
		deploy.Spec.Template.Spec.Containers[0].Resources = *tl.Spec.ChiaConfig.Resources
	}

	if tl.Spec.ChiaExporterConfig.Enabled {
		exporterContainer := kube.GetChiaExporterContainer(ctx, tl.Spec.ChiaExporterConfig.Image, containerSecurityContext, tl.Spec.ImagePullPolicy, containerResorces)
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, exporterContainer)
	}

	if tl.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = tl.Spec.PodSecurityContext
	}

	// TODO add pod affinity, tolerations

	return deploy
}
