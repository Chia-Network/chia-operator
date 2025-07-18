# Chia Operator Documentation

This documentation contains comprehensive guides for deploying and managing Chia services on Kubernetes.

## Quick Start

- **[Start a Farm](start-a-farm.md)** - Complete guide to setting up a Chia farm with all necessary components

## Core Components

### Farming Services

- **[ChiaCA](chiaca.md)** - Certificate Authority for secure communication between Chia services
- **[ChiaNode](chianode.md)** - Full node for blockchain synchronization and peer networking
- **[ChiaFarmer](chiafarmer.md)** - Farming service that creates blocks and earns rewards
- **[ChiaHarvester](chiaharvester.md)** - Plot management and harvesting service
- **[ChiaWallet](chiawallet.md)** - Wallet service for managing XCH and transactions

### Network Services

- **[ChiaCrawler](chiacrawler.md)** - Network crawler for peer discovery
- **[ChiaIntroducer](chiaintroducer.md)** - Introduction service for connecting new nodes to the network
- **[ChiaSeeder](chiaseeder.md)** - DNS seeder for providing initial peer connections
- **[ChiaTimelord](chiatimelord.md)** - Verifiable delay function service

### Data Services

- **[ChiaDataLayer](chiadatalayer.md)** - Data layer service for decentralized data storage

### Infrastructure

- **[ChiaCertificates](chiacertificates.md)** - Certificate management for secure communications
- **[ChiaNetwork](chianetwork.md)** - Network configuration and management

## Configuration Guides

### General Configuration

- **[All Components](all.md)** - Common configuration options applicable to all Chia resources
- **[Services and Networking](services-networking.md)** - Service configuration, load balancing, and networking options
- **[Storage](storage.md)** - Persistent volume and storage configuration
- **[Advanced](advanced.md)** - Advanced configurations including sidecars and init containers

### Monitoring and Health

- **[Chia Exporter](chia-exporter.md)** - Prometheus metrics collection and monitoring
- **[Chia Healthcheck](chia-healthcheck.md)** - Health checking and readiness probes

## Getting Help

Check out our **[Troubleshooting](troubleshooting.md)** guide for instant help with common issues.

If you still need help:

1. Check the [chia-operator GitHub Issues](https://github.com/Chia-Network/chia-operator/issues)
2. Review the [Chia documentation](https://docs.chia.net/)
3. Join the [Chia community Discord](https://discord.gg/chia)

When reporting issues, please include:

- Kubernetes version
- chia-operator version
- Resource definitions (with sensitive data removed)
- Pod logs and events
- Output of `kubectl describe` for affected resources
