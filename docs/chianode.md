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

## Chia configuration

Some of Chia's configuration can be changed from within the CR.

```yaml
spec:
  chia:
    testnet: true # Switches to the default testnet in the Chia configuration file.
    timezone: "UTC" # Switches the tzdata timezone in the container.
    logLevel: "INFO" # Sets the Chia log level.
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

## chia-exporter sidecar

[chia-exporter](https://github.com/chia-network/chia-exporter) is a Prometheus exporter that surfaces scrape-able metrics to a Prometheus server. chia-exporter runs as a sidecar container to all Chia services ran by this operator by default.

### Add labels to chia-exporter service

You may want to add some labels to your chia-exporter Service that get added as labels to your Prometheus metrics.

```yaml
spec:
  chiaExporter:
    serviceLabels:
      network: "mainnet"
```

### Disable chia-exporter

```yaml
spec:
  chiaExporter:
    enabled: false
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
