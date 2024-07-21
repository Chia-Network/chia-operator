/*
Copyright 2023 Chia Network Inc.
*/

package chiatimelord

import (
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
func assemblePeerService(tl k8schianetv1.ChiaTimelord) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:           fmt.Sprintf(chiatimelordNamePattern, tl.Name),
		Namespace:      tl.Namespace,
		OwnerReference: getOwnerReference(tl),
		Ports: []corev1.ServicePort{
			{
				Port:       consts.TimelordPort,
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
		},
	}

	if tl.Spec.ChiaConfig.PeerService != nil {
		inputs.ServiceType = tl.Spec.ChiaConfig.PeerService.ServiceType
		inputs.IPFamilyPolicy = tl.Spec.ChiaConfig.PeerService.IPFamilyPolicy
		inputs.IPFamilies = tl.Spec.ChiaConfig.PeerService.IPFamilies

		// Labels
		var additionalServiceLabels = make(map[string]string)
		if tl.Spec.ChiaConfig.PeerService.Labels != nil {
			additionalServiceLabels = tl.Spec.ChiaConfig.PeerService.Labels
		}
		inputs.Labels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
		inputs.SelectorLabels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels)

		// Annotations
		var additionalServiceAnnotations = make(map[string]string)
		if tl.Spec.ChiaConfig.PeerService.Annotations != nil {
			additionalServiceAnnotations = tl.Spec.ChiaConfig.PeerService.Annotations
		}
		inputs.Annotations = kube.CombineMaps(tl.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)
	}

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaTimelord CR
func assembleDaemonService(tl k8schianetv1.ChiaTimelord) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:           fmt.Sprintf(chiatimelordNamePattern, tl.Name) + "-daemon",
		Namespace:      tl.Namespace,
		OwnerReference: getOwnerReference(tl),
		Ports:          kube.GetChiaDaemonServicePorts(),
	}

	if tl.Spec.ChiaConfig.DaemonService != nil {
		inputs.ServiceType = tl.Spec.ChiaConfig.DaemonService.ServiceType
		inputs.IPFamilyPolicy = tl.Spec.ChiaConfig.DaemonService.IPFamilyPolicy
		inputs.IPFamilies = tl.Spec.ChiaConfig.DaemonService.IPFamilies

		// Labels
		var additionalServiceLabels = make(map[string]string)
		if tl.Spec.ChiaConfig.DaemonService.Labels != nil {
			additionalServiceLabels = tl.Spec.ChiaConfig.DaemonService.Labels
		}
		inputs.Labels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
		inputs.SelectorLabels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels)

		// Annotations
		var additionalServiceAnnotations = make(map[string]string)
		if tl.Spec.ChiaConfig.DaemonService.Annotations != nil {
			additionalServiceAnnotations = tl.Spec.ChiaConfig.DaemonService.Annotations
		}
		inputs.Annotations = kube.CombineMaps(tl.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)
	}

	return kube.AssembleCommonService(inputs)
}

// assembleRPCService assembles the RPC Service resource for a ChiaTimelord CR
func assembleRPCService(tl k8schianetv1.ChiaTimelord) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:           fmt.Sprintf(chiatimelordNamePattern, tl.Name) + "-rpc",
		Namespace:      tl.Namespace,
		OwnerReference: getOwnerReference(tl),
		Ports: []corev1.ServicePort{
			{
				Port:       consts.TimelordRPCPort,
				TargetPort: intstr.FromString("rpc"),
				Protocol:   "TCP",
				Name:       "rpc",
			},
		},
	}

	if tl.Spec.ChiaConfig.RPCService != nil {
		inputs.ServiceType = tl.Spec.ChiaConfig.RPCService.ServiceType
		inputs.IPFamilyPolicy = tl.Spec.ChiaConfig.RPCService.IPFamilyPolicy
		inputs.IPFamilies = tl.Spec.ChiaConfig.RPCService.IPFamilies

		// Labels
		var additionalServiceLabels = make(map[string]string)
		if tl.Spec.ChiaConfig.RPCService.Labels != nil {
			additionalServiceLabels = tl.Spec.ChiaConfig.RPCService.Labels
		}
		inputs.Labels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
		inputs.SelectorLabels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels)

		// Annotations
		var additionalServiceAnnotations = make(map[string]string)
		if tl.Spec.ChiaConfig.RPCService.Annotations != nil {
			additionalServiceAnnotations = tl.Spec.ChiaConfig.RPCService.Annotations
		}
		inputs.Annotations = kube.CombineMaps(tl.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)
	}

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaTimelord CR
func assembleChiaExporterService(tl k8schianetv1.ChiaTimelord) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:           fmt.Sprintf(chiatimelordNamePattern, tl.Name) + "-metrics",
		Namespace:      tl.Namespace,
		OwnerReference: getOwnerReference(tl),
		Ports:          kube.GetChiaExporterServicePorts(),
	}

	if tl.Spec.ChiaExporterConfig.Service != nil {
		inputs.ServiceType = tl.Spec.ChiaExporterConfig.Service.ServiceType
		inputs.IPFamilyPolicy = tl.Spec.ChiaExporterConfig.Service.IPFamilyPolicy
		inputs.IPFamilies = tl.Spec.ChiaExporterConfig.Service.IPFamilies

		// Labels
		var additionalServiceLabels = make(map[string]string)
		if tl.Spec.ChiaExporterConfig.Service.Labels != nil {
			additionalServiceLabels = tl.Spec.ChiaExporterConfig.Service.Labels
		}
		inputs.Labels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
		inputs.SelectorLabels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels)

		// Annotations
		var additionalServiceAnnotations = make(map[string]string)
		if tl.Spec.ChiaExporterConfig.Service.Annotations != nil {
			additionalServiceAnnotations = tl.Spec.ChiaExporterConfig.Service.Annotations
		}
		inputs.Annotations = kube.CombineMaps(tl.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)
	}

	return kube.AssembleCommonService(inputs)
}

