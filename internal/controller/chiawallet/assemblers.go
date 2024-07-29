/*
Copyright 2023 Chia Network Inc.
*/

package chiawallet

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const chiawalletNamePattern = "%s-wallet"

// assemblePeerService assembles the peer Service resource for a ChiaWallet CR
func assemblePeerService(wallet k8schianetv1.ChiaWallet) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:           fmt.Sprintf(chiawalletNamePattern, wallet.Name),
		Namespace:      wallet.Namespace,
		OwnerReference: getOwnerReference(wallet),
		Ports: []corev1.ServicePort{
			{
				Port:       consts.WalletPort,
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
		},
	}

	if wallet.Spec.ChiaConfig.PeerService != nil {
		inputs.ServiceType = wallet.Spec.ChiaConfig.PeerService.ServiceType
		inputs.IPFamilyPolicy = wallet.Spec.ChiaConfig.PeerService.IPFamilyPolicy
		inputs.IPFamilies = wallet.Spec.ChiaConfig.PeerService.IPFamilies

		// Labels
		var additionalServiceLabels = make(map[string]string)
		if wallet.Spec.ChiaConfig.PeerService.Labels != nil {
			additionalServiceLabels = wallet.Spec.ChiaConfig.PeerService.Labels
		}
		inputs.Labels = kube.GetCommonLabels(wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
		inputs.SelectorLabels = kube.GetCommonLabels(wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels)

		// Annotations
		var additionalServiceAnnotations = make(map[string]string)
		if wallet.Spec.ChiaConfig.PeerService.Annotations != nil {
			additionalServiceAnnotations = wallet.Spec.ChiaConfig.PeerService.Annotations
		}
		inputs.Annotations = kube.CombineMaps(wallet.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)
	}

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaWallet CR
func assembleDaemonService(wallet k8schianetv1.ChiaWallet) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:           fmt.Sprintf(chiawalletNamePattern, wallet.Name) + "-daemon",
		Namespace:      wallet.Namespace,
		OwnerReference: getOwnerReference(wallet),
		Ports:          kube.GetChiaDaemonServicePorts(),
	}

	if wallet.Spec.ChiaConfig.DaemonService != nil {
		inputs.ServiceType = wallet.Spec.ChiaConfig.DaemonService.ServiceType
		inputs.IPFamilyPolicy = wallet.Spec.ChiaConfig.DaemonService.IPFamilyPolicy
		inputs.IPFamilies = wallet.Spec.ChiaConfig.DaemonService.IPFamilies

		// Labels
		var additionalServiceLabels = make(map[string]string)
		if wallet.Spec.ChiaConfig.DaemonService.Labels != nil {
			additionalServiceLabels = wallet.Spec.ChiaConfig.DaemonService.Labels
		}
		inputs.Labels = kube.GetCommonLabels(wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
		inputs.SelectorLabels = kube.GetCommonLabels(wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels)

		// Annotations
		var additionalServiceAnnotations = make(map[string]string)
		if wallet.Spec.ChiaConfig.DaemonService.Annotations != nil {
			additionalServiceAnnotations = wallet.Spec.ChiaConfig.DaemonService.Annotations
		}
		inputs.Annotations = kube.CombineMaps(wallet.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)
	}

	return kube.AssembleCommonService(inputs)
}

// assembleRPCService assembles the RPC Service resource for a ChiaWallet CR
func assembleRPCService(wallet k8schianetv1.ChiaWallet) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:           fmt.Sprintf(chiawalletNamePattern, wallet.Name) + "-rpc",
		Namespace:      wallet.Namespace,
		OwnerReference: getOwnerReference(wallet),
		Ports: []corev1.ServicePort{
			{
				Port:       consts.WalletRPCPort,
				TargetPort: intstr.FromString("rpc"),
				Protocol:   "TCP",
				Name:       "rpc",
			},
		},
	}

	if wallet.Spec.ChiaConfig.RPCService != nil {
		inputs.ServiceType = wallet.Spec.ChiaConfig.RPCService.ServiceType
		inputs.IPFamilyPolicy = wallet.Spec.ChiaConfig.RPCService.IPFamilyPolicy
		inputs.IPFamilies = wallet.Spec.ChiaConfig.RPCService.IPFamilies

		// Labels
		var additionalServiceLabels = make(map[string]string)
		if wallet.Spec.ChiaConfig.RPCService.Labels != nil {
			additionalServiceLabels = wallet.Spec.ChiaConfig.RPCService.Labels
		}
		inputs.Labels = kube.GetCommonLabels(wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
		inputs.SelectorLabels = kube.GetCommonLabels(wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels)

		// Annotations
		var additionalServiceAnnotations = make(map[string]string)
		if wallet.Spec.ChiaConfig.RPCService.Annotations != nil {
			additionalServiceAnnotations = wallet.Spec.ChiaConfig.RPCService.Annotations
		}
		inputs.Annotations = kube.CombineMaps(wallet.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)
	}

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaWallet CR
func assembleChiaExporterService(wallet k8schianetv1.ChiaWallet) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:           fmt.Sprintf(chiawalletNamePattern, wallet.Name) + "-metrics",
		Namespace:      wallet.Namespace,
		OwnerReference: getOwnerReference(wallet),
		Ports:          kube.GetChiaExporterServicePorts(),
	}

	if wallet.Spec.ChiaExporterConfig.Service != nil {
		inputs.ServiceType = wallet.Spec.ChiaExporterConfig.Service.ServiceType
		inputs.IPFamilyPolicy = wallet.Spec.ChiaExporterConfig.Service.IPFamilyPolicy
		inputs.IPFamilies = wallet.Spec.ChiaExporterConfig.Service.IPFamilies

		// Labels
		var additionalServiceLabels = make(map[string]string)
		if wallet.Spec.ChiaExporterConfig.Service.Labels != nil {
			additionalServiceLabels = wallet.Spec.ChiaExporterConfig.Service.Labels
		}
		inputs.Labels = kube.GetCommonLabels(wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
		inputs.SelectorLabels = kube.GetCommonLabels(wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels)

		// Annotations
		var additionalServiceAnnotations = make(map[string]string)
		if wallet.Spec.ChiaExporterConfig.Service.Annotations != nil {
			additionalServiceAnnotations = wallet.Spec.ChiaExporterConfig.Service.Annotations
		}
		inputs.Annotations = kube.CombineMaps(wallet.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)
	}

	return kube.AssembleCommonService(inputs)
}

