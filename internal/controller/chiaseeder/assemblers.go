/*
Copyright 2023 Chia Network Inc.
*/

package chiaseeder

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

const chiaseederNamePattern = "%s-seeder"

// assemblePeerService assembles the peer Service resource for a ChiaSeeder CR
func assemblePeerService(seeder k8schianetv1.ChiaSeeder, fullNodePort int32) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiaseederNamePattern, seeder.Name),
		Namespace: seeder.Namespace,
		Ports: []corev1.ServicePort{
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
				Port:       fullNodePort,
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
		},
	}

	inputs.ServiceType = seeder.Spec.ChiaConfig.PeerService.ServiceType
	inputs.ExternalTrafficPolicy = seeder.Spec.ChiaConfig.PeerService.ExternalTrafficPolicy
	inputs.SessionAffinity = seeder.Spec.ChiaConfig.PeerService.SessionAffinity
	inputs.SessionAffinityConfig = seeder.Spec.ChiaConfig.PeerService.SessionAffinityConfig
	inputs.IPFamilyPolicy = seeder.Spec.ChiaConfig.PeerService.IPFamilyPolicy
	inputs.IPFamilies = seeder.Spec.ChiaConfig.PeerService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if seeder.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = seeder.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if seeder.Spec.ChiaConfig.PeerService.Annotations != nil {
		additionalServiceAnnotations = seeder.Spec.ChiaConfig.PeerService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(seeder.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Handle the Service rollup feature
	if kube.ShouldMakeService(seeder.Spec.ChiaHealthcheckConfig.Service, false) && kube.ShouldRollIntoMainPeerService(seeder.Spec.ChiaHealthcheckConfig.Service) {
		inputs.Ports = append(inputs.Ports, kube.GetChiaHealthcheckServicePorts()...)
	}

	return kube.AssembleCommonService(inputs)
}

// assembleAllService assembles the all-port Service resource for a ChiaSeeder CR
func assembleAllService(seeder k8schianetv1.ChiaSeeder, fullNodePort int32) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiaseederNamePattern, seeder.Name) + "-all",
		Namespace: seeder.Namespace,
		Ports: []corev1.ServicePort{
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

	inputs.ServiceType = seeder.Spec.ChiaConfig.AllService.ServiceType
	inputs.ExternalTrafficPolicy = seeder.Spec.ChiaConfig.AllService.ExternalTrafficPolicy
	inputs.SessionAffinity = seeder.Spec.ChiaConfig.PeerService.SessionAffinity
	inputs.SessionAffinityConfig = seeder.Spec.ChiaConfig.PeerService.SessionAffinityConfig
	inputs.IPFamilyPolicy = seeder.Spec.ChiaConfig.AllService.IPFamilyPolicy
	inputs.IPFamilies = seeder.Spec.ChiaConfig.AllService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if seeder.Spec.ChiaConfig.AllService.Labels != nil {
		additionalServiceLabels = seeder.Spec.ChiaConfig.AllService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if seeder.Spec.ChiaConfig.AllService.Annotations != nil {
		additionalServiceAnnotations = seeder.Spec.ChiaConfig.AllService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(seeder.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaSeeder CR
func assembleDaemonService(seeder k8schianetv1.ChiaSeeder) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiaseederNamePattern, seeder.Name) + "-daemon",
		Namespace: seeder.Namespace,
		Ports:     kube.GetChiaDaemonServicePorts(),
	}

	inputs.ServiceType = seeder.Spec.ChiaConfig.DaemonService.ServiceType
	inputs.ExternalTrafficPolicy = seeder.Spec.ChiaConfig.DaemonService.ExternalTrafficPolicy
	inputs.SessionAffinity = seeder.Spec.ChiaConfig.PeerService.SessionAffinity
	inputs.SessionAffinityConfig = seeder.Spec.ChiaConfig.PeerService.SessionAffinityConfig
	inputs.IPFamilyPolicy = seeder.Spec.ChiaConfig.DaemonService.IPFamilyPolicy
	inputs.IPFamilies = seeder.Spec.ChiaConfig.DaemonService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if seeder.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = seeder.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if seeder.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = seeder.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(seeder.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleRPCService assembles the RPC Service resource for a ChiaSeeder CR
func assembleRPCService(seeder k8schianetv1.ChiaSeeder) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiaseederNamePattern, seeder.Name) + "-rpc",
		Namespace: seeder.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       consts.CrawlerRPCPort,
				TargetPort: intstr.FromString("rpc"),
				Protocol:   "TCP",
				Name:       "rpc",
			},
		},
	}

	inputs.ServiceType = seeder.Spec.ChiaConfig.RPCService.ServiceType
	inputs.ExternalTrafficPolicy = seeder.Spec.ChiaConfig.RPCService.ExternalTrafficPolicy
	inputs.SessionAffinity = seeder.Spec.ChiaConfig.PeerService.SessionAffinity
	inputs.SessionAffinityConfig = seeder.Spec.ChiaConfig.PeerService.SessionAffinityConfig
	inputs.IPFamilyPolicy = seeder.Spec.ChiaConfig.RPCService.IPFamilyPolicy
	inputs.IPFamilies = seeder.Spec.ChiaConfig.RPCService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if seeder.Spec.ChiaConfig.RPCService.Labels != nil {
		additionalServiceLabels = seeder.Spec.ChiaConfig.RPCService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if seeder.Spec.ChiaConfig.RPCService.Annotations != nil {
		additionalServiceAnnotations = seeder.Spec.ChiaConfig.RPCService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(seeder.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaSeeder CR
func assembleChiaExporterService(seeder k8schianetv1.ChiaSeeder) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiaseederNamePattern, seeder.Name) + "-metrics",
		Namespace: seeder.Namespace,
		Ports:     kube.GetChiaExporterServicePorts(),
	}

	inputs.ServiceType = seeder.Spec.ChiaExporterConfig.Service.ServiceType
	inputs.ExternalTrafficPolicy = seeder.Spec.ChiaExporterConfig.Service.ExternalTrafficPolicy
	inputs.SessionAffinity = seeder.Spec.ChiaConfig.PeerService.SessionAffinity
	inputs.SessionAffinityConfig = seeder.Spec.ChiaConfig.PeerService.SessionAffinityConfig
	inputs.IPFamilyPolicy = seeder.Spec.ChiaExporterConfig.Service.IPFamilyPolicy
	inputs.IPFamilies = seeder.Spec.ChiaExporterConfig.Service.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if seeder.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = seeder.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if seeder.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = seeder.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(seeder.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleChiaHealthcheckService assembles the chia-healthcheck Service resource for a ChiaSeeder CR
func assembleChiaHealthcheckService(seeder k8schianetv1.ChiaSeeder) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiaseederNamePattern, seeder.Name) + "-healthcheck",
		Namespace: seeder.Namespace,
		Ports:     kube.GetChiaHealthcheckServicePorts(),
	}

	inputs.ServiceType = seeder.Spec.ChiaHealthcheckConfig.Service.ServiceType
	inputs.ExternalTrafficPolicy = seeder.Spec.ChiaHealthcheckConfig.Service.ExternalTrafficPolicy
	inputs.SessionAffinity = seeder.Spec.ChiaConfig.PeerService.SessionAffinity
	inputs.SessionAffinityConfig = seeder.Spec.ChiaConfig.PeerService.SessionAffinityConfig
	inputs.IPFamilyPolicy = seeder.Spec.ChiaHealthcheckConfig.Service.IPFamilyPolicy
	inputs.IPFamilies = seeder.Spec.ChiaHealthcheckConfig.Service.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if seeder.Spec.ChiaHealthcheckConfig.Service.Labels != nil {
		additionalServiceLabels = seeder.Spec.ChiaHealthcheckConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if seeder.Spec.ChiaHealthcheckConfig.Service.Annotations != nil {
		additionalServiceAnnotations = seeder.Spec.ChiaHealthcheckConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(seeder.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleVolumeClaim assembles the PVC resource for a ChiaSeeder CR
func assembleVolumeClaim(seeder k8schianetv1.ChiaSeeder) (*corev1.PersistentVolumeClaim, error) {
	if seeder.Spec.Storage == nil || seeder.Spec.Storage.ChiaRoot == nil || seeder.Spec.Storage.ChiaRoot.PersistentVolumeClaim == nil {
		return nil, nil
	}

	resourceReq, err := resource.ParseQuantity(seeder.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return nil, err
	}

	accessModes := []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	if len(seeder.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = seeder.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes
	}

	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(chiaseederNamePattern, seeder.Name),
			Namespace: seeder.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      accessModes,
			StorageClassName: &seeder.Spec.Storage.ChiaRoot.PersistentVolumeClaim.StorageClass,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resourceReq,
				},
			},
		},
	}

	return &pvc, nil
}

