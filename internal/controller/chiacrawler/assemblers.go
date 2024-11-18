/*
Copyright 2024 Chia Network Inc.
*/

package chiacrawler

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

const chiacrawlerNamePattern = "%s-crawler"

// assemblePeerService assembles the peer Service resource for a ChiaCrawler CR
func assemblePeerService(crawler k8schianetv1.ChiaCrawler, fullNodePort int32) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiacrawlerNamePattern, crawler.Name),
		Namespace: crawler.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       fullNodePort,
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
		},
	}

	inputs.ServiceType = crawler.Spec.ChiaConfig.PeerService.ServiceType
	inputs.ExternalTrafficPolicy = crawler.Spec.ChiaConfig.PeerService.ExternalTrafficPolicy
	inputs.IPFamilyPolicy = crawler.Spec.ChiaConfig.PeerService.IPFamilyPolicy
	inputs.IPFamilies = crawler.Spec.ChiaConfig.PeerService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if crawler.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = crawler.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if crawler.Spec.ChiaConfig.PeerService.Annotations != nil {
		additionalServiceAnnotations = crawler.Spec.ChiaConfig.PeerService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(crawler.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleAllService assembles the all-port Service resource for a ChiaCrawler CR
func assembleAllService(crawler k8schianetv1.ChiaCrawler, fullNodePort int32) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiacrawlerNamePattern, crawler.Name) + "-all",
		Namespace: crawler.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       fullNodePort,
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
			{
				Port:       consts.CrawlerRPCPort,
				TargetPort: intstr.FromString("rpc"),
				Protocol:   "TCP",
				Name:       "rpc",
			},
		},
	}
	inputs.Ports = append(inputs.Ports, kube.GetChiaDaemonServicePorts()...)

	inputs.ServiceType = crawler.Spec.ChiaConfig.AllService.ServiceType
	inputs.ExternalTrafficPolicy = crawler.Spec.ChiaConfig.AllService.ExternalTrafficPolicy
	inputs.IPFamilyPolicy = crawler.Spec.ChiaConfig.AllService.IPFamilyPolicy
	inputs.IPFamilies = crawler.Spec.ChiaConfig.AllService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if crawler.Spec.ChiaConfig.AllService.Labels != nil {
		additionalServiceLabels = crawler.Spec.ChiaConfig.AllService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if crawler.Spec.ChiaConfig.AllService.Annotations != nil {
		additionalServiceAnnotations = crawler.Spec.ChiaConfig.AllService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(crawler.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaCrawler CR
func assembleDaemonService(crawler k8schianetv1.ChiaCrawler) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiacrawlerNamePattern, crawler.Name) + "-daemon",
		Namespace: crawler.Namespace,
		Ports:     kube.GetChiaDaemonServicePorts(),
	}

	inputs.ServiceType = crawler.Spec.ChiaConfig.DaemonService.ServiceType
	inputs.ExternalTrafficPolicy = crawler.Spec.ChiaConfig.DaemonService.ExternalTrafficPolicy
	inputs.IPFamilyPolicy = crawler.Spec.ChiaConfig.DaemonService.IPFamilyPolicy
	inputs.IPFamilies = crawler.Spec.ChiaConfig.DaemonService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if crawler.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = crawler.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if crawler.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = crawler.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(crawler.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleRPCService assembles the RPC Service resource for a ChiaCrawler CR
func assembleRPCService(crawler k8schianetv1.ChiaCrawler) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiacrawlerNamePattern, crawler.Name) + "-rpc",
		Namespace: crawler.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       consts.CrawlerRPCPort,
				TargetPort: intstr.FromString("rpc"),
				Protocol:   "TCP",
				Name:       "rpc",
			},
		},
	}

	inputs.ServiceType = crawler.Spec.ChiaConfig.RPCService.ServiceType
	inputs.ExternalTrafficPolicy = crawler.Spec.ChiaConfig.RPCService.ExternalTrafficPolicy
	inputs.IPFamilyPolicy = crawler.Spec.ChiaConfig.RPCService.IPFamilyPolicy
	inputs.IPFamilies = crawler.Spec.ChiaConfig.RPCService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if crawler.Spec.ChiaConfig.RPCService.Labels != nil {
		additionalServiceLabels = crawler.Spec.ChiaConfig.RPCService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if crawler.Spec.ChiaConfig.RPCService.Annotations != nil {
		additionalServiceAnnotations = crawler.Spec.ChiaConfig.RPCService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(crawler.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaCrawler CR
func assembleChiaExporterService(crawler k8schianetv1.ChiaCrawler) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiacrawlerNamePattern, crawler.Name) + "-metrics",
		Namespace: crawler.Namespace,
		Ports:     kube.GetChiaExporterServicePorts(),
	}

	inputs.ServiceType = crawler.Spec.ChiaExporterConfig.Service.ServiceType
	inputs.ExternalTrafficPolicy = crawler.Spec.ChiaExporterConfig.Service.ExternalTrafficPolicy
	inputs.IPFamilyPolicy = crawler.Spec.ChiaExporterConfig.Service.IPFamilyPolicy
	inputs.IPFamilies = crawler.Spec.ChiaExporterConfig.Service.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if crawler.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = crawler.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if crawler.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = crawler.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(crawler.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleVolumeClaim assembles the PVC resource for a ChiaCrawler CR
func assembleVolumeClaim(crawler k8schianetv1.ChiaCrawler) (*corev1.PersistentVolumeClaim, error) {
	if crawler.Spec.Storage == nil || crawler.Spec.Storage.ChiaRoot == nil || crawler.Spec.Storage.ChiaRoot.PersistentVolumeClaim == nil {
		return nil, nil
	}

	resourceReq, err := resource.ParseQuantity(crawler.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return nil, err
	}

	accessModes := []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	if len(crawler.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = crawler.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes
	}

	pvc := corev1.PersistentVolumeClaim{
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
	}

	return &pvc, nil
}

// assembleDeployment assembles the crawler Deployment resource for a ChiaCrawler CR
func assembleDeployment(crawler k8schianetv1.ChiaCrawler, fullNodePort int32, networkData *map[string]string) (appsv1.Deployment, error) {
	var deploy = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(chiacrawlerNamePattern, crawler.Name),
			Namespace:   crawler.Namespace,
			Labels:      kube.GetCommonLabels(crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels),
			Annotations: crawler.Spec.AdditionalMetadata.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(crawler.Kind, crawler.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(crawler.Kind, crawler.ObjectMeta, crawler.Spec.AdditionalMetadata.Labels),
					Annotations: crawler.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					Affinity:     crawler.Spec.Affinity,
					NodeSelector: crawler.Spec.NodeSelector,
					Volumes:      getChiaVolumes(crawler),
				},
			},
		},
	}

	if crawler.Spec.ServiceAccountName != nil && *crawler.Spec.ServiceAccountName != "" {
		deploy.Spec.Template.Spec.ServiceAccountName = *crawler.Spec.ServiceAccountName
	}

	chiaContainer, err := assembleChiaContainer(crawler, fullNodePort, networkData)
	if err != nil {
		return appsv1.Deployment{}, err
	}
	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, chiaContainer)

	// Get Init Containers
	deploy.Spec.Template.Spec.InitContainers = kube.GetExtraContainers(crawler.Spec.InitContainers, chiaContainer)
	// Add Init Container Volumes
	for _, init := range crawler.Spec.InitContainers {
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, init.Volumes...)
	}

	// Get Sidecar Containers
	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, kube.GetExtraContainers(crawler.Spec.Sidecars, chiaContainer)...)
	// Add Sidecar Container Volumes
	for _, sidecar := range crawler.Spec.Sidecars {
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, sidecar.Volumes...)
	}

	if crawler.Spec.ImagePullSecrets != nil && len(*crawler.Spec.ImagePullSecrets) != 0 {
		deploy.Spec.Template.Spec.ImagePullSecrets = *crawler.Spec.ImagePullSecrets
	}

	if crawler.Spec.ChiaExporterConfig.Enabled {
		chiaExporterContainer := assembleChiaExporterContainer(crawler)
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, chiaExporterContainer)
	}

	if crawler.Spec.Strategy != nil {
		deploy.Spec.Strategy = *crawler.Spec.Strategy
	}

	if crawler.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = crawler.Spec.PodSecurityContext
	}

	// TODO add pod tolerations

	return deploy, nil
}

