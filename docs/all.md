# All

This documentation is meant to be applicable generally to all resources the operator manages.

## Chia configuration

Some of Chia's configuration can be changed from within the CR.

```yaml
spec:
  chia:
    testnet: true # Switches to the default testnet in the Chia configuration file.
    timezone: "UTC" # Switches the tzdata timezone in the container.
    logLevel: "INFO" # Sets the Chia log level.
    selfHostname: "0.0.0.0" # Sets the self_hostname setting in the config, which affects what network interfaces chia services are bound to.
```

### Selecting a network

You can select a network from your chia configuration with the following options:

```yaml
spec:
  chia:
    network: "testnetZZ" # Switches to the network with a matching name in the chia config file.
    networkPort: 58445 # Switches the default network port full_nodes connect with.
    introducerAddress: "introducer.default.svc.cluster.local" # Sets the introducer address used in the chia config file.
    dnsIntroducerAddress: "dns-introducer.default.svc.cluster.local" # Sets the DNS introducer address used in the chia config file.
```

## Chia container resource requests and limits

You can set resource requests and limits for the chia container deployed from a custom resource with the following (note that these are just example values, and not to be taken as recommendations for your deployments):

```yaml
spec:
  resources:
    requests:
      memory: "256Mi"
      cpu: "500m"
    limits:
      memory: "1028Mi"
      cpu: "1000m"
```

Before setting these, ensure you have an idea of how much memory and cpu the chia service being deployed tends to use under normal circumstances. If too low of a limit is specified, the chia container may restart often. If given too great of requests, you may be wasting some of the scheduling capabilities of a kubernetes node.

## Pod Affinity

You can set Pod affinity and anti-affinity rules for any custom resource like so (this is just an example):

```yaml
spec:
  affinity:
    # These are just example anti-affinity rules
    podAntiAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        - labelSelector:
            matchExpressions:
              - key: app
                operator: In
                values:
                  - my-app
          topologyKey: "kubernetes.io/hostname"
```

## Pod Security Contexts

This sets the securityContext for a pod.

View the Kubernetes [documentation](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/) for examples of applicable Pod securityContext fields.

```yaml
spec:
  podSecurityContext:
    fsGroup: 2000
```

## Container Security Contexts

This sets the securityContext for the chia container in a Pod.

View the Kubernetes [documentation](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/) for examples of applicable container securityContext fields.

```yaml
spec:
  chia:
    securityContext:
      allowPrivilegeEscalation: false
```

## Node Selectors

You can pin a Pod to a specific node using labels from that node in a nodeSelector.

```yaml
spec:
  nodeSelctor:
    "kubernetes.io/hostname": "worker1"
```

## Update Strategy

You can set a custom update strategy using [kubernetes Deployment update strategy](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy) definitions.

```yaml
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
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
