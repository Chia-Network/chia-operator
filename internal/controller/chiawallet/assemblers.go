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
func (r *ChiaWalletReconciler) assemblePeerService(ctx context.Context, wallet k8schianetv1.ChiaWallet) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiawalletNamePattern, wallet.Name)
	inputs.Namespace = wallet.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, wallet)

	// Service Type
	if wallet.Spec.ChiaConfig.PeerService != nil && wallet.Spec.ChiaConfig.PeerService.ServiceType != nil {
		inputs.ServiceType = *wallet.Spec.ChiaConfig.PeerService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if wallet.Spec.ChiaConfig.PeerService != nil && wallet.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = wallet.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if wallet.Spec.ChiaConfig.PeerService != nil && wallet.Spec.ChiaConfig.PeerService.Annotations != nil {
		additionalServiceAnnotations = wallet.Spec.ChiaConfig.PeerService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(wallet.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       consts.WalletPort,
			TargetPort: intstr.FromString("peers"),
			Protocol:   "TCP",
			Name:       "peers",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaWallet CR
func (r *ChiaWalletReconciler) assembleDaemonService(ctx context.Context, wallet k8schianetv1.ChiaWallet) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiawalletNamePattern, wallet.Name) + "-daemon"
	inputs.Namespace = wallet.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, wallet)

	// Service Type
	if wallet.Spec.ChiaConfig.DaemonService != nil && wallet.Spec.ChiaConfig.DaemonService.ServiceType != nil {
		inputs.ServiceType = *wallet.Spec.ChiaConfig.DaemonService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if wallet.Spec.ChiaConfig.DaemonService != nil && wallet.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = wallet.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if wallet.Spec.ChiaConfig.DaemonService != nil && wallet.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = wallet.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(wallet.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       consts.DaemonPort,
			TargetPort: intstr.FromString("daemon"),
			Protocol:   "TCP",
			Name:       "daemon",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleRPCService assembles the RPC Service resource for a ChiaWallet CR
func (r *ChiaWalletReconciler) assembleRPCService(ctx context.Context, wallet k8schianetv1.ChiaWallet) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiawalletNamePattern, wallet.Name) + "-rpc"
	inputs.Namespace = wallet.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, wallet)

	// Service Type
	if wallet.Spec.ChiaConfig.RPCService != nil && wallet.Spec.ChiaConfig.RPCService.ServiceType != nil {
		inputs.ServiceType = *wallet.Spec.ChiaConfig.RPCService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if wallet.Spec.ChiaConfig.RPCService != nil && wallet.Spec.ChiaConfig.RPCService.Labels != nil {
		additionalServiceLabels = wallet.Spec.ChiaConfig.RPCService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if wallet.Spec.ChiaConfig.RPCService != nil && wallet.Spec.ChiaConfig.RPCService.Annotations != nil {
		additionalServiceAnnotations = wallet.Spec.ChiaConfig.RPCService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(wallet.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       consts.WalletRPCPort,
			TargetPort: intstr.FromString("rpc"),
			Protocol:   "TCP",
			Name:       "rpc",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaWallet CR
func (r *ChiaWalletReconciler) assembleChiaExporterService(ctx context.Context, wallet k8schianetv1.ChiaWallet) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chiawalletNamePattern, wallet.Name) + "-metrics"
	inputs.Namespace = wallet.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, wallet)

	// Service Type
	if wallet.Spec.ChiaExporterConfig.Service != nil && wallet.Spec.ChiaExporterConfig.Service.ServiceType != nil {
		inputs.ServiceType = *wallet.Spec.ChiaExporterConfig.Service.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if wallet.Spec.ChiaExporterConfig.Service != nil && wallet.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = wallet.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if wallet.Spec.ChiaExporterConfig.Service != nil && wallet.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = wallet.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(wallet.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       consts.ChiaExporterPort,
			TargetPort: intstr.FromString("metrics"),
			Protocol:   "TCP",
			Name:       "metrics",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleVolumeClaim assembles the PVC resource for a ChiaWallet CR
func (r *ChiaWalletReconciler) assembleVolumeClaim(ctx context.Context, wallet k8schianetv1.ChiaWallet) (corev1.PersistentVolumeClaim, error) {
	resourceReq, err := resource.ParseQuantity(wallet.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return corev1.PersistentVolumeClaim{}, err
	}

	var accessModes []corev1.PersistentVolumeAccessMode
	if len(wallet.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = wallet.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes
	} else {
		accessModes = []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
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

// assembleDeployment reconciles the wallet Deployment resource for a ChiaWallet CR
func (r *ChiaWalletReconciler) assembleDeployment(ctx context.Context, wallet k8schianetv1.ChiaWallet) appsv1.Deployment {
	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiawalletNamePattern, wallet.Name),
			Namespace:       wallet.Namespace,
			Labels:          kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels),
			Annotations:     wallet.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, wallet),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(ctx, wallet.Kind, wallet.ObjectMeta, wallet.Spec.AdditionalMetadata.Labels),
					Annotations: wallet.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					// TODO add: imagePullSecret, serviceAccountName config
					Containers: []corev1.Container{
						{
							Name:            "chia",
							Image:           wallet.Spec.ChiaConfig.Image,
							ImagePullPolicy: wallet.Spec.ImagePullPolicy,
							Env:             r.getChiaEnv(ctx, wallet),
							Ports: []corev1.ContainerPort{
								{
									Name:          "daemon",
									ContainerPort: consts.DaemonPort,
									Protocol:      "TCP",
								},
								{
									Name:          "peers",
									ContainerPort: consts.WalletPort,
									Protocol:      "TCP",
								},
								{
									Name:          "rpc",
									ContainerPort: consts.WalletRPCPort,
									Protocol:      "TCP",
								},
							},
							VolumeMounts: r.getChiaVolumeMounts(ctx, wallet),
						},
					},
					NodeSelector: wallet.Spec.NodeSelector,
					Volumes:      r.getChiaVolumes(ctx, wallet),
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
				cont.Container.VolumeMounts = r.getChiaVolumeMounts(ctx, wallet)
			}

			// Share chia env if enabled
			if cont.ShareEnv {
				cont.Container.Env = append(cont.Container.Env, r.getChiaEnv(ctx, wallet)...)
			}

			deploy.Spec.Template.Spec.InitContainers = append(deploy.Spec.Template.Spec.InitContainers, cont.Container)
		}
	}

	if wallet.Spec.Strategy != nil {
		deploy.Spec.Strategy = *wallet.Spec.Strategy
	}

	var containerSecurityContext *corev1.SecurityContext
	if wallet.Spec.ChiaConfig.SecurityContext != nil {
		containerSecurityContext = wallet.Spec.ChiaConfig.SecurityContext
		deploy.Spec.Template.Spec.Containers[0].SecurityContext = wallet.Spec.ChiaConfig.SecurityContext
	}

	if wallet.Spec.ChiaConfig.LivenessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].LivenessProbe = wallet.Spec.ChiaConfig.LivenessProbe
	}

	if wallet.Spec.ChiaConfig.ReadinessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].ReadinessProbe = wallet.Spec.ChiaConfig.ReadinessProbe
	}

	if wallet.Spec.ChiaConfig.StartupProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].StartupProbe = wallet.Spec.ChiaConfig.StartupProbe
	}

	var containerResorces corev1.ResourceRequirements
	if wallet.Spec.ChiaConfig.Resources != nil {
		containerResorces = *wallet.Spec.ChiaConfig.Resources
		deploy.Spec.Template.Spec.Containers[0].Resources = *wallet.Spec.ChiaConfig.Resources
	}

	if wallet.Spec.ChiaExporterConfig.Enabled {
		exporterContainer := kube.GetChiaExporterContainer(ctx, wallet.Spec.ChiaExporterConfig.Image, containerSecurityContext, wallet.Spec.ImagePullPolicy, containerResorces)
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, exporterContainer)
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
