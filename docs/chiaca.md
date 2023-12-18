# ChiaCA

Chia components require a common certificate authority to talk to each other securely. It is also a hard requirement in some situations such as between a harvester and a farmer.

The ChiaCA custom resource (CR) was created out of convenience to generate a certificate authority for you and put it in a kubernetes Secret.

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaCA
metadata:
  name: my-ca
spec:
  secret: my-ca
```

This will create a kubernetes Secret in the same namespace that this CR is applied named `my-ca`. If you have your own pre-existing CA that you would like to continue using instead, you can also [create a kubernetes Secret manually, documented in this section of the readme.](https://github.com/Chia-Network/chia-operator/blob/main/README.md#ssl-ca).

You can then supply this CA Secret to other Chia custom resources like so:

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaFarmer
metadata:
  name: my-farmer
spec:
  chia:
    caSecretName: my-ca
```
