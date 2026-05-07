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

## Trusted Peers

You can optionally specify a list of [CIDRs](https://aws.amazon.com/what-is/cidr/) that this ful_node should trust other full_node peers from. View the [Chia documentation on trusted peers](https://docs.chia.net/faq/?_highlight=trust#what-are-trusted-peers-and-how-do-i-add-them) to understand whether you should use this feature or not.

Here's an example ChiaNode that specifies trusted CIDRs:

```yaml
spec:
  chia:
    trustedCIDRs:
      - "192.168.1.0/24"
      - "10.0.0/8"
```

This specifies two trusted CIDRs, where if the IP address of a full_node peer is discovered to be within one of these two CIDR ranges, chia will consider that a trusted peer.

## chia-db-pull init container

ChiaNode supports an optional first-class `chia-db-pull` init container that downloads a chia blockchain database from an S3-compatible bucket into `CHIA_ROOT` before the chia container starts. This can dramatically reduce sync time for fresh nodes.

A minimal example:

```yaml
spec:
  chiaDBPull:
    enabled: true
    s3Prefix: "s3://chia-blockchain-sqlite-backups/testnet11/"
    network: "testnet11"
```

A more full example using a Secret for AWS credentials and a min-height threshold:

```yaml
spec:
  chiaDBPull:
    enabled: true
    s3Prefix: "s3://chia-blockchain-sqlite-backups/testnet11/"
    network: "testnet11"
    minHeight: 123456
    awsCredentialsSecret: aws-creds
```

The Secret referenced by `awsCredentialsSecret` is mounted via `envFrom`, so its keys (e.g. `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, optionally `AWS_SESSION_TOKEN`) become environment variables on the init container without needing to be in the CR in plaintext.

If you are running on EKS with IRSA (or another mechanism where the pod's ServiceAccount provides AWS credentials), simply omit `awsCredentialsSecret` and set the appropriate `serviceAccountName` on the ChiaNode.

When `chiaDBPull.enabled` is `true`, `chiaDBPull.s3Prefix` is required; the controller will refuse to reconcile and emit an event otherwise. The S3 bucket + path specified by `chiaDBPull.s3Prefix` must contain the following files: `blockchain_v2_${NETWORK}.sqlite`, `height-to-hash`, and `sub-epoch-summaries`. If any of those files are missing within the `s3Prefix` the init container will fail and the node won't start.

### Note on ordering with `spec.initContainers`

The first-class `chia-db-pull` container is appended to the StatefulSet's init container list **after** any containers defined in `spec.initContainers`. This means manually-defined init containers run first (good for things like clearing peer caches), and `chia-db-pull` is the last init container before the main chia container starts.

## More Info

This page contains documentation specific to this resource. Please see the rest of the documentation for information on more available configurations.

* [Generic options for all chia-operator resources.](all.md)
* [chia-exporter configuration](chia-exporter.md)
* [chia-healthcheck configuration](chia-healthcheck.md)
* [Services and networking](services-networking.md)
* [Storage](storage.md)
