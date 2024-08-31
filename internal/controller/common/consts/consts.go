/*
Copyright 2023 Chia Network Inc.
*/

package consts

// ChiaKind enumerates the list of Chia component custom resources this operator controls
type ChiaKind string

const (
	// ChiaCAKind is the API Kind for Chia certificate authorities
	ChiaCAKind ChiaKind = "ChiaCA"

	// ChiaCrawlerKind is the API Kind for Chia crawlers
	ChiaCrawlerKind ChiaKind = "ChiaCrawler"

	// ChiaFarmerKind is the API Kind for Chia farmers
	ChiaFarmerKind ChiaKind = "ChiaFarmer"

	// ChiaHarvesterKind is the API Kind for Chia harvesters
	ChiaHarvesterKind ChiaKind = "ChiaHarvester"

	// ChiaIntroducerKind is the API Kind for Chia introducers
	ChiaIntroducerKind ChiaKind = "ChiaIntroducer"

	// ChiaNodeKind is the API Kind for Chia full_nodes
	ChiaNodeKind ChiaKind = "ChiaNode"

	// ChiaSeederKind is the API Kind for Chia seeders / dns-introducers
	ChiaSeederKind ChiaKind = "ChiaSeeder"

	// ChiaWalletKind is the API Kind for Chia wallets
	ChiaWalletKind ChiaKind = "ChiaWallet"
)

// API default image constants
const (
	// DefaultChiaCAImageName contains the default image name for the ca-gen image
	DefaultChiaCAImageName = "ghcr.io/chia-network/chia-operator/ca-gen"

	// DefaultChiaCAImageTag contains the default tag name for the ca-gen image
	DefaultChiaCAImageTag = "0.7.5"

	// DefaultChiaImageName contains the default image name for the chia-docker image
	DefaultChiaImageName = "ghcr.io/chia-network/chia"

	// DefaultChiaImageTag contains the default tag name for the chia-docker image
	DefaultChiaImageTag = "2.4.3"

	// DefaultChiaExporterImageName contains the default image name for the chia-exporter image
	DefaultChiaExporterImageName = "ghcr.io/chia-network/chia-exporter"

	// DefaultChiaExporterImageTag contains the default tag name for the chia-exporter image
	DefaultChiaExporterImageTag = "0.15.3"

	// DefaultChiaHealthcheckImageName contains the default image name for the chia-healthcheck image
	DefaultChiaHealthcheckImageName = "ghcr.io/chia-network/chia-healthcheck"

	// DefaultChiaHealthcheckImageTag contains the default tag name for the chia-healthcheck image
	DefaultChiaHealthcheckImageTag = "0.2.1"
)

const (
	// DaemonPort defines the port for the Chia daemon
	DaemonPort = 55400

	// FarmerPort defines the port for farmer instances
	FarmerPort = 8447

	// FarmerRPCPort defines the port for the farmer RPC
	FarmerRPCPort = 8559

	// HarvesterPort defines the port for harvester instances
	HarvesterPort = 8448

	// HarvesterRPCPort defines the port for the harvester RPC
	HarvesterRPCPort = 8560

	// MainnetNodePort defines the port for mainnet nodes
	MainnetNodePort = 8444

	// TestnetNodePort defines the port for testnet nodes
	TestnetNodePort = 58444

	// NodeRPCPort defines the port for the full_node RPC
	NodeRPCPort = 8555

	// CrawlerRPCPort defines the port for the crawler RPC
	CrawlerRPCPort = 8561

	// TimelordPort defines the port for timelord
	TimelordPort = 8446

	// TimelordRPCPort defines the port for the timelord RPC
	TimelordRPCPort = 8557

	// WalletPort defines the port for wallet instances
	WalletPort = 8449

	// WalletRPCPort defines the port for the wallet RPC
	WalletRPCPort = 9256

	// ChiaExporterPort defines the port for Chia Exporter instances
	ChiaExporterPort = 9914

	// ChiaHealthcheckPort defines the port for Chia Healthcheck instances
	ChiaHealthcheckPort = 9950
)
