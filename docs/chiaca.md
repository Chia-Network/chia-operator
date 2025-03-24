# ChiaCA

Chia components require a common certificate authority to talk to each other securely. It is also a hard requirement in some situations such as between a harvester and a farmer.

The ChiaCA custom resource (CR) was created out of convenience to generate a certificate authority for you and put it in a kubernetes Secret.

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaCA
metadata:
  name: my-ca
spec:
  secret: my-ca-secret # optional: name of the Secret to create (defaults to the name of the ChiaCA resource)
```

This will create a kubernetes Secret in the same namespace that this CR is applied named `my-ca-secret`. If you have your own pre-existing CA that you would like to continue using instead, you can also [create a kubernetes Secret manually, documented in this section of the readme](https://github.com/Chia-Network/chia-operator/blob/main/README.md#ssl-ca).

You can then supply this CA Secret to other Chia custom resources like so:

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaFarmer
metadata:
  name: my-farmer
spec:
  chia:
    caSecretName: my-ca-secret
```

## Manually create a CA Secret

The ChiaCA custom resource (CR) exists as an option of convenience, but if you have your own CA you'd like to use instead, you'll need to create a Secret that contains all the files in the `$CHIA_ROOT/config/ssl/ca` directory, like so:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-ca
stringData:
  chia_ca.crt: |
    <redacted file output>
  chia_ca.key: |
    <redacted file output>
  private_ca.crt: |
    <redacted file output>
  private_ca.key: |
    <redacted file output>
type: Opaque
```

You only need to do this if you don't want to use the ChiaCA CR to make it for you.
