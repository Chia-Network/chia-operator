# Advanced

This documentation describes some advanced (or uncommon) usages of chia-operator. "Advanced" or "uncommon" will be defined as not necessary to run chia, and more for running services that utilize chia installations in some manner (via network requests or otherwise.)

## Sidecar containers

You can run a container in the same kubernetes Pod as your chia components utilizing the `spec.sidecars` segment of all Chia resources supported by this operator except ChiaCAs and ChiaNetworks.

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
    - container:
        name: my-sidecar-container
        image: nginx:latest
        ports:
          - containerPort: 80
            name: http
        env:
          - name: SIDECAR_CONTAINER_VAR
            value: "sidecar_container_value"
        volumeMounts:
          - name: sidecar-data
            mountPath: /usr/share/nginx/html
      volumes:
        - name: sidecar-data
          emptyDir: {}
      shareVolumeMounts: true
      shareEnv: true
```

`sidecars` is a list that contains a few keys: `container`, `volumes`, `shareVolumeMounts` and `shareEnv`.

* `container` is just a normal kubernetes container specification which can contain a name, image, environment variables, volumeMounts, etc.
* `volumes` is a list of kubernetes Volumes that should be added to the Pod specification. You will need to add these to the sidecar container's volumeMounts as shown above.
* `shareVolumeMounts` if set to true, gives the sidecar container the same volume mounts as the chia container, in the same mountpoints as the chia container.
* `shareEnv` if set to true, gives the sidecar container the same environment variables as the chia container.

You may specify multiple sidecar containers, additional volumes, and toggle the shareVolumtMounts/shareEnv options separately for each of them in this way.

## Init containers

You can run a container as an init container in the same kubernetes Pod as your chia components utilizing the `spec.initContainers` segment of all Chia resources supported by this operator except ChiaCAs and ChiaNetworks.

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
        volumeMounts:
          - name: init-data
            mountPath: /usr/share/nginx/html
      volumes:
        - name: init-data
          emptyDir: {}
      shareVolumeMounts: true
      shareEnv: true
```

`initContainers` is a list that contains a few keys: `container`, `volumes`, `shareVolumeMounts` and `shareEnv`.

* `container` is just a normal kubernetes container specification which can contain a name, image, environment variables, volumeMounts, etc.
* `volumes` is a list of kubernetes Volumes that should be added to the Pod specification. You will need to add these to the init container's volumeMounts as shown above.
* `shareVolumeMounts` if set to true, gives the init container the same volume mounts as the chia container, in the same mountpoints as the chia container.
* `shareEnv` if set to true, gives the init container the same environment variables as the chia container.

You may specify multiple init containers, additional volumes, and toggle the shareVolumtMounts/shareEnv options separately for each of them in this way.

## Specify a Chia version

Operator releases tend to pin to the current latest version of chia (at the time the release was published) but if you'd like to manage the version of chia ran yourself, there's a field to do so:

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaNode
metadata:
  name: my-node
spec:
  chia:
    image: ghcr.io/chia-network/chia:2.4.3
```

The example shows a ChiaNode (full_node) resource on v2.4.3 of chia, but this field is also available on other resources.

Since this is an image field, you can point to any OCI image containing chia, but note that this operator makes heavy use of the [chia-docker](https://github.com/Chia-Network/chia-docker) entrypoint script for setting a lot of the chia configuration, so it should be compatible with that script to ensure your Chia services start up properly. Using an image that isn't at least based on the official chia-docker image will likely result in a broken installation.
