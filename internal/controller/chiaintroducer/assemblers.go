/*
Copyright 2024 Chia Network Inc.
*/

package chiaintroducer

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

const chiaintroducerNamePattern = "%s-introducer"

// assemblePeerService assembles the peer Service resource for a ChiaIntroducer CR
func assemblePeerService(introducer k8schianetv1.ChiaIntroducer, fullNodePort int32) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiaintroducerNamePattern, introducer.Name),
		Namespace: introducer.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       fullNodePort,
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
		},
	}

	inputs.ServiceType = introducer.Spec.ChiaConfig.PeerService.ServiceType
	inputs.ExternalTrafficPolicy = introducer.Spec.ChiaConfig.PeerService.ExternalTrafficPolicy
	inputs.SessionAffinity = introducer.Spec.ChiaConfig.PeerService.SessionAffinity
	inputs.SessionAffinityConfig = introducer.Spec.ChiaConfig.PeerService.SessionAffinityConfig
	inputs.IPFamilyPolicy = introducer.Spec.ChiaConfig.PeerService.IPFamilyPolicy
	inputs.IPFamilies = introducer.Spec.ChiaConfig.PeerService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if introducer.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = introducer.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if introducer.Spec.ChiaConfig.PeerService.Annotations != nil {
		additionalServiceAnnotations = introducer.Spec.ChiaConfig.PeerService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(introducer.Spec.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleAllService assembles the all-port Service resource for a ChiaIntroducer CR
func assembleAllService(introducer k8schianetv1.ChiaIntroducer, fullNodePort int32) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiaintroducerNamePattern, introducer.Name) + "-all",
		Namespace: introducer.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       fullNodePort,
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
		},
	}
	inputs.Ports = append(inputs.Ports, kube.GetChiaDaemonServicePorts()...)

	inputs.ServiceType = introducer.Spec.ChiaConfig.AllService.ServiceType
	inputs.ExternalTrafficPolicy = introducer.Spec.ChiaConfig.AllService.ExternalTrafficPolicy
	inputs.SessionAffinity = introducer.Spec.ChiaConfig.AllService.SessionAffinity
	inputs.SessionAffinityConfig = introducer.Spec.ChiaConfig.AllService.SessionAffinityConfig
	inputs.IPFamilyPolicy = introducer.Spec.ChiaConfig.AllService.IPFamilyPolicy
	inputs.IPFamilies = introducer.Spec.ChiaConfig.AllService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if introducer.Spec.ChiaConfig.AllService.Labels != nil {
		additionalServiceLabels = introducer.Spec.ChiaConfig.AllService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if introducer.Spec.ChiaConfig.AllService.Annotations != nil {
		additionalServiceAnnotations = introducer.Spec.ChiaConfig.AllService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(introducer.Spec.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaIntroducer CR
func assembleDaemonService(introducer k8schianetv1.ChiaIntroducer) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiaintroducerNamePattern, introducer.Name) + "-daemon",
		Namespace: introducer.Namespace,
		Ports:     kube.GetChiaDaemonServicePorts(),
	}

	inputs.ServiceType = introducer.Spec.ChiaConfig.DaemonService.ServiceType
	inputs.ExternalTrafficPolicy = introducer.Spec.ChiaConfig.DaemonService.ExternalTrafficPolicy
	inputs.SessionAffinity = introducer.Spec.ChiaConfig.DaemonService.SessionAffinity
	inputs.SessionAffinityConfig = introducer.Spec.ChiaConfig.DaemonService.SessionAffinityConfig
	inputs.IPFamilyPolicy = introducer.Spec.ChiaConfig.DaemonService.IPFamilyPolicy
	inputs.IPFamilies = introducer.Spec.ChiaConfig.DaemonService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if introducer.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = introducer.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if introducer.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = introducer.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(introducer.Spec.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaIntroducer CR
func assembleChiaExporterService(introducer k8schianetv1.ChiaIntroducer) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiaintroducerNamePattern, introducer.Name) + "-metrics",
		Namespace: introducer.Namespace,
		Ports:     kube.GetChiaExporterServicePorts(),
	}

	inputs.ServiceType = introducer.Spec.ChiaExporterConfig.Service.ServiceType
	inputs.ExternalTrafficPolicy = introducer.Spec.ChiaExporterConfig.Service.ExternalTrafficPolicy
	inputs.SessionAffinity = introducer.Spec.ChiaExporterConfig.Service.SessionAffinity
	inputs.SessionAffinityConfig = introducer.Spec.ChiaExporterConfig.Service.SessionAffinityConfig
	inputs.IPFamilyPolicy = introducer.Spec.ChiaExporterConfig.Service.IPFamilyPolicy
	inputs.IPFamilies = introducer.Spec.ChiaExporterConfig.Service.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if introducer.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = introducer.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if introducer.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = introducer.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(introducer.Spec.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleVolumeClaim assembles the PVC resource for a ChiaIntroducer CR
func assembleVolumeClaim(introducer k8schianetv1.ChiaIntroducer) (*corev1.PersistentVolumeClaim, error) {
	if introducer.Spec.Storage == nil || introducer.Spec.Storage.ChiaRoot == nil || introducer.Spec.Storage.ChiaRoot.PersistentVolumeClaim == nil {
		return nil, nil
	}

	resourceReq, err := resource.ParseQuantity(introducer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return nil, err
	}

	accessModes := []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	if len(introducer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = introducer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes
	}

	pvc := corev1.PersistentVolumeClaim{
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
	}

	return &pvc, nil
}

// assembleDeployment assembles the introducer Deployment resource for a ChiaIntroducer CR
func assembleDeployment(introducer k8schianetv1.ChiaIntroducer, fullNodePort int32, networkData *map[string]string) (appsv1.Deployment, error) {
	var deploy = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(chiaintroducerNamePattern, introducer.Name),
			Namespace:   introducer.Namespace,
			Labels:      kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.Labels),
			Annotations: introducer.Spec.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(introducer.Kind, introducer.ObjectMeta, introducer.Spec.Labels),
					Annotations: introducer.Spec.Annotations,
				},
				Spec: corev1.PodSpec{
					Affinity:                  introducer.Spec.Affinity,
					TopologySpreadConstraints: introducer.Spec.TopologySpreadConstraints,
					NodeSelector:              introducer.Spec.NodeSelector,
					Volumes:                   getChiaVolumes(introducer),
				},
			},
		},
	}

	if introducer.Spec.ServiceAccountName != nil && *introducer.Spec.ServiceAccountName != "" {
		deploy.Spec.Template.Spec.ServiceAccountName = *introducer.Spec.ServiceAccountName
	}

	chiaContainer, err := assembleChiaContainer(introducer, fullNodePort, networkData)
	if err != nil {
		return appsv1.Deployment{}, err
	}
	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, chiaContainer)

	// Get Init Containers
	deploy.Spec.Template.Spec.InitContainers = kube.GetExtraContainers(introducer.Spec.InitContainers, chiaContainer)
	// Add Init Container Volumes
	for _, init := range introducer.Spec.InitContainers {
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, init.Volumes...)
	}

	// Get Sidecar Containers
	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, kube.GetExtraContainers(introducer.Spec.Sidecars, chiaContainer)...)
	// Add Sidecar Container Volumes
	for _, sidecar := range introducer.Spec.Sidecars {
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, sidecar.Volumes...)
	}

	if introducer.Spec.ImagePullSecrets != nil && len(*introducer.Spec.ImagePullSecrets) != 0 {
		deploy.Spec.Template.Spec.ImagePullSecrets = *introducer.Spec.ImagePullSecrets
	}

	if kube.ChiaExporterEnabled(introducer.Spec.ChiaExporterConfig) {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, assembleChiaExporterContainer(introducer))
	}

	if introducer.Spec.Strategy != nil {
		deploy.Spec.Strategy = *introducer.Spec.Strategy
	}

	if introducer.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = introducer.Spec.PodSecurityContext
	}

	// TODO add pod tolerations

	return deploy, nil
}

