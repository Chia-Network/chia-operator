/*
Copyright 2023 Chia Network Inc.
*/

package chiatimelord

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

const chiatimelordNamePattern = "%s-timelord"

// assemblePeerService assembles the peer Service resource for a ChiaTimelord CR
func (r *ChiaTimelordReconciler) assemblePeerService(ctx context.Context, tl k8schianetv1.ChiaTimelord) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiatimelordNamePattern, tl.Name)
	inputs.Namespace = tl.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, tl)

	// Service Type
	if tl.Spec.ChiaConfig.PeerService != nil && tl.Spec.ChiaConfig.PeerService.ServiceType != nil {
		inputs.ServiceType = *tl.Spec.ChiaConfig.PeerService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if tl.Spec.ChiaConfig.PeerService != nil && tl.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = tl.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if tl.Spec.ChiaConfig.PeerService != nil && tl.Spec.ChiaConfig.PeerService.Annotations != nil {
		additionalServiceAnnotations = tl.Spec.ChiaConfig.PeerService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(tl.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       consts.TimelordPort,
			TargetPort: intstr.FromString("peers"),
			Protocol:   "TCP",
			Name:       "peers",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaTimelord CR
func (r *ChiaTimelordReconciler) assembleDaemonService(ctx context.Context, tl k8schianetv1.ChiaTimelord) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiatimelordNamePattern, tl.Name) + "-daemon"
	inputs.Namespace = tl.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, tl)

	// Service Type
	if tl.Spec.ChiaConfig.DaemonService != nil && tl.Spec.ChiaConfig.DaemonService.ServiceType != nil {
		inputs.ServiceType = *tl.Spec.ChiaConfig.DaemonService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if tl.Spec.ChiaConfig.DaemonService != nil && tl.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = tl.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if tl.Spec.ChiaConfig.DaemonService != nil && tl.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = tl.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(tl.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

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

// assembleRPCService assembles the RPC Service resource for a ChiaTimelord CR
func (r *ChiaTimelordReconciler) assembleRPCService(ctx context.Context, tl k8schianetv1.ChiaTimelord) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiatimelordNamePattern, tl.Name) + "-rpc"
	inputs.Namespace = tl.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, tl)

	// Service Type
	if tl.Spec.ChiaConfig.RPCService != nil && tl.Spec.ChiaConfig.RPCService.ServiceType != nil {
		inputs.ServiceType = *tl.Spec.ChiaConfig.RPCService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if tl.Spec.ChiaConfig.RPCService != nil && tl.Spec.ChiaConfig.RPCService.Labels != nil {
		additionalServiceLabels = tl.Spec.ChiaConfig.RPCService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if tl.Spec.ChiaConfig.RPCService != nil && tl.Spec.ChiaConfig.RPCService.Annotations != nil {
		additionalServiceAnnotations = tl.Spec.ChiaConfig.RPCService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(tl.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       consts.TimelordRPCPort,
			TargetPort: intstr.FromString("rpc"),
			Protocol:   "TCP",
			Name:       "rpc",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaTimelord CR
func (r *ChiaTimelordReconciler) assembleChiaExporterService(ctx context.Context, tl k8schianetv1.ChiaTimelord) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiatimelordNamePattern, tl.Name) + "-metrics"
	inputs.Namespace = tl.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, tl)

	// Service Type
	if tl.Spec.ChiaExporterConfig.Service != nil && tl.Spec.ChiaExporterConfig.Service.ServiceType != nil {
		inputs.ServiceType = *tl.Spec.ChiaExporterConfig.Service.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if tl.Spec.ChiaExporterConfig.Service != nil && tl.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = tl.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if tl.Spec.ChiaExporterConfig.Service != nil && tl.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = tl.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(tl.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

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

// assembleVolumeClaim assembles the PVC resource for a ChiaTimelord CR
func (r *ChiaTimelordReconciler) assembleVolumeClaim(ctx context.Context, tl k8schianetv1.ChiaTimelord) (corev1.PersistentVolumeClaim, error) {
	resourceReq, err := resource.ParseQuantity(tl.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return corev1.PersistentVolumeClaim{}, err
	}

	var accessModes []corev1.PersistentVolumeAccessMode
	if len(tl.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = tl.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes
	} else {
		accessModes = []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	}

	return corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(chiatimelordNamePattern, tl.Name),
			Namespace: tl.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      accessModes,
			StorageClassName: &tl.Spec.Storage.ChiaRoot.PersistentVolumeClaim.StorageClass,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resourceReq,
				},
			},
		},
	}, nil
}

// assembleDeployment assembles the tl Deployment resource for a ChiaTimelord CR
func (r *ChiaTimelordReconciler) assembleDeployment(ctx context.Context, tl k8schianetv1.ChiaTimelord) appsv1.Deployment {
	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiatimelordNamePattern, tl.Name),
			Namespace:       tl.Namespace,
			Labels:          kube.GetCommonLabels(ctx, tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels),
			Annotations:     tl.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, tl),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(ctx, tl.Kind, tl.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(ctx, tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels),
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
							VolumeMounts: r.getChiaVolumeMounts(ctx, tl),
						},
					},
					NodeSelector: tl.Spec.NodeSelector,
					Volumes:      r.getChiaVolumes(ctx, tl),
				},
			},
		},
	}

	if len(tl.Spec.InitContainers) != 0 {
		// Overwrite any volumeMounts specified in init containers. Not currently supported.
		for _, cont := range tl.Spec.InitContainers {
			cont.Container.VolumeMounts = []corev1.VolumeMount{}

			// Share chia volume mounts if enabled
			if cont.ShareVolumeMounts {
				cont.Container.VolumeMounts = r.getChiaVolumeMounts(ctx, tl)
			}

			// Share chia env if enabled
			if cont.ShareEnv {
				cont.Container.Env = append(cont.Container.Env, r.getChiaEnv(ctx, tl)...)
			}

			deploy.Spec.Template.Spec.InitContainers = append(deploy.Spec.Template.Spec.InitContainers, cont.Container)
		}
	}

	if tl.Spec.Strategy != nil {
		deploy.Spec.Strategy = *tl.Spec.Strategy
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

	if len(tl.Spec.Sidecars.Containers) > 0 {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, tl.Spec.Sidecars.Containers...)
	}

	// TODO add pod affinity, tolerations

	return deploy
}
