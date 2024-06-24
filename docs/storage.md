# Storage

You can have the operator create persistent volume claims for your Chia resource Deployments.

This applies to all Chia components except for ChiaNodes which use Statefulsets and are configured a slightly different way. See the ChiaNode documentation for more details.

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
