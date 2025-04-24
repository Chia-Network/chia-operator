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
func assemblePeerService(tl k8schianetv1.ChiaTimelord) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiatimelordNamePattern, tl.Name),
		Namespace: tl.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       consts.TimelordPort,
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
		},
	}

	inputs.ServiceType = tl.Spec.ChiaConfig.PeerService.ServiceType
	inputs.ExternalTrafficPolicy = tl.Spec.ChiaConfig.PeerService.ExternalTrafficPolicy
	inputs.SessionAffinity = tl.Spec.ChiaConfig.PeerService.SessionAffinity
	inputs.SessionAffinityConfig = tl.Spec.ChiaConfig.PeerService.SessionAffinityConfig
	inputs.IPFamilyPolicy = tl.Spec.ChiaConfig.PeerService.IPFamilyPolicy
	inputs.IPFamilies = tl.Spec.ChiaConfig.PeerService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if tl.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = tl.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if tl.Spec.ChiaConfig.PeerService.Annotations != nil {
		additionalServiceAnnotations = tl.Spec.ChiaConfig.PeerService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(tl.Spec.Annotations, additionalServiceAnnotations)

	// Handle the Service rollup feature
	if kube.ShouldMakeService(tl.Spec.ChiaHealthcheckConfig.Service, false) && kube.ShouldRollIntoMainPeerService(tl.Spec.ChiaHealthcheckConfig.Service) {
		inputs.Ports = append(inputs.Ports, kube.GetChiaHealthcheckServicePorts()...)
	}

	return kube.AssembleCommonService(inputs)
}

