/*
Copyright 2023 Chia Network Inc.
*/

package v1

import corev1 "k8s.io/api/core/v1"

// CommonSpec represents the common configuration options for controller APIs at the top-spec level
type CommonSpec struct {
	AdditionalMetadata `json:",inline"`

	// ChiaExporterConfig defines the configuration options available to Chia component containers
	// +optional
	ChiaExporterConfig SpecChiaExporter `json:"chiaExporter,omitempty"`

	// InitContainers allows defining a list of containers that will run as init containers in the kubernetes Pods this resource creates
	// +optional
	InitContainers []ExtraContainer `json:"initContainers,omitempty"`

	// Sidecars allows defining a list of containers and volumes that will share the kubernetes Pod alongside a Chia container
	// +optional
	Sidecars []ExtraContainer `json:"sidecars,omitempty"`

	//StorageConfig defines the Chia container's CHIA_ROOT storage config
	// +optional
	Storage *StorageConfig `json:"storage,omitempty"`

	// ImagePullPolicy is the pull policy for containers in the pod
	// +optional
	// +kubebuilder:default="Always"
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// ImagePullSecrets is a local object reference list to some image pull secrets for pod templates
	// +optional
	ImagePullSecrets *[]corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// ServiceAccountName is an optional name of a Service Account in the target namespace to use for this Chia deployment
	// +optional
	ServiceAccountName *string `json:"serviceAccountName,omitempty"`

	// NodeSelector selects a node by key value pairs
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// PodSecurityContext defines the security context for the pod
	// +optional
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`

	// Affinity defines a group of affinity or anti-affinity rules
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// TopologySpreadConstraints describes how a group of pods ought to spread across topology domains.
	// +optional
	TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
}

// ExtraContainer allows defining a container spec that will share the kubernetes Pod alongside a Chia container, or run as an init container, along with some additional Pod spec configuration
type ExtraContainer struct {
	// Container allows defining a container spec that will share the kubernetes Pod alongside a Chia container
	// +optional
	Container corev1.Container `json:"container,omitempty"`

	// Volumes allows defining a list of volumes that can be mounted by this container
	// +optional
	Volumes []corev1.Volume `json:"volumes,omitempty"`

	// ShareVolumeMounts if set to true, shares any volume mounts from the main chia container to this container
	// +optional
	ShareVolumeMounts bool `json:"shareVolumeMounts,omitempty"`

	// ShareEnv if set to true, shares the environment variables from the main chia container to this container
	// +optional
	ShareEnv bool `json:"shareEnv,omitempty"`
}

// CommonSpecChia represents the common configuration options for a chia spec
type CommonSpecChia struct {
	// Image defines the image to use for the chia component containers
	// +optional
	Image *string `json:"image,omitempty"`

	// Testnet is set to true if the Chia container should switch to the latest default testnet's settings
	// +optional
	Testnet *bool `json:"testnet,omitempty"`

	// ChiaNetwork is the name of a ChiaNetwork resource in the same namespace as this resource
	// +optional
	ChiaNetwork *string `json:"chiaNetwork,omitempty"`

	// Network can be set to a network name in the chia configuration file to switch to
	// +optional
	Network *string `json:"network,omitempty"`

	// NetworkPort can be set to the port that full_nodes will use in the selected network.
	// This implies specification of the Network setting.
	// +optional
	NetworkPort *uint16 `json:"networkPort,omitempty"`

	// IntroducerAddress can be set to the hostname or IP address of an introducer to set in the chia config.
	// No port should be specified, it's taken from the value of the NetworkPort setting.
	// +optional
	IntroducerAddress *string `json:"introducerAddress,omitempty"`

	// DNSIntroducerAddress can be set to a hostname to a DNS Introducer server.
	// +optional
	DNSIntroducerAddress *string `json:"dnsIntroducerAddress,omitempty"`

	// Timezone can be set to your local timezone for accurate timestamps. Defaults to UTC
	// +optional
	Timezone *string `json:"timezone,omitempty"`

	// SourceRef is set to the desired ref of the chia-blockchain repository to install from. Defaults to unset (uses the installation already in the chia image.)
	// +optional
	SourceRef *string `json:"sourceRef,omitempty"`

	// LogLevel is set to the desired chia config log_level
	// +optional
	LogLevel *string `json:"logLevel,omitempty"`

	// SelfHostname defines the bind address of chia services in the container
	// Setting to `0.0.0.0` binds chia services to all interfaces
	// +optional
	SelfHostname *string `json:"selfHostname,omitempty"`

	// PeerService defines settings for the default Service installed with any Chia component resource.
	// This Service usually contains ports for peer connections, or in the case of seeders port 53.
	// This Service will default to being enabled with a ClusterIP Service type.
	// +optional
	PeerService Service `json:"peerService,omitempty"`

	// DaemonService defines settings for the daemon Service installed with any Chia component resource.
	// This Service usually contains the port for the Chia daemon that runs alongside any Chia instance.
	// This Service will default to being enabled with a ClusterIP Service type.
	// +optional
	DaemonService Service `json:"daemonService,omitempty"`

	// RPCService defines settings for the RPC Service installed with any Chia component resource.
	// This Service contains the port for the Chia RPC API.
	// This Service will default to being enabled with a ClusterIP Service type.
	// +optional
	RPCService Service `json:"rpcService,omitempty"`

	// AllService defines settings for a Service that contains all the ports from the peer, daemon, and RPC Services installed with any Chia component resource.
	// This Service will default to being enabled with a ClusterIP Service type.
	// +optional
	AllService Service `json:"allService,omitempty"`

	// AdditionalEnv contain a list of additional environment variables to be supplied to the chia container.
	// These variables will be placed at the end of the environment variable list in the resulting container, this means they overwrite variables of the same name created by the operator in the container env.
	// +optional
	AdditionalEnv *[]corev1.EnvVar `json:"additionalEnv,omitempty"`

	// LivenessProbe used to determine if a container is running properly and will restart the container if the probe fails
	// +optional
	LivenessProbe *corev1.Probe `json:"livenessProbe,omitempty"`

	// ReadinessProbe used to indicate when a container is ready to accept traffic and prevent traffic from being sent to pods that aren't ready.
	// +optional
	ReadinessProbe *corev1.Probe `json:"readinessProbe,omitempty"`

	// StartupProbe used to give applications time to initialize fully before liveness and readiness probes begin checking, preventing premature restarts of slow-starting containers.
	// +optional
	StartupProbe *corev1.Probe `json:"startupProbe,omitempty"`

	// Resources defines the compute resources (limits/requests) for the chia container.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// SecurityContext defines the security context for the chia container
	// +optional
	SecurityContext *corev1.SecurityContext `json:"securityContext,omitempty"`
}

// SpecChiaExporter defines the desired state of Chia exporter configuration
type SpecChiaExporter struct {
	// Enabled defines whether a chia-exporter sidecar container should run with the chia container
	// Defaults to enabled
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Image defines the image to use for the chia exporter containers
	// +optional
	Image *string `json:"image,omitempty"`

	// Service defines settings for the Service installed with any chia-exporter resource.
	// This Service contains the port for chia-exporter's web exporter.
	// This Service will default to being enabled with a ClusterIP Service type if chia-exporter is enabled.
	// +optional
	Service Service `json:"service,omitempty"`

	// ConfigSecretName is the name of an optional Secret that contains the environment variables that will be mounted in the chia-exporter container.
	// +optional
	ConfigSecretName *string `json:"configSecretName,omitempty"`

	// Resources defines the compute resources (limits/requests) for the chia container.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// SecurityContext defines the security context for the chia container
	// +optional
	SecurityContext *corev1.SecurityContext `json:"securityContext,omitempty"`
}

// SpecChiaHealthcheck defines the desired state of Chia healthcheck configuration
type SpecChiaHealthcheck struct {
	// Enabled defines whether a chia-exporter sidecar container should run with the chia container
	// Defaults to enabled on services that support it
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Image defines the image to use for the chia exporter containers
	// +optional
	Image *string `json:"image,omitempty"`

	// DNSHostname is the hostname to check for DNS responses. Disabled if not provided.
	// +optional
	DNSHostname *string `json:"dnsHostname,omitempty"`

	// Service defines settings for the Service installed with any chia-healthcheck resource.
	// This Service contains the port for chia-healthcheck's web server.
	// This Service will default to being disabled.
	// +optional
	Service Service `json:"service,omitempty"`

	// Resources defines the compute resources (limits/requests) for the chia container.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// SecurityContext defines the security context for the chia container
	// +optional
	SecurityContext *corev1.SecurityContext `json:"securityContext,omitempty"`
}

// ChiaSecretKey defines the name of a kubernetes secret and key in that namespace that contains the Chia mnemonic
type ChiaSecretKey struct {
	// SecretName is the name of the kubernetes secret containing a mnemonic key
	Name string `json:"name"`

	// Key is the key of the data item in the Secret
	Key string `json:"key"`
}

// AdditionalMetadata contains labels and annotations to attach to created objects
type AdditionalMetadata struct {
	// Labels is a map of string keys and values to attach to created objects
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations is a map of string keys and values to attach to created objects
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// Service contains kubernetes Service related configuration options
type Service struct {
	AdditionalMetadata `json:",inline"`

	// Enabled is a boolean selector for a Service if it should be generated.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// ServiceType is the Type of the Service. Defaults to ClusterIP
	// +optional
	ServiceType *corev1.ServiceType `json:"type,omitempty"`

	// IPFamilyPolicy represents the dual-stack-ness requested or required by a Service
	// +optional
	IPFamilyPolicy *corev1.IPFamilyPolicy `json:"ipFamilyPolicy,omitempty"`

	// IPFamilies represents a list of IP families (IPv4 and/or IPv6) required by a Service
	// +optional
	IPFamilies *[]corev1.IPFamily `json:"ipFamilies,omitempty"`

	// ExternalTrafficPolicy sets the external traffic policy for the service
	// +optional
	ExternalTrafficPolicy *corev1.ServiceExternalTrafficPolicy `json:"externalTrafficPolicy,omitempty"`

	// SessionAffinity can be set to "ClientIP" to enable session affinity based on client IP
	// +optional
	SessionAffinity *corev1.ServiceAffinity `json:"sessionAffinity,omitempty"`

	// SessionAffinityConfig allows configuring the settings for sessionAffinity
	// +optional
	SessionAffinityConfig *corev1.SessionAffinityConfig `json:"sessionAffinityConfig,omitempty"`

	// RollIntoPeerService tells the controller to not actually generate this Service, but instead roll the Service ports of this Service into the peer Service.
	// The peer Service is often considered the primary Service generated for a chia resource, as it is the most likely Service to expose publicly.
	// This option is default, and only provides its functionality on chia-healthcheck Services. It may be included to other Services someday if a use case arises.
	// +optional
	RollIntoPeerService *bool `json:"rollIntoPeerService,omitempty"`
}

// StorageConfig contains storage configuration settings
type StorageConfig struct {
	// Storage configuration for CHIA_ROOT
	// +optional
	ChiaRoot *ChiaRootConfig `json:"chiaRoot,omitempty"`

	// Storage configuration for harvester plots
	// +optional
	Plots *PlotsConfig `json:"plots,omitempty"`

	// Storage configuration for data_layer server files
	// +optional
	DataLayerServerFiles *DataLayerServerFilesConfig `json:"dataLayerServerFiles,omitempty"`
}

// ChiaRootConfig optional config for CHIA_ROOT persistent storage, likely only needed for Chia full_nodes, but may help in startup time for other components.
// Both options may be specified but only one can be used, therefore PersistentVolumeClaims will be respected over HostPath volumes if both are specified.
type ChiaRootConfig struct {
	// PersistentVolumeClaim use an existing persistent volume claim to store CHIA_ROOT data
	// +optional
	PersistentVolumeClaim *PersistentVolumeClaimConfig `json:"persistentVolumeClaim,omitempty"`

	// HostPathVolume use an existing directory on the host to store CHIA_ROOT data
	// +optional
	HostPathVolume *HostPathVolumeConfig `json:"hostPathVolume,omitempty"`
}

// PlotsConfig optional config for harvester plots persistent storage, only needed for Chia harvesters.
// Supports adding both PVCs and hostPath volumes.
type PlotsConfig struct {
	// PersistentVolumeClaim use an existing persistent volume claim to mount plot directories
	// +optional
	PersistentVolumeClaim []*PersistentVolumeClaimConfig `json:"persistentVolumeClaim,omitempty"`

	// HostPathVolume use an existing directory on the host to mount plot directories
	// +optional
	HostPathVolume []*HostPathVolumeConfig `json:"hostPathVolume,omitempty"`
}

// DataLayerServerFilesConfig optional config for data_layer server file persistent storage.
// Both options may be specified but only one can be used, therefore PersistentVolumeClaims will be respected over HostPath volumes if both are specified.
type DataLayerServerFilesConfig struct {
	// PersistentVolumeClaim use an existing persistent volume claim to store server files
	// +optional
	PersistentVolumeClaim *PersistentVolumeClaimConfig `json:"persistentVolumeClaim,omitempty"`

	// HostPathVolume use an existing directory on the host to store server files
	// +optional
	HostPathVolume *HostPathVolumeConfig `json:"hostPathVolume,omitempty"`
}

// PersistentVolumeClaimConfig config for PVC volumes in kubernetes
type PersistentVolumeClaimConfig struct {
	// ClaimName is the name of an existing PersistentVolumeClaim in the target namespace
	// This field does nothing on ChiaNode resources.
	// This field does nothing when GenerateVolumeClaims is set to true.
	// +optional
	ClaimName string `json:"claimName,omitempty"`

	// GenerateVolumeClaims is mutually exclusive with the ClaimName field, and overrides that field if set.
	// Instead, an operator generated PVC name will be made, and the operator will provision a volume claim for you.
	// This field does nothing on ChiaNode resources.
	// +optional
	GenerateVolumeClaims bool `json:"generateVolumeClaims,omitempty"`

	// StorageClass is the name of a storage class for the PVC. Only relevant for ChiaNodes and use with the GenerateVolumeClaims option.
	// +optional
	StorageClass string `json:"storageClass,omitempty"`

	// AccessModes are the volume access modes. Only relevant for ChiaNodes and use with the GenerateVolumeClaims option.
	// Defaults to RWO if unspecified.
	// +optional
	AccessModes []corev1.PersistentVolumeAccessMode `json:"accessModes"`

	// ResourceRequest is the amount of storage requested. Only relevant for ChiaNodes and use with the GenerateVolumeClaims option.
	// +optional
	ResourceRequest string `json:"resourceRequest,omitempty"`
}

// HostPathVolumeConfig config for hostPath volumes in kubernetes
type HostPathVolumeConfig struct {
	// Path use an existing directory on your Pod's host to mount in the Pod's containers.
	// If a HostPath is used, it is highly recommended that a NodeSelector is used to keep the Pod on the host that has the directory to mount.
	// +optional
	Path string `json:"path,omitempty"`
}

// Peer config for a peer - host and port
type Peer struct {
	// Host is the IP address or hostname to a full_node peer.
	Host string `json:"host"`

	// Port is the port number the full_node's peer port is listening on.
	Port uint16 `json:"port"`
}
