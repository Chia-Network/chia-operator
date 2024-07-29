# Storage

## Chia Root

This section defines the configuration for storing your CHIA_ROOT data somewhere persistently.

### Persistent Volumes

You can have the operator create persistent volume claims for your Chia resource Deployments.

This uses a ChiaFarmer as an example but the same applies to all Chia components except for ChiaNodes which use Statefulsets and are configured a slightly different way. See the ChiaNode documentation for more details.

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaFarmer
metadata:
  name: my-farmer
spec:
  storage:
    chiaRoot:
      persistentVolumeClaim:
        generateVolumeClaims: true
        resourceRequest: 2Gi
        storageClass: "local-path"
        accessModes:
          - "ReadWriteOnce"
```

The `accessModes` field is optional and will default to ReadWriteOnce if unspecified.

If you have a pre-existing persistent volume claim that you would like to use, simply specify `claimName` instead, like so:

```yaml
spec:
  storage:
    chiaRoot:
      persistentVolumeClaim:
        claimName: "chiaroot-data"
```

### Hostpath Volumes

Sometimes your persistent data is just on a particular host, rather than in a kubernetes persistent volume. In that case, usually you need to define the host's path and a NodeSelector which pins the deployed Pods to the node you select. This uses a ChiaFarmer as an example but the same goes for any other resource:

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaFarmer
metadata:
  name: my-farmer
spec:
  nodeSelector:
    kubernetes.io/hostname: "node-with-hostpath"
  storage:
    chiaRoot:
      hostPathVolume:
        path: "/home/user/storage/chiaroot"
```

The `.spec.nodeSelector` field defines a label that exists on the particular kubernetes node to pin the Pod to. And the `.spec.storage.chiaRoot.hostPathVolume.path` field defines the path on the host to a directory containing your CHIA_ROOT data.
