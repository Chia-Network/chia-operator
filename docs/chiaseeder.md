# ChiaSeeder

Specifying a ChiaSeeder will create a kubernetes Deployment and some Services for a Chia seeder. It also requires a specified [Chia certificate authority](chiaca.md).

The majority of people do not need to run a seeder. Seeders in Chia serve the purpose of introducing full_nodes in a network to other full_node peers on that network. See the [seeder documentation](https://docs.chia.net/guides/seeder-user-guide/) for more information.

Seeders have some pre-requisites that you will normally configure outside of a kubernetes cluster. This operator doesn't do any of that configuration on your behalf, so in short you will need:

* A DNS `A` record that points to your server's IP address. In this instance the A record will probably be your public IP address if you intend on the DNS server to be reachable publicly, or an internal address if you're reserving the seeder's DNS server for your use.
* A DNS `AAAA` record is not strictly needed, but is often preferred if your network is IPv6 enabled.
* A DNS `NS` record that points to your `A`/`AAAA` record(s).
* Networking fixtures between the public internet and your seeder server. This may be a NodePort Service that points to your ChiaSeeder kubernetes Pod. And port forwards on your firewall for port 53 to your NodePort Service. Seeder servers respond to queries on both TCP and UDP, but other full_nodes will only make contact using the UDP protocol.

ChiaSeeder Deployments add the `NET_BIND_SERVICE` linux capability to bind to privileged ports, as is typical of any DNS server ran on linux. See the [linux man pages](https://man7.org/linux/man-pages/man7/capabilities.7.html) for more information.

Here's a ChiaSeeder example custom resource (CR):

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaSeeder
metadata:
  name: my-seeder
spec:
  chia:
    caSecretName: chiaca-secret # A kubernetes Secret containing certificate authority files
    domainName: "seeder.example.com." # name of the NS record for your server with a trailing period. (ex. "seeder.example.com.")
    nameserver: "seeder-mainnet-1.example.com." # name of the A record for your server with a trailing period. (ex. "seeder-us-west-2.example.com.")
    rname: "admin.example.com." # an administrator's email address with '@' replaced with '.' and a trailing period.
```

## Chia configuration

Some of Chia's configuration can be changed from within the CR.

```yaml
spec:
  chia:
    testnet: true # Switches to the default testnet in the Chia configuration file.
    timezone: "UTC" # Switches the tzdata timezone in the container.
    logLevel: "INFO" # Sets the Chia log level.
    
    # Seeder config settings
    minimumHeight: 240000 # Only consider nodes synced at least to this height
    bootstrapPeer: "mainnet-node.chia.svc.cluster.local" # Peers used for the initial crawler run to find peers
    ttl: 900 # field on DNS records that controls the length of time that a record is considered valid
```

## CHIA_ROOT storage

`CHIA_ROOT` is an environment variable that tells chia services where to expect a data directory to be for local chia state. You can store your chia state persistently a couple of different ways: either with a host mount or a persistent volume claim.

To use a persistent volume claim, first create one in the same namespace and then give its name in the CR like the following:

```yaml
spec:
  storage:
    chiaRoot:
      persistentVolumeClaim:
        claimName: "chiaroot-data"
```

To use a hostPath volume, first create a directory on the host and specify the path in the CR like the following:

```yaml
spec:
  storage:
    chiaRoot:
      hostPathVolume:
        path: "/home/user/storage/chiaroot"
```

If using a hostPath, you may want to pin the pod to a specific kubernetes node using a NodeSelector:

```yaml
spec:
  nodeSelector:
    kubernetes.io/hostname: "node-with-hostpath"
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
