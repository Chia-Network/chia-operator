# ChiaDataLayer

Specifying a ChiaDataLayer will create a Kubernetes Deployment and Services for a Chia DataLayer server that connects to a local [full_node](chianode.md). It also requires a specified [Chia certificate authority](chiaca.md).

It is also expected you have a pre-existing Chia key to import, likely one that you generated locally in a Chia GUI installation.

Here's a minimal ChiaDataLayer example custom resource (CR):

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaDataLayer
metadata:
  name: my-datalayer
spec:
  chia:
    caSecretName: chiaca-secret # A kubernetes Secret containing certificate authority files
    # A local full_node using kubernetes DNS names
    fullNodePeers:
      - host: "node.default.svc.cluster.local"
        port: 8444
    # A kubernetes Secret named chiakey-secret containing a key.txt file with your mnemonic key
    secretKey:
      name: "chiakey-secret"
      key: "key.txt"
```

## Secret key

The `secretKey` field in the ChiaDataLayer's spec defines the name of a Kubernetes Secret that contains your mnemonic. Only Wallets and Farmers need your mnemonic key to function. You can create your Kubernetes Secret like so:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: chiakey-secret
stringData:
  key.txt: "your mnemonic goes here"
type: Opaque
```

Replace the text value for `key.txt` with your mnemonic, and then reference it in your ChiaDataLayer resource in the way shown above.

## Fileserver Configuration

The ChiaDataLayer can optionally run a fileserver sidecar container to serve the data_layer server files. This is disabled by default but can be enabled with the following configuration:

```yaml
spec:
  fileserver:
    enabled: true
    # Optional custom image for the fileserver
    image: "custom/fileserver:tag"
    # Optional custom mount path for server files
    serverFileMountpath: "/custom/path"
    # Optional custom container port
    containerPort: 8080
    # Optional service configuration
    service:
      enabled: true
      type: ClusterIP
      externalTrafficPolicy: Local
```

### Additional Environment Variables

You can add custom environment variables to the fileserver container using the `additionalEnv` field:

```yaml
spec:
  fileserver:
    enabled: true
    additionalEnv:
      - name: CUSTOM_VAR
        value: "custom-value"
      - name: SECRET_VAR
        valueFrom:
          secretKeyRef:
            name: my-secret
            key: secret-key
```

### Container Health Checks

The fileserver container supports standard Kubernetes probes for health checking:

```yaml
spec:
  fileserver:
    enabled: true
    # Liveness probe to check if container is running properly
    livenessProbe:
      httpGet:
        path: /health
        port: 8080
      initialDelaySeconds: 30
      periodSeconds: 10
    # Readiness probe to check if container is ready to accept traffic
    readinessProbe:
      httpGet:
        path: /ready
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 5
    # Startup probe to give container time to initialize
    startupProbe:
      httpGet:
        path: /startup
        port: 8080
      failureThreshold: 30
      periodSeconds: 10
```

### Resource Requirements

You can specify resource limits and requests for the fileserver container:

```yaml
spec:
  fileserver:
    enabled: true
    resources:
      requests:
        memory: "256Mi"
        cpu: "250m"
      limits:
        memory: "512Mi"
        cpu: "500m"
```

### Security Context

You can configure the security context for the fileserver container:

```yaml
spec:
  fileserver:
    enabled: true
    securityContext:
      runAsNonRoot: true
      runAsUser: 1000
      allowPrivilegeEscalation: false
      capabilities:
        drop:
          - ALL
```

## Trusted Peers

You can specify trusted CIDRs for your DataLayer server using the `trustedCIDRs` field:

```yaml
spec:
  chia:
    trustedCIDRs:
      - "192.168.1.0/24"
      - "10.0.0.0/8"
```

## More Info

This page contains documentation specific to this resource. Please see the rest of the documentation for information on more available configurations.

* [Generic options for all chia-operator resources.](all.md)
* [chia-exporter configuration](chia-exporter.md)
* [Services and networking](services-networking.md)
* [Storage](storage.md)
