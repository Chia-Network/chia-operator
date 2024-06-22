/*
Copyright 2024 Chia Network Inc.
*/

package chiacrawler

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

const chiacrawlerNamePattern = "%s-crawler"

// assemblePeerService assembles the peer Service resource for a ChiaCrawler CR
func (r *ChiaCrawlerReconciler) assemblePeerService(ctx context.Context, crawler k8schianetv1.ChiaCrawler) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiacrawlerNamePattern, crawler.Name)
	inputs.Namespace = crawler.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, crawler)

	// Service Type
	if crawler.Spec.ChiaConfig.PeerService != nil && crawler.Spec.ChiaConfig.PeerService.ServiceType != nil {
		inputs.ServiceType = *crawler.Spec.ChiaConfig.PeerService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if crawler.Spec.ChiaConfig.PeerService != nil && crawler.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = crawler.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if crawler.Spec.ChiaConfig.PeerService != nil && crawler.Spec.ChiaConfig.PeerService.Annotations != nil {
		additionalServiceAnnotations = crawler.Spec.ChiaConfig.PeerService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(crawler.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       r.getFullNodePort(ctx, crawler),
			TargetPort: intstr.FromString("peers"),
			Protocol:   "TCP",
			Name:       "peers",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaCrawler CR
func (r *ChiaCrawlerReconciler) assembleDaemonService(ctx context.Context, crawler k8schianetv1.ChiaCrawler) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiacrawlerNamePattern, crawler.Name) + "-daemon"
	inputs.Namespace = crawler.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, crawler)

	// Service Type
	if crawler.Spec.ChiaConfig.DaemonService != nil && crawler.Spec.ChiaConfig.DaemonService.ServiceType != nil {
		inputs.ServiceType = *crawler.Spec.ChiaConfig.DaemonService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if crawler.Spec.ChiaConfig.DaemonService != nil && crawler.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = crawler.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if crawler.Spec.ChiaConfig.DaemonService != nil && crawler.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = crawler.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(crawler.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

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

// assembleRPCService assembles the RPC Service resource for a ChiaCrawler CR
func (r *ChiaCrawlerReconciler) assembleRPCService(ctx context.Context, crawler k8schianetv1.ChiaCrawler) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiacrawlerNamePattern, crawler.Name) + "-rpc"
	inputs.Namespace = crawler.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, crawler)

	// Service Type
	if crawler.Spec.ChiaConfig.RPCService != nil && crawler.Spec.ChiaConfig.RPCService.ServiceType != nil {
		inputs.ServiceType = *crawler.Spec.ChiaConfig.RPCService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if crawler.Spec.ChiaConfig.RPCService != nil && crawler.Spec.ChiaConfig.RPCService.Labels != nil {
		additionalServiceLabels = crawler.Spec.ChiaConfig.RPCService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if crawler.Spec.ChiaConfig.RPCService != nil && crawler.Spec.ChiaConfig.RPCService.Annotations != nil {
		additionalServiceAnnotations = crawler.Spec.ChiaConfig.RPCService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(crawler.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

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

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaCrawler CR
func (r *ChiaCrawlerReconciler) assembleChiaExporterService(ctx context.Context, crawler k8schianetv1.ChiaCrawler) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiacrawlerNamePattern, crawler.Name) + "-metrics"
	inputs.Namespace = crawler.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, crawler)

	// Service Type
	if crawler.Spec.ChiaExporterConfig.Service != nil && crawler.Spec.ChiaExporterConfig.Service.ServiceType != nil {
		inputs.ServiceType = *crawler.Spec.ChiaExporterConfig.Service.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if crawler.Spec.ChiaExporterConfig.Service != nil && crawler.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = crawler.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if crawler.Spec.ChiaExporterConfig.Service != nil && crawler.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = crawler.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(crawler.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

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

// assembleVolumeClaim assembles the PVC resource for a ChiaCrawler CR
func (r *ChiaCrawlerReconciler) assembleVolumeClaim(ctx context.Context, crawler k8schianetv1.ChiaCrawler) (corev1.PersistentVolumeClaim, error) {
	resourceReq, err := resource.ParseQuantity(crawler.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return corev1.PersistentVolumeClaim{}, err
	}

	var accessModes []corev1.PersistentVolumeAccessMode
	if len(crawler.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = crawler.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes
	} else {
		accessModes = []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	}

	return corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(chiacrawlerNamePattern, crawler.Name),
			Namespace: crawler.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      accessModes,
			StorageClassName: &crawler.Spec.Storage.ChiaRoot.PersistentVolumeClaim.StorageClass,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resourceReq,
				},
			},
		},
	}, nil
}

// assembleDeployment assembles the Deployment resource for a ChiaCrawler CR
func (r *ChiaCrawlerReconciler) assembleDeployment(ctx context.Context, crawler k8schianetv1.ChiaCrawler) appsv1.Deployment {
	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiacrawlerNamePattern, crawler.Name),
			Namespace:       crawler.Namespace,
			Labels:          kube.GetCommonLabels(ctx, crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels),
			Annotations:     crawler.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, crawler),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(ctx, crawler.Kind, crawler.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(ctx, crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels),
					Annotations: crawler.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					// TODO add: imagePullSecret, serviceAccountName config
					Containers: []corev1.Container{
						{
							Name:            "chia",
							Image:           crawler.Spec.ChiaConfig.Image,
							ImagePullPolicy: crawler.Spec.ImagePullPolicy,
							Env:             r.getChiaEnv(ctx, crawler),
							Ports: []corev1.ContainerPort{
								{
									Name:          "daemon",
									ContainerPort: consts.DaemonPort,
									Protocol:      "TCP",
								},
								{
									Name:          "peers",
									ContainerPort: r.getFullNodePort(ctx, crawler),
									Protocol:      "TCP",
								},
								{
									Name:          "rpc",
									ContainerPort: consts.CrawlerRPCPort,
									Protocol:      "TCP",
								},
							},
							VolumeMounts: r.getChiaVolumeMounts(ctx, crawler),
						},
					},
					NodeSelector: crawler.Spec.NodeSelector,
					Volumes:      r.getChiaVolumes(ctx, crawler),
				},
			},
		},
	}

	if len(crawler.Spec.InitContainers) != 0 {
		// Overwrite any volumeMounts specified in init containers. Not currently supported.
		for _, cont := range crawler.Spec.InitContainers {
			cont.Container.VolumeMounts = []corev1.VolumeMount{}

			// Share chia volume mounts if enabled
			if cont.ShareVolumeMounts {
				cont.Container.VolumeMounts = r.getChiaVolumeMounts(ctx, crawler)
			}

			// Share chia env if enabled
			if cont.ShareEnv {
				cont.Container.Env = append(cont.Container.Env, r.getChiaEnv(ctx, crawler)...)
			}

			deploy.Spec.Template.Spec.InitContainers = append(deploy.Spec.Template.Spec.InitContainers, cont.Container)
		}
	}

	if crawler.Spec.Strategy != nil {
		deploy.Spec.Strategy = *crawler.Spec.Strategy
	}

	var containerSecurityContext *corev1.SecurityContext
	if crawler.Spec.ChiaConfig.SecurityContext != nil {
		containerSecurityContext = crawler.Spec.ChiaConfig.SecurityContext
		deploy.Spec.Template.Spec.Containers[0].SecurityContext = crawler.Spec.ChiaConfig.SecurityContext
	}

	if crawler.Spec.ChiaConfig.LivenessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].LivenessProbe = crawler.Spec.ChiaConfig.LivenessProbe
	}

	if crawler.Spec.ChiaConfig.ReadinessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].ReadinessProbe = crawler.Spec.ChiaConfig.ReadinessProbe
	}

	if crawler.Spec.ChiaConfig.StartupProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].StartupProbe = crawler.Spec.ChiaConfig.StartupProbe
	}

	var containerResorces corev1.ResourceRequirements
	if crawler.Spec.ChiaConfig.Resources != nil {
		containerResorces = *crawler.Spec.ChiaConfig.Resources
		deploy.Spec.Template.Spec.Containers[0].Resources = *crawler.Spec.ChiaConfig.Resources
	}

	if crawler.Spec.ChiaExporterConfig.Enabled {
		exporterContainer := kube.GetChiaExporterContainer(ctx, crawler.Spec.ChiaExporterConfig.Image, containerSecurityContext, crawler.Spec.ImagePullPolicy, containerResorces)
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, exporterContainer)
	}

	if crawler.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = crawler.Spec.PodSecurityContext
	}

	if len(crawler.Spec.Sidecars.Containers) > 0 {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, crawler.Spec.Sidecars.Containers...)
	}

	// TODO add pod affinity, tolerations

	return deploy
}
