package chiafarmer

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

const chiafarmerNamePattern = "%s-farmer"

// assembleBaseService assembles the main Service resource for a Chiafarmer CR
func (r *ChiaFarmerReconciler) assembleBaseService(ctx context.Context, farmer k8schianetv1.ChiaFarmer) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiafarmerNamePattern, farmer.Name),
			Namespace:       farmer.Namespace,
			Labels:          r.getLabels(ctx, farmer, farmer.Spec.AdditionalMetadata.Labels),
			Annotations:     farmer.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, farmer),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(farmer.Spec.ServiceType),
			Ports: []corev1.ServicePort{
				{
					Port:       consts.DaemonPort,
					TargetPort: intstr.FromString("daemon"),
					Protocol:   "TCP",
					Name:       "daemon",
				},
				{
					Port:       consts.FarmerPort,
					TargetPort: intstr.FromString("peers"),
					Protocol:   "TCP",
					Name:       "peers",
				},
				{
					Port:       consts.FarmerRPCPort,
					TargetPort: intstr.FromString("rpc"),
					Protocol:   "TCP",
					Name:       "rpc",
				},
			},
			Selector: r.getLabels(ctx, farmer, farmer.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaFarmer CR
func (r *ChiaFarmerReconciler) assembleChiaExporterService(ctx context.Context, farmer k8schianetv1.ChiaFarmer) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiafarmerNamePattern, farmer.Name) + "-metrics",
			Namespace:       farmer.Namespace,
			Labels:          r.getLabels(ctx, farmer, farmer.Spec.AdditionalMetadata.Labels, farmer.Spec.ChiaExporterConfig.ServiceLabels),
			Annotations:     farmer.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, farmer),
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
			Selector: r.getLabels(ctx, farmer, farmer.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleDeployment assembles the farmer Deployment resource for a ChiaFarmer CR
func (r *ChiaFarmerReconciler) assembleDeployment(ctx context.Context, farmer k8schianetv1.ChiaFarmer) appsv1.Deployment {
	var chiaSecContext *corev1.SecurityContext
	if farmer.Spec.ChiaConfig.SecurityContext != nil {
		chiaSecContext = farmer.Spec.ChiaConfig.SecurityContext
	}

	var chiaLivenessProbe *corev1.Probe
	if farmer.Spec.ChiaConfig.LivenessProbe != nil {
		chiaLivenessProbe = farmer.Spec.ChiaConfig.LivenessProbe
	}

	var chiaReadinessProbe *corev1.Probe
	if farmer.Spec.ChiaConfig.ReadinessProbe != nil {
		chiaReadinessProbe = farmer.Spec.ChiaConfig.ReadinessProbe
	}

	var chiaStartupProbe *corev1.Probe
	if farmer.Spec.ChiaConfig.StartupProbe != nil {
		chiaStartupProbe = farmer.Spec.ChiaConfig.StartupProbe
	}

	var chiaResources corev1.ResourceRequirements
	if farmer.Spec.ChiaConfig.Resources != nil {
		chiaResources = *farmer.Spec.ChiaConfig.Resources
	}

	var imagePullPolicy corev1.PullPolicy
	if farmer.Spec.ImagePullPolicy != nil {
		imagePullPolicy = *farmer.Spec.ImagePullPolicy
	}

	var chiaExporterImage = farmer.Spec.ChiaExporterConfig.Image
	if chiaExporterImage == "" {
		chiaExporterImage = consts.DefaultChiaExporterImage
	}

	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiafarmerNamePattern, farmer.Name),
			Namespace:       farmer.Namespace,
			Labels:          r.getLabels(ctx, farmer, farmer.Spec.AdditionalMetadata.Labels),
			Annotations:     farmer.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, farmer),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: r.getLabels(ctx, farmer),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      r.getLabels(ctx, farmer, farmer.Spec.AdditionalMetadata.Labels),
					Annotations: farmer.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					// TODO add: imagePullSecret, serviceAccountName config
					Containers: []corev1.Container{
						{
							Name:            "chia",
							SecurityContext: chiaSecContext,
							Image:           farmer.Spec.ChiaConfig.Image,
							ImagePullPolicy: imagePullPolicy,
							Env:             r.getChiaEnv(ctx, farmer),
							Ports: []corev1.ContainerPort{
								{
									Name:          "daemon",
									ContainerPort: consts.DaemonPort,
									Protocol:      "TCP",
								},
								{
									Name:          "peers",
									ContainerPort: consts.FarmerPort,
									Protocol:      "TCP",
								},
								{
									Name:          "rpc",
									ContainerPort: consts.FarmerRPCPort,
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
					NodeSelector: farmer.Spec.NodeSelector,
					Volumes:      r.getChiaVolumes(ctx, farmer),
				},
			},
		},
	}

	exporterContainer := kube.GetChiaExporterContainer(ctx, chiaExporterImage, chiaSecContext, imagePullPolicy, chiaResources)
	deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, exporterContainer)

	if farmer.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = farmer.Spec.PodSecurityContext
	}

	// TODO add pod affinity, tolerations

	return deploy
}
