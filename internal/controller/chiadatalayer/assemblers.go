/*
Copyright 2024 Chia Network Inc.
*/

package chiadatalayer

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/chiadatalayer/fileserver"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
)

const chiadatalayerNamePattern = "%s-datalayer"

// assembleDaemonService assembles the daemon Service resource for a ChiaDataLayer CR
func assembleDaemonService(datalayer k8schianetv1.ChiaDataLayer) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiadatalayerNamePattern, datalayer.Name) + "-daemon",
		Namespace: datalayer.Namespace,
		Ports:     kube.GetChiaDaemonServicePorts(),
	}

	inputs.ServiceType = datalayer.Spec.ChiaConfig.DaemonService.ServiceType
	inputs.ExternalTrafficPolicy = datalayer.Spec.ChiaConfig.DaemonService.ExternalTrafficPolicy
	inputs.SessionAffinity = datalayer.Spec.ChiaConfig.DaemonService.SessionAffinity
	inputs.SessionAffinityConfig = datalayer.Spec.ChiaConfig.DaemonService.SessionAffinityConfig
	inputs.IPFamilyPolicy = datalayer.Spec.ChiaConfig.DaemonService.IPFamilyPolicy
	inputs.IPFamilies = datalayer.Spec.ChiaConfig.DaemonService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if datalayer.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = datalayer.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(datalayer.Kind, datalayer.ObjectMeta, datalayer.Spec.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(datalayer.Kind, datalayer.ObjectMeta, datalayer.Spec.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if datalayer.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = datalayer.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(datalayer.Spec.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleRPCService assembles the RPC Service resource for a ChiaDataLayer CR
func assembleRPCService(datalayer k8schianetv1.ChiaDataLayer) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiadatalayerNamePattern, datalayer.Name) + "-rpc",
		Namespace: datalayer.Namespace,
		Ports: []corev1.ServicePort{
			{
				Port:       consts.DataLayerRPCPort,
				TargetPort: intstr.FromString("rpc"),
				Protocol:   "TCP",
				Name:       "rpc",
			},
			{
				Port:       consts.WalletRPCPort,
				TargetPort: intstr.FromString("wallet-rpc"),
				Protocol:   "TCP",
				Name:       "wallet-rpc",
			},
		},
	}

	inputs.ServiceType = datalayer.Spec.ChiaConfig.RPCService.ServiceType
	inputs.ExternalTrafficPolicy = datalayer.Spec.ChiaConfig.RPCService.ExternalTrafficPolicy
	inputs.SessionAffinity = datalayer.Spec.ChiaConfig.RPCService.SessionAffinity
	inputs.SessionAffinityConfig = datalayer.Spec.ChiaConfig.RPCService.SessionAffinityConfig
	inputs.IPFamilyPolicy = datalayer.Spec.ChiaConfig.RPCService.IPFamilyPolicy
	inputs.IPFamilies = datalayer.Spec.ChiaConfig.RPCService.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if datalayer.Spec.ChiaConfig.RPCService.Labels != nil {
		additionalServiceLabels = datalayer.Spec.ChiaConfig.RPCService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(datalayer.Kind, datalayer.ObjectMeta, datalayer.Spec.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(datalayer.Kind, datalayer.ObjectMeta, datalayer.Spec.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if datalayer.Spec.ChiaConfig.RPCService.Annotations != nil {
		additionalServiceAnnotations = datalayer.Spec.ChiaConfig.RPCService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(datalayer.Spec.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaDataLayer CR
func assembleChiaExporterService(datalayer k8schianetv1.ChiaDataLayer) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:      fmt.Sprintf(chiadatalayerNamePattern, datalayer.Name) + "-metrics",
		Namespace: datalayer.Namespace,
		Ports:     kube.GetChiaExporterServicePorts(),
	}

	inputs.ServiceType = datalayer.Spec.ChiaExporterConfig.Service.ServiceType
	inputs.ExternalTrafficPolicy = datalayer.Spec.ChiaExporterConfig.Service.ExternalTrafficPolicy
	inputs.SessionAffinity = datalayer.Spec.ChiaExporterConfig.Service.SessionAffinity
	inputs.SessionAffinityConfig = datalayer.Spec.ChiaExporterConfig.Service.SessionAffinityConfig
	inputs.IPFamilyPolicy = datalayer.Spec.ChiaExporterConfig.Service.IPFamilyPolicy
	inputs.IPFamilies = datalayer.Spec.ChiaExporterConfig.Service.IPFamilies

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if datalayer.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = datalayer.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(datalayer.Kind, datalayer.ObjectMeta, datalayer.Spec.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(datalayer.Kind, datalayer.ObjectMeta, datalayer.Spec.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if datalayer.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = datalayer.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(datalayer.Spec.Annotations, additionalServiceAnnotations)

	return kube.AssembleCommonService(inputs)
}

// assembleChiaRootVolumeClaim assembles the CHIA_ROOT PVC resource for a ChiaDataLayer CR
func assembleChiaRootVolumeClaim(datalayer k8schianetv1.ChiaDataLayer) (*corev1.PersistentVolumeClaim, error) {
	if datalayer.Spec.Storage == nil || datalayer.Spec.Storage.ChiaRoot == nil || datalayer.Spec.Storage.ChiaRoot.PersistentVolumeClaim == nil {
		return nil, nil
	}

	resourceReq, err := resource.ParseQuantity(datalayer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return nil, err
	}

	accessModes := []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	if len(datalayer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = datalayer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes
	}

	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(chiadatalayerNamePattern, datalayer.Name),
			Namespace: datalayer.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      accessModes,
			StorageClassName: &datalayer.Spec.Storage.ChiaRoot.PersistentVolumeClaim.StorageClass,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resourceReq,
				},
			},
		},
	}

	return &pvc, nil
}

// assembleDataLayerFilesVolumeClaim assembles the data_layer server files PVC resource for a ChiaDataLayer CR
func assembleDataLayerFilesVolumeClaim(datalayer k8schianetv1.ChiaDataLayer) (*corev1.PersistentVolumeClaim, error) {
	if datalayer.Spec.Storage == nil || datalayer.Spec.Storage.DataLayerServerFiles == nil || datalayer.Spec.Storage.DataLayerServerFiles.PersistentVolumeClaim == nil {
		return nil, nil
	}

	resourceReq, err := resource.ParseQuantity(datalayer.Spec.Storage.DataLayerServerFiles.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return nil, err
	}

	accessModes := []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	if len(datalayer.Spec.Storage.DataLayerServerFiles.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = datalayer.Spec.Storage.DataLayerServerFiles.PersistentVolumeClaim.AccessModes
	}

	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(chiadatalayerNamePattern, datalayer.Name) + "-server",
			Namespace: datalayer.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      accessModes,
			StorageClassName: &datalayer.Spec.Storage.DataLayerServerFiles.PersistentVolumeClaim.StorageClass,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resourceReq,
				},
			},
		},
	}

	return &pvc, nil
}

