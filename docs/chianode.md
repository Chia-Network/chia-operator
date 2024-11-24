# ChiaNode

Specifying a ChiaNode will create a kubernetes Statefulset and some Services for a Chia full_node. It also requires a specified [Chia certificate authority](chiaca.md).

Here's a minimal ChiaNode example custom resource (CR):

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaNode
metadata:
  name: my-node
spec:
  chia:
    caSecretName: chiaca-secret # A kubernetes Secret containing certificate authority files
```

## Replicas

To specify the number of replicas that are in the resulting StatefulSet, you can update `.spec.replicas` with an integer number of replicas.

```yaml
spec:
  replicas: 1
```

If you would like to ensure your replicas get scheduled on different kubernetes nodes, view the [Pod Affinity documentation.](all.md#pod-affinity)

## Full Node Peers

You may optionally specify a list of full_nodes peer(s) that your node will always try to remain connected to.

```yaml
spec:
  chia:
    # A local full_node using kubernetes DNS names
    fullNodePeers:
      - host: "node.default.svc.cluster.local"
        port: 8444
```

## CHIA_ROOT storage

`CHIA_ROOT` is an environment variable that tells chia services where to expect a data directory to be for local chia state. You can store your chia state persistently a couple of different ways: either with a host mount or a persistent volume claim.

To use a persistent volume claim, first create one in the same namespace and then give its name in the CR like the following:

```yaml
spec:
  storage:
    chiaRoot:
      persistentVolumeClaim:
        storageClass: ""
        resourceRequest: "300Gi"
```

To use a hostPath volume, first create a directory on the host and specify the path in the CR like the following:

```yaml
spec:
  storage:
    chiaRoot:
      hostPathVolume:
        path: "/home/user/storage/chiaroot"
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
* [chia-healthcheck configuration](chia-healthcheck.md)
* [Services and networking](services-networking.md)
* [Storage](storage.md)
