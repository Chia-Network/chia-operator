package chiawallet

import (
	"context"
	"fmt"

	k8schianetv1 "github.com/chia-network/chia-operator/api/v1"
	"github.com/chia-network/chia-operator/internal/controller/common/consts"
	"github.com/chia-network/chia-operator/internal/controller/common/kube"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const chiawalletNamePattern = "%s-wallet"

// assembleBaseService reconciles the main Service resource for a ChiaWallet CR
func (r *ChiaWalletReconciler) assembleBaseService(ctx context.Context, wallet k8schianetv1.ChiaWallet) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiawalletNamePattern, wallet.Name),
			Namespace:       wallet.Namespace,
			Labels:          r.getLabels(ctx, wallet, wallet.Spec.AdditionalMetadata.Labels),
			Annotations:     wallet.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, wallet),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(wallet.Spec.ServiceType),
			Ports: []corev1.ServicePort{
				{
					Port:       consts.DaemonPort,
					TargetPort: intstr.FromString("daemon"),
					Protocol:   "TCP",
					Name:       "daemon",
				},
				{
					Port:       consts.WalletPort,
					TargetPort: intstr.FromString("peers"),
					Protocol:   "TCP",
					Name:       "peers",
				},
				{
					Port:       consts.WalletRPCPort,
					TargetPort: intstr.FromString("rpc"),
					Protocol:   "TCP",
					Name:       "rpc",
				},
			},
			Selector: r.getLabels(ctx, wallet, wallet.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaWallet CR
func (r *ChiaWalletReconciler) assembleChiaExporterService(ctx context.Context, wallet k8schianetv1.ChiaWallet) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiawalletNamePattern, wallet.Name) + "-metrics",
			Namespace:       wallet.Namespace,
			Labels:          r.getLabels(ctx, wallet, wallet.Spec.AdditionalMetadata.Labels, wallet.Spec.ChiaExporterConfig.ServiceLabels),
			Annotations:     wallet.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, wallet),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType("ClusterIP"),
			Ports: []corev1.ServicePort{
				{
					Port:       consts.ChiaExporterPort,
					TargetPort: intstr.FromString("metrics"),
					Protocol:   "TCP",
					Name:       "metrics",
				},
			},
			Selector: r.getLabels(ctx, wallet, wallet.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleDeployment reconciles the wallet Deployment resource for a ChiaWallet CR
func (r *ChiaWalletReconciler) assembleDeployment(ctx context.Context, wallet k8schianetv1.ChiaWallet) appsv1.Deployment {
	var chiaSecContext *corev1.SecurityContext
	if wallet.Spec.ChiaConfig.SecurityContext != nil {
		chiaSecContext = wallet.Spec.ChiaConfig.SecurityContext
	}

	var chiaLivenessProbe *corev1.Probe
	if wallet.Spec.ChiaConfig.LivenessProbe != nil {
		chiaLivenessProbe = wallet.Spec.ChiaConfig.LivenessProbe
	}

	var chiaReadinessProbe *corev1.Probe
	if wallet.Spec.ChiaConfig.ReadinessProbe != nil {
		chiaReadinessProbe = wallet.Spec.ChiaConfig.ReadinessProbe
	}

	var chiaStartupProbe *corev1.Probe
	if wallet.Spec.ChiaConfig.StartupProbe != nil {
		chiaStartupProbe = wallet.Spec.ChiaConfig.StartupProbe
	}

	var chiaResources corev1.ResourceRequirements
	if wallet.Spec.ChiaConfig.Resources != nil {
		chiaResources = *wallet.Spec.ChiaConfig.Resources
	}

	var imagePullPolicy corev1.PullPolicy
	if wallet.Spec.ImagePullPolicy != nil {
		imagePullPolicy = *wallet.Spec.ImagePullPolicy
	}

	var chiaExporterImage = wallet.Spec.ChiaExporterConfig.Image
	if chiaExporterImage == "" {
		chiaExporterImage = consts.DefaultChiaExporterImage
	}

	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiawalletNamePattern, wallet.Name),
			Namespace:       wallet.Namespace,
			Labels:          r.getLabels(ctx, wallet, wallet.Spec.AdditionalMetadata.Labels),
			Annotations:     wallet.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, wallet),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: r.getLabels(ctx, wallet),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      r.getLabels(ctx, wallet, wallet.Spec.AdditionalMetadata.Labels),
					Annotations: wallet.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					// TODO add: imagePullSecret, serviceAccountName config
					Containers: []corev1.Container{
						{
							Name:            "chia",
							SecurityContext: chiaSecContext,
							Image:           wallet.Spec.ChiaConfig.Image,
							ImagePullPolicy: imagePullPolicy,
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
							LivenessProbe:  chiaLivenessProbe,
							ReadinessProbe: chiaReadinessProbe,
							StartupProbe:   chiaStartupProbe,
							Resources:      chiaResources,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "secret-ca",
									MountPath: "/chia-ca",
								},
								{
									Name:      "key",
									MountPath: "/key",
								},
								{
									Name:      "chiaroot",
									MountPath: "/chia-data",
								},
							},
						},
					},
					NodeSelector: wallet.Spec.NodeSelector,
					Volumes:      r.getChiaVolumes(ctx, wallet),
				},
			},
		},
	}

	exporterContainer := kube.GetChiaExporterContainer(ctx, chiaExporterImage, chiaSecContext, imagePullPolicy, chiaResources)
	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, exporterContainer)

	if wallet.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = wallet.Spec.PodSecurityContext
	}

	// TODO add pod affinity, tolerations

	return deploy
}
