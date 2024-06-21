/*
Copyright 2023 Chia Network Inc.
*/

package chiafarmer

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

const chiafarmerNamePattern = "%s-farmer"

// assemblePeerService assembles the peer Service resource for a Chiafarmer CR
func (r *ChiaFarmerReconciler) assemblePeerService(ctx context.Context, farmer k8schianetv1.ChiaFarmer) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiafarmerNamePattern, farmer.Name)
	inputs.Namespace = farmer.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, farmer)

	// Service Type
	if farmer.Spec.ChiaConfig.PeerService != nil && farmer.Spec.ChiaConfig.PeerService.ServiceType != nil {
		inputs.ServiceType = *farmer.Spec.ChiaConfig.PeerService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if farmer.Spec.ChiaConfig.PeerService != nil && farmer.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = farmer.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if farmer.Spec.ChiaConfig.PeerService != nil && farmer.Spec.ChiaConfig.PeerService.Annotations != nil {
		additionalServiceAnnotations = farmer.Spec.ChiaConfig.PeerService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(farmer.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       consts.FarmerPort,
			TargetPort: intstr.FromString("peers"),
			Protocol:   "TCP",
			Name:       "peers",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a Chiafarmer CR
func (r *ChiaFarmerReconciler) assembleDaemonService(ctx context.Context, farmer k8schianetv1.ChiaFarmer) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiafarmerNamePattern, farmer.Name) + "-daemon"
	inputs.Namespace = farmer.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, farmer)

	// Service Type
	if farmer.Spec.ChiaConfig.DaemonService != nil && farmer.Spec.ChiaConfig.DaemonService.ServiceType != nil {
		inputs.ServiceType = *farmer.Spec.ChiaConfig.DaemonService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if farmer.Spec.ChiaConfig.DaemonService != nil && farmer.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = farmer.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if farmer.Spec.ChiaConfig.DaemonService != nil && farmer.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = farmer.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(farmer.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

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

// assembleRPCService assembles the RPC Service resource for a Chiafarmer CR
func (r *ChiaFarmerReconciler) assembleRPCService(ctx context.Context, farmer k8schianetv1.ChiaFarmer) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiafarmerNamePattern, farmer.Name) + "-rpc"
	inputs.Namespace = farmer.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, farmer)

	// Service Type
	if farmer.Spec.ChiaConfig.RPCService != nil && farmer.Spec.ChiaConfig.RPCService.ServiceType != nil {
		inputs.ServiceType = *farmer.Spec.ChiaConfig.RPCService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if farmer.Spec.ChiaConfig.RPCService != nil && farmer.Spec.ChiaConfig.RPCService.Labels != nil {
		additionalServiceLabels = farmer.Spec.ChiaConfig.RPCService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if farmer.Spec.ChiaConfig.RPCService != nil && farmer.Spec.ChiaConfig.RPCService.Annotations != nil {
		additionalServiceAnnotations = farmer.Spec.ChiaConfig.RPCService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(farmer.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       consts.FarmerRPCPort,
			TargetPort: intstr.FromString("rpc"),
			Protocol:   "TCP",
			Name:       "rpc",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a Chiafarmer CR
func (r *ChiaFarmerReconciler) assembleChiaExporterService(ctx context.Context, farmer k8schianetv1.ChiaFarmer) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiafarmerNamePattern, farmer.Name) + "-metrics"
	inputs.Namespace = farmer.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, farmer)

	// Service Type
	if farmer.Spec.ChiaExporterConfig.Service != nil && farmer.Spec.ChiaExporterConfig.Service.ServiceType != nil {
		inputs.ServiceType = *farmer.Spec.ChiaExporterConfig.Service.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if farmer.Spec.ChiaExporterConfig.Service != nil && farmer.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = farmer.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if farmer.Spec.ChiaExporterConfig.Service != nil && farmer.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = farmer.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(farmer.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

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

// assembleDeployment assembles the farmer Deployment resource for a ChiaFarmer CR
func (r *ChiaFarmerReconciler) assembleDeployment(ctx context.Context, farmer k8schianetv1.ChiaFarmer) appsv1.Deployment {
	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiafarmerNamePattern, farmer.Name),
			Namespace:       farmer.Namespace,
			Labels:          kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels),
			Annotations:     farmer.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, farmer),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels),
					Annotations: farmer.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					// TODO add: imagePullSecret, serviceAccountName config
					Containers: []corev1.Container{
						{
							Name:            "chia",
							Image:           farmer.Spec.ChiaConfig.Image,
							ImagePullPolicy: farmer.Spec.ImagePullPolicy,
							Env:             r.getChiaEnv(ctx, farmer),
							Ports: []corev1.ContainerPort{
								{
									Name:          "daemon",
									ContainerPort: consts.DaemonPort,
									Protocol:      "TCP",
								},
								{
									Name:          "peers",
									ContainerPort: consts.FarmerPort,
									Protocol:      "TCP",
								},
								{
									Name:          "rpc",
									ContainerPort: consts.FarmerRPCPort,
									Protocol:      "TCP",
								},
							},
							VolumeMounts: r.getChiaVolumeMounts(ctx, farmer),
						},
					},
					NodeSelector: farmer.Spec.NodeSelector,
					Volumes:      r.getChiaVolumes(ctx, farmer),
				},
			},
		},
	}

	if len(farmer.Spec.InitContainers) != 0 {
		// Overwrite any volumeMounts specified in init containers. Not currently supported.
		for _, cont := range farmer.Spec.InitContainers {
			cont.Container.VolumeMounts = []corev1.VolumeMount{}

			// Share chia volume mounts if enabled
			if cont.ShareVolumeMounts {
				cont.Container.VolumeMounts = r.getChiaVolumeMounts(ctx, farmer)
			}

			// Share chia env if enabled
			if cont.ShareEnv {
				cont.Container.Env = append(cont.Container.Env, r.getChiaEnv(ctx, farmer)...)
			}

			deploy.Spec.Template.Spec.InitContainers = append(deploy.Spec.Template.Spec.InitContainers, cont.Container)
		}
	}

	var containerSecurityContext *corev1.SecurityContext
	if farmer.Spec.ChiaConfig.SecurityContext != nil {
		containerSecurityContext = farmer.Spec.ChiaConfig.SecurityContext
		deploy.Spec.Template.Spec.Containers[0].SecurityContext = farmer.Spec.ChiaConfig.SecurityContext
	}

	if farmer.Spec.ChiaConfig.LivenessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].LivenessProbe = farmer.Spec.ChiaConfig.LivenessProbe
	}

	if farmer.Spec.ChiaConfig.ReadinessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].ReadinessProbe = farmer.Spec.ChiaConfig.ReadinessProbe
	}

	if farmer.Spec.ChiaConfig.StartupProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].StartupProbe = farmer.Spec.ChiaConfig.StartupProbe
	}

	var containerResorces corev1.ResourceRequirements
	if farmer.Spec.ChiaConfig.Resources != nil {
		containerResorces = *farmer.Spec.ChiaConfig.Resources
		deploy.Spec.Template.Spec.Containers[0].Resources = *farmer.Spec.ChiaConfig.Resources
	}

	if farmer.Spec.ChiaExporterConfig.Enabled {
		exporterContainer := kube.GetChiaExporterContainer(ctx, farmer.Spec.ChiaExporterConfig.Image, containerSecurityContext, farmer.Spec.ImagePullPolicy, containerResorces)
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, exporterContainer)
	}

	if farmer.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = farmer.Spec.PodSecurityContext
	}

	if len(farmer.Spec.Sidecars.Containers) > 0 {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, farmer.Spec.Sidecars.Containers...)
	}

	// TODO add pod affinity, tolerations

	return deploy
}
