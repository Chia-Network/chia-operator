/*
Copyright 2023 Chia Network Inc.
*/

package chiaseeder

import (
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

const chiaseederNamePattern = "%s-seeder"

// assemblePeerService assembles the peer Service resource for a ChiaSeeder CR
func assemblePeerService(seeder k8schianetv1.ChiaSeeder) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:           fmt.Sprintf(chiaseederNamePattern, seeder.Name),
		Namespace:      seeder.Namespace,
		OwnerReference: getOwnerReference(seeder),
		Ports: []corev1.ServicePort{
			{
				Port:       53,
				TargetPort: intstr.FromString("dns"),
				Protocol:   "UDP",
				Name:       "dns",
			},
			{
				Port:       53,
				TargetPort: intstr.FromString("dns-tcp"),
				Protocol:   "TCP",
				Name:       "dns-tcp",
			},
			{
				Port:       getFullNodePort(seeder),
				TargetPort: intstr.FromString("peers"),
				Protocol:   "TCP",
				Name:       "peers",
			},
		},
	}

	if seeder.Spec.ChiaConfig.PeerService != nil {
		inputs.ServiceType = seeder.Spec.ChiaConfig.PeerService.ServiceType
		inputs.IPFamilyPolicy = seeder.Spec.ChiaConfig.PeerService.IPFamilyPolicy
		inputs.IPFamilies = seeder.Spec.ChiaConfig.PeerService.IPFamilies

		// Labels
		var additionalServiceLabels = make(map[string]string)
		if seeder.Spec.ChiaConfig.PeerService.Labels != nil {
			additionalServiceLabels = seeder.Spec.ChiaConfig.PeerService.Labels
		}
		inputs.Labels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
		inputs.SelectorLabels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels)

		// Annotations
		var additionalServiceAnnotations = make(map[string]string)
		if seeder.Spec.ChiaConfig.PeerService.Annotations != nil {
			additionalServiceAnnotations = seeder.Spec.ChiaConfig.PeerService.Annotations
		}
		inputs.Annotations = kube.CombineMaps(seeder.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)
	}

	return kube.AssembleCommonService(inputs)
}

// assembleDaemonService assembles the daemon Service resource for a ChiaSeeder CR
func assembleDaemonService(seeder k8schianetv1.ChiaSeeder) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:           fmt.Sprintf(chiaseederNamePattern, seeder.Name) + "-daemon",
		Namespace:      seeder.Namespace,
		OwnerReference: getOwnerReference(seeder),
		Ports:          kube.GetChiaDaemonServicePorts(),
	}

	if seeder.Spec.ChiaConfig.DaemonService != nil {
		inputs.ServiceType = seeder.Spec.ChiaConfig.DaemonService.ServiceType
		inputs.IPFamilyPolicy = seeder.Spec.ChiaConfig.DaemonService.IPFamilyPolicy
		inputs.IPFamilies = seeder.Spec.ChiaConfig.DaemonService.IPFamilies

		// Labels
		var additionalServiceLabels = make(map[string]string)
		if seeder.Spec.ChiaConfig.DaemonService.Labels != nil {
			additionalServiceLabels = seeder.Spec.ChiaConfig.DaemonService.Labels
		}
		inputs.Labels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
		inputs.SelectorLabels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels)

		// Annotations
		var additionalServiceAnnotations = make(map[string]string)
		if seeder.Spec.ChiaConfig.DaemonService.Annotations != nil {
			additionalServiceAnnotations = seeder.Spec.ChiaConfig.DaemonService.Annotations
		}
		inputs.Annotations = kube.CombineMaps(seeder.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)
	}

	return kube.AssembleCommonService(inputs)
}

// assembleRPCService assembles the RPC Service resource for a ChiaSeeder CR
func assembleRPCService(seeder k8schianetv1.ChiaSeeder) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:           fmt.Sprintf(chiaseederNamePattern, seeder.Name) + "-rpc",
		Namespace:      seeder.Namespace,
		OwnerReference: getOwnerReference(seeder),
		Ports: []corev1.ServicePort{
			{
				Port:       consts.CrawlerRPCPort,
				TargetPort: intstr.FromString("rpc"),
				Protocol:   "TCP",
				Name:       "rpc",
			},
		},
	}

	if seeder.Spec.ChiaConfig.RPCService != nil {
		inputs.ServiceType = seeder.Spec.ChiaConfig.RPCService.ServiceType
		inputs.IPFamilyPolicy = seeder.Spec.ChiaConfig.RPCService.IPFamilyPolicy
		inputs.IPFamilies = seeder.Spec.ChiaConfig.RPCService.IPFamilies

		// Labels
		var additionalServiceLabels = make(map[string]string)
		if seeder.Spec.ChiaConfig.RPCService.Labels != nil {
			additionalServiceLabels = seeder.Spec.ChiaConfig.RPCService.Labels
		}
		inputs.Labels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
		inputs.SelectorLabels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels)

		// Annotations
		var additionalServiceAnnotations = make(map[string]string)
		if seeder.Spec.ChiaConfig.RPCService.Annotations != nil {
			additionalServiceAnnotations = seeder.Spec.ChiaConfig.RPCService.Annotations
		}
		inputs.Annotations = kube.CombineMaps(seeder.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)
	}

	return kube.AssembleCommonService(inputs)
}

