/*
Copyright 2023 Chia Network Inc.
*/

package chiaharvester

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/resource"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
)

const chiaharvesterNamePattern = "%s-harvester"

// assemblePeerService assembles the peer Service resource for a ChiaHarvester CR
func (r *ChiaHarvesterReconciler) assemblePeerService(ctx context.Context, harvester k8schianetv1.ChiaHarvester) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiaharvesterNamePattern, harvester.Name)
	inputs.Namespace = harvester.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, harvester)

	// Service Type
	if harvester.Spec.ChiaConfig.PeerService != nil {
		inputs.ServiceType = harvester.Spec.ChiaConfig.PeerService.ServiceType
		inputs.IPFamilyPolicy = harvester.Spec.ChiaConfig.PeerService.IPFamilyPolicy
		inputs.IPFamilies = harvester.Spec.ChiaConfig.PeerService.IPFamilies
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if harvester.Spec.ChiaConfig.PeerService != nil && harvester.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = harvester.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if harvester.Spec.ChiaConfig.PeerService != nil && harvester.Spec.ChiaConfig.PeerService.Annotations != nil {
		additionalServiceAnnotations = harvester.Spec.ChiaConfig.PeerService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(harvester.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       consts.HarvesterPort,
			TargetPort: intstr.FromString("peers"),
			Protocol:   "TCP",
			Name:       "peers",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaHarvester CR
func (r *ChiaHarvesterReconciler) assembleDaemonService(ctx context.Context, harvester k8schianetv1.ChiaHarvester) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiaharvesterNamePattern, harvester.Name) + "-daemon"
	inputs.Namespace = harvester.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, harvester)

	// Service Type
	if harvester.Spec.ChiaConfig.DaemonService != nil {
		inputs.ServiceType = harvester.Spec.ChiaConfig.DaemonService.ServiceType
		inputs.IPFamilyPolicy = harvester.Spec.ChiaConfig.DaemonService.IPFamilyPolicy
		inputs.IPFamilies = harvester.Spec.ChiaConfig.DaemonService.IPFamilies
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if harvester.Spec.ChiaConfig.DaemonService != nil && harvester.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = harvester.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if harvester.Spec.ChiaConfig.DaemonService != nil && harvester.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = harvester.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(harvester.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

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

// assembleRPCService assembles the RPC Service resource for a ChiaHarvester CR
func (r *ChiaHarvesterReconciler) assembleRPCService(ctx context.Context, harvester k8schianetv1.ChiaHarvester) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiaharvesterNamePattern, harvester.Name) + "-rpc"
	inputs.Namespace = harvester.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, harvester)

	// Service Type
	if harvester.Spec.ChiaConfig.RPCService != nil {
		inputs.ServiceType = harvester.Spec.ChiaConfig.RPCService.ServiceType
		inputs.IPFamilyPolicy = harvester.Spec.ChiaConfig.RPCService.IPFamilyPolicy
		inputs.IPFamilies = harvester.Spec.ChiaConfig.RPCService.IPFamilies
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if harvester.Spec.ChiaConfig.RPCService != nil && harvester.Spec.ChiaConfig.RPCService.Labels != nil {
		additionalServiceLabels = harvester.Spec.ChiaConfig.RPCService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if harvester.Spec.ChiaConfig.RPCService != nil && harvester.Spec.ChiaConfig.RPCService.Annotations != nil {
		additionalServiceAnnotations = harvester.Spec.ChiaConfig.RPCService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(harvester.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       consts.HarvesterRPCPort,
			TargetPort: intstr.FromString("rpc"),
			Protocol:   "TCP",
			Name:       "rpc",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaHarvester CR
func (r *ChiaHarvesterReconciler) assembleChiaExporterService(ctx context.Context, harvester k8schianetv1.ChiaHarvester) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiaharvesterNamePattern, harvester.Name) + "-metrics"
	inputs.Namespace = harvester.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, harvester)

	// Service Type
	if harvester.Spec.ChiaExporterConfig.Service != nil {
		inputs.ServiceType = harvester.Spec.ChiaExporterConfig.Service.ServiceType
		inputs.IPFamilyPolicy = harvester.Spec.ChiaExporterConfig.Service.IPFamilyPolicy
		inputs.IPFamilies = harvester.Spec.ChiaExporterConfig.Service.IPFamilies
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if harvester.Spec.ChiaExporterConfig.Service != nil && harvester.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = harvester.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if harvester.Spec.ChiaExporterConfig.Service != nil && harvester.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = harvester.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(harvester.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

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

// assembleVolumeClaim assembles the PVC resource for a ChiaHarvester CR
func (r *ChiaHarvesterReconciler) assembleVolumeClaim(ctx context.Context, harvester k8schianetv1.ChiaHarvester) (corev1.PersistentVolumeClaim, error) {
	resourceReq, err := resource.ParseQuantity(harvester.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return corev1.PersistentVolumeClaim{}, err
	}

	var accessModes []corev1.PersistentVolumeAccessMode
	if len(harvester.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = harvester.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes
	} else {
		accessModes = []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	}

	return corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(chiaharvesterNamePattern, harvester.Name),
			Namespace: harvester.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      accessModes,
			StorageClassName: &harvester.Spec.Storage.ChiaRoot.PersistentVolumeClaim.StorageClass,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resourceReq,
				},
			},
		},
	}, nil
}

// assembleDeployment assembles the harvester Deployment resource for a ChiaHarvester CR
func (r *ChiaHarvesterReconciler) assembleDeployment(ctx context.Context, harvester k8schianetv1.ChiaHarvester) appsv1.Deployment {
	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(chiaharvesterNamePattern, harvester.Name),
			Namespace:   harvester.Namespace,
			Labels:      kube.GetCommonLabels(ctx, harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels),
			Annotations: harvester.Spec.AdditionalMetadata.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(ctx, harvester.Kind, harvester.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(ctx, harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels),
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

	if len(harvester.Spec.InitContainers) != 0 {
		// Overwrite any volumeMounts specified in init containers. Not currently supported.
		for _, cont := range harvester.Spec.InitContainers {
			cont.Container.VolumeMounts = []corev1.VolumeMount{}

			// Share chia volume mounts if enabled
			if cont.ShareVolumeMounts {
				cont.Container.VolumeMounts = r.getChiaVolumeMounts(ctx, harvester)
			}

			// Share chia env if enabled
			if cont.ShareEnv {
				cont.Container.Env = append(cont.Container.Env, r.getChiaEnv(ctx, harvester)...)
			}

			deploy.Spec.Template.Spec.InitContainers = append(deploy.Spec.Template.Spec.InitContainers, cont.Container)
		}
	}

	if harvester.Spec.Strategy != nil {
		deploy.Spec.Strategy = *harvester.Spec.Strategy
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
		exporterContainer := kube.AssembleChiaExporterContainer(kube.AssembleChiaExporterContainerInputs{
			Image:                harvester.Spec.ChiaExporterConfig.Image,
			ConfigSecretName:     harvester.Spec.ChiaExporterConfig.ConfigSecretName,
			SecurityContext:      containerSecurityContext,
			PullPolicy:           harvester.Spec.ImagePullPolicy,
			ResourceRequirements: containerResorces,
		})
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, exporterContainer)
	}

	if harvester.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = harvester.Spec.PodSecurityContext
	}

	if len(harvester.Spec.Sidecars.Containers) > 0 {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, harvester.Spec.Sidecars.Containers...)
	}

	// TODO add pod affinity, tolerations

	return deploy
}
