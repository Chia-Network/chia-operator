package chianode

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

const chianodeNamePattern = "%s-node"

// assembleBaseService assembles the main Service resource for a ChiaNode CR
func (r *ChiaNodeReconciler) assembleBaseService(ctx context.Context, node k8schianetv1.ChiaNode) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chianodeNamePattern, node.Name),
			Namespace:       node.Namespace,
			Labels:          r.getLabels(ctx, node, node.Spec.AdditionalMetadata.Labels),
			Annotations:     node.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, node),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(node.Spec.ServiceType),
			Ports: []corev1.ServicePort{
				{
					Port:       consts.DaemonPort,
					TargetPort: intstr.FromString("daemon"),
					Protocol:   "TCP",
					Name:       "daemon",
				},
				{
					Port:       r.getFullNodePort(ctx, node),
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
			Selector: r.getLabels(ctx, node, node.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleInternalService assembles the internal Service resource for a ChiaNode CR
func (r *ChiaNodeReconciler) assembleInternalService(ctx context.Context, node k8schianetv1.ChiaNode) corev1.Service {
	local := corev1.ServiceInternalTrafficPolicyLocal
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chianodeNamePattern, node.Name) + "-internal",
			Namespace:       node.Namespace,
			Labels:          r.getLabels(ctx, node, node.Spec.AdditionalMetadata.Labels),
			Annotations:     node.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, node),
		},
		Spec: corev1.ServiceSpec{
			Type:                  corev1.ServiceType("ClusterIP"),
			InternalTrafficPolicy: &local,
			Ports: []corev1.ServicePort{
				{
					Port:       consts.DaemonPort,
					TargetPort: intstr.FromString("daemon"),
					Protocol:   "TCP",
					Name:       "daemon",
				},
				{
					Port:       r.getFullNodePort(ctx, node),
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
			Selector: r.getLabels(ctx, node, node.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleHeadlessService assembles the headless Service resource for a ChiaNode CR
func (r *ChiaNodeReconciler) assembleHeadlessService(ctx context.Context, node k8schianetv1.ChiaNode) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chianodeNamePattern, node.Name) + "-headless",
			Namespace:       node.Namespace,
			Labels:          r.getLabels(ctx, node, node.Spec.AdditionalMetadata.Labels),
			Annotations:     node.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, node),
		},
		Spec: corev1.ServiceSpec{
			Type:      corev1.ServiceType("ClusterIP"),
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				{
					Port:       consts.DaemonPort,
					TargetPort: intstr.FromString("daemon"),
					Protocol:   "TCP",
					Name:       "daemon",
				},
				{
					Port:       r.getFullNodePort(ctx, node),
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
			Selector: r.getLabels(ctx, node, node.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaNode CR
func (r *ChiaNodeReconciler) assembleChiaExporterService(ctx context.Context, node k8schianetv1.ChiaNode) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chianodeNamePattern, node.Name) + "-metrics",
			Namespace:       node.Namespace,
			Labels:          r.getLabels(ctx, node, node.Spec.AdditionalMetadata.Labels, node.Spec.ChiaExporterConfig.ServiceLabels),
			Annotations:     node.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, node),
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
			Selector: r.getLabels(ctx, node, node.Spec.AdditionalMetadata.Labels),
		},
	}
}

// assembleStatefulset assembles the node StatefulSet resource for a ChiaNode CR
func (r *ChiaNodeReconciler) assembleStatefulset(ctx context.Context, node k8schianetv1.ChiaNode) appsv1.StatefulSet {
	var chiaSecContext *corev1.SecurityContext
	if node.Spec.ChiaConfig.SecurityContext != nil {
		chiaSecContext = node.Spec.ChiaConfig.SecurityContext
	}

	var chiaLivenessProbe *corev1.Probe
	if node.Spec.ChiaConfig.LivenessProbe != nil {
		chiaLivenessProbe = node.Spec.ChiaConfig.LivenessProbe
	}

	var chiaReadinessProbe *corev1.Probe
	if node.Spec.ChiaConfig.ReadinessProbe != nil {
		chiaReadinessProbe = node.Spec.ChiaConfig.ReadinessProbe
	}

	var chiaStartupProbe *corev1.Probe
	if node.Spec.ChiaConfig.StartupProbe != nil {
		chiaStartupProbe = node.Spec.ChiaConfig.StartupProbe
	}

	var chiaResources corev1.ResourceRequirements
	if node.Spec.ChiaConfig.Resources != nil {
		chiaResources = *node.Spec.ChiaConfig.Resources
	}

	var imagePullPolicy corev1.PullPolicy
	if node.Spec.ImagePullPolicy != nil {
		imagePullPolicy = *node.Spec.ImagePullPolicy
	}

	var chiaExporterImage = node.Spec.ChiaExporterConfig.Image
	if chiaExporterImage == "" {
		chiaExporterImage = consts.DefaultChiaExporterImage
	}

	vols, volClaimTemplates := r.getChiaVolumesAndTemplates(ctx, node)

	var stateful appsv1.StatefulSet = appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chianodeNamePattern, node.Name),
			Namespace:       node.Namespace,
			Labels:          r.getLabels(ctx, node, node.Spec.AdditionalMetadata.Labels),
			Annotations:     node.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, node),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: node.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: r.getLabels(ctx, node),
			},
			ServiceName: fmt.Sprintf(chianodeNamePattern, node.Name) + "-headless",
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      r.getLabels(ctx, node, node.Spec.AdditionalMetadata.Labels),
					Annotations: node.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					// TODO add: imagePullSecret, serviceAccountName config
					Containers: []corev1.Container{
						{
							Name:            "chia",
							SecurityContext: chiaSecContext,
							Image:           node.Spec.ChiaConfig.Image,
							ImagePullPolicy: imagePullPolicy,
							Env:             r.getChiaNodeEnv(ctx, node),
							Ports: []corev1.ContainerPort{
								{
									Name:          "daemon",
									ContainerPort: consts.DaemonPort,
									Protocol:      "TCP",
								},
								{
									Name:          "peers",
									ContainerPort: r.getFullNodePort(ctx, node),
									Protocol:      "TCP",
								},
								{
									Name:          "rpc",
									ContainerPort: consts.NodeRPCPort,
									Protocol:      "TCP",
								},
							},
							LivenessProbe:  chiaLivenessProbe,
							ReadinessProbe: chiaReadinessProbe,
							StartupProbe:   chiaStartupProbe,
							Resources:      chiaResources,
							VolumeMounts:   r.getChiaVolumeMounts(ctx, node),
						},
					},
					NodeSelector: node.Spec.NodeSelector,
					Volumes:      vols,
				},
			},
			VolumeClaimTemplates: volClaimTemplates,
		},
	}

	exporterContainer := kube.GetChiaExporterContainer(ctx, chiaExporterImage, chiaSecContext, imagePullPolicy, chiaResources)
	stateful.Spec.Template.Spec.Containers = append(stateful.Spec.Template.Spec.Containers, exporterContainer)

	if node.Spec.PodSecurityContext != nil {
		stateful.Spec.Template.Spec.SecurityContext = node.Spec.PodSecurityContext
	}

	// TODO add pod affinity, tolerations

	return stateful
}