// assembleDeployment assembles the seeder Deployment resource for a ChiaSeeder CR
func assembleDeployment(seeder k8schianetv1.ChiaSeeder, fullNodePort int32, networkData *map[string]string) (appsv1.Deployment, error) {
	var deploy = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(chiaseederNamePattern, seeder.Name),
			Namespace:   seeder.Namespace,
			Labels:      kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels),
			Annotations: seeder.Spec.AdditionalMetadata.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels),
					Annotations: seeder.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					Affinity:                  seeder.Spec.Affinity,
					TopologySpreadConstraints: seeder.Spec.TopologySpreadConstraints,
					NodeSelector:              seeder.Spec.NodeSelector,
					Volumes:                   getChiaVolumes(seeder),
				},
			},
		},
	}

	if seeder.Spec.ServiceAccountName != nil && *seeder.Spec.ServiceAccountName != "" {
		deploy.Spec.Template.Spec.ServiceAccountName = *seeder.Spec.ServiceAccountName
	}

	chiaContainer, err := assembleChiaContainer(seeder, fullNodePort, networkData)
	if err != nil {
		return appsv1.Deployment{}, err
	}
	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, chiaContainer)

	// Get Init Containers
	deploy.Spec.Template.Spec.InitContainers = kube.GetExtraContainers(seeder.Spec.InitContainers, chiaContainer)
	// Add Init Container Volumes
	for _, init := range seeder.Spec.InitContainers {
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, init.Volumes...)
	}
	// Get Sidecar Containers
	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, kube.GetExtraContainers(seeder.Spec.Sidecars, chiaContainer)...)
	// Add Sidecar Container Volumes
	for _, sidecar := range seeder.Spec.Sidecars {
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, sidecar.Volumes...)
	}

	if seeder.Spec.ImagePullSecrets != nil && len(*seeder.Spec.ImagePullSecrets) != 0 {
		deploy.Spec.Template.Spec.ImagePullSecrets = *seeder.Spec.ImagePullSecrets
	}

	if seeder.Spec.ChiaExporterConfig.Enabled {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, assembleChiaExporterContainer(seeder))
	}

	if seeder.Spec.ChiaHealthcheckConfig.Enabled && seeder.Spec.ChiaHealthcheckConfig.DNSHostname != nil {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, assembleChiaHealthcheckContainer(seeder))
	}

	if seeder.Spec.Strategy != nil {
		deploy.Spec.Strategy = *seeder.Spec.Strategy
	}

	if seeder.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = seeder.Spec.PodSecurityContext
	}

	// TODO add pod tolerations

	return deploy, nil
}

