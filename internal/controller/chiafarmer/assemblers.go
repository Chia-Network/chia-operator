/*
Copyright 2023 Chia Network Inc.
*/

package chiafarmer

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

const chiafarmerNamePattern = "%s-farmer"

// assembleBaseService assembles the main Service resource for a Chiafarmer CR
func (r *ChiaFarmerReconciler) assembleBaseService(ctx context.Context, farmer k8schianetv1.ChiaFarmer) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiafarmerNamePattern, farmer.Name),
			Namespace:       farmer.Namespace,
			Labels:          kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels),
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
			Selector: kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaFarmer CR
func (r *ChiaFarmerReconciler) assembleChiaExporterService(ctx context.Context, farmer k8schianetv1.ChiaFarmer) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiafarmerNamePattern, farmer.Name) + "-metrics",
			Namespace:       farmer.Namespace,
			Labels:          kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels, farmer.Spec.ChiaExporterConfig.ServiceLabels),
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
			Selector: kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleDeployment assembles the farmer Deployment resource for a ChiaFarmer CR
func (r *ChiaFarmerReconciler) assembleDeployment(ctx context.Context, farmer k8schianetv1.ChiaFarmer) appsv1.Deployment {
	var deploy appsv1.Deployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chiafarmerNamePattern, farmer.Name),
			Namespace:       farmer.Namespace,
			Labels:          kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels),
			Annotations:     farmer.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, farmer),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(ctx, farmer.Kind, farmer.ObjectMeta, farmer.Spec.AdditionalMetadata.Labels),
					Annotations: farmer.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					// TODO add: imagePullSecret, serviceAccountName config
					Containers: []corev1.Container{
						{
							Name:            "chia",
							Image:           farmer.Spec.ChiaConfig.Image,
							ImagePullPolicy: farmer.Spec.ImagePullPolicy,
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

	var containerSecurityContext *corev1.SecurityContext
	if farmer.Spec.ChiaConfig.SecurityContext != nil {
		containerSecurityContext = farmer.Spec.ChiaConfig.SecurityContext
		deploy.Spec.Template.Spec.Containers[0].SecurityContext = farmer.Spec.ChiaConfig.SecurityContext
	}

	if farmer.Spec.ChiaConfig.LivenessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].LivenessProbe = farmer.Spec.ChiaConfig.LivenessProbe
	}

	if farmer.Spec.ChiaConfig.ReadinessProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].ReadinessProbe = farmer.Spec.ChiaConfig.ReadinessProbe
	}

	if farmer.Spec.ChiaConfig.StartupProbe != nil {
		deploy.Spec.Template.Spec.Containers[0].StartupProbe = farmer.Spec.ChiaConfig.StartupProbe
	}

	var containerResorces corev1.ResourceRequirements
	if farmer.Spec.ChiaConfig.Resources != nil {
		containerResorces = *farmer.Spec.ChiaConfig.Resources
		deploy.Spec.Template.Spec.Containers[0].Resources = *farmer.Spec.ChiaConfig.Resources
	}

	if farmer.Spec.ChiaExporterConfig.Enabled {
		exporterContainer := kube.GetChiaExporterContainer(ctx, farmer.Spec.ChiaExporterConfig.Image, containerSecurityContext, farmer.Spec.ImagePullPolicy, containerResorces)
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, exporterContainer)
	}

	if farmer.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = farmer.Spec.PodSecurityContext
	}

	// TODO add pod affinity, tolerations

	return deploy
}
