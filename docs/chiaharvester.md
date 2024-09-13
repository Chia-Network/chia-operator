# ChiaHarvester

Specifying a ChiaHarvester will create a kubernetes Deployment and some Services for a Chia harvester that connects to a local [farmer](chiafarmer.md). It also requires a specified [Chia certificate authority](chiaca.md).

It is also expected you have some pre-existing plots in persistent volumes or mounted to a host path on one of your k8s nodes.

Here's a minimal ChiaHarvester example custom resource (CR):

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaHarvester
metadata:
  name: my-harvester
spec:
  chia:
    caSecretName: chiaca-secret # A kubernetes Secret containing certificate authority files
    farmerAddress: "farmer.default.svc.cluster.local" # A local farmer using kubernetes DNS names
```

## Plot storage

You can mount hostPath volumes or persistent volumes in a harvester pod using the following syntax. All claims/hostPaths get mounted as subdirectories of `/plots` in the container, and are mounted as read-only volumes. Harvesters ran with this operator set the `recursive_plot_scan` option to true.

```yaml
spec:
  storage:
    plots:
      persistentVolumeClaim:
        - claimName: "plot1"
        - claimName: "plot2"
      hostPathVolume:
        - path: "/home/user/storage/plots3"
        - path: "/home/user/storage/plots4"
```

If using a hostPath, you may want to pin the pod to a specific kubernetes node using a NodeSelector:

```yaml
spec:
  nodeSelector:
    kubernetes.io/hostname: "node-with-hostpath"
```

## More Info

This page contains documentation specific to this resource. Please see the rest of the documentation for information on more available configurations.

* [Generic options for all chia-operator resources.](all.md)
* [chia-exporter configuration](chia-exporter.md)
* [Services and networking](services-networking.md)
* [Storage](storage.md)
