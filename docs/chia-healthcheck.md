# Chia Healthcheck

[Chia Healthcheck](https://github.com/Chia-Network/chia-healthcheck) is an optional component to certain Chia resources installed by this operator. It will run as a sidecar container in the Pod to Chia components if enabled.

Supported components:

- ChiaNodes
- ChiaSeeders
- ChiaTimelords

## Enable

You can enable the healthcheck sidecar with the following:

```yaml
spec:
  chiaHealthcheck:
    enabled: true
```

## Liveness/Readiness/Startup Probes

When enabled, chia-healthcheck will be configured as the liveness, readiness, and startup probes for the main chia container by default. You can optionally specify your own probes which will take precedence over the automatic ones:

```yaml
spec:
  chia:
    livenessProbe:
      httpGet:
        path: /full_node
        port: 9950
        scheme: HTTP
    readinessProbe:
      httpGet:
        path: /full_node
        port: 9950
        scheme: HTTP
    startupProbe:
      httpGet:
        path: /full_node
        port: 9950
        scheme: HTTP
```

The path specified in your probe config will change depending on the component you're installing:

-`/full_node` for ChiaNodes
-`/seeder` for ChiaSeeders
-`/timelord` for ChiaSeeders

## DNS Hostnames

To configure the ChiaSeeder healthcheck, you need to specify a configuration option in your custom resource file:

```yaml
spec:
  chiaHealthcheck:
    enabled: true
    dnsHostname: "my-seeder.chia.domain.url"
```

This configures the hostname used to check for DNS responses, and corresponds to the `--dns-hostname` flag from chia-healthcheck. The seeder healthcheck will be disabled if not provided.

## Roll Healthcheck Service ports into the Peer Service

Often times there is a need to expose the peer Service publicly (so peers outside your network can connect to you.) In some deployments, users may use load balancers of some variety to expose the Service publicly, which will be fronted by a public IP address. Similarly, it may be desired in some deployments to expose the chia-healthcheck Service publicly, for example if you wish to use an external health monitoring tool like Uptime Robot or Uptime Kuma. Rather than keep the peer and healthcheck Services separately in this scenario, you may want to create one Service that exposes your peer ports and healthcheck port on the same Service.

For that use case, there is an option that works for the chia-healthcheck Service's configuration to "roll up" the chia-healthcheck Service port into the main peer Service's ports. To do that, set the following in your custom resource that supports chia-healthcheck:

```yaml
spec:
  chiaHealthcheck:
    enabled: true
    service:
      enabled: true
      rollIntoPeerService: true
```

NOTE: If you had custom labels/annotations for your healthcheck Service, you should add them to the Peer Service configuration instead.

## Specify the version of chia-healthcheck

Operator releases tend to pin to the current latest version of chia-healthcheck (at the time the release was published) but if you'd like to manage the version of chia-healthcheck yourself, there's a field to do so:

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaSeeder
metadata:
  name: my-seeder
spec:
  chiaHealthcheck:
    enabled: true
    image: ghcr.io/chia-network/chia-healthcheck:0.2.1
```

The example shows a ChiaSeeder (seeder) resource on 0.2.1 of chia, but this field is also available on other resources that support chia-healthcheck.
