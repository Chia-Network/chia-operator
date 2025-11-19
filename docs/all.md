# All

This documentation is meant to be applicable generally to all resources the operator manages.

## Table of Contents

- [Chia Configuration](#chia-configuration)
  - [Network Selection](#selecting-a-network)
  - [Specify a Chia image](#specify-a-chia-image)
  - [Install from Specific Ref](#install-chia-from-a-specific-ref)
- [Requests and Limits](#chia-container-resource-requests-and-limits)
- [Environment Variables](#chia-container-additional-environment-variables)
- [Pod Affinity](#pod-affinity)
- [Topology Spread Constraints](#topology-spread-constraints)
- [Pod Security Contexts](#pod-security-contexts)
- [Container Security Contexts](#container-security-contexts)
- [Node Selectors](#node-selectors)
- [Update Strategies](#update-strategy)
- [Health Checks](#configure-readiness-liveness-and-startup-probes)
- [Image Pull Secret](#specify-image-pull-secrets)
- [Image Pull Policy](#specify-image-pull-policy)
- [Service Account](#specify-a-service-account)

## Chia configuration

Some of Chia's configuration can be changed from within the CR.

```yaml
spec:
  chia:
    testnet: true # Switches to the default testnet in the Chia configuration file.
    timezone: "UTC" # Switches the tzdata timezone in the container.
    logLevel: "INFO" # Sets the Chia log level.
    selfHostname: "127.0.0.1" # Sets the self_hostname setting in the config, which affects what network interfaces chia services are bound to (defaults to 0.0.0.0) 
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

### Specify a Chia image

Operator releases tend to pin to the latest version of chia (at the time the release was published) but if you'd like to manage the version of chia ran yourself, there's a field to do so:

```yaml
spec:
  chia:
    image: ghcr.io/chia-network/chia:2.5.7
```

Since this is an image field, you can point to any OCI image containing chia, but note that this operator makes heavy use of the [chia-docker](https://github.com/Chia-Network/chia-docker) entrypoint script for setting a lot of the chia configuration, so it should be compatible with that script to ensure your Chia services start up properly. Using an image that isn't at least based on the official chia-docker image will likely result in a broken installation.

### Install chia from a specific ref

You can select a specific ref (commit sha or branch) from the chia-blockchain repository to install chia from. This is unnecessary the majority of the time as the image has a chia installation by default, but this may be useful for testing specific versions of chia:

```yaml
spec:
  chia:
    sourceRef: "a6a27bfe8e8d3e3db16701e7a33182ac11ce0723" # commit sha of github.com/Chia-Network/chia-blockchain to install from
```

Note that if you use this configuration, the tag of the chia image running in your Pods may still specify a version of chia-blockchain, but is no longer the version of chia installed in the image.

## Chia container resource requests and limits

You can set resource requests and limits for the chia container deployed from a custom resource with the following (note that these are just example values, and not to be taken as recommendations for your deployments):

```yaml
spec:
  chia:
    resources:
      requests:
        memory: "256Mi"
        cpu: "500m"
      limits:
        memory: "1028Mi"
        cpu: "1000m"
```

Before setting these, ensure you have an idea of how much memory and cpu the chia service being deployed tends to use under normal circumstances. If too low of a limit is specified, the chia container may restart often. If given too great of requests, you may be wasting some of the scheduling capabilities of a kubernetes node.

## Chia container additional environment variables

WARNING: This is a dangerous feature for advanced use cases. This can override variables set by the operator, which may cause issues deploying chia services. The vast majority of uses for additional environment variables should be solved with settings within the CRD. But under some circumstances you may want to set some chia container env for niche settings for chia-tools to plant in your chia configuration. If you have a use for setting additional container env, please make an Issue so we can discuss potentially adding it to the CRD.

```yaml
spec:
  chia:
    additionalEnv:
      # Set normal environment variable key-value
      - name: FOO
        value: BAR
      # Set using valueFrom
      - name: HELLO
        valueFrom:
          configMapKeyRef:
            name: my-configmap
            key: WORLD
```

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

## Topology Spread Constraints

You can also use `TopologySpreadConstraints` to describe how a group of pods ought to spread across topology domains. The Scheduler will schedule pods in a way which abides by the constraints. All topology spread constraints are ANDed.

```yaml
spec:
  topologySpreadConstraints:
    - maxSkew: 1
      topologyKey: "kubernetes.io/hostname"
      whenUnsatisfiable: DoNotSchedule
      labelSelector:
        matchExpressions:
          - key: app
            operator: In
            values:
              - my-app
```

## Tolerations

You can also use `Tolerations` to allow scheduling the pod to a node with a matching taint.

```yaml
spec:
  tolerations:
    - key: "key1"
      operator: "Exists"
      effect: "NoSchedule"
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
  nodeSelector:
    "kubernetes.io/hostname": "worker1"
```

## Update Strategy

You can set a custom update strategy using [kubernetes Deployment update strategy](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy) definitions.

NOTE: This applies to all resources that deploy Pods except for ChiaNodes.

```yaml
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
```

### ChiaNode Update Strategies

ChiaNodes deploy StatefulSet resources which use a different update strategy definition. See the documentation for [kubernetes StatefulSet update strategy](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#update-strategies) definitions.

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaNode
metadata:
  name: my-node
spec:
  updateStrategy:
    type: RollingUpdate
```

## Configure Readiness, Liveness, and Startup probes

By default, if running a service supported by [chia-healthcheck](chia-healthcheck.md), and chia-healthcheck is enabled (it is enabled by default), then some startup, readiness, and liveness probes will be configured for the chia container using endpoints from the chia-healthcheck sidecar.

If chia-healthcheck is not running as a sidecar to the chia container, then no readiness, liveness, or startup probes are configured by default. These can be configured in the chia container's spec though. The following example configures the simple healthcheck script contained in the chia-docker image:

```yaml
spec:
  chia:
    livenessProbe:
      exec:
        command:
          - /bin/sh
          - '-c'
          - /usr/local/bin/docker-healthcheck.sh || exit 1
    readinessProbe:
      exec:
        command:
          - /bin/sh
          - '-c'
          - /usr/local/bin/docker-healthcheck.sh || exit 1
    startupProbe:
      exec:
        command:
          - /bin/sh
          - '-c'
          - /usr/local/bin/docker-healthcheck.sh || exit 1
```

The chia-exporter and chia-healthcheck containers come with their own readiness, liveness, and startup probes which are enabled by default.

## Specify Image Pull Secrets

Most of the time you won't need to specify imagePullSecrets when using this operator, but if you specify custom images from your own registries for init containers, sidecar containers, or custom versions for any of the default containers this operator supports, you can specify imagePullSecrets to pull those images. These are dockerconfigjson Secrets you would have created in the cluster yourself. Here's an example of how you would specify imagePullSecrets in your chia custom resource:

```yaml
spec:
  imagePullSecrets:
    - name: my-registry-secret
```

And here is an example dockerconfigjson Secret, make sure this is installed in the cluster in the same namespace as your chia custom resource:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-registry-secret
type: kubernetes.io/dockerconfigjson
data:
  .dockerconfigjson: <base64-encoded-auth-secret>
```

## Specify Image Pull Policy

If you need to specify your image pull policy for container images:

```yaml
spec:
  imagePullPolicy: "IfNotPresent"
```

## Specify a Service Account

If you need to specify an existing ServiceAccount for your chia deployments, you can do so. This assumes the ServiceAccount already exists in the same namespace as this Chia resource, it won't create one for you.

```yaml
spec:
  serviceAccountName: "my-service-account"
```
