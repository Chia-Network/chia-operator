# ChiaDataLayer

ChiaDataLayers run the data_layer Chia component, which comes bundled with a Chia wallet. In a future time, the wallet may be able to be run separately, but it is not currently possible.

The data_layer_http server runs as an optional sidecar. In a future release, it may be possible to run the HTTP server separately from the data_layer server, but it is not currently implemented.

Here's a minimal ChiaDataLayer example custom resource (CR):

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaDataLayer
metadata:
  name: my-datalayer
spec:
  chia:
    caSecretName: chiaca-secret # A kubernetes Secret containing certificate authority files
```