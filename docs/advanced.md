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

## Change arbitrary Chia configuration fields

Some pieces of the [Chia configuration file](https://github.com/Chia-Network/chia-blockchain/blob/main/chia/util/initial-config.yaml) do not have first-class settings available in the custom resources that this operator manages. First-class settings are ones that have a dedicated place in the Chia custom resource configuration. It's possible to change any piece of the config via the `additionalEnv` field, however.

```yaml
spec:
  chia:
    additionalEnv:
      # Sets a field in the yaml Chia config at the path 'full_node.enable_upnp' to 'False'
      # The variable must be prefixed with 'chia.'
      - name: "chia.full_node.enable_upnp"
        value: "False"
```

In this example, we disabled UPNP in the full_node config. If you later unset this, you may assume the setting would change back to "True" (the default in the Chia config.) This is not the case, assuming you mount your CHIA_ROOT in a persistent volume. First-class settings supported by chia-operator try to set defaults in the config back for you if you later undefine them. This is not the case for config settings changed through `additionalEnv`. If, for example, you want to re-enable UPNP in the future, you would need to set `chia.full_node.enable_upnp: True`, rather than undefine it. After that change is applied, however, you can undefine it from the custom resource if you would like.

You should only use `additionalEnv` to specify Chia config fields that don't have first-class settings support. A setting defined in `additionalEnv` should take precedence over the same first-class setting, but there's no guarantee of this, and you're asking for headaches for no good reason.
