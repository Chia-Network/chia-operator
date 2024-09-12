# ChiaSeeder

Specifying a ChiaSeeder will create a kubernetes Deployment and some Services for a Chia seeder.

The majority of people do not need to run a seeder. Seeders in Chia serve the purpose of introducing full_nodes in a network to other full_node peers on that network. See the [seeder documentation](https://docs.chia.net/guides/seeder-user-guide/) for more information.

Seeders have some pre-requisites that you will normally configure outside a kubernetes cluster. This operator doesn't do any of that configuration on your behalf, so in short you will need:

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
    domainName: "seeder.example.com." # name of the NS record for your server with a trailing period. (ex. "seeder.example.com.")
    nameserver: "seeder-mainnet-1.example.com." # name of the A record for your server with a trailing period. (ex. "seeder-us-west-2.example.com.")
    rname: "admin.example.com." # an administrator's email address with '@' replaced with '.' and a trailing period.
```

## Chia configuration

Some of Chia's configuration can be changed from within the CR.

```yaml
spec:
  chia:
    minimumHeight: 240000 # Only consider nodes synced at least to this height
    bootstrapPeer: "mainnet-node.chia.svc.cluster.local" # Peers used for the initial crawler run to find peers
    ttl: 900 # field on DNS records that controls the length of time that a record is considered valid
```

## More Info

This page contains documentation specific to this resource. Please see the rest of the documentation for information on more available configurations.

* [Generic options for all chia-operator resources.](all.md)
* [chia-exporter configuration](chia-exporter.md)
* [chia-healthcheck configuration](chia-healthcheck.md)
* [Services and networking](services-networking.md)
* [Storage](storage.md)