// assembleAllService assembles the all-port Service resource for a ChiaTimelord CR
func assembleAllService(timelord k8schianetv1.ChiaTimelord) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiatimelordNamePattern, timelord.Name) + "-all",
		Namespace: timelord.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       consts.TimelordPort,
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
			{
				Port:       consts.TimelordRPCPort,
				TargetPort: intstr.FromString("rpc"),
				Protocol:   "TCP",
				Name:       "rpc",
			},
		},
	}
	inputs.Ports = append(inputs.Ports, kube.GetChiaDaemonServicePorts()...)

	inputs.ServiceType = timelord.Spec.ChiaConfig.AllService.ServiceType
	inputs.ExternalTrafficPolicy = timelord.Spec.ChiaConfig.AllService.ExternalTrafficPolicy
	inputs.SessionAffinity = timelord.Spec.ChiaConfig.AllService.SessionAffinity
	inputs.SessionAffinityConfig = timelord.Spec.ChiaConfig.AllService.SessionAffinityConfig
	inputs.IPFamilyPolicy = timelord.Spec.ChiaConfig.AllService.IPFamilyPolicy
	inputs.IPFamilies = timelord.Spec.ChiaConfig.AllService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if timelord.Spec.ChiaConfig.AllService.Labels != nil {
		additionalServiceLabels = timelord.Spec.ChiaConfig.AllService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(timelord.Kind, timelord.ObjectMeta, timelord.Spec.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(timelord.Kind, timelord.ObjectMeta, timelord.Spec.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if timelord.Spec.ChiaConfig.AllService.Annotations != nil {
		additionalServiceAnnotations = timelord.Spec.ChiaConfig.AllService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(timelord.Spec.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaTimelord CR
func assembleDaemonService(tl k8schianetv1.ChiaTimelord) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiatimelordNamePattern, tl.Name) + "-daemon",
		Namespace: tl.Namespace,
		Ports:     kube.GetChiaDaemonServicePorts(),
	}

	inputs.ServiceType = tl.Spec.ChiaConfig.DaemonService.ServiceType
	inputs.ExternalTrafficPolicy = tl.Spec.ChiaConfig.DaemonService.ExternalTrafficPolicy
	inputs.SessionAffinity = tl.Spec.ChiaConfig.DaemonService.SessionAffinity
	inputs.SessionAffinityConfig = tl.Spec.ChiaConfig.DaemonService.SessionAffinityConfig
	inputs.IPFamilyPolicy = tl.Spec.ChiaConfig.DaemonService.IPFamilyPolicy
	inputs.IPFamilies = tl.Spec.ChiaConfig.DaemonService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if tl.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = tl.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if tl.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = tl.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(tl.Spec.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleRPCService assembles the RPC Service resource for a ChiaTimelord CR
func assembleRPCService(tl k8schianetv1.ChiaTimelord) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiatimelordNamePattern, tl.Name) + "-rpc",
		Namespace: tl.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       consts.TimelordRPCPort,
				TargetPort: intstr.FromString("rpc"),
				Protocol:   "TCP",
				Name:       "rpc",
			},
		},
	}

	inputs.ServiceType = tl.Spec.ChiaConfig.RPCService.ServiceType
	inputs.ExternalTrafficPolicy = tl.Spec.ChiaConfig.RPCService.ExternalTrafficPolicy
	inputs.SessionAffinity = tl.Spec.ChiaConfig.RPCService.SessionAffinity
	inputs.SessionAffinityConfig = tl.Spec.ChiaConfig.RPCService.SessionAffinityConfig
	inputs.IPFamilyPolicy = tl.Spec.ChiaConfig.RPCService.IPFamilyPolicy
	inputs.IPFamilies = tl.Spec.ChiaConfig.RPCService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if tl.Spec.ChiaConfig.RPCService.Labels != nil {
		additionalServiceLabels = tl.Spec.ChiaConfig.RPCService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if tl.Spec.ChiaConfig.RPCService.Annotations != nil {
		additionalServiceAnnotations = tl.Spec.ChiaConfig.RPCService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(tl.Spec.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaTimelord CR
func assembleChiaExporterService(tl k8schianetv1.ChiaTimelord) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiatimelordNamePattern, tl.Name) + "-metrics",
		Namespace: tl.Namespace,
		Ports:     kube.GetChiaExporterServicePorts(),
	}

	inputs.ServiceType = tl.Spec.ChiaExporterConfig.Service.ServiceType
	inputs.ExternalTrafficPolicy = tl.Spec.ChiaExporterConfig.Service.ExternalTrafficPolicy
	inputs.SessionAffinity = tl.Spec.ChiaExporterConfig.Service.SessionAffinity
	inputs.SessionAffinityConfig = tl.Spec.ChiaExporterConfig.Service.SessionAffinityConfig
	inputs.IPFamilyPolicy = tl.Spec.ChiaExporterConfig.Service.IPFamilyPolicy
	inputs.IPFamilies = tl.Spec.ChiaExporterConfig.Service.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if tl.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = tl.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if tl.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = tl.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(tl.Spec.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleChiaHealthcheckService assembles the chia-healthcheck Service resource
func assembleChiaHealthcheckService(tl k8schianetv1.ChiaTimelord) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiatimelordNamePattern, tl.Name) + "-healthcheck",
		Namespace: tl.Namespace,
		Ports:     kube.GetChiaHealthcheckServicePorts(),
	}

	inputs.ServiceType = tl.Spec.ChiaHealthcheckConfig.Service.ServiceType
	inputs.ExternalTrafficPolicy = tl.Spec.ChiaHealthcheckConfig.Service.ExternalTrafficPolicy
	inputs.SessionAffinity = tl.Spec.ChiaHealthcheckConfig.Service.SessionAffinity
	inputs.SessionAffinityConfig = tl.Spec.ChiaHealthcheckConfig.Service.SessionAffinityConfig
	inputs.IPFamilyPolicy = tl.Spec.ChiaHealthcheckConfig.Service.IPFamilyPolicy
	inputs.IPFamilies = tl.Spec.ChiaHealthcheckConfig.Service.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if tl.Spec.ChiaHealthcheckConfig.Service.Labels != nil {
		additionalServiceLabels = tl.Spec.ChiaHealthcheckConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if tl.Spec.ChiaHealthcheckConfig.Service.Annotations != nil {
		additionalServiceAnnotations = tl.Spec.ChiaHealthcheckConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(tl.Spec.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleVolumeClaim assembles the PVC resource for a ChiaTimelord CR
func assembleVolumeClaim(tl k8schianetv1.ChiaTimelord) (*corev1.PersistentVolumeClaim, error) {
	if tl.Spec.Storage == nil || tl.Spec.Storage.ChiaRoot == nil || tl.Spec.Storage.ChiaRoot.PersistentVolumeClaim == nil {
		return nil, nil
	}

	resourceReq, err := resource.ParseQuantity(tl.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return nil, err
	}

	accessModes := []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	if len(tl.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = tl.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes
	}

	pvc := corev1.PersistentVolumeClaim{
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
	}

	return &pvc, nil
}

// assembleDeployment assembles the tl Deployment resource for a ChiaTimelord CR
func assembleDeployment(ctx context.Context, tl k8schianetv1.ChiaTimelord, networkData *map[string]string) (appsv1.Deployment, error) {
	var deploy = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(chiatimelordNamePattern, tl.Name),
			Namespace:   tl.Namespace,
			Labels:      kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.Labels),
			Annotations: tl.Spec.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(tl.Kind, tl.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(tl.Kind, tl.ObjectMeta, tl.Spec.Labels),
					Annotations: tl.Spec.Annotations,
				},
				Spec: corev1.PodSpec{
					Affinity:                  tl.Spec.Affinity,
					TopologySpreadConstraints: tl.Spec.TopologySpreadConstraints,
					NodeSelector:              tl.Spec.NodeSelector,
					Volumes:                   getChiaVolumes(tl),
				},
			},
		},
	}

	if tl.Spec.ServiceAccountName != nil && *tl.Spec.ServiceAccountName != "" {
		deploy.Spec.Template.Spec.ServiceAccountName = *tl.Spec.ServiceAccountName
	}

	chiaContainer, err := assembleChiaContainer(ctx, tl, networkData)
	if err != nil {
		return appsv1.Deployment{}, err
	}
	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, chiaContainer)

	// Get Init Containers
	deploy.Spec.Template.Spec.InitContainers = kube.GetExtraContainers(tl.Spec.InitContainers, chiaContainer)
	// Add Init Container Volumes
	for _, init := range tl.Spec.InitContainers {
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, init.Volumes...)
	}
	// Get Sidecar Containers
	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, kube.GetExtraContainers(tl.Spec.Sidecars, chiaContainer)...)
	// Add Sidecar Container Volumes
	for _, sidecar := range tl.Spec.Sidecars {
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, sidecar.Volumes...)
	}

	if tl.Spec.ImagePullSecrets != nil && len(*tl.Spec.ImagePullSecrets) != 0 {
		deploy.Spec.Template.Spec.ImagePullSecrets = *tl.Spec.ImagePullSecrets
	}

	if kube.ChiaExporterEnabled(tl.Spec.ChiaExporterConfig) {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, assembleChiaExporterContainer(tl))
	}

	if kube.ChiaHealthcheckEnabled(tl.Spec.ChiaHealthcheckConfig) {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, assembleChiaHealthcheckContainer(tl))
	}

	if tl.Spec.Strategy != nil {
		deploy.Spec.Strategy = *tl.Spec.Strategy
	}

	if tl.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = tl.Spec.PodSecurityContext
	}

	// TODO add pod tolerations

	return deploy, nil
}

