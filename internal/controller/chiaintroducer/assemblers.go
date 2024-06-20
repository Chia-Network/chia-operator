/*
Copyright 2024 Chia Network Inc.
*/

package chiaintroducer

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

const chiaintroducerNamePattern = "%s-introducer"

// assemblePeerService assembles the peer Service resource for a ChiaIntroducer CR
func (r *ChiaIntroducerReconciler) assemblePeerService(ctx context.Context, introducer k8schianetv1.ChiaIntroducer) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiaintroducerNamePattern, introducer.Name)
	inputs.Namespace = introducer.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, introducer)

	// Service Type
	if introducer.Spec.ChiaConfig.PeerService != nil && introducer.Spec.ChiaConfig.PeerService.ServiceType != nil {
		inputs.ServiceType = *introducer.Spec.ChiaConfig.PeerService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if introducer.Spec.ChiaConfig.PeerService != nil && introducer.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = introducer.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if introducer.Spec.ChiaConfig.PeerService != nil && introducer.Spec.ChiaConfig.PeerService.Annotations != nil {
		additionalServiceAnnotations = introducer.Spec.ChiaConfig.PeerService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(introducer.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       r.getFullNodePort(ctx, introducer),
			TargetPort: intstr.FromString("peers"),
			Protocol:   "TCP",
			Name:       "peers",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaIntroducer CR
func (r *ChiaIntroducerReconciler) assembleDaemonService(ctx context.Context, introducer k8schianetv1.ChiaIntroducer) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiaintroducerNamePattern, introducer.Name) + "-daemon"
	inputs.Namespace = introducer.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, introducer)

	// Service Type
	if introducer.Spec.ChiaConfig.DaemonService != nil && introducer.Spec.ChiaConfig.DaemonService.ServiceType != nil {
		inputs.ServiceType = *introducer.Spec.ChiaConfig.DaemonService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if introducer.Spec.ChiaConfig.DaemonService != nil && introducer.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = introducer.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if introducer.Spec.ChiaConfig.DaemonService != nil && introducer.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = introducer.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(introducer.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       consts.DaemonPort,
			TargetPort: intstr.FromString("daemon"),
			Protocol:   "TCP",
			Name:       "daemon",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaIntroducer CR
func (r *ChiaIntroducerReconciler) assembleChiaExporterService(ctx context.Context, introducer k8schianetv1.ChiaIntroducer) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiaintroducerNamePattern, introducer.Name) + "-metrics"
	inputs.Namespace = introducer.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, introducer)

	// Service Type
	if introducer.Spec.ChiaExporterConfig.Service != nil && introducer.Spec.ChiaExporterConfig.Service.ServiceType != nil {
		inputs.ServiceType = *introducer.Spec.ChiaExporterConfig.Service.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if introducer.Spec.ChiaExporterConfig.Service != nil && introducer.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = introducer.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if introducer.Spec.ChiaExporterConfig.Service != nil && introducer.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = introducer.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(introducer.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       consts.ChiaExporterPort,
			TargetPort: intstr.FromString("metrics"),
			Protocol:   "TCP",
			Name:       "metrics",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleDeployment assembles the Deployment resource for a ChiaIntroducer CR
func (r *ChiaIntroducerReconciler) assembleDeployment(ctx context.Context, introducer k8schianetv1.ChiaIntroducer) appsv1.Deployment {
	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiaintroducerNamePattern, introducer.Name),
			Namespace:       introducer.Namespace,
			Labels:          kube.GetCommonLabels(ctx, introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels),
			Annotations:     introducer.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, introducer),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(ctx, introducer.Kind, introducer.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(ctx, introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels),
					Annotations: introducer.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					// TODO add: imagePullSecret, serviceAccountName config
					Containers: []corev1.Container{
						{
							Name:            "chia",
							Image:           introducer.Spec.ChiaConfig.Image,
							ImagePullPolicy: introducer.Spec.ImagePullPolicy,
							Env:             r.getChiaEnv(ctx, introducer),
							Ports: []corev1.ContainerPort{
								{
									Name:          "daemon",
									ContainerPort: consts.DaemonPort,
									Protocol:      "TCP",
								},
								{
									Name:          "peers",
									ContainerPort: r.getFullNodePort(ctx, introducer),
									Protocol:      "TCP",
								},
							},
							VolumeMounts: r.getChiaVolumeMounts(ctx, introducer),
						},
					},
					NodeSelector: introducer.Spec.NodeSelector,
					Volumes:      r.getChiaVolumes(ctx, introducer),
				},
			},
		},
	}

	if introducer.Spec.Strategy != nil {
		deploy.Spec.Strategy = *introducer.Spec.Strategy
	}

	var containerSecurityContext *corev1.SecurityContext
	if introducer.Spec.ChiaConfig.SecurityContext != nil {
		containerSecurityContext = introducer.Spec.ChiaConfig.SecurityContext
		deploy.Spec.Template.Spec.Containers[0].SecurityContext = introducer.Spec.ChiaConfig.SecurityContext
	}

	if introducer.Spec.ChiaConfig.LivenessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].LivenessProbe = introducer.Spec.ChiaConfig.LivenessProbe
	}

	if introducer.Spec.ChiaConfig.ReadinessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].ReadinessProbe = introducer.Spec.ChiaConfig.ReadinessProbe
	}

	if introducer.Spec.ChiaConfig.StartupProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].StartupProbe = introducer.Spec.ChiaConfig.StartupProbe
	}

	var containerResorces corev1.ResourceRequirements
	if introducer.Spec.ChiaConfig.Resources != nil {
		containerResorces = *introducer.Spec.ChiaConfig.Resources
		deploy.Spec.Template.Spec.Containers[0].Resources = *introducer.Spec.ChiaConfig.Resources
	}

	if introducer.Spec.ChiaExporterConfig.Enabled {
		exporterContainer := kube.GetChiaExporterContainer(ctx, introducer.Spec.ChiaExporterConfig.Image, containerSecurityContext, introducer.Spec.ImagePullPolicy, containerResorces)
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, exporterContainer)
	}

	if introducer.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = introducer.Spec.PodSecurityContext
	}

	if len(introducer.Spec.Sidecars.Containers) > 0 {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, introducer.Spec.Sidecars.Containers...)
	}

	// TODO add pod affinity, tolerations

	return deploy
}
