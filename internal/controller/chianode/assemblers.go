/*
Copyright 2023 Chia Network Inc.
*/

package chianode

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
)

const chianodeNamePattern = "%s-node"

// assemblePeerService assembles the peer Service resource for a ChiaNode CR
func assemblePeerService(node k8schianetv1.ChiaNode, fullNodePort int32) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chianodeNamePattern, node.Name),
		Namespace: node.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       fullNodePort,
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
		},
	}

	inputs.ServiceType = node.Spec.ChiaConfig.PeerService.ServiceType
	inputs.ExternalTrafficPolicy = node.Spec.ChiaConfig.PeerService.ExternalTrafficPolicy
	inputs.SessionAffinity = node.Spec.ChiaConfig.PeerService.SessionAffinity
	inputs.SessionAffinityConfig = node.Spec.ChiaConfig.PeerService.SessionAffinityConfig
	inputs.IPFamilyPolicy = node.Spec.ChiaConfig.PeerService.IPFamilyPolicy
	inputs.IPFamilies = node.Spec.ChiaConfig.PeerService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if node.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = node.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if node.Spec.ChiaConfig.PeerService.Annotations != nil {
		additionalServiceAnnotations = node.Spec.ChiaConfig.PeerService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(node.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Handle the Service rollup feature
	if kube.ShouldMakeService(node.Spec.ChiaHealthcheckConfig.Service, false) && kube.ShouldRollIntoMainPeerService(node.Spec.ChiaHealthcheckConfig.Service) {
		inputs.Ports = append(inputs.Ports, kube.GetChiaHealthcheckServicePorts()...)
	}

	return kube.AssembleCommonService(inputs)
}

// assembleAllService assembles the all-port Service resource for a ChiaNode CR
func assembleAllService(node k8schianetv1.ChiaNode, fullNodePort int32) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chianodeNamePattern, node.Name) + "-all",
		Namespace: node.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       fullNodePort,
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
			{
				Port:       consts.NodeRPCPort,
				TargetPort: intstr.FromString("rpc"),
				Protocol:   "TCP",
				Name:       "rpc",
			},
		},
	}
	inputs.Ports = append(inputs.Ports, kube.GetChiaDaemonServicePorts()...)

	inputs.ServiceType = node.Spec.ChiaConfig.AllService.ServiceType
	inputs.ExternalTrafficPolicy = node.Spec.ChiaConfig.AllService.ExternalTrafficPolicy
	inputs.SessionAffinity = node.Spec.ChiaConfig.PeerService.SessionAffinity
	inputs.SessionAffinityConfig = node.Spec.ChiaConfig.PeerService.SessionAffinityConfig
	inputs.IPFamilyPolicy = node.Spec.ChiaConfig.AllService.IPFamilyPolicy
	inputs.IPFamilies = node.Spec.ChiaConfig.AllService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if node.Spec.ChiaConfig.AllService.Labels != nil {
		additionalServiceLabels = node.Spec.ChiaConfig.AllService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if node.Spec.ChiaConfig.AllService.Annotations != nil {
		additionalServiceAnnotations = node.Spec.ChiaConfig.AllService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(node.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaNode CR
func assembleDaemonService(node k8schianetv1.ChiaNode) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chianodeNamePattern, node.Name) + "-daemon",
		Namespace: node.Namespace,
		Ports:     kube.GetChiaDaemonServicePorts(),
	}

	inputs.ServiceType = node.Spec.ChiaConfig.DaemonService.ServiceType
	inputs.ExternalTrafficPolicy = node.Spec.ChiaConfig.DaemonService.ExternalTrafficPolicy
	inputs.SessionAffinity = node.Spec.ChiaConfig.PeerService.SessionAffinity
	inputs.SessionAffinityConfig = node.Spec.ChiaConfig.PeerService.SessionAffinityConfig
	inputs.IPFamilyPolicy = node.Spec.ChiaConfig.DaemonService.IPFamilyPolicy
	inputs.IPFamilies = node.Spec.ChiaConfig.DaemonService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if node.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = node.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if node.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = node.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(node.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleRPCService assembles the RPC Service resource for a ChiaNode CR
func assembleRPCService(node k8schianetv1.ChiaNode) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chianodeNamePattern, node.Name) + "-rpc",
		Namespace: node.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       consts.NodeRPCPort,
				TargetPort: intstr.FromString("rpc"),
				Protocol:   "TCP",
				Name:       "rpc",
			},
		},
	}

	inputs.ServiceType = node.Spec.ChiaConfig.RPCService.ServiceType
	inputs.ExternalTrafficPolicy = node.Spec.ChiaConfig.RPCService.ExternalTrafficPolicy
	inputs.SessionAffinity = node.Spec.ChiaConfig.PeerService.SessionAffinity
	inputs.SessionAffinityConfig = node.Spec.ChiaConfig.PeerService.SessionAffinityConfig
	inputs.IPFamilyPolicy = node.Spec.ChiaConfig.RPCService.IPFamilyPolicy
	inputs.IPFamilies = node.Spec.ChiaConfig.RPCService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if node.Spec.ChiaConfig.RPCService.Labels != nil {
		additionalServiceLabels = node.Spec.ChiaConfig.RPCService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if node.Spec.ChiaConfig.RPCService.Annotations != nil {
		additionalServiceAnnotations = node.Spec.ChiaConfig.RPCService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(node.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaNode CR
func assembleChiaExporterService(node k8schianetv1.ChiaNode) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chianodeNamePattern, node.Name) + "-metrics",
		Namespace: node.Namespace,
		Ports:     kube.GetChiaExporterServicePorts(),
	}

	inputs.ServiceType = node.Spec.ChiaExporterConfig.Service.ServiceType
	inputs.ExternalTrafficPolicy = node.Spec.ChiaExporterConfig.Service.ExternalTrafficPolicy
	inputs.SessionAffinity = node.Spec.ChiaConfig.PeerService.SessionAffinity
	inputs.SessionAffinityConfig = node.Spec.ChiaConfig.PeerService.SessionAffinityConfig
	inputs.IPFamilyPolicy = node.Spec.ChiaExporterConfig.Service.IPFamilyPolicy
	inputs.IPFamilies = node.Spec.ChiaExporterConfig.Service.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if node.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = node.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if node.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = node.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(node.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleChiaHealthcheckService assembles the chia-healthcheck Service resource for a ChiaNode CR
func assembleChiaHealthcheckService(node k8schianetv1.ChiaNode) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chianodeNamePattern, node.Name) + "-healthcheck",
		Namespace: node.Namespace,
		Ports:     kube.GetChiaHealthcheckServicePorts(),
	}

	inputs.ServiceType = node.Spec.ChiaHealthcheckConfig.Service.ServiceType
	inputs.ExternalTrafficPolicy = node.Spec.ChiaHealthcheckConfig.Service.ExternalTrafficPolicy
	inputs.SessionAffinity = node.Spec.ChiaConfig.PeerService.SessionAffinity
	inputs.SessionAffinityConfig = node.Spec.ChiaConfig.PeerService.SessionAffinityConfig
	inputs.IPFamilyPolicy = node.Spec.ChiaHealthcheckConfig.Service.IPFamilyPolicy
	inputs.IPFamilies = node.Spec.ChiaHealthcheckConfig.Service.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if node.Spec.ChiaHealthcheckConfig.Service.Labels != nil {
		additionalServiceLabels = node.Spec.ChiaHealthcheckConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if node.Spec.ChiaHealthcheckConfig.Service.Annotations != nil {
		additionalServiceAnnotations = node.Spec.ChiaHealthcheckConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(node.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleHeadlessPeerService assembles the headless peer Service for a Chianode CR
func assembleHeadlessPeerService(node k8schianetv1.ChiaNode, fullNodePort int32) corev1.Service {
	srv := assemblePeerService(node, fullNodePort)

	srv.Name = srv.Name + "-headless"
	srv.Annotations = node.Spec.AdditionalMetadata.Annotations // Overwrites the annotations from the peer Service, since those may contain some related to tools like external-dns
	srv.Spec.Type = corev1.ServiceTypeClusterIP
	srv.Spec.ClusterIP = "None"
	srv.Spec.ExternalTrafficPolicy = ""

	return srv
}

// assembleHeadlessPeerService assembles the headless peer Service for a Chianode CR
func assembleLocalPeerService(node k8schianetv1.ChiaNode, fullNodePort int32) corev1.Service {
	srv := assemblePeerService(node, fullNodePort)

	srv.Name = srv.Name + "-internal"
	srv.Annotations = node.Spec.AdditionalMetadata.Annotations // Overwrites the annotations from the peer Service, since those may contain some related to tools like external-dns
	srv.Spec.Type = corev1.ServiceTypeClusterIP
	local := corev1.ServiceInternalTrafficPolicyLocal
	srv.Spec.InternalTrafficPolicy = &local
	srv.Spec.ExternalTrafficPolicy = ""

	return srv
}

// assembleStatefulset assembles the node StatefulSet resource for a ChiaNode CR
func assembleStatefulset(ctx context.Context, node k8schianetv1.ChiaNode, fullNodePort int32, networkData *map[string]string) (appsv1.StatefulSet, error) {
	vols, volClaimTemplates := getChiaVolumesAndTemplates(node)

	stateful := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(chianodeNamePattern, node.Name),
			Namespace:   node.Namespace,
			Labels:      kube.GetCommonLabels(node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels),
			Annotations: node.Spec.AdditionalMetadata.Annotations,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &node.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(node.Kind, node.ObjectMeta),
			},
			ServiceName: fmt.Sprintf(chianodeNamePattern, node.Name) + "-headless",
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels),
					Annotations: node.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					Affinity:     node.Spec.Affinity,
					NodeSelector: node.Spec.NodeSelector,
					Volumes:      vols,
				},
			},
			VolumeClaimTemplates: volClaimTemplates,
		},
	}

	if node.Spec.ServiceAccountName != nil && *node.Spec.ServiceAccountName != "" {
		stateful.Spec.Template.Spec.ServiceAccountName = *node.Spec.ServiceAccountName
	}

	chiaContainer, err := assembleChiaContainer(ctx, node, fullNodePort, networkData)
	if err != nil {
		return appsv1.StatefulSet{}, err
	}
	stateful.Spec.Template.Spec.Containers = append(stateful.Spec.Template.Spec.Containers, chiaContainer)

	// Get Init Containers
	stateful.Spec.Template.Spec.InitContainers = kube.GetExtraContainers(node.Spec.InitContainers, chiaContainer)
	// Add Init Container Volumes
	for _, init := range node.Spec.InitContainers {
		stateful.Spec.Template.Spec.Volumes = append(stateful.Spec.Template.Spec.Volumes, init.Volumes...)
	}

	// Get Sidecar Containers
	stateful.Spec.Template.Spec.Containers = append(stateful.Spec.Template.Spec.Containers, kube.GetExtraContainers(node.Spec.Sidecars, chiaContainer)...)
	// Add Sidecar Container Volumes
	for _, sidecar := range node.Spec.Sidecars {
		stateful.Spec.Template.Spec.Volumes = append(stateful.Spec.Template.Spec.Volumes, sidecar.Volumes...)
	}

	if node.Spec.ImagePullSecrets != nil && len(*node.Spec.ImagePullSecrets) != 0 {
		stateful.Spec.Template.Spec.ImagePullSecrets = *node.Spec.ImagePullSecrets
	}

	if node.Spec.ChiaExporterConfig.Enabled {
		stateful.Spec.Template.Spec.Containers = append(stateful.Spec.Template.Spec.Containers, assembleChiaExporterContainer(node))
	}

	if node.Spec.ChiaHealthcheckConfig.Enabled {
		stateful.Spec.Template.Spec.Containers = append(stateful.Spec.Template.Spec.Containers, assembleChiaHealthcheckContainer(node))
	}

	if node.Spec.UpdateStrategy != nil {
		stateful.Spec.UpdateStrategy = *node.Spec.UpdateStrategy
	}

	if node.Spec.PodSecurityContext != nil {
		stateful.Spec.Template.Spec.SecurityContext = node.Spec.PodSecurityContext
	}

	// TODO add pod tolerations

	return stateful, nil
}

