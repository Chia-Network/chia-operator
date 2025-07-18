# ChiaCrawler

Specifying a ChiaCrawler will create a kubernetes Deployment and some Services for a Chia crawler service. It optionally requires a specified [Chia certificate authority](chiaca.md) for secure communication.

A Chia crawler service helps discover and connect to other nodes in the Chia network by crawling the peer-to-peer network.

Here's a minimal ChiaCrawler example custom resource (CR):

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaCrawler
metadata:
  name: my-crawler
spec:
  chia:
    testnet: true # Optional: set to true for testnet, false for mainnet
```

## Certificate Authority (Optional)

While not required, you can specify a CA secret for secure communication:

```yaml
spec:
  chia:
    caSecretName: chiaca-secret # A kubernetes Secret containing certificate authority files
```

## CHIA_ROOT storage

`CHIA_ROOT` is an environment variable that tells chia services where to expect a data directory to be for local chia state. You can store your chia state persistently using either a host mount or a persistent volume claim.

To use a persistent volume claim:

```yaml
spec:
  storage:
    chiaRoot:
      persistentVolumeClaim:
        generateVolumeClaims: true
        resourceRequest: "1Gi"
        storageClass: "standard"
        accessModes:
          - "ReadWriteOnce"
```

To use a hostPath volume:

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
