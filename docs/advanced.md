# Advanced

This documentation describes some advanced (or uncommon) usages of chia-operator. "Advanced" or "uncommon" will be defined as not necessary to run chia, and more for running services that utilize chia installations in some manner (via network requests or otherwise.)

## Sidecar containers

You can run a container in the same kubernetes Pod as your chia components utilizing the `spec.sidecars` segment of all Chia resources supported by this operator except ChiaCAs.

To create a Chia component that runs some sidecar container, do something like the following:

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaNode
metadata:
  name: my-node
spec:
  chia:
    [...]
  chiaExporter:
    [...]
  sidecars:
    containers:
    - name: sidecar
      image: nginx:latest
      ports:
      - containerPort: 80
        name: http
      env:
      - name: SIDECAR_CONTAINER_VAR
        value: "sidecar_container_value"
      volumeMounts:
      - name: chiaroot
        mountPath: /data
      - name: sidecar-data
        mountPath: /usr/share/nginx/html
    volumes:
    - name: sidecar-data
      emptyDir: {}
```

If you were to apply this to a cluster, it would create a Statefulset with 3 containers per Pod replica. The container names would be `chia`, `chia-exporter`, and `nginx`. The `nginx` container would expose containerPort 80, an environment variable named `SIDECAR_VAR`, and it would mount the main CHIA_ROOT volume as well as an emptydir volume that we specified for this sidecar that neither the `chia` or `chia-exporter` containers would mount.

## Init containers

You can run a container as an init container in the same kubernetes Pod as your chia components utilizing the `spec.initContainer` segment of all Chia resources supported by this operator except ChiaCAs.

To create a Chia component that runs an init container, do something like the following:

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaNode
metadata:
  name: my-node
spec:
  chia:
    [...]
  chiaExporter:
    [...]
  initContainers:
    - container:
        name: my-init-container
        image: nginx:latest
        ports:
        - containerPort: 80
          name: http
        env:
        - name: INIT_CONTAINER_VAR
          value: "init_container_value"
```

`initContainers` is a list of the normal kubernetes container specification. It does not support or respect setting the volumeMounts field in the container, however. Any volumeMounts specified will be overwritten.

### Share Chia Volumes

You can share volumes from the main chia container to your init containers using the following:

```yaml
spec:
  initContainers:
    - shareVolumeMounts: true # Option to share the volume mounts from the chia container, useful if you specified a CHIA_ROOT volume and want to add some data to it before chia starts
      container:
        name: my-init-container
        image: nginx:latest
        ports:
        - containerPort: 80
          name: http
        env:
        - name: INIT_CONTAINER_VAR
          value: "init_container_value"
```

### Share Chia Env

You can share environment variable from the main chia container to your init containers using the following:

```yaml
spec:
  initContainers:
    - shareEnv: true # Option to share the environment variables from the chia container, useful if the init container image is a derivative of the chia-docker image
      container:
        name: my-init-container
        image: nginx:latest
        ports:
        - containerPort: 80
          name: http
        env:
        - name: INIT_CONTAINER_VAR
          value: "init_container_value"
```