func assembleChiaContainer(ctx context.Context, node k8schianetv1.ChiaNode, fullNodePort int32, networkData *map[string]string) (corev1.Container, error) {
	input := kube.AssembleChiaContainerInputs{
		Image:           node.Spec.ChiaConfig.Image,
		ImagePullPolicy: node.Spec.ImagePullPolicy,
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
				ContainerPort: consts.NodeRPCPort,
				Protocol:      "TCP",
			},
		},
		VolumeMounts: getChiaVolumeMounts(),
	}

	env, err := getChiaEnv(ctx, node, networkData)
	if err != nil {
		return corev1.Container{}, err
	}
	input.Env = env

	if node.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = node.Spec.ChiaConfig.SecurityContext
	}

	if node.Spec.ChiaConfig.LivenessProbe != nil {
		input.LivenessProbe = node.Spec.ChiaConfig.LivenessProbe
	} else if node.Spec.ChiaHealthcheckConfig.Enabled {
		input.LivenessProbe = kube.AssembleChiaHealthcheckProbe(kube.AssembleChiaHealthcheckProbeInputs{
			Path: "/full_node",
		})
	}

	if node.Spec.ChiaConfig.ReadinessProbe != nil {
		input.ReadinessProbe = node.Spec.ChiaConfig.ReadinessProbe
	} else if node.Spec.ChiaHealthcheckConfig.Enabled {
		input.ReadinessProbe = kube.AssembleChiaHealthcheckProbe(kube.AssembleChiaHealthcheckProbeInputs{
			Path: "/full_node/readiness",
		})
	}

	if node.Spec.ChiaConfig.StartupProbe != nil {
		input.StartupProbe = node.Spec.ChiaConfig.StartupProbe
	} else if node.Spec.ChiaHealthcheckConfig.Enabled {
		failThresh := int32(30)
		periodSec := int32(10)
		input.StartupProbe = kube.AssembleChiaHealthcheckProbe(kube.AssembleChiaHealthcheckProbeInputs{
			Path:             "/full_node/readiness",
			FailureThreshold: &failThresh,
			PeriodSeconds:    &periodSec,
		})
	}

	if node.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = node.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaContainer(input), nil
}

func assembleChiaExporterContainer(node k8schianetv1.ChiaNode) corev1.Container {
	input := kube.AssembleChiaExporterContainerInputs{
		Image:            node.Spec.ChiaExporterConfig.Image,
		ConfigSecretName: node.Spec.ChiaExporterConfig.ConfigSecretName,
		ImagePullPolicy:  node.Spec.ImagePullPolicy,
	}

	if node.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = node.Spec.ChiaConfig.SecurityContext
	}

	if node.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = *node.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaExporterContainer(input)
}

func assembleChiaHealthcheckContainer(node k8schianetv1.ChiaNode) corev1.Container {
	input := kube.AssembleChiaHealthcheckContainerInputs{
		Image:           node.Spec.ChiaHealthcheckConfig.Image,
		ImagePullPolicy: node.Spec.ImagePullPolicy,
	}

	if node.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = node.Spec.ChiaConfig.SecurityContext
	}

	if node.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = *node.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaHealthcheckContainer(input)
}