// assembleVolumeClaim assembles the PVC resource for a ChiaWallet CR
func assembleVolumeClaim(wallet k8schianetv1.ChiaWallet) (corev1.PersistentVolumeClaim, error) {
	resourceReq, err := resource.ParseQuantity(wallet.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return corev1.PersistentVolumeClaim{}, err
	}

	accessModes := []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	if len(wallet.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = wallet.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes
	}

	return corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(chiawalletNamePattern, wallet.Name),
			Namespace: wallet.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      accessModes,
			StorageClassName: &wallet.Spec.Storage.ChiaRoot.PersistentVolumeClaim.StorageClass,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resourceReq,
				},
			},
		},
	}, nil
}

// assembleDeployment assembles the wallet Deployment resource for a ChiaWallet CR
func assembleDeployment(ctx context.Context, wallet k8schianetv1.ChiaWallet) appsv1.Deployment {
	var deploy = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(chiawalletNamePattern, wallet.Name),
			Namespace:   wallet.Namespace,
			Labels:      kube.GetCommonLabels(wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels),
			Annotations: wallet.Spec.AdditionalMetadata.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(wallet.Kind, wallet.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels),
					Annotations: wallet.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					Containers:   []corev1.Container{assembleChiaContainer(ctx, wallet)},
					NodeSelector: wallet.Spec.NodeSelector,
					Volumes:      getChiaVolumes(wallet),
				},
			},
		},
	}

	if len(wallet.Spec.InitContainers) != 0 {
		// Overwrite any volumeMounts specified in init containers. Not currently supported.
		for _, cont := range wallet.Spec.InitContainers {
			cont.Container.VolumeMounts = []corev1.VolumeMount{}

			// Share chia volume mounts if enabled
			if cont.ShareVolumeMounts {
				cont.Container.VolumeMounts = getChiaVolumeMounts()
			}

			// Share chia env if enabled
			if cont.ShareEnv {
				cont.Container.Env = append(cont.Container.Env, getChiaEnv(ctx, wallet)...)
			}

			deploy.Spec.Template.Spec.InitContainers = append(deploy.Spec.Template.Spec.InitContainers, cont.Container)
		}
	}

	if wallet.Spec.ChiaExporterConfig.Enabled {
		chiaExporterContainer := assembleChiaExporterContainer(wallet)
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, chiaExporterContainer)
	}

	if wallet.Spec.Strategy != nil {
		deploy.Spec.Strategy = *wallet.Spec.Strategy
	}

	if wallet.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = wallet.Spec.PodSecurityContext
	}

	if len(wallet.Spec.Sidecars.Containers) > 0 {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, wallet.Spec.Sidecars.Containers...)
	}

	// TODO add pod affinity, tolerations

	return deploy
}

func assembleChiaContainer(ctx context.Context, wallet k8schianetv1.ChiaWallet) corev1.Container {
	input := kube.AssembleChiaContainerInputs{
		Image:           wallet.Spec.ChiaConfig.Image,
		ImagePullPolicy: wallet.Spec.ImagePullPolicy,
		Env:             getChiaEnv(ctx, wallet),
		Ports:           getChiaPorts(),
		VolumeMounts:    getChiaVolumeMounts(),
	}

	if wallet.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = wallet.Spec.ChiaConfig.SecurityContext
	}

	if wallet.Spec.ChiaConfig.LivenessProbe != nil {
		input.LivenessProbe = wallet.Spec.ChiaConfig.LivenessProbe
	}

	if wallet.Spec.ChiaConfig.ReadinessProbe != nil {
		input.ReadinessProbe = wallet.Spec.ChiaConfig.ReadinessProbe
	}

	if wallet.Spec.ChiaConfig.StartupProbe != nil {
		input.StartupProbe = wallet.Spec.ChiaConfig.StartupProbe
	}

	if wallet.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = wallet.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaContainer(input)
}

func assembleChiaExporterContainer(wallet k8schianetv1.ChiaWallet) corev1.Container {
	input := kube.AssembleChiaExporterContainerInputs{
		Image:            wallet.Spec.ChiaExporterConfig.Image,
		ConfigSecretName: wallet.Spec.ChiaExporterConfig.ConfigSecretName,
		ImagePullPolicy:  wallet.Spec.ImagePullPolicy,
	}

	if wallet.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = wallet.Spec.ChiaConfig.SecurityContext
	}

	if wallet.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = *wallet.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaExporterContainer(input)
}
