/*
Copyright 2023 Chia Network Inc.
*/

package consts

// ControllerOwner bool to help set the controller owner for a create kubernetes Kind
var ControllerOwner = true

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

	// walletRPCPort defines the port for the wallet RPC
	WalletRPCPort = 9256

	// ChiaExporterPort defines the port for Chia Exporter instances
	ChiaExporterPort = 9914

	// DefaultChiaExporterImage is the default image name and tag of the chia-exporter image
	DefaultChiaExporterImage = "ghcr.io/chia-network/chia-exporter:latest"
)
