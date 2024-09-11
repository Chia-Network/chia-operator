/*
Copyright 2023 Chia Network Inc.
*/

package chiaharvester

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

const chiaharvesterNamePattern = "%s-harvester"

// assemblePeerService assembles the peer Service resource for a ChiaHarvester CR
func assemblePeerService(harvester k8schianetv1.ChiaHarvester) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiaharvesterNamePattern, harvester.Name),
		Namespace: harvester.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       consts.HarvesterPort,
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
		},
	}

	inputs.ServiceType = harvester.Spec.ChiaConfig.PeerService.ServiceType
	inputs.IPFamilyPolicy = harvester.Spec.ChiaConfig.PeerService.IPFamilyPolicy
	inputs.IPFamilies = harvester.Spec.ChiaConfig.PeerService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if harvester.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = harvester.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if harvester.Spec.ChiaConfig.PeerService.Annotations != nil {
		additionalServiceAnnotations = harvester.Spec.ChiaConfig.PeerService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(harvester.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleAllService assembles the all-port Service resource for a ChiaHarvester CR
func assembleAllService(harvester k8schianetv1.ChiaHarvester) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiaharvesterNamePattern, harvester.Name) + "-all",
		Namespace: harvester.Namespace,
		Ports: []corev1.ServicePort{
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
	}
	inputs.Ports = append(inputs.Ports, kube.GetChiaDaemonServicePorts()...)

	inputs.ServiceType = harvester.Spec.ChiaConfig.AllService.ServiceType
	inputs.IPFamilyPolicy = harvester.Spec.ChiaConfig.AllService.IPFamilyPolicy
	inputs.IPFamilies = harvester.Spec.ChiaConfig.AllService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if harvester.Spec.ChiaConfig.AllService.Labels != nil {
		additionalServiceLabels = harvester.Spec.ChiaConfig.AllService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if harvester.Spec.ChiaConfig.AllService.Annotations != nil {
		additionalServiceAnnotations = harvester.Spec.ChiaConfig.AllService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(harvester.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaHarvester CR
func assembleDaemonService(harvester k8schianetv1.ChiaHarvester) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiaharvesterNamePattern, harvester.Name) + "-daemon",
		Namespace: harvester.Namespace,
		Ports:     kube.GetChiaDaemonServicePorts(),
	}

	inputs.ServiceType = harvester.Spec.ChiaConfig.DaemonService.ServiceType
	inputs.IPFamilyPolicy = harvester.Spec.ChiaConfig.DaemonService.IPFamilyPolicy
	inputs.IPFamilies = harvester.Spec.ChiaConfig.DaemonService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if harvester.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = harvester.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if harvester.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = harvester.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(harvester.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleRPCService assembles the RPC Service resource for a ChiaHarvester CR
func assembleRPCService(harvester k8schianetv1.ChiaHarvester) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiaharvesterNamePattern, harvester.Name) + "-rpc",
		Namespace: harvester.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       consts.HarvesterRPCPort,
				TargetPort: intstr.FromString("rpc"),
				Protocol:   "TCP",
				Name:       "rpc",
			},
		},
	}

	inputs.ServiceType = harvester.Spec.ChiaConfig.RPCService.ServiceType
	inputs.IPFamilyPolicy = harvester.Spec.ChiaConfig.RPCService.IPFamilyPolicy
	inputs.IPFamilies = harvester.Spec.ChiaConfig.RPCService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if harvester.Spec.ChiaConfig.RPCService.Labels != nil {
		additionalServiceLabels = harvester.Spec.ChiaConfig.RPCService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if harvester.Spec.ChiaConfig.RPCService.Annotations != nil {
		additionalServiceAnnotations = harvester.Spec.ChiaConfig.RPCService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(harvester.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaHarvester CR
func assembleChiaExporterService(harvester k8schianetv1.ChiaHarvester) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiaharvesterNamePattern, harvester.Name) + "-metrics",
		Namespace: harvester.Namespace,
		Ports:     kube.GetChiaExporterServicePorts(),
	}

	inputs.ServiceType = harvester.Spec.ChiaExporterConfig.Service.ServiceType
	inputs.IPFamilyPolicy = harvester.Spec.ChiaExporterConfig.Service.IPFamilyPolicy
	inputs.IPFamilies = harvester.Spec.ChiaExporterConfig.Service.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if harvester.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = harvester.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if harvester.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = harvester.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(harvester.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleVolumeClaim assembles the PVC resource for a ChiaHarvester CR
func assembleVolumeClaim(harvester k8schianetv1.ChiaHarvester) (*corev1.PersistentVolumeClaim, error) {
	if harvester.Spec.Storage == nil || harvester.Spec.Storage.ChiaRoot == nil || harvester.Spec.Storage.ChiaRoot.PersistentVolumeClaim == nil {
		return nil, nil
	}

	resourceReq, err := resource.ParseQuantity(harvester.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return nil, err
	}

	accessModes := []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	if len(harvester.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = harvester.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes
	}

	pvc := corev1.PersistentVolumeClaim{
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
	}

	return &pvc, nil
}

// assembleDeployment assembles the harvester Deployment resource for a ChiaHarvester CR
func assembleDeployment(harvester k8schianetv1.ChiaHarvester) appsv1.Deployment {
	var deploy = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(chiaharvesterNamePattern, harvester.Name),
			Namespace:   harvester.Namespace,
			Labels:      kube.GetCommonLabels(harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels),
			Annotations: harvester.Spec.AdditionalMetadata.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(harvester.Kind, harvester.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(harvester.Kind, harvester.ObjectMeta, harvester.Spec.AdditionalMetadata.Labels),
					Annotations: harvester.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					Containers:   []corev1.Container{assembleChiaContainer(harvester)},
					Affinity:     harvester.Spec.Affinity,
					NodeSelector: harvester.Spec.NodeSelector,
					Volumes:      getChiaVolumes(harvester),
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
				cont.Container.VolumeMounts = getChiaVolumeMounts(harvester)
			}

			// Share chia env if enabled
			if cont.ShareEnv {
				cont.Container.Env = append(cont.Container.Env, getChiaEnv(harvester)...)
			}

			deploy.Spec.Template.Spec.InitContainers = append(deploy.Spec.Template.Spec.InitContainers, cont.Container)
		}
	}

	if harvester.Spec.ImagePullSecrets != nil && len(*harvester.Spec.ImagePullSecrets) != 0 {
		deploy.Spec.Template.Spec.ImagePullSecrets = *harvester.Spec.ImagePullSecrets
	}

	if harvester.Spec.ChiaExporterConfig.Enabled {
		chiaExporterContainer := assembleChiaExporterContainer(harvester)
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, chiaExporterContainer)
	}

	if harvester.Spec.Strategy != nil {
		deploy.Spec.Strategy = *harvester.Spec.Strategy
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

func assembleChiaContainer(harvester k8schianetv1.ChiaHarvester) corev1.Container {
	input := kube.AssembleChiaContainerInputs{
		Image:           harvester.Spec.ChiaConfig.Image,
		ImagePullPolicy: harvester.Spec.ImagePullPolicy,
		Env:             getChiaEnv(harvester),
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
		VolumeMounts: getChiaVolumeMounts(harvester),
	}

	if harvester.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = harvester.Spec.ChiaConfig.SecurityContext
	}

	if harvester.Spec.ChiaConfig.LivenessProbe != nil {
		input.LivenessProbe = harvester.Spec.ChiaConfig.LivenessProbe
	}

	if harvester.Spec.ChiaConfig.ReadinessProbe != nil {
		input.ReadinessProbe = harvester.Spec.ChiaConfig.ReadinessProbe
	}

	if harvester.Spec.ChiaConfig.StartupProbe != nil {
		input.StartupProbe = harvester.Spec.ChiaConfig.StartupProbe
	}

	if harvester.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = harvester.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaContainer(input)
}

func assembleChiaExporterContainer(harvester k8schianetv1.ChiaHarvester) corev1.Container {
	input := kube.AssembleChiaExporterContainerInputs{
		Image:            harvester.Spec.ChiaExporterConfig.Image,
		ConfigSecretName: harvester.Spec.ChiaExporterConfig.ConfigSecretName,
		ImagePullPolicy:  harvester.Spec.ImagePullPolicy,
	}

	if harvester.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = harvester.Spec.ChiaConfig.SecurityContext
	}

	if harvester.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = *harvester.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaExporterContainer(input)
}
