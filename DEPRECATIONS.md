# Deprecations

## ChiaFarmer

### fullNodePeer

* Deprecated in: 0.12.6
* Expected to be removed in: 1.0.0

The `fullNodePeer` API field was deprecated in favor of the more versatile `fullNodePeers` list field. The latter being the preferred way to specify one or multiple full_node peers for your configuration.

```yaml
spec:
  chia:
    # fullNodePeer: peer1:8444
    fullNodePeers:
      - host: peer1
        port: 8444
```

## ChiaTimelord

### fullNodePeer

* Deprecated in: 0.12.6
* Expected to be removed in: 1.0.0

The `fullNodePeer` API field was deprecated in favor of the more versatile `fullNodePeers` list field. The latter being the preferred way to specify one or multiple full_node peers for your configuration.

```yaml
spec:
  chia:
    # fullNodePeer: peer1:8444
    fullNodePeers:
      - host: peer1
        port: 8444
```

## ChiaWallet

### fullNodePeer

* Deprecated in: 0.12.6
* Expected to be removed in: 1.0.0

The `fullNodePeer` API field was deprecated in favor of the more versatile `fullNodePeers` list field. The latter being the preferred way to specify one or multiple full_node peers for your configuration.

```yaml
spec:
  chia:
    # fullNodePeer: peer1:8444
    fullNodePeers:
      - host: peer1
        port: 8444
```

## ChiaSeeder

### bootstrapPeer

* Deprecated in: 0.12.5
* Expected to be removed in: 1.0.0

The `bootstrapPeer` API field was deprecated in favor of `bootstrapPeers`. The latter being the preferred way to specify multiple bootstrap peers for a seeder installation, which still allows for specifying a single peer.

Switch to bootstrapPeers by specifying a yaml string list instead of a string:

```yaml
spec:
  chia:
    # bootstrapPeer: "mainnet-node.chia.svc.cluster.local"
    bootstrapPeers:
      - "mainnet-node.chia.svc.cluster.local"
```
