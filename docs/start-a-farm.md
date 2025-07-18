# Start a farm

This guide provides a complete walkthrough for setting up a Chia farm using the chia-operator. It covers all the essential components needed for farming Chia.

This guide installs everything in the default namespace, but you can of course install them in any namespace. These are also all fairly minimal examples with just enough config to be helpful. Other options are supported.

## Table of Contents

- [SSL CA](#ssl-ca)
- [Full Node](#full_node)
- [Farmer](#farmer)
- [Harvester](#harvester)
- [Wallet](#wallet)

## SSL CA

First thing you'll need is a CA Secret. Chia components all communicate with each other over TLS with signed certificates all using the same certificate authority. This presents a problem in k8s, because each chia-docker container will try to generate their own CAs if none are declared, and all your components will refuse to communicate with each other. This operator contains a ChiaCA CRD that will generate a new CA and set it as a kubernetes Secret for you, or you can make your own Secret with a pre-existing ssl/ca directory. In this guide, we'll show the ChiaCA method first, and then the pre-made Secret second.

Create a file named `ca.yaml`:

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaCA
metadata:
  name: mainnet-ca
spec:
  secret: mainnet-ca
```

The `spec.secret` key specifies the name of the k8s Secret that will be created. The Secret will be created in the same namespace that the ChiaCA CR was created in. Apply this with `kubectl apply -f ca.yaml`

You can also specify a CA Secret without using the ChiaCA custom resource helper. [See the chiaca documentation.](chiaca.md#manually-create-a-ca-secret)

## full_node

Next we need a full_node. Create a file named `node.yaml`:

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaNode
metadata:
  name: mainnet
spec:
  replicas: 1
  chia:
    caSecretName: mainnet-ca
    timezone: "UTC"
  storage:
    chiaRoot:
      hostPathVolume:
        path: "/home/user/.chia/mainnetk8s"
  nodeSelector:
    kubernetes.io/hostname: "node-with-hostpath"
```

As you can see, we used a hostPath volume for CHIA_ROOT. We also specified a nodeSelector for the full_node pod that will be brought up in this Statefulset, this is because a hostPath won't move over to other nodes in the cluster if it gets scheduled elsewhere. You can configure a persistent volume claim instead. ChiaNode objects create StatefulSets which allow each pod to generate their very own PersistentVolumeClaim with this as your storage config:

```yaml
storage:
  chiaRoot:
    persistentVolumeClaim:
      storageClass: ""
      resourceRequest: "300Gi"
```

Finally, apply your ChiaNode with: `kubectl apply -f node.yaml`

## farmer

Now we can create a farmer that talks to our full_node. Create a file named `farmer.yaml`:

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaFarmer
metadata:
  name: mainnet
spec:
  chia:
    caSecretName: mainnet-ca
    timezone: "UTC"
    fullNodePeer: "mainnet-node.default.svc.cluster.local:8444"
    secretKey:
      name: "chiakey"
      key: "key.txt"
```

A couple of things going on here. First, we configured the fullNodePeer address, which we'll use kubernetes internal DNS names for services, targeting port 8444 on the `mainnet-node` service, in the `default` namespace, using the default cluster domain `cluster.local`. If your cluster uses a non-default domain name, switch it to that. Also switch `default` to whatever namespace your ChiaNode is deployed to.

We also have a `secretKey` in the chia config spec. That defines a k8s Secret in the same namespace as this ChiaFarmer, named `chiakey` which contains one data key `key.txt` which contains your Chia mnemonic.

Finally, apply this ChiaFarmer with `kubectl apply -f farmer.yaml`

## harvester

Now we can create a harvester that talks to our farmer. Create a file named `harvester.yaml`:

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaHarvester
metadata:
  name: mainnet
spec:
  chia:
    caSecretName: mainnet-ca
    timezone: "UTC"
    farmerAddress: "mainnet-farmer.default.svc.cluster.local"
  storage:
    plots:
      hostPathVolume:
        - path: "/mnt/plot1"
        - path: "/mnt/plot2"
  nodeSelector:
    kubernetes.io/hostname: "node-with-hostpaths"
```

The config here is very similar to the other components we already made, but we're specifying the farmerAddress, which tells the harvester where to look for the farmer. The farmer port is inferred. And in the storage config, we're specifying two plot directories that are mounted to a particular host. And we're pinning this harvester pod to that node using a nodeSelector with a label that exists on that particular node.

## wallet

Now we can create a wallet that talks to our full_node. Create a file named `wallet.yaml`:

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaWallet
metadata:
  name: mainnet
spec:
  chia:
    caSecretName: mainnet-ca
    timezone: "UTC"
    fullNodePeers:
      - host: "mainnet-node.default.svc.cluster.local"
        port: 8444
    secretKey:
      name: "chiakey"
      key: "key.txt"
```

The config here is very similar to the farmer we already made since it also requires your mnemonic key and a full_node peer.

Finally, apply this ChiaWallet with `kubectl apply -f wallet.yaml`

## Next Steps

After deploying all components, you should have a complete Chia farm running in Kubernetes with:

- A Certificate Authority for secure communications
- A full node syncing with the Chia blockchain
- A farmer connected to your full node
- A harvester managing your plots
- A wallet for managing your XCH

### Additional Configuration

For more advanced configurations, see:

- [Generic options for all chia resources](all.md)
- [Services and networking](services-networking.md)
- [Storage](storage.md)
- [chia-exporter configuration](chia-exporter.md) for Prometheus metrics
- [chia-healthcheck configuration](chia-healthcheck.md) for health monitoring