func assembleChiaContainer(introducer k8schianetv1.ChiaIntroducer, fullNodePort int32, networkData *map[string]string) (corev1.Container, error) {
	input := kube.AssembleChiaContainerInputs{
		Image:           introducer.Spec.ChiaConfig.Image,
		ImagePullPolicy: introducer.Spec.ImagePullPolicy,
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
		},
		VolumeMounts: getChiaVolumeMounts(introducer),
	}

	env, err := getChiaEnv(introducer, networkData)
	if err != nil {
		return corev1.Container{}, err
	}
	input.Env = env

	if introducer.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = introducer.Spec.ChiaConfig.SecurityContext
	}

	if introducer.Spec.ChiaConfig.LivenessProbe != nil {
		input.LivenessProbe = introducer.Spec.ChiaConfig.LivenessProbe
	}

	if introducer.Spec.ChiaConfig.ReadinessProbe != nil {
		input.ReadinessProbe = introducer.Spec.ChiaConfig.ReadinessProbe
	}

	if introducer.Spec.ChiaConfig.StartupProbe != nil {
		input.StartupProbe = introducer.Spec.ChiaConfig.StartupProbe
	}

	if introducer.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = introducer.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaContainer(input), nil
}

func assembleChiaExporterContainer(introducer k8schianetv1.ChiaIntroducer) corev1.Container {
	input := kube.AssembleChiaExporterContainerInputs{
		Image:            introducer.Spec.ChiaExporterConfig.Image,
		ConfigSecretName: introducer.Spec.ChiaExporterConfig.ConfigSecretName,
		ImagePullPolicy:  introducer.Spec.ImagePullPolicy,
	}

	if introducer.Spec.ChiaExporterConfig.SecurityContext != nil {
		input.SecurityContext = introducer.Spec.ChiaExporterConfig.SecurityContext
	}

	if introducer.Spec.ChiaExporterConfig.Resources != nil {
		input.ResourceRequirements = *introducer.Spec.ChiaExporterConfig.Resources
	}

	return kube.AssembleChiaExporterContainer(input)
}
