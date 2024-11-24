# ChiaTimelord

Specifying a ChiaTimelord will create a kubernetes Deployment and some Services for a Chia timelord that connects to a local [full_node](chianode.md). It also requires a specified [Chia certificate authority](chiaca.md).

Here's a minimal ChiaTimelord example custom resource (CR):

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaTimelord
metadata:
  name: my-timelord
spec:
  chia:
    caSecretName: chiaca-secret # A kubernetes Secret containing certificate authority files
    # A full_node using kubernetes DNS names
    fullNodePeers:
      - host: "node.default.svc.cluster.local"
        port: 8444
```

## More Info

This page contains documentation specific to this resource. Please see the rest of the documentation for information on more available configurations.

* [Generic options for all chia-operator resources.](all.md)
* [chia-exporter configuration](chia-exporter.md)
* [Services and networking](services-networking.md)
* [Storage](storage.md)
