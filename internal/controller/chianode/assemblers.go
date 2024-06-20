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

// assemblePeerService assembles the peer Service for a Chianode CR
func (r *ChiaNodeReconciler) assemblePeerService(ctx context.Context, node k8schianetv1.ChiaNode) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chianodeNamePattern, node.Name)
	inputs.Namespace = node.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, node)

	// Service Type
	if node.Spec.ChiaConfig.PeerService != nil && node.Spec.ChiaConfig.PeerService.ServiceType != nil {
		inputs.ServiceType = *node.Spec.ChiaConfig.PeerService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if node.Spec.ChiaConfig.PeerService != nil && node.Spec.ChiaConfig.PeerService.Labels != nil {
		additionalServiceLabels = node.Spec.ChiaConfig.PeerService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if node.Spec.ChiaConfig.PeerService != nil && node.Spec.ChiaConfig.PeerService.Annotations != nil {
		additionalServiceAnnotations = node.Spec.ChiaConfig.PeerService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(node.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       r.getFullNodePort(ctx, node),
			TargetPort: intstr.FromString("peers"),
			Protocol:   "TCP",
			Name:       "peers",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleHeadlessPeerService assembles the headless peer Service for a Chianode CR
func (r *ChiaNodeReconciler) assembleHeadlessPeerService(ctx context.Context, node k8schianetv1.ChiaNode) corev1.Service {
	srv := r.assemblePeerService(ctx, node)

	srv.Name = srv.Name + "-headless"
	srv.Spec.Type = "ClusterIP"
	srv.Spec.ClusterIP = "None"

	return srv
}

// assembleHeadlessPeerService assembles the headless peer Service for a Chianode CR
func (r *ChiaNodeReconciler) assembleLocalPeerService(ctx context.Context, node k8schianetv1.ChiaNode) corev1.Service {
	srv := r.assemblePeerService(ctx, node)

	srv.Name = srv.Name + "-internal"
	local := corev1.ServiceInternalTrafficPolicyLocal
	srv.Spec.InternalTrafficPolicy = &local

	return srv
}

// assembleDaemonService assembles the daemon Service resource for a Chianode CR
func (r *ChiaNodeReconciler) assembleDaemonService(ctx context.Context, node k8schianetv1.ChiaNode) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chianodeNamePattern, node.Name) + "-daemon"
	inputs.Namespace = node.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, node)

	// Service Type
	if node.Spec.ChiaConfig.DaemonService != nil && node.Spec.ChiaConfig.DaemonService.ServiceType != nil {
		inputs.ServiceType = *node.Spec.ChiaConfig.DaemonService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if node.Spec.ChiaConfig.DaemonService != nil && node.Spec.ChiaConfig.DaemonService.Labels != nil {
		additionalServiceLabels = node.Spec.ChiaConfig.DaemonService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if node.Spec.ChiaConfig.DaemonService != nil && node.Spec.ChiaConfig.DaemonService.Annotations != nil {
		additionalServiceAnnotations = node.Spec.ChiaConfig.DaemonService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(node.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

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

// assembleRPCService assembles the RPC Service resource for a Chianode CR
func (r *ChiaNodeReconciler) assembleRPCService(ctx context.Context, node k8schianetv1.ChiaNode) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chianodeNamePattern, node.Name) + "-rpc"
	inputs.Namespace = node.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, node)

	// Service Type
	if node.Spec.ChiaConfig.RPCService != nil && node.Spec.ChiaConfig.RPCService.ServiceType != nil {
		inputs.ServiceType = *node.Spec.ChiaConfig.RPCService.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if node.Spec.ChiaConfig.RPCService != nil && node.Spec.ChiaConfig.RPCService.Labels != nil {
		additionalServiceLabels = node.Spec.ChiaConfig.RPCService.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if node.Spec.ChiaConfig.RPCService != nil && node.Spec.ChiaConfig.RPCService.Annotations != nil {
		additionalServiceAnnotations = node.Spec.ChiaConfig.RPCService.Annotations
	}
	inputs.Annotations = kube.CombineMaps(node.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

	// Ports
	inputs.Ports = []corev1.ServicePort{
		{
			Port:       consts.NodeRPCPort,
			TargetPort: intstr.FromString("rpc"),
			Protocol:   "TCP",
			Name:       "rpc",
		},
	}

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a Chianode CR
func (r *ChiaNodeReconciler) assembleChiaExporterService(ctx context.Context, node k8schianetv1.ChiaNode) corev1.Service {
	var inputs kube.AssembleCommonServiceInputs

	// Service Metadata
	inputs.Name = fmt.Sprintf(chianodeNamePattern, node.Name) + "-metrics"
	inputs.Namespace = node.Namespace
	inputs.OwnerReference = r.getOwnerReference(ctx, node)

	// Service Type
	if node.Spec.ChiaExporterConfig.Service != nil && node.Spec.ChiaExporterConfig.Service.ServiceType != nil {
		inputs.ServiceType = *node.Spec.ChiaExporterConfig.Service.ServiceType
	} else {
		inputs.ServiceType = corev1.ServiceTypeClusterIP
	}

	// Labels
	var additionalServiceLabels = make(map[string]string)
	if node.Spec.ChiaExporterConfig.Service != nil && node.Spec.ChiaExporterConfig.Service.Labels != nil {
		additionalServiceLabels = node.Spec.ChiaExporterConfig.Service.Labels
	}
	inputs.Labels = kube.GetCommonLabels(ctx, node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
	inputs.SelectorLabels = kube.GetCommonLabels(ctx, node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels)

	// Annotations
	var additionalServiceAnnotations = make(map[string]string)
	if node.Spec.ChiaExporterConfig.Service != nil && node.Spec.ChiaExporterConfig.Service.Annotations != nil {
		additionalServiceAnnotations = node.Spec.ChiaExporterConfig.Service.Annotations
	}
	inputs.Annotations = kube.CombineMaps(node.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)

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

// assembleStatefulset assembles the node StatefulSet resource for a ChiaNode CR
func (r *ChiaNodeReconciler) assembleStatefulset(ctx context.Context, node k8schianetv1.ChiaNode) appsv1.StatefulSet {
	vols, volClaimTemplates := r.getChiaVolumesAndTemplates(ctx, node)

	var stateful appsv1.StatefulSet = appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:            fmt.Sprintf(chianodeNamePattern, node.Name),
			Namespace:       node.Namespace,
			Labels:          kube.GetCommonLabels(ctx, node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels),
			Annotations:     node.Spec.AdditionalMetadata.Annotations,
			OwnerReferences: r.getOwnerReference(ctx, node),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &node.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(ctx, node.Kind, node.ObjectMeta),
			},
			ServiceName: fmt.Sprintf(chianodeNamePattern, node.Name) + "-headless",
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(ctx, node.Kind, node.ObjectMeta, node.Spec.AdditionalMetadata.Labels),
					Annotations: node.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					// TODO add: imagePullSecret, serviceAccountName config
					Containers: []corev1.Container{
						{
							Name:            "chia",
							Image:           node.Spec.ChiaConfig.Image,
							ImagePullPolicy: node.Spec.ImagePullPolicy,
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
							VolumeMounts: r.getChiaVolumeMounts(ctx, node),
						},
					},
					NodeSelector: node.Spec.NodeSelector,
					Volumes:      vols,
				},
			},
			VolumeClaimTemplates: volClaimTemplates,
		},
	}

	var containerSecurityContext *corev1.SecurityContext
	if node.Spec.ChiaConfig.SecurityContext != nil {
		containerSecurityContext = node.Spec.ChiaConfig.SecurityContext
		stateful.Spec.Template.Spec.Containers[0].SecurityContext = node.Spec.ChiaConfig.SecurityContext
	}

	if node.Spec.ChiaConfig.LivenessProbe != nil {
		stateful.Spec.Template.Spec.Containers[0].LivenessProbe = node.Spec.ChiaConfig.LivenessProbe
	}

	if node.Spec.ChiaConfig.ReadinessProbe != nil {
		stateful.Spec.Template.Spec.Containers[0].ReadinessProbe = node.Spec.ChiaConfig.ReadinessProbe
	}

	if node.Spec.ChiaConfig.StartupProbe != nil {
		stateful.Spec.Template.Spec.Containers[0].StartupProbe = node.Spec.ChiaConfig.StartupProbe
	}

	var containerResorces corev1.ResourceRequirements
	if node.Spec.ChiaConfig.Resources != nil {
		containerResorces = *node.Spec.ChiaConfig.Resources
		stateful.Spec.Template.Spec.Containers[0].Resources = *node.Spec.ChiaConfig.Resources
	}

	if node.Spec.ChiaExporterConfig.Enabled {
		exporterContainer := kube.GetChiaExporterContainer(ctx, node.Spec.ChiaExporterConfig.Image, containerSecurityContext, node.Spec.ImagePullPolicy, containerResorces)
		stateful.Spec.Template.Spec.Containers = append(stateful.Spec.Template.Spec.Containers, exporterContainer)
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