func assembleChiaContainer(ctx context.Context, tl k8schianetv1.ChiaTimelord, networkData *map[string]string) (corev1.Container, error) {
	input := kube.AssembleChiaContainerInputs{
		Image:           tl.Spec.ChiaConfig.Image,
		ImagePullPolicy: tl.Spec.ImagePullPolicy,
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

	env, err := getChiaEnv(ctx, tl, networkData)
	if err != nil {
		return corev1.Container{}, err
	}
	input.Env = env

	if tl.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = tl.Spec.ChiaConfig.SecurityContext
	}

	if tl.Spec.ChiaConfig.LivenessProbe != nil {
		input.LivenessProbe = tl.Spec.ChiaConfig.LivenessProbe
	} else if kube.ChiaHealthcheckEnabled(tl.Spec.ChiaHealthcheckConfig) {
		input.LivenessProbe = kube.AssembleChiaHealthcheckProbe(kube.AssembleChiaHealthcheckProbeInputs{
			Path: "/timelord",
		})
	}

	if tl.Spec.ChiaConfig.ReadinessProbe != nil {
		input.ReadinessProbe = tl.Spec.ChiaConfig.ReadinessProbe
	} else if kube.ChiaHealthcheckEnabled(tl.Spec.ChiaHealthcheckConfig) {
		input.ReadinessProbe = kube.AssembleChiaHealthcheckProbe(kube.AssembleChiaHealthcheckProbeInputs{
			Path: "/timelord/readiness",
		})
	}

	if tl.Spec.ChiaConfig.StartupProbe != nil {
		input.StartupProbe = tl.Spec.ChiaConfig.StartupProbe
	} else if kube.ChiaHealthcheckEnabled(tl.Spec.ChiaHealthcheckConfig) {
		failThresh := int32(30)
		periodSec := int32(10)
		input.StartupProbe = kube.AssembleChiaHealthcheckProbe(kube.AssembleChiaHealthcheckProbeInputs{
			Path:             "/timelord/readiness",
			FailureThreshold: &failThresh,
			PeriodSeconds:    &periodSec,
		})
	}

	if tl.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = tl.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaContainer(input), nil
}

func assembleChiaExporterContainer(tl k8schianetv1.ChiaTimelord) corev1.Container {
	input := kube.AssembleChiaExporterContainerInputs{
		Image:            tl.Spec.ChiaExporterConfig.Image,
		ConfigSecretName: tl.Spec.ChiaExporterConfig.ConfigSecretName,
		ImagePullPolicy:  tl.Spec.ImagePullPolicy,
	}

	if tl.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = tl.Spec.ChiaConfig.SecurityContext
	}

	if tl.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = *tl.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaExporterContainer(input)
}

func assembleChiaHealthcheckContainer(tl k8schianetv1.ChiaTimelord) corev1.Container {
	input := kube.AssembleChiaHealthcheckContainerInputs{
		Image:           tl.Spec.ChiaHealthcheckConfig.Image,
		ImagePullPolicy: tl.Spec.ImagePullPolicy,
	}

	if tl.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = tl.Spec.ChiaConfig.SecurityContext
	}

	if tl.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = *tl.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaHealthcheckContainer(input)
}
