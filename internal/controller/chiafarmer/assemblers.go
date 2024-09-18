/*
Copyright 2023 Chia Network Inc.
*/

package chiafarmer

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

const chiafarmerNamePattern = "%s-farmer"

// assemblePeerService assembles the peer Service resource for a ChiaFarmer CR
func assemblePeerService(farmer k8schianetv1.ChiaFarmer) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiafarmerNamePattern, farmer.Name),
		Namespace: farmer.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       consts.FarmerPort,
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
		},
	}

	inputs.ServiceType = farmer.Spec.ChiaConfig.PeerService.ServiceType
	inputs.IPFamilyPolicy = farmer.Spec.ChiaConfig.PeerService.IPFamilyPolicy
	inputs.IPFamilies = farmer.Spec.ChiaConfig.PeerService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if farmer.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = farmer.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if farmer.Spec.ChiaConfig.PeerService.Annotations != nil {
		additionalServiceAnnotations = farmer.Spec.ChiaConfig.PeerService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(farmer.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleAllService assembles the all-port Service resource for a ChiaFarmer CR
func assembleAllService(farmer k8schianetv1.ChiaFarmer) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiafarmerNamePattern, farmer.Name) + "-all",
		Namespace: farmer.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       consts.FarmerPort,
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
			{
				Port:       consts.FarmerRPCPort,
				TargetPort: intstr.FromString("rpc"),
				Protocol:   "TCP",
				Name:       "rpc",
			},
		},
	}
	inputs.Ports = append(inputs.Ports, kube.GetChiaDaemonServicePorts()...)

	inputs.ServiceType = farmer.Spec.ChiaConfig.AllService.ServiceType
	inputs.IPFamilyPolicy = farmer.Spec.ChiaConfig.AllService.IPFamilyPolicy
	inputs.IPFamilies = farmer.Spec.ChiaConfig.AllService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if farmer.Spec.ChiaConfig.AllService.Labels != nil {
		additionalServiceLabels = farmer.Spec.ChiaConfig.AllService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if farmer.Spec.ChiaConfig.AllService.Annotations != nil {
		additionalServiceAnnotations = farmer.Spec.ChiaConfig.AllService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(farmer.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaFarmer CR
func assembleDaemonService(farmer k8schianetv1.ChiaFarmer) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiafarmerNamePattern, farmer.Name) + "-daemon",
		Namespace: farmer.Namespace,
		Ports:     kube.GetChiaDaemonServicePorts(),
	}

	inputs.ServiceType = farmer.Spec.ChiaConfig.DaemonService.ServiceType
	inputs.IPFamilyPolicy = farmer.Spec.ChiaConfig.DaemonService.IPFamilyPolicy
	inputs.IPFamilies = farmer.Spec.ChiaConfig.DaemonService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if farmer.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = farmer.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if farmer.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = farmer.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(farmer.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleRPCService assembles the RPC Service resource for a ChiaFarmer CR
func assembleRPCService(farmer k8schianetv1.ChiaFarmer) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiafarmerNamePattern, farmer.Name) + "-rpc",
		Namespace: farmer.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       consts.FarmerRPCPort,
				TargetPort: intstr.FromString("rpc"),
				Protocol:   "TCP",
				Name:       "rpc",
			},
		},
	}

	inputs.ServiceType = farmer.Spec.ChiaConfig.RPCService.ServiceType
	inputs.IPFamilyPolicy = farmer.Spec.ChiaConfig.RPCService.IPFamilyPolicy
	inputs.IPFamilies = farmer.Spec.ChiaConfig.RPCService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if farmer.Spec.ChiaConfig.RPCService.Labels != nil {
		additionalServiceLabels = farmer.Spec.ChiaConfig.RPCService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if farmer.Spec.ChiaConfig.RPCService.Annotations != nil {
		additionalServiceAnnotations = farmer.Spec.ChiaConfig.RPCService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(farmer.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaFarmer CR
func assembleChiaExporterService(farmer k8schianetv1.ChiaFarmer) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiafarmerNamePattern, farmer.Name) + "-metrics",
		Namespace: farmer.Namespace,
		Ports:     kube.GetChiaExporterServicePorts(),
	}

	inputs.ServiceType = farmer.Spec.ChiaExporterConfig.Service.ServiceType
	inputs.IPFamilyPolicy = farmer.Spec.ChiaExporterConfig.Service.IPFamilyPolicy
	inputs.IPFamilies = farmer.Spec.ChiaExporterConfig.Service.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if farmer.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = farmer.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if farmer.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = farmer.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(farmer.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleVolumeClaim assembles the PVC resource for a ChiaFarmer CR
func assembleVolumeClaim(farmer k8schianetv1.ChiaFarmer) (*corev1.PersistentVolumeClaim, error) {
	if farmer.Spec.Storage == nil || farmer.Spec.Storage.ChiaRoot == nil || farmer.Spec.Storage.ChiaRoot.PersistentVolumeClaim == nil {
		return nil, nil
	}

	resourceReq, err := resource.ParseQuantity(farmer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return nil, err
	}

	accessModes := []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	if len(farmer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = farmer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes
	}

	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(chiafarmerNamePattern, farmer.Name),
			Namespace: farmer.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      accessModes,
			StorageClassName: &farmer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.StorageClass,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resourceReq,
				},
			},
		},
	}

	return &pvc, nil
}