// assembleChiaExporterService assembles the chia-exporter Service resource for a ChiaSeeder CR
func assembleChiaExporterService(seeder k8schianetv1.ChiaSeeder) corev1.Service {
	inputs := kube.AssembleCommonServiceInputs{
		Name:           fmt.Sprintf(chiaseederNamePattern, seeder.Name) + "-metrics",
		Namespace:      seeder.Namespace,
		OwnerReference: getOwnerReference(seeder),
		Ports:          kube.GetChiaExporterServicePorts(),
	}

	if seeder.Spec.ChiaExporterConfig.Service != nil {
		inputs.ServiceType = seeder.Spec.ChiaExporterConfig.Service.ServiceType
		inputs.IPFamilyPolicy = seeder.Spec.ChiaExporterConfig.Service.IPFamilyPolicy
		inputs.IPFamilies = seeder.Spec.ChiaExporterConfig.Service.IPFamilies

		// Labels
		var additionalServiceLabels = make(map[string]string)
		if seeder.Spec.ChiaExporterConfig.Service.Labels != nil {
			additionalServiceLabels = seeder.Spec.ChiaExporterConfig.Service.Labels
		}
		inputs.Labels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels, additionalServiceLabels)
		inputs.SelectorLabels = kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels)

		// Annotations
		var additionalServiceAnnotations = make(map[string]string)
		if seeder.Spec.ChiaExporterConfig.Service.Annotations != nil {
			additionalServiceAnnotations = seeder.Spec.ChiaExporterConfig.Service.Annotations
		}
		inputs.Annotations = kube.CombineMaps(seeder.Spec.AdditionalMetadata.Annotations, additionalServiceAnnotations)
	}

	return kube.AssembleCommonService(inputs)
}

// assembleVolumeClaim assembles the PVC resource for a ChiaSeeder CR
func assembleVolumeClaim(seeder k8schianetv1.ChiaSeeder) (corev1.PersistentVolumeClaim, error) {
	resourceReq, err := resource.ParseQuantity(seeder.Spec.Storage.ChiaRoot.PersistentVolumeClaim.ResourceRequest)
	if err != nil {
		return corev1.PersistentVolumeClaim{}, err
	}

	accessModes := []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	if len(seeder.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes) != 0 {
		accessModes = seeder.Spec.Storage.ChiaRoot.PersistentVolumeClaim.AccessModes
	}

	return corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(chiaseederNamePattern, seeder.Name),
			Namespace: seeder.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      accessModes,
			StorageClassName: &seeder.Spec.Storage.ChiaRoot.PersistentVolumeClaim.StorageClass,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resourceReq,
				},
			},
		},
	}, nil
}

// assembleDeployment assembles the seeder Deployment resource for a ChiaSeeder CR
func assembleDeployment(seeder k8schianetv1.ChiaSeeder) appsv1.Deployment {
	var deploy = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf(chiaseederNamePattern, seeder.Name),
			Namespace:   seeder.Namespace,
			Labels:      kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels),
			Annotations: seeder.Spec.AdditionalMetadata.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      kube.GetCommonLabels(seeder.Kind, seeder.ObjectMeta, seeder.Spec.AdditionalMetadata.Labels),
					Annotations: seeder.Spec.AdditionalMetadata.Annotations,
				},
				Spec: corev1.PodSpec{
					Containers:   []corev1.Container{assembleChiaContainer(seeder)},
					NodeSelector: seeder.Spec.NodeSelector,
					Volumes:      getChiaVolumes(seeder),
				},
			},
		},
	}

	if len(seeder.Spec.InitContainers) != 0 {
		// Overwrite any volumeMounts specified in init containers. Not currently supported.
		for _, cont := range seeder.Spec.InitContainers {
			cont.Container.VolumeMounts = []corev1.VolumeMount{}

			// Share chia volume mounts if enabled
			if cont.ShareVolumeMounts {
				cont.Container.VolumeMounts = getChiaVolumeMounts(seeder)
			}

			// Share chia env if enabled
			if cont.ShareEnv {
				cont.Container.Env = append(cont.Container.Env, getChiaEnv(seeder)...)
			}

			deploy.Spec.Template.Spec.InitContainers = append(deploy.Spec.Template.Spec.InitContainers, cont.Container)
		}
	}

	if seeder.Spec.ChiaExporterConfig.Enabled {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, assembleChiaExporterContainer(seeder))
	}

	if seeder.Spec.ChiaHealthcheckConfig.Enabled && seeder.Spec.ChiaHealthcheckConfig.DNSHostname != nil {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, assembleChiaHealthcheckContainer(seeder))
	}

	if seeder.Spec.Strategy != nil {
		deploy.Spec.Strategy = *seeder.Spec.Strategy
	}

	if seeder.Spec.PodSecurityContext != nil {
		deploy.Spec.Template.Spec.SecurityContext = seeder.Spec.PodSecurityContext
	}

	if len(seeder.Spec.Sidecars.Containers) > 0 {
		deploy.Spec.Template.Spec.Containers = append(deploy.Spec.Template.Spec.Containers, seeder.Spec.Sidecars.Containers...)
	}

	// TODO add pod affinity, tolerations

	return deploy
}

