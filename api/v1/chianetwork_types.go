/*
Copyright 2024 Chia Network Inc.
*/

package v1

import (
	"github.com/chia-network/go-chia-libs/pkg/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChiaNetworkSpec defines the desired state of ChiaNetwork
type ChiaNetworkSpec struct {
	// NetworkConstants specifies the network constants for this network in the config
	// +optional
	NetworkConstants *NetworkConstants `json:"constants"`

	// NetworkConfig is the config for the network (address prefix and default full_node port)
	// +optional
	NetworkConfig *config.NetworkConfig `json:"config"`

	// NetworkName is the name of the selected network in the config, and will also be used as the key for related network config and constants.
	// If specified on a ChiaNetwork, and passed to a chia-deploying resource, this will override any value specified for `.spec.chia.network` on that resource.
	// This field is optional, and network name will default to the ChiaNetwork name if unspecified.
	// +optional
	NetworkName *string `json:"networkName,omitempty"`

	// NetworkPort can be set to the port that full_nodes will use in the selected network.
	// If specified on a ChiaNetwork, and passed to a chia-deploying resource, this will override any value specified for `.spec.chia.networkPort` on that resource.
	// +optional
	NetworkPort *uint16 `json:"networkPort,omitempty"`

	// IntroducerAddress can be set to the hostname or IP address of an introducer to set in the chia config.
	// No port should be specified, it's taken from the value of the NetworkPort setting.
	// If specified on a ChiaNetwork, and passed to a chia-deploying resource, this will override any value specified for `.spec.chia.introducerAddress` on that resource.
	// +optional
	IntroducerAddress *string `json:"introducerAddress,omitempty"`

	// DNSIntroducerAddress can be set to a hostname to a DNS Introducer server.
	// If specified on a ChiaNetwork, and passed to a chia-deploying resource, this will override any value specified for `.spec.chia.dnsIntroducerAddress` on that resource.
	// +optional
	DNSIntroducerAddress *string `json:"dnsIntroducerAddress,omitempty"`
}

// NetworkConstants the constants for each network
type NetworkConstants struct {
	GenesisChallenge               string `json:"GENESIS_CHALLENGE"`
	GenesisPreFarmPoolPuzzleHash   string `json:"GENESIS_PRE_FARM_POOL_PUZZLE_HASH"`
	GenesisPreFarmFarmerPuzzleHash string `json:"GENESIS_PRE_FARM_FARMER_PUZZLE_HASH"`

	// +optional
	AggSigMeAdditionalData string `json:"AGG_SIG_ME_ADDITIONAL_DATA,omitempty"`

	// +optional
	DifficultyConstantFactor uint64 `json:"DIFFICULTY_CONSTANT_FACTOR,omitempty"`

	// +optional
	DifficultyStarting uint64 `json:"DIFFICULTY_STARTING,omitempty"`

	// +optional
	EpochBlocks uint32 `json:"EPOCH_BLOCKS,omitempty"`

	// +optional
	MempoolBlockBuffer uint8 `json:"MEMPOOL_BLOCK_BUFFER,omitempty"`

	// +optional
	MinPlotSize uint8 `json:"MIN_PLOT_SIZE,omitempty"`

	// +optional
	NetworkType uint8 `json:"NETWORK_TYPE,omitempty"`

	// +optional
	SubSlotItersStarting uint64 `json:"SUB_SLOT_ITERS_STARTING,omitempty"`

	// +optional
	HardForkHeight uint32 `json:"HARD_FORK_HEIGHT,omitempty"`

	// +optional
	SoftFork4Height uint32 `json:"SOFT_FORK4_HEIGHT,omitempty"`

	// +optional
	SoftFork5Height uint32 `json:"SOFT_FORK5_HEIGHT,omitempty"`

	// +optional
	PlotFilter128Height uint32 `json:"PLOT_FILTER_128_HEIGHT,omitempty"`

	// +optional
	PlotFilter64Height uint32 `json:"PLOT_FILTER_64_HEIGHT,omitempty"`

	// +optional
	PlotFilter32Height uint32 `json:"PLOT_FILTER_32_HEIGHT,omitempty"`
}

// ChiaNetworkStatus defines the observed state of ChiaNetwork
type ChiaNetworkStatus struct {
	// Ready says whether the ChiaNetwork is ready, which should be true when the ConfigMap is created
	// +kubebuilder:default=false
	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ChiaNetwork is the Schema for the chianetworks API
type ChiaNetwork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChiaNetworkSpec   `json:"spec,omitempty"`
	Status ChiaNetworkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ChiaNetworkList contains a list of ChiaNetwork
type ChiaNetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChiaNetwork `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChiaNetwork{}, &ChiaNetworkList{})
}
