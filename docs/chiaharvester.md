# ChiaHarvester

Specifying a ChiaHarvester will create a kubernetes Deployment and some Services for a Chia harvester that connects to a local [farmer](chiafarmer.md). It also requires a specified [Chia certificate authority](chiaca.md).

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

## Chia configuration

Some of Chia's configuration can be changed from within the CR.

```yaml
spec:
  chia:
    testnet: true # Switches to the default testnet in the Chia configuration file.
    timezone: "UTC" # Switches the tzdata timezone in the container.
    logLevel: "INFO" # Sets the Chia log level.
```

## Plot storage

You can mount hostPath volumes or persistent volumes in a harvester pod using the following syntax. All claims/hostPaths get mounted as sub-directories of `/plots` in the container. Harvesters ran with this operator set the `recursive_plot_scan` option to true.

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

## CHIA_ROOT storage

`CHIA_ROOT` is an environment variable that tells chia services where to expect a data directory to be for local chia state. You can store your chia state persistently a couple of different ways: either with a host mount or a persistent volume claim.

To use a persistent volume claim, first create one in the same namespace and then give its name in the CR like the following:

```yaml
spec:
  storage:
    chiaRoot:
      persistentVolumeClaim:
        claimName: "chiaroot-data"
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

## Selecting a network

You can select a network from your chia configuration with the following options:

```yaml
spec:
  chia:
    network: "testnetZZ" # Switches to the network with a matching name in the chia config file.
    networkPort: 58445 # Switches the default network port full_nodes connect with.
    introducerAddress: "introducer.default.svc.cluster.local" # Sets the introducer address used in the chia config file.
    dnsIntroducerAddress: "dns-introducer.default.svc.cluster.local" # Sets the DNS introducer address used in the chia config file.
```

## Configure Readiness, Liveness, and Startup probes

By default, if chia-exporter is enabled it comes with its own readiness and liveness probes. But you can configure readiness, liveness, and startup probes for the chia container in your deployed Pods, too:

```yaml
spec:
  chia:
    livenessProbe:
      exec:
        command:
          - /bin/sh
          - '-c'
          - /usr/local/bin/docker-healthcheck.sh || exit 1
      initialDelaySeconds: 30
    readinessProbe:
      exec:
        command:
          - /bin/sh
          - '-c'
          - /usr/local/bin/docker-healthcheck.sh || exit 1
      initialDelaySeconds: 30
    startupProbe:
      exec:
        command:
          - /bin/sh
          - '-c'
          - /usr/local/bin/docker-healthcheck.sh || exit 1
```

## Update Strategy

You can set a custom update strategy using [kubernetes Deployment update strategy](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy) definitions.

Example:

```yaml
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
```