// assembleDeployment assembles the datalayer Deployment resource for a ChiaDataLayer CR
func assembleDeployment(ctx context.Context, datalayer k8schianetv1.ChiaDataLayer, networkData *map[string]string) (appsv1.Deployment, error) {
	var deploy = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(chiadatalayerNamePattern, datalayer.Name),
			Namespace:   datalayer.Namespace,
			Labels:      kube.GetCommonLabels(datalayer.Kind, datalayer.ObjectMeta, datalayer.Spec.Labels),
			Annotations: datalayer.Spec.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(datalayer.Kind, datalayer.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(datalayer.Kind, datalayer.ObjectMeta, datalayer.Spec.Labels),
					Annotations: datalayer.Spec.Annotations,
				},
				Spec: corev1.PodSpec{
					Affinity:                  datalayer.Spec.Affinity,
					TopologySpreadConstraints: datalayer.Spec.TopologySpreadConstraints,
					NodeSelector:              datalayer.Spec.NodeSelector,
					Volumes:                   getChiaVolumes(datalayer),
				},
			},
		},
	}

	if datalayer.Spec.ServiceAccountName != nil && *datalayer.Spec.ServiceAccountName != "" {
		deploy.Spec.Template.Spec.ServiceAccountName = *datalayer.Spec.ServiceAccountName
	}

	chiaContainer, err := assembleChiaContainer(ctx, datalayer, networkData)
	if err != nil {
		return appsv1.Deployment{}, err
	}
	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, chiaContainer)

	// Get Init Containers
	deploy.Spec.Template.Spec.InitContainers = kube.GetExtraContainers(datalayer.Spec.InitContainers, chiaContainer)
	// Add Init Container Volumes
	for _, init := range datalayer.Spec.InitContainers {
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, init.Volumes...)
	}
	// Get Sidecar Containers
	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, kube.GetExtraContainers(datalayer.Spec.Sidecars, chiaContainer)...)
	// Add Sidecar Container Volumes
	for _, sidecar := range datalayer.Spec.Sidecars {
		deploy.Spec.Template.Spec.Volumes = append(deploy.Spec.Template.Spec.Volumes, sidecar.Volumes...)
	}

	if datalayer.Spec.ImagePullSecrets != nil && len(*datalayer.Spec.ImagePullSecrets) != 0 {
		deploy.Spec.Template.Spec.ImagePullSecrets = *datalayer.Spec.ImagePullSecrets
	}

	if kube.ChiaExporterEnabled(datalayer.Spec.ChiaExporterConfig) {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, assembleChiaExporterContainer(datalayer))
	}

	if datalayer.Spec.FileserverConfig.Enabled != nil && *datalayer.Spec.FileserverConfig.Enabled {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, fileserver.AssembleContainer(datalayer))
	}

	if datalayer.Spec.Strategy != nil {
		deploy.Spec.Strategy = *datalayer.Spec.Strategy
	}

	if datalayer.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = datalayer.Spec.PodSecurityContext
	}

	// TODO add pod tolerations

	return deploy, nil
}

