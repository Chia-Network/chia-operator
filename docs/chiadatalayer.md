# ChiaDataLayer

ChiaDataLayers run the data_layer Chia component, which comes bundled with a Chia wallet. In a future time, the wallet may be able to be run separately, but it is not currently possible.

The data_layer_http server runs as an optional sidecar. In a future release, it may be possible to run the HTTP server separately from the data_layer server, but it is not currently implemented.

Here's a minimal ChiaDataLayer example custom resource (CR):

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaDataLayer
metadata:
  name: mainnet
spec:
  chia:
    caSecretName: "chiaca-secret" # A kubernetes Secret containing certificate authority files
    # A kubernetes Secret named chiakey-secret containing a key.txt file with your mnemonic key
    secretKey:
      name: "chiakey-secret"
      key: "key.txt"
  dataLayerHTTP:
    enabled: true # Enabled the data_layer_http sidecar service
```

## Full Node Peers

You may optionally specify a list of full_nodes peer(s) that your wallet will always try to remain connected to.

```yaml
spec:
  chia:
    # A local full_node using kubernetes DNS names
    fullNodePeers:
      - host: "node.default.svc.cluster.local"
        port: 8444
```

## Trusted Peers

You can optionally specify a list of [CIDRs](https://aws.amazon.com/what-is/cidr/) that the wallet ran alongside data_layer should trust full_node peers from. View the [Chia documentation on trusted peers](https://docs.chia.net/faq/?_highlight=trust#what-are-trusted-peers-and-how-do-i-add-them) to understand whether you should use this feature or not.

Here's an example ChiaDataLayer that specifies trusted CIDRs:

```yaml
spec:
  chia:
    trustedCIDRs:
      - "192.168.1.0/24"
      - "10.0.0/8"
```

## Server files storage

Datalayer stores its server files in `/datalayer/server` inside the container. You can set a persistent volume for this directory by adding the following:

```yaml
spec:
  storage:
    dataLayerServerFiles:
      persistentVolumeClaim:
        generateVolumeClaims: true
        storageClass: ""
        resourceRequest: "10Gi"
```

## HTTP File Servers

The `ChiaDataLayer` resource supports two options for serving files over HTTP:

### DataLayer HTTP

The built-in `data_layer_http` sidecar is a Chia component that provides HTTP access to the data layer server files. This is the official Chia implementation for serving files.

```yaml
apiVersion: chia.network/v1
kind: ChiaDataLayer
metadata:
  name: my-datalayer
spec:
  dataLayerHTTP:
    enabled: true  # Enable the data_layer_http sidecar
    service:
      type: ClusterIP  # Optional - defaults to ClusterIP
      # Available service options:
      # externalTrafficPolicy
      # sessionAffinity
      # sessionAffinityConfig
      # ipFamilyPolicy
      # ipFamilies
      # labels
      # annotations
```

The `data_layer_http` sidecar:

- Is the official Chia implementation for serving files
- Runs on port 8575 (configurable via `consts.DataLayerHTTPPort`)
- Serves files from the data layer server files directory
- Inherits Chia-specific configuration options from `CommonSpecChia`

### Nginx Sidecar

The `ChiaDataLayer` resource supports an optional nginx sidecar container that can be used to serve static files from the data layer server files directory. This is useful for serving files over HTTP without exposing the data layer HTTP service.

```yaml
apiVersion: chia.network/v1
kind: ChiaDataLayer
metadata:
  name: my-datalayer
spec:
  nginx:
    enabled: true  # Enable the nginx sidecar
    image: "nginx:1.25"  # Optional - defaults to nginx:latest
    service:
      type: ClusterIP  # Optional - defaults to ClusterIP
      # Available service options:
      # externalTrafficPolicy
      # sessionAffinity
      # sessionAffinityConfig
      # ipFamilyPolicy
      # ipFamilies
      # labels
      # annotations
    # Available container configuration options:
    # securityContext
    # livenessProbe
    # readinessProbe
    # startupProbe
    # resources
```

The nginx sidecar:

- Runs on port 8575 (configurable via `consts.NginxPort`)
- Serves files from `/datalayer/server` (mounted from the data layer server files volume)
- Uses a simple nginx configuration that serves static files
- Can be configured with custom container settings like security context, probes, and resource limits
- Creates a service to expose the nginx server (configurable via the `service` field)

The nginx sidecar is particularly useful when you want to:

- Serve static files without exposing the data layer HTTP service
- Use nginx's caching and serving capabilities for better performance
- Apply custom nginx configurations for serving files

### Choosing Between HTTP Servers

When deciding between the `data_layer_http` sidecar and the nginx sidecar, consider:

- Use `data_layer_http` if you want the official Chia implementation and don't need additional HTTP server features
- Use nginx if you need:
  - Better performance through nginx's caching
  - Custom nginx configurations
  - Additional HTTP server features
  - Separation of concerns between data layer and HTTP serving

Note: You can enable both sidecars simultaneously, but they will both try to use port 8575. In this case, you should configure one of them to use a different port.

## More Info

This page contains documentation specific to this resource. Please see the rest of the documentation for information on more available configurations.

-[Generic options for all chia-operator resources.](all.md)
-[chia-exporter configuration](chia-exporter.md)
-[chia-healthcheck configuration](chia-healthcheck.md)
-[Services and networking](services-networking.md)
-[Storage](storage.md)