func assembleChiaContainer(crawler k8schianetv1.ChiaCrawler, fullNodePort int32, networkData *map[string]string) (corev1.Container, error) {
	input := kube.AssembleChiaContainerInputs{
		Image:           crawler.Spec.ChiaConfig.Image,
		ImagePullPolicy: crawler.Spec.ImagePullPolicy,
		Ports: []corev1.ContainerPort{
			{
				Name:          "daemon",
				ContainerPort: consts.DaemonPort,
				Protocol:      "TCP",
			},
			{
				Name:          "peers",
				ContainerPort: fullNodePort,
				Protocol:      "TCP",
			},
			{
				Name:          "rpc",
				ContainerPort: consts.CrawlerRPCPort,
				Protocol:      "TCP",
			},
		},
		VolumeMounts: getChiaVolumeMounts(crawler),
	}

	env, err := getChiaEnv(crawler, networkData)
	if err != nil {
		return corev1.Container{}, err
	}
	input.Env = env

	if crawler.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = crawler.Spec.ChiaConfig.SecurityContext
	}

	if crawler.Spec.ChiaConfig.LivenessProbe != nil {
		input.LivenessProbe = crawler.Spec.ChiaConfig.LivenessProbe
	}

	if crawler.Spec.ChiaConfig.ReadinessProbe != nil {
		input.ReadinessProbe = crawler.Spec.ChiaConfig.ReadinessProbe
	}

	if crawler.Spec.ChiaConfig.StartupProbe != nil {
		input.StartupProbe = crawler.Spec.ChiaConfig.StartupProbe
	}

	if crawler.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = crawler.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaContainer(input), nil
}

func assembleChiaExporterContainer(crawler k8schianetv1.ChiaCrawler) corev1.Container {
	input := kube.AssembleChiaExporterContainerInputs{
		Image:            crawler.Spec.ChiaExporterConfig.Image,
		ConfigSecretName: crawler.Spec.ChiaExporterConfig.ConfigSecretName,
		ImagePullPolicy:  crawler.Spec.ImagePullPolicy,
	}

	if crawler.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = crawler.Spec.ChiaConfig.SecurityContext
	}

	if crawler.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = *crawler.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaExporterContainer(input)
}