func assembleChiaContainer(ctx context.Context, datalayer k8schianetv1.ChiaDataLayer, networkData *map[string]string) (corev1.Container, error) {
	input := kube.AssembleChiaContainerInputs{
		Image:           datalayer.Spec.ChiaConfig.Image,
		ImagePullPolicy: datalayer.Spec.ImagePullPolicy,
		Ports:           getChiaPorts(),
		VolumeMounts:    getChiaVolumeMounts(datalayer),
	}

	env, err := getChiaEnv(ctx, datalayer, networkData)
	if err != nil {
		return corev1.Container{}, err
	}
	input.Env = env

	if datalayer.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = datalayer.Spec.ChiaConfig.SecurityContext
	}

	if datalayer.Spec.ChiaConfig.LivenessProbe != nil {
		input.LivenessProbe = datalayer.Spec.ChiaConfig.LivenessProbe
	}

	if datalayer.Spec.ChiaConfig.ReadinessProbe != nil {
		input.ReadinessProbe = datalayer.Spec.ChiaConfig.ReadinessProbe
	}

	if datalayer.Spec.ChiaConfig.StartupProbe != nil {
		input.StartupProbe = datalayer.Spec.ChiaConfig.StartupProbe
	}

	if datalayer.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = datalayer.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaContainer(input), nil
}

func assembleChiaExporterContainer(datalayer k8schianetv1.ChiaDataLayer) corev1.Container {
	input := kube.AssembleChiaExporterContainerInputs{
		Image:            datalayer.Spec.ChiaExporterConfig.Image,
		ConfigSecretName: datalayer.Spec.ChiaExporterConfig.ConfigSecretName,
		ImagePullPolicy:  datalayer.Spec.ImagePullPolicy,
	}

	if datalayer.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = datalayer.Spec.ChiaConfig.SecurityContext
	}

	if datalayer.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = *datalayer.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaExporterContainer(input)
}