func assembleChiaContainer(seeder k8schianetv1.ChiaSeeder) corev1.Container {
	input := kube.AssembleChiaContainerInputs{
		Image:           seeder.Spec.ChiaConfig.Image,
		ImagePullPolicy: seeder.Spec.ImagePullPolicy,
		Env:             getChiaEnv(seeder),
		Ports:           getChiaPorts(seeder),
		VolumeMounts:    getChiaVolumeMounts(seeder),
	}

	if seeder.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = seeder.Spec.ChiaConfig.SecurityContext
	} else {
		input.SecurityContext = &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{
					"NET_BIND_SERVICE",
				},
			},
		}
	}

	if seeder.Spec.ChiaConfig.LivenessProbe != nil {
		input.LivenessProbe = seeder.Spec.ChiaConfig.LivenessProbe
	} else if seeder.Spec.ChiaHealthcheckConfig.Enabled && seeder.Spec.ChiaHealthcheckConfig.DNSHostname != nil {
		input.ReadinessProbe = kube.AssembleChiaHealthcheckProbe(kube.AssembleChiaHealthcheckProbeInputs{
			Kind: consts.ChiaSeederKind,
		})
	}

	if seeder.Spec.ChiaConfig.ReadinessProbe != nil {
		input.ReadinessProbe = seeder.Spec.ChiaConfig.ReadinessProbe
	} else if seeder.Spec.ChiaHealthcheckConfig.Enabled && seeder.Spec.ChiaHealthcheckConfig.DNSHostname != nil {
		input.ReadinessProbe = kube.AssembleChiaHealthcheckProbe(kube.AssembleChiaHealthcheckProbeInputs{
			Kind: consts.ChiaSeederKind,
		})
	}

	if seeder.Spec.ChiaConfig.StartupProbe != nil {
		input.StartupProbe = seeder.Spec.ChiaConfig.StartupProbe
	} else if seeder.Spec.ChiaHealthcheckConfig.Enabled && seeder.Spec.ChiaHealthcheckConfig.DNSHostname != nil {
		failThresh := int32(30)
		periodSec := int32(10)
		input.StartupProbe = kube.AssembleChiaHealthcheckProbe(kube.AssembleChiaHealthcheckProbeInputs{
			Kind:             consts.ChiaSeederKind,
			FailureThreshold: &failThresh,
			PeriodSeconds:    &periodSec,
		})
	}

	if seeder.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = seeder.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaContainer(input)
}

func assembleChiaExporterContainer(seeder k8schianetv1.ChiaSeeder) corev1.Container {
	input := kube.AssembleChiaExporterContainerInputs{
		Image:            seeder.Spec.ChiaExporterConfig.Image,
		ConfigSecretName: seeder.Spec.ChiaExporterConfig.ConfigSecretName,
		PullPolicy:       seeder.Spec.ImagePullPolicy,
	}

	if seeder.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = seeder.Spec.ChiaConfig.SecurityContext
	}

	if seeder.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = *seeder.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaExporterContainer(input)
}

func assembleChiaHealthcheckContainer(seeder k8schianetv1.ChiaSeeder) corev1.Container {
	input := kube.AssembleChiaHealthcheckContainerInputs{
		Image:       seeder.Spec.ChiaHealthcheckConfig.Image,
		DNSHostname: seeder.Spec.ChiaHealthcheckConfig.DNSHostname,
		PullPolicy:  seeder.Spec.ImagePullPolicy,
	}

	if seeder.Spec.ChiaConfig.SecurityContext != nil {
		input.SecurityContext = seeder.Spec.ChiaConfig.SecurityContext
	}

	if seeder.Spec.ChiaConfig.Resources != nil {
		input.ResourceRequirements = *seeder.Spec.ChiaConfig.Resources
	}

	return kube.AssembleChiaHealthcheckContainer(input)
}
