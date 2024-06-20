/*
Copyright 2023 Chia Network Inc.
*/

package chiaseeder

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

const chiaseederNamePattern = "%s-seeder"

// assemblePeerService assembles the peer Service resource for a ChiaSeeder CR
func (r *ChiaSeederReconciler) assemblePeerService(ctx context.Context, seeder k8schianetv1.ChiaSeeder) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiaseederNamePattern, seeder.Name)
	inputs.Namespace = seeder.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, seeder)

	// Service Type
	if seeder.Spec.ChiaConfig.PeerService != nil && seeder.Spec.ChiaConfig.PeerService.ServiceType != nil {
		inputs.ServiceType = *seeder.Spec.ChiaConfig.PeerService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if seeder.Spec.ChiaConfig.PeerService != nil && seeder.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = seeder.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if seeder.Spec.ChiaConfig.PeerService != nil && seeder.Spec.ChiaConfig.PeerService.Annotations != nil {
		additionalServiceAnnotations = seeder.Spec.ChiaConfig.PeerService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(seeder.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
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
			Port:       r.getFullNodePort(ctx, seeder),
			TargetPort: intstr.FromString("peers"),
			Protocol:   "TCP",
			Name:       "peers",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaSeeder CR
func (r *ChiaSeederReconciler) assembleDaemonService(ctx context.Context, seeder k8schianetv1.ChiaSeeder) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiaseederNamePattern, seeder.Name) + "-daemon"
	inputs.Namespace = seeder.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, seeder)

	// Service Type
	if seeder.Spec.ChiaConfig.DaemonService != nil && seeder.Spec.ChiaConfig.DaemonService.ServiceType != nil {
		inputs.ServiceType = *seeder.Spec.ChiaConfig.DaemonService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if seeder.Spec.ChiaConfig.DaemonService != nil && seeder.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = seeder.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if seeder.Spec.ChiaConfig.DaemonService != nil && seeder.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = seeder.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(seeder.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

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

// assembleRPCService assembles the RPC Service resource for a ChiaSeeder CR
func (r *ChiaSeederReconciler) assembleRPCService(ctx context.Context, seeder k8schianetv1.ChiaSeeder) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiaseederNamePattern, seeder.Name) + "-rpc"
	inputs.Namespace = seeder.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, seeder)

	// Service Type
	if seeder.Spec.ChiaConfig.RPCService != nil && seeder.Spec.ChiaConfig.RPCService.ServiceType != nil {
		inputs.ServiceType = *seeder.Spec.ChiaConfig.RPCService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if seeder.Spec.ChiaConfig.RPCService != nil && seeder.Spec.ChiaConfig.RPCService.Labels != nil {
		additionalServiceLabels = seeder.Spec.ChiaConfig.RPCService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if seeder.Spec.ChiaConfig.RPCService != nil && seeder.Spec.ChiaConfig.RPCService.Annotations != nil {
		additionalServiceAnnotations = seeder.Spec.ChiaConfig.RPCService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(seeder.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       consts.CrawlerRPCPort,
			TargetPort: intstr.FromString("rpc"),
			Protocol:   "TCP",
			Name:       "rpc",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaSeeder CR
func (r *ChiaSeederReconciler) assembleChiaExporterService(ctx context.Context, seeder k8schianetv1.ChiaSeeder) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiaseederNamePattern, seeder.Name) + "-metrics"
	inputs.Namespace = seeder.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, seeder)

	// Service Type
	if seeder.Spec.ChiaExporterConfig.Service != nil && seeder.Spec.ChiaExporterConfig.Service.ServiceType != nil {
		inputs.ServiceType = *seeder.Spec.ChiaExporterConfig.Service.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if seeder.Spec.ChiaExporterConfig.Service != nil && seeder.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = seeder.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if seeder.Spec.ChiaExporterConfig.Service != nil && seeder.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = seeder.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(seeder.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

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

// assembleDeployment assembles the Deployment resource for a ChiaSeeder CR
func (r *ChiaSeederReconciler) assembleDeployment(ctx context.Context, seeder k8schianetv1.ChiaSeeder) appsv1.Deployment {
	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiaseederNamePattern, seeder.Name),
			Namespace:       seeder.Namespace,
			Labels:          kube.GetCommonLabels(ctx, seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels),
			Annotations:     seeder.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, seeder),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(ctx, seeder.Kind, seeder.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(ctx, seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels),
					Annotations: seeder.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					// TODO add: imagePullSecret, serviceAccountName config
					Containers: []corev1.Container{
						{
							Name:            "chia",
							Image:           seeder.Spec.ChiaConfig.Image,
							ImagePullPolicy: seeder.Spec.ImagePullPolicy,
							Env:             r.getChiaEnv(ctx, seeder),
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
									ContainerPort: r.getFullNodePort(ctx, seeder),
									Protocol:      "TCP",
								},
								{
									Name:          "rpc",
									ContainerPort: consts.CrawlerRPCPort,
									Protocol:      "TCP",
								},
							},
							VolumeMounts: r.getChiaVolumeMounts(ctx, seeder),
							SecurityContext: &corev1.SecurityContext{
								Capabilities: &corev1.Capabilities{
									Add: []corev1.Capability{
										"NET_BIND_SERVICE",
									},
								},
							},
						},
					},
					NodeSelector: seeder.Spec.NodeSelector,
					Volumes:      r.getChiaVolumes(ctx, seeder),
				},
			},
		},
	}

	if seeder.Spec.Strategy != nil {
		deploy.Spec.Strategy = *seeder.Spec.Strategy
	}

	var containerSecurityContext *corev1.SecurityContext
	if seeder.Spec.ChiaConfig.SecurityContext != nil {
		containerSecurityContext = seeder.Spec.ChiaConfig.SecurityContext
		deploy.Spec.Template.Spec.Containers[0].SecurityContext = seeder.Spec.ChiaConfig.SecurityContext
	}

	if seeder.Spec.ChiaConfig.LivenessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].LivenessProbe = seeder.Spec.ChiaConfig.LivenessProbe
	}

	if seeder.Spec.ChiaConfig.ReadinessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].ReadinessProbe = seeder.Spec.ChiaConfig.ReadinessProbe
	}

	if seeder.Spec.ChiaConfig.StartupProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].StartupProbe = seeder.Spec.ChiaConfig.StartupProbe
	}

	var containerResorces corev1.ResourceRequirements
	if seeder.Spec.ChiaConfig.Resources != nil {
		containerResorces = *seeder.Spec.ChiaConfig.Resources
		deploy.Spec.Template.Spec.Containers[0].Resources = *seeder.Spec.ChiaConfig.Resources
	}

	if seeder.Spec.ChiaExporterConfig.Enabled {
		exporterContainer := kube.GetChiaExporterContainer(ctx, seeder.Spec.ChiaExporterConfig.Image, containerSecurityContext, seeder.Spec.ImagePullPolicy, containerResorces)
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, exporterContainer)
	}

	if seeder.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = seeder.Spec.PodSecurityContext
	}

	if len(seeder.Spec.Sidecars.Containers) > 0 {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, seeder.Spec.Sidecars.Containers...)
	}

	// TODO add pod affinity, tolerations

	return deploy
}
