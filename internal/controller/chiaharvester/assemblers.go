/*
Copyright 2023 Chia Network Inc.
*/

package chiaharvester

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

const chiaharvesterNamePattern = "%s-harvester"

// assembleBaseService reconciles the main Service resource for a ChiaHarvester CR
func (r *ChiaHarvesterReconciler) assembleBaseService(ctx context.Context, harvester k8schianetv1.ChiaHarvester) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiaharvesterNamePattern, harvester.Name),
			Namespace:       harvester.Namespace,
			Labels:          kube.GetCommonLabels(ctx, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels),
			Annotations:     harvester.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, harvester),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(harvester.Spec.ServiceType),
			Ports: []corev1.ServicePort{
				{
					Port:       consts.DaemonPort,
					TargetPort: intstr.FromString("daemon"),
					Protocol:   "TCP",
					Name:       "daemon",
				},
				{
					Port:       consts.HarvesterPort,
					TargetPort: intstr.FromString("peers"),
					Protocol:   "TCP",
					Name:       "peers",
				},
				{
					Port:       consts.HarvesterRPCPort,
					TargetPort: intstr.FromString("rpc"),
					Protocol:   "TCP",
					Name:       "rpc",
				},
			},
			Selector: kube.GetCommonLabels(ctx, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaHarvester CR
func (r *ChiaHarvesterReconciler) assembleChiaExporterService(ctx context.Context, harvester k8schianetv1.ChiaHarvester) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiaharvesterNamePattern, harvester.Name) + "-metrics",
			Namespace:       harvester.Namespace,
			Labels:          kube.GetCommonLabels(ctx, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels, harvester.Spec.ChiaExporterConfig.ServiceLabels),
			Annotations:     harvester.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, harvester),
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
			Selector: kube.GetCommonLabels(ctx, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleDeployment assembles the harvester Deployment resource for a ChiaHarvester CR
func (r *ChiaHarvesterReconciler) assembleDeployment(ctx context.Context, harvester k8schianetv1.ChiaHarvester) appsv1.Deployment {
	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiaharvesterNamePattern, harvester.Name),
			Namespace:       harvester.Namespace,
			Labels:          kube.GetCommonLabels(ctx, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels),
			Annotations:     harvester.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, harvester),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(ctx, harvester.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(ctx, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels),
					Annotations: harvester.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					// TODO add: imagePullSecret, serviceAccountName config
					Containers: []corev1.Container{
						{
							Name:            "chia",
							Image:           harvester.Spec.ChiaConfig.Image,
							ImagePullPolicy: harvester.Spec.ImagePullPolicy,
							Env:             r.getChiaEnv(ctx, harvester),
							Ports: []corev1.ContainerPort{
								{
									Name:          "daemon",
									ContainerPort: consts.DaemonPort,
									Protocol:      "TCP",
								},
								{
									Name:          "peers",
									ContainerPort: consts.HarvesterPort,
									Protocol:      "TCP",
								},
								{
									Name:          "rpc",
									ContainerPort: consts.HarvesterRPCPort,
									Protocol:      "TCP",
								},
							},
							VolumeMounts: r.getChiaVolumeMounts(ctx, harvester),
						},
					},
					NodeSelector: harvester.Spec.NodeSelector,
					Volumes:      r.getChiaVolumes(ctx, harvester),
				},
			},
		},
	}

	var containerSecurityContext *corev1.SecurityContext
	if harvester.Spec.ChiaConfig.SecurityContext != nil {
		containerSecurityContext = harvester.Spec.ChiaConfig.SecurityContext
		deploy.Spec.Template.Spec.Containers[0].SecurityContext = harvester.Spec.ChiaConfig.SecurityContext
	}

	if harvester.Spec.ChiaConfig.LivenessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].LivenessProbe = harvester.Spec.ChiaConfig.LivenessProbe
	}

	if harvester.Spec.ChiaConfig.ReadinessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].ReadinessProbe = harvester.Spec.ChiaConfig.ReadinessProbe
	}

	if harvester.Spec.ChiaConfig.StartupProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].StartupProbe = harvester.Spec.ChiaConfig.StartupProbe
	}

	var containerResorces corev1.ResourceRequirements
	if harvester.Spec.ChiaConfig.Resources != nil {
		containerResorces = *harvester.Spec.ChiaConfig.Resources
		deploy.Spec.Template.Spec.Containers[0].Resources = *harvester.Spec.ChiaConfig.Resources
	}

	if harvester.Spec.ChiaExporterConfig.Enabled {
		exporterContainer := kube.GetChiaExporterContainer(ctx, harvester.Spec.ChiaExporterConfig.Image, containerSecurityContext, harvester.Spec.ImagePullPolicy, containerResorces)
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, exporterContainer)
	}

	if harvester.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = harvester.Spec.PodSecurityContext
	}

	// TODO add pod affinity, tolerations

	return deploy
}
