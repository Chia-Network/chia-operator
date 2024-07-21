/*
Copyright 2024 Chia Network Inc.
*/

package chiaintroducer

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

const chiaintroducerNamePattern = "%s-introducer"

// assemblePeerService assembles the peer Service resource for a ChiaIntroducer CR
func (r *ChiaIntroducerReconciler) assemblePeerService(ctx context.Context, introducer k8schianetv1.ChiaIntroducer) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiaintroducerNamePattern, introducer.Name)
	inputs.Namespace = introducer.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, introducer)

	// Service Type
	if introducer.Spec.ChiaConfig.PeerService != nil {
		inputs.ServiceType = introducer.Spec.ChiaConfig.PeerService.ServiceType
		inputs.IPFamilyPolicy = introducer.Spec.ChiaConfig.PeerService.IPFamilyPolicy
		inputs.IPFamilies = introducer.Spec.ChiaConfig.PeerService.IPFamilies
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if introducer.Spec.ChiaConfig.PeerService != nil && introducer.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = introducer.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels)

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
	if introducer.Spec.ChiaConfig.DaemonService != nil {
		inputs.ServiceType = introducer.Spec.ChiaConfig.DaemonService.ServiceType
		inputs.IPFamilyPolicy = introducer.Spec.ChiaConfig.DaemonService.IPFamilyPolicy
		inputs.IPFamilies = introducer.Spec.ChiaConfig.DaemonService.IPFamilies
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if introducer.Spec.ChiaConfig.DaemonService != nil && introducer.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = introducer.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels)

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
	if introducer.Spec.ChiaExporterConfig.Service != nil {
		inputs.ServiceType = introducer.Spec.ChiaExporterConfig.Service.ServiceType
		inputs.IPFamilyPolicy = introducer.Spec.ChiaExporterConfig.Service.IPFamilyPolicy
		inputs.IPFamilies = introducer.Spec.ChiaExporterConfig.Service.IPFamilies
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if introducer.Spec.ChiaExporterConfig.Service != nil && introducer.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = introducer.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels)

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

// assembleVolumeClaim assembles the PVC resource for a ChiaIntroducer CR
func (r *ChiaIntroducerReconciler) assembleVolumeClaim(ctx context.Context, introducer k8schianetv1.ChiaIntroducer) (corev1.PersistentVolumeClaim, error) {
	resourceReq, err := resource.ParseQuantity(introducer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return corev1.PersistentVolumeClaim{}, err
	}

	var accessModes []corev1.PersistentVolumeAccessMode
	if len(introducer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = introducer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes
	} else {
		accessModes = []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	}

	return corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(chiaintroducerNamePattern, introducer.Name),
			Namespace: introducer.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      accessModes,
			StorageClassName: &introducer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.StorageClass,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resourceReq,
				},
			},
		},
	}, nil
}

// assembleDeployment assembles the Deployment resource for a ChiaIntroducer CR
func (r *ChiaIntroducerReconciler) assembleDeployment(ctx context.Context, introducer k8schianetv1.ChiaIntroducer) appsv1.Deployment {
	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(chiaintroducerNamePattern, introducer.Name),
			Namespace:   introducer.Namespace,
			Labels:      kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels),
			Annotations: introducer.Spec.AdditionalMetadata.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.AdditionalMetadata.Labels),
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

	if len(introducer.Spec.InitContainers) != 0 {
		// Overwrite any volumeMounts specified in init containers. Not currently supported.
		for _, cont := range introducer.Spec.InitContainers {
			cont.Container.VolumeMounts = []corev1.VolumeMount{}

			// Share chia volume mounts if enabled
			if cont.ShareVolumeMounts {
				cont.Container.VolumeMounts = r.getChiaVolumeMounts(ctx, introducer)
			}

			// Share chia env if enabled
			if cont.ShareEnv {
				cont.Container.Env = append(cont.Container.Env, r.getChiaEnv(ctx, introducer)...)
			}

			deploy.Spec.Template.Spec.InitContainers = append(deploy.Spec.Template.Spec.InitContainers, cont.Container)
		}
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
		exporterContainer := kube.AssembleChiaExporterContainer(kube.AssembleChiaExporterContainerInputs{
			Image:                introducer.Spec.ChiaExporterConfig.Image,
			ConfigSecretName:     introducer.Spec.ChiaExporterConfig.ConfigSecretName,
			SecurityContext:      containerSecurityContext,
			PullPolicy:           introducer.Spec.ImagePullPolicy,
			ResourceRequirements: containerResorces,
		})
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
