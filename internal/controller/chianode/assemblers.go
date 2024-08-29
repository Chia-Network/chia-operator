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
func assemblePeerService(node k8schianetv1.ChiaNode) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chianodeNamePattern, node.Name),
		Namespace: node.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       kube.GetFullNodePort(node.Spec.ChiaConfig.CommonSpecChia),
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
		},
	}

	inputs.ServiceType = node.Spec.ChiaConfig.PeerService.ServiceType
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

// assembleDaemonService assembles the daemon Service resource for a ChiaNode CR
func assembleDaemonService(node k8schianetv1.ChiaNode) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chianodeNamePattern, node.Name) + "-daemon",
		Namespace: node.Namespace,
		Ports:     kube.GetChiaDaemonServicePorts(),
	}

	inputs.ServiceType = node.Spec.ChiaConfig.DaemonService.ServiceType
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
func assembleHeadlessPeerService(node k8schianetv1.ChiaNode) corev1.Service {
	srv := assemblePeerService(node)

	srv.Name = srv.Name + "-headless"
	srv.Spec.Type = "ClusterIP"
	srv.Spec.ClusterIP = "None"

	return srv
}

// assembleHeadlessPeerService assembles the headless peer Service for a Chianode CR
func assembleLocalPeerService(node k8schianetv1.ChiaNode) corev1.Service {
	srv := assemblePeerService(node)

	srv.Name = srv.Name + "-internal"
	local := corev1.ServiceInternalTrafficPolicyLocal
	srv.Spec.InternalTrafficPolicy = &local

	return srv
}

// assembleStatefulset assembles the node StatefulSet resource for a ChiaNode CR
func assembleStatefulset(ctx context.Context, node k8schianetv1.ChiaNode) appsv1.StatefulSet {
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
					// TODO add: imagePullSecret, serviceAccountName config
					Containers:   []corev1.Container{assembleChiaContainer(ctx, node)},
					Affinity:     node.Spec.Affinity,
					NodeSelector: node.Spec.NodeSelector,
					Volumes:      vols,
				},
			},
			VolumeClaimTemplates: volClaimTemplates,
		},
	}

	if len(node.Spec.InitContainers) != 0 {
		// Overwrite any volumeMounts specified in init containers. Not currently supported.
		for _, cont := range node.Spec.InitContainers {
			cont.Container.VolumeMounts = []corev1.VolumeMount{}

			// Share chia volume mounts if enabled
			if cont.ShareVolumeMounts {
				cont.Container.VolumeMounts = getChiaVolumeMounts()
			}

			// Share chia env if enabled
			if cont.ShareEnv {
				cont.Container.Env = append(cont.Container.Env, getChiaEnv(ctx, node)...)
			}

			stateful.Spec.Template.Spec.InitContainers = append(stateful.Spec.Template.Spec.InitContainers, cont.Container)
		}
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

	if len(node.Spec.Sidecars.Containers) > 0 {
		stateful.Spec.Template.Spec.Containers = append(stateful.Spec.Template.Spec.Containers, node.Spec.Sidecars.Containers...)
	}

	// TODO add pod affinity, tolerations

	return stateful
}

func assembleChiaContainer(ctx context.Context, node k8schianetv1.ChiaNode) corev1.Container {
	input := kube.AssembleChiaContainerInputs{
		Image:           node.Spec.ChiaConfig.Image,
		ImagePullPolicy: node.Spec.ImagePullPolicy,
		Env:             getChiaEnv(ctx, node),
		Ports: []corev1.ContainerPort{
			{
				Name:          "daemon",
				ContainerPort: consts.DaemonPort,
				Protocol:      "TCP",
			},
			{
				Name:          "peers",
				ContainerPort: kube.GetFullNodePort(node.Spec.ChiaConfig.CommonSpecChia),
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

	if node.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = node.Spec.ChiaConfig.SecurityContext
	}

	if node.Spec.ChiaConfig.LivenessProbe != nil {
		input.LivenessProbe = node.Spec.ChiaConfig.LivenessProbe
	} else if node.Spec.ChiaHealthcheckConfig.Enabled {
		input.LivenessProbe = kube.AssembleChiaHealthcheckProbe(kube.AssembleChiaHealthcheckProbeInputs{
			Kind: consts.ChiaNodeKind,
		})
	}

	if node.Spec.ChiaConfig.ReadinessProbe != nil {
		input.ReadinessProbe = node.Spec.ChiaConfig.ReadinessProbe
	} else if node.Spec.ChiaHealthcheckConfig.Enabled {
		input.ReadinessProbe = kube.AssembleChiaHealthcheckProbe(kube.AssembleChiaHealthcheckProbeInputs{
			Kind: consts.ChiaNodeKind,
		})
	}

	if node.Spec.ChiaConfig.StartupProbe != nil {
		input.StartupProbe = node.Spec.ChiaConfig.StartupProbe
	} else if node.Spec.ChiaHealthcheckConfig.Enabled {
		failThresh := int32(30)
		periodSec := int32(10)
		input.StartupProbe = kube.AssembleChiaHealthcheckProbe(kube.AssembleChiaHealthcheckProbeInputs{
			Kind:             consts.ChiaNodeKind,
			FailureThreshold: &failThresh,
			PeriodSeconds:    &periodSec,
		})
	}

	if node.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = node.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaContainer(input)
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