// assembleVolumeClaim assembles the PVC resource for a ChiaTimelord CR
func assembleVolumeClaim(tl k8schianetv1.ChiaTimelord) (corev1.PersistentVolumeClaim, error) {
	resourceReq, err := resource.ParseQuantity(tl.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return corev1.PersistentVolumeClaim{}, err
	}

	accessModes := []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	if len(tl.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = tl.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes
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
func assembleDeployment(tl k8schianetv1.ChiaTimelord) appsv1.Deployment {
	chiaContainer := assembleChiaContainer(tl)
	chiaExporterContainer := assembleChiaExporterContainer(tl)

	var deploy = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(chiatimelordNamePattern, tl.Name),
			Namespace:   tl.Namespace,
			Labels:      kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels),
			Annotations: tl.Spec.AdditionalMetadata.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(tl.Kind, tl.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.AdditionalMetadata.Labels),
					Annotations: tl.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					Containers:   []corev1.Container{chiaContainer, chiaExporterContainer},
					NodeSelector: tl.Spec.NodeSelector,
					Volumes:      getChiaVolumes(tl),
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
				cont.Container.VolumeMounts = getChiaVolumeMounts()
			}

			// Share chia env if enabled
			if cont.ShareEnv {
				cont.Container.Env = append(cont.Container.Env, getChiaEnv(tl)...)
			}

			deploy.Spec.Template.Spec.InitContainers = append(deploy.Spec.Template.Spec.InitContainers, cont.Container)
		}
	}

	if tl.Spec.Strategy != nil {
		deploy.Spec.Strategy = *tl.Spec.Strategy
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

func assembleChiaContainer(tl k8schianetv1.ChiaTimelord) corev1.Container {
	input := kube.AssembleChiaContainerInputs{
		Image:           tl.Spec.ChiaConfig.Image,
		ImagePullPolicy: tl.Spec.ImagePullPolicy,
		Env:             getChiaEnv(tl),
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
		VolumeMounts: getChiaVolumeMounts(),
	}

	if tl.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = tl.Spec.ChiaConfig.SecurityContext
	}

	if tl.Spec.ChiaConfig.LivenessProbe != nil {
		input.LivenessProbe = tl.Spec.ChiaConfig.LivenessProbe
	}

	if tl.Spec.ChiaConfig.ReadinessProbe != nil {
		input.ReadinessProbe = tl.Spec.ChiaConfig.ReadinessProbe
	}

	if tl.Spec.ChiaConfig.StartupProbe != nil {
		input.StartupProbe = tl.Spec.ChiaConfig.StartupProbe
	}

	if tl.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = tl.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaContainer(input)
}

func assembleChiaExporterContainer(tl k8schianetv1.ChiaTimelord) corev1.Container {
	input := kube.AssembleChiaExporterContainerInputs{
		Image:            tl.Spec.ChiaExporterConfig.Image,
		ConfigSecretName: tl.Spec.ChiaExporterConfig.ConfigSecretName,
		PullPolicy:       tl.Spec.ImagePullPolicy,
	}

	if tl.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = tl.Spec.ChiaConfig.SecurityContext
	}

	if tl.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = *tl.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaExporterContainer(input)
}
