# Chia Healthcheck

[Chia Healthcheck](https://github.com/Chia-Network/chia-healthcheck) is an optional component to certain Chia resources installed by this operator that can be used as a startup, liveness, and readiness probe. Chia-Healthcheck provides more intelligent healthchecking logic to ensure your chia services are healthy.

Supported components:

- ChiaNodes
- ChiaSeeders
- ChiaTimelords

## Enable

The chia-healthcheck sidecar will be enabled by default for all services that support it. But you can explicitly enable or disable it with the following:

```yaml
spec:
  chiaHealthcheck:
    enabled: false
```

NOTE: ChiaSeeders require an additional parameter for chia-healthcheck:

```yaml
spec:
  chiaHealthcheck:
    enabled: true
    dnsHostname: seeder.example.com
```

The `dnsHostname` setting is only required for seeders. If you enable chia-healthcheck on a ChiaSeeder, but omit this setting, the operator will override the value of `enabled` by disabling chia-healthcheck.

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

Operator releases tend to pin to the current latest version of chia-healthcheck (at the time the release was published.) If you would like to manage the version of chia-healthcheck yourself, you can specify the version of the image to use:

```yaml
spec:
  chiaHealthcheck:
    enabled: true
    image: ghcr.io/chia-network/chia-healthcheck:0.2.1
```
