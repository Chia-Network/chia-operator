# ChiaDataLayer

ChiaDataLayers run the data_layer Chia component, which comes bundled with a Chia wallet. In a future time, the wallet may be able to be run separately, but it is not currently possible.

The data_layer_http server runs as an optional sidecar. In a future release, it may be possible to run the HTTP server separately from the data_layer server, but it is not currently implemented.

Here's a minimal ChiaDataLayer example custom resource (CR):

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaDataLayer
metadata:
  name: mainnet
spec:
  chia:
    caSecretName: "chiaca-secret" # A kubernetes Secret containing certificate authority files
    # A kubernetes Secret named chiakey-secret containing a key.txt file with your mnemonic key
    secretKey:
      name: "chiakey-secret"
      key: "key.txt"
  dataLayerHTTP:
    enabled: true # Enabled the data_layer_http sidecar service
```

## Full Node Peers

You may optionally specify a list of full_nodes peer(s) that your wallet will always try to remain connected to.

```yaml
spec:
  chia:
    # A local full_node using kubernetes DNS names
    fullNodePeers:
      - host: "node.default.svc.cluster.local"
        port: 8444
```

## Server files storage

Datalayer stores its server files in `/datalayer/server` inside the container. You can set a persistent volume for this directory by adding the following:

```yaml
spec:
  storage:
    dataLayerServerFiles:
      persistentVolumeClaim:
        generateVolumeClaims: true
        storageClass: ""
        resourceRequest: "10Gi"
```

## More Info

This page contains documentation specific to this resource. Please see the rest of the documentation for information on more available configurations.

* [Generic options for all chia-operator resources.](all.md)
* [chia-exporter configuration](chia-exporter.md)
* [chia-healthcheck configuration](chia-healthcheck.md)
* [Services and networking](services-networking.md)
* [Storage](storage.md)