func assembleChiaContainer(seeder k8schianetv1.ChiaSeeder, fullNodePort int32, networkData *map[string]string) (corev1.Container, error) {
	input := kube.AssembleChiaContainerInputs{
		Image:           seeder.Spec.ChiaConfig.Image,
		ImagePullPolicy: seeder.Spec.ImagePullPolicy,
		Ports:           getChiaPorts(fullNodePort),
		VolumeMounts:    getChiaVolumeMounts(seeder),
	}

	env, err := getChiaEnv(seeder, networkData)
	if err != nil {
		return corev1.Container{}, err
	}
	input.Env = env

	if seeder.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = seeder.Spec.ChiaConfig.SecurityContext
	} else {
		input.SecurityContext = &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{
					"NET_BIND_SERVICE",
				},
			},
		}
	}

	if seeder.Spec.ChiaConfig.LivenessProbe != nil {
		input.LivenessProbe = seeder.Spec.ChiaConfig.LivenessProbe
	} else if seeder.Spec.ChiaHealthcheckConfig.Enabled && seeder.Spec.ChiaHealthcheckConfig.DNSHostname != nil {
		input.ReadinessProbe = kube.AssembleChiaHealthcheckProbe(kube.AssembleChiaHealthcheckProbeInputs{
			Kind: consts.ChiaSeederKind,
		})
	}

	if seeder.Spec.ChiaConfig.ReadinessProbe != nil {
		input.ReadinessProbe = seeder.Spec.ChiaConfig.ReadinessProbe
	} else if seeder.Spec.ChiaHealthcheckConfig.Enabled && seeder.Spec.ChiaHealthcheckConfig.DNSHostname != nil {
		input.ReadinessProbe = kube.AssembleChiaHealthcheckProbe(kube.AssembleChiaHealthcheckProbeInputs{
			Kind: consts.ChiaSeederKind,
		})
	}

	if seeder.Spec.ChiaConfig.StartupProbe != nil {
		input.StartupProbe = seeder.Spec.ChiaConfig.StartupProbe
	} else if seeder.Spec.ChiaHealthcheckConfig.Enabled && seeder.Spec.ChiaHealthcheckConfig.DNSHostname != nil {
		failThresh := int32(30)
		periodSec := int32(10)
		input.StartupProbe = kube.AssembleChiaHealthcheckProbe(kube.AssembleChiaHealthcheckProbeInputs{
			Kind:             consts.ChiaSeederKind,
			FailureThreshold: &failThresh,
			PeriodSeconds:    &periodSec,
		})
	}

	if seeder.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = seeder.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaContainer(input), nil
}

func assembleChiaExporterContainer(seeder k8schianetv1.ChiaSeeder) corev1.Container {
	input := kube.AssembleChiaExporterContainerInputs{
		Image:            seeder.Spec.ChiaExporterConfig.Image,
		ConfigSecretName: seeder.Spec.ChiaExporterConfig.ConfigSecretName,
		ImagePullPolicy:  seeder.Spec.ImagePullPolicy,
	}

	if seeder.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = seeder.Spec.ChiaConfig.SecurityContext
	}

	if seeder.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = *seeder.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaExporterContainer(input)
}

func assembleChiaHealthcheckContainer(seeder k8schianetv1.ChiaSeeder) corev1.Container {
	input := kube.AssembleChiaHealthcheckContainerInputs{
		Image:           seeder.Spec.ChiaHealthcheckConfig.Image,
		DNSHostname:     seeder.Spec.ChiaHealthcheckConfig.DNSHostname,
		ImagePullPolicy: seeder.Spec.ImagePullPolicy,
	}

	if seeder.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = seeder.Spec.ChiaConfig.SecurityContext
	}

	if seeder.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = *seeder.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaHealthcheckContainer(input)
}
