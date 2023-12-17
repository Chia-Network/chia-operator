/*
Copyright 2023 Chia Network Inc.
*/

package chiaharvester

import (
	"context"
	"fmt"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const chiaharvesterNamePattern = "%s-harvester"

// assembleBaseService reconciles the main Service resource for a ChiaHarvester CR
func (r *ChiaHarvesterReconciler) assembleBaseService(ctx context.Context, harvester k8schianetv1.ChiaHarvester) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiaharvesterNamePattern, harvester.Name),
			Namespace:       harvester.Namespace,
			Labels:          r.getLabels(ctx, harvester, harvester.Spec.AdditionalMetadata.Labels),
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
			Selector: r.getLabels(ctx, harvester, harvester.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaHarvester CR
func (r *ChiaHarvesterReconciler) assembleChiaExporterService(ctx context.Context, harvester k8schianetv1.ChiaHarvester) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiaharvesterNamePattern, harvester.Name) + "-metrics",
			Namespace:       harvester.Namespace,
			Labels:          r.getLabels(ctx, harvester, harvester.Spec.AdditionalMetadata.Labels, harvester.Spec.ChiaExporterConfig.ServiceLabels),
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
			Selector: r.getLabels(ctx, harvester, harvester.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleDeployment assembles the harvester Deployment resource for a ChiaHarvester CR
func (r *ChiaHarvesterReconciler) assembleDeployment(ctx context.Context, harvester k8schianetv1.ChiaHarvester) appsv1.Deployment {
	var chiaSecContext *corev1.SecurityContext
	if harvester.Spec.ChiaConfig.SecurityContext != nil {
		chiaSecContext = harvester.Spec.ChiaConfig.SecurityContext
	}

	var chiaLivenessProbe *corev1.Probe
	if harvester.Spec.ChiaConfig.LivenessProbe != nil {
		chiaLivenessProbe = harvester.Spec.ChiaConfig.LivenessProbe
	}

	var chiaReadinessProbe *corev1.Probe
	if harvester.Spec.ChiaConfig.ReadinessProbe != nil {
		chiaReadinessProbe = harvester.Spec.ChiaConfig.ReadinessProbe
	}

	var chiaStartupProbe *corev1.Probe
	if harvester.Spec.ChiaConfig.StartupProbe != nil {
		chiaStartupProbe = harvester.Spec.ChiaConfig.StartupProbe
	}

	var chiaResources corev1.ResourceRequirements
	if harvester.Spec.ChiaConfig.Resources != nil {
		chiaResources = *harvester.Spec.ChiaConfig.Resources
	}

	var chiaExporterImage = harvester.Spec.ChiaExporterConfig.Image
	if chiaExporterImage == "" {
		chiaExporterImage = consts.DefaultChiaExporterImage
	}

	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiaharvesterNamePattern, harvester.Name),
			Namespace:       harvester.Namespace,
			Labels:          r.getLabels(ctx, harvester, harvester.Spec.AdditionalMetadata.Labels),
			Annotations:     harvester.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, harvester),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: r.getLabels(ctx, harvester),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      r.getLabels(ctx, harvester, harvester.Spec.AdditionalMetadata.Labels),
					Annotations: harvester.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					// TODO add: imagePullSecret, serviceAccountName config
					Containers: []corev1.Container{
						{
							Name:            "chia",
							SecurityContext: chiaSecContext,
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
							LivenessProbe:  chiaLivenessProbe,
							ReadinessProbe: chiaReadinessProbe,
							StartupProbe:   chiaStartupProbe,
							Resources:      chiaResources,
							VolumeMounts:   r.getChiaVolumeMounts(ctx, harvester),
						},
					},
					NodeSelector: harvester.Spec.NodeSelector,
					Volumes:      r.getChiaVolumes(ctx, harvester),
				},
			},
		},
	}

	exporterContainer := kube.GetChiaExporterContainer(ctx, chiaExporterImage, chiaSecContext, harvester.Spec.ImagePullPolicy, chiaResources)
	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, exporterContainer)

	if harvester.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = harvester.Spec.PodSecurityContext
	}

	// TODO add pod affinity, tolerations

	return deploy
}