// assembleDeployment assembles the farmer Deployment resource for a ChiaFarmer CR
func assembleDeployment(farmer k8schianetv1.ChiaFarmer, networkData *map[string]string) (appsv1.Deployment, error) {
	var deploy = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(chiafarmerNamePattern, farmer.Name),
			Namespace:   farmer.Namespace,
			Labels:      kube.GetCommonLabels(farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels),
			Annotations: farmer.Spec.AdditionalMetadata.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(farmer.Kind, farmer.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels),
					Annotations: farmer.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					Affinity:     farmer.Spec.Affinity,
					NodeSelector: farmer.Spec.NodeSelector,
					Volumes:      getChiaVolumes(farmer),
				},
			},
		},
	}

	chiaContainer, err := assembleChiaContainer(farmer, networkData)
	if err != nil {
		return appsv1.Deployment{}, err
	}
	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, chiaContainer)

	if len(farmer.Spec.InitContainers) != 0 {
		// Overwrite any volumeMounts specified in init containers. Not currently supported.
		for _, cont := range farmer.Spec.InitContainers {
			cont.Container.VolumeMounts = []corev1.VolumeMount{}

			// Share chia volume mounts if enabled
			if cont.ShareVolumeMounts {
				cont.Container.VolumeMounts = getChiaVolumeMounts()
			}

			// Share chia env if enabled
			if cont.ShareEnv {
				cont.Container.Env = append(cont.Container.Env, chiaContainer.Env...)
			}

			deploy.Spec.Template.Spec.InitContainers = append(deploy.Spec.Template.Spec.InitContainers, cont.Container)
		}
	}

	if farmer.Spec.ImagePullSecrets != nil && len(*farmer.Spec.ImagePullSecrets) != 0 {
		deploy.Spec.Template.Spec.ImagePullSecrets = *farmer.Spec.ImagePullSecrets
	}

	if farmer.Spec.ChiaExporterConfig.Enabled {
		chiaExporterContainer := assembleChiaExporterContainer(farmer)
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, chiaExporterContainer)
	}

	if farmer.Spec.Strategy != nil {
		deploy.Spec.Strategy = *farmer.Spec.Strategy
	}

	if farmer.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = farmer.Spec.PodSecurityContext
	}

	if len(farmer.Spec.Sidecars.Containers) > 0 {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, farmer.Spec.Sidecars.Containers...)
	}

	// TODO add pod affinity, tolerations

	return deploy, nil
}

func assembleChiaContainer(farmer k8schianetv1.ChiaFarmer, networkData *map[string]string) (corev1.Container, error) {
	input := kube.AssembleChiaContainerInputs{
		Image:           farmer.Spec.ChiaConfig.Image,
		ImagePullPolicy: farmer.Spec.ImagePullPolicy,
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
		VolumeMounts: getChiaVolumeMounts(),
	}

	env, err := getChiaEnv(farmer, networkData)
	if err != nil {
		return corev1.Container{}, err
	}
	input.Env = env

	if farmer.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = farmer.Spec.ChiaConfig.SecurityContext
	}

	if farmer.Spec.ChiaConfig.LivenessProbe != nil {
		input.LivenessProbe = farmer.Spec.ChiaConfig.LivenessProbe
	}

	if farmer.Spec.ChiaConfig.ReadinessProbe != nil {
		input.ReadinessProbe = farmer.Spec.ChiaConfig.ReadinessProbe
	}

	if farmer.Spec.ChiaConfig.StartupProbe != nil {
		input.StartupProbe = farmer.Spec.ChiaConfig.StartupProbe
	}

	if farmer.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = farmer.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaContainer(input), nil
}

func assembleChiaExporterContainer(farmer k8schianetv1.ChiaFarmer) corev1.Container {
	input := kube.AssembleChiaExporterContainerInputs{
		Image:            farmer.Spec.ChiaExporterConfig.Image,
		ConfigSecretName: farmer.Spec.ChiaExporterConfig.ConfigSecretName,
		ImagePullPolicy:  farmer.Spec.ImagePullPolicy,
	}

	if farmer.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = farmer.Spec.ChiaConfig.SecurityContext
	}

	if farmer.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = *farmer.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaExporterContainer(input)
}
