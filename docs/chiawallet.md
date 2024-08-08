# ChiaWallet

Specifying a ChiaWallet will create a kubernetes Deployment and some Services for a Chia wallet that optionally connects to a local [full_node](chianode.md). It also requires a specified [Chia certificate authority](chiaca.md).

It is also expected you have a pre-existing Chia key to import, likely one that you generated locally in a Chia GUI installation.

Here's a minimal ChiaWallet example custom resource (CR):

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaWallet
metadata:
  name: my-wallet
spec:
  chia:
    caSecretName: chiaca-secret # A kubernetes Secret containing certificate authority files
    # A kubernetes Secret named chiakey-secret containing a key.txt file with your mnemonic key
    secretKey:
      name: "chiakey-secret"
      key: "key.txt"
```

## Secret key

The `secretKey` field in the ChiaWallet's spec defines the name of a Kubernetes Secret that contains your mnemonic. Only Wallets and Farmers need your mnemonic key to function. You can create your Kubernetes Secret like so:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: chiakey-secret
stringData:
  key.txt: your mnemonic goes here
type: Opaque
```

Replace the text value for `key.txt` with your mnemonic, and then reference it in your ChiaWallet resource in the way shown above.

## More Info

This page contains documentation specific to this resource. Please see the rest of the doucmentation for information on more available configurations.

* [Generic options for all chia-operator resources.](all.md)
* [chia-exporter configuration](chia-exporter.md)
* [Services and networking](services-networking.md)
* [Storage](storage.md)
