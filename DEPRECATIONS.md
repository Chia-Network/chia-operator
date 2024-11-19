# Deprecations

## ChiaSeeder

### bootstrapPeer

Deprecated in: 0.12.5
Expected to be removed in: 1.0.0

The `bootstrapPeer` API field was deprecated in favor of `bootstrapPeers`. The latter being the preferred way to specify multiple bootstrap peers for a seeder installation, which still allows for specifying a single peer.

Switch to bootstrapPeers by specifying a yaml string list instead of a string:

```yaml
spec:
  chia:
    # bootstrapPeer: "mainnet-node.chia.svc.cluster.local"
    bootstrapPeers:
      - "mainnet-node.chia.svc.cluster.local"
```
