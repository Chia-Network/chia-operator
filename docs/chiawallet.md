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

## Chia configuration

Some of Chia's configuration can be changed from within the CR.

```yaml
spec:
  chia:
    testnet: true # Switches to the default testnet in the Chia configuration file.
    timezone: "UTC" # Switches the tzdata timezone in the container.
    logLevel: "INFO" # Sets the Chia log level.
```

You can also give a local full_node as a peer for your wallet. This does not currently automate the process of setting the full_node for trusted sync, but could still speed up the wallet syncing process because a full_node running in your cluster is probably geographically closer than other full_nodes around the world. Set a full_node peer like so:

```yaml
spec:
  chia:
    fullNodePeer: "node.default.svc.cluster.local:8444" # A local full_node using kubernetes DNS names
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
