# Advanced

This documentation describes some advanced (or uncommon) usages of chia-operator. "Advanced" or "uncommon" will be defined as not necessary to run chia, and more for running services that utilize chia installations in some manner (via network requests or otherwise.)

# Sidecar containers

You can run a container in the same kubernetes Pod as your chia components utilizing the `sidecarContainer` segment of most of the custom resources managed by this operator.

Supported custom resources (CRs):

- ChiaFarmer
- ChiaHarvester
- ChiaNode
- ChiaSeeder
- ChiaTimelord
- ChiaWallet

## Usage

To create a Chia component that runs some sidecar container, do something like the following:

```
apiVersion: k8s.chia.net/v1
kind: ChiaNode
metadata:
  name: my-node
spec:
  chia:
    [...]
  chiaExporter:
    [...]
  sidecarContainers:
    - name: nginx
      image: nginx:latest
      ports:
      - containerPort: 80
        name: http
      env:
      - name: SIDECAR_VAR
        value: "sidecar_value"
      volumeMounts:
      - name: chiaroot
        mountPath: /data
```

If you were to apply this to a cluster, it would create a Statefulset with 3 containers per Pod replica. The container names would be `chia`, `chia-exporter`, and `nginx`. The `nginx` container would expose containerPort 80, an environment variable named `SIDECAR_VAR`, and it would mount the main CHIA_ROOT volume.
