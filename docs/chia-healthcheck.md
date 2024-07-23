# Chia Healthcheck

[Chia Healthcheck](https://github.com/Chia-Network/chia-healthcheck) is an optional component to certain Chia resources installed by this operator. It will run as a sidecar container in the Pod to Chia components if enabled.

Supported components:

- ChiaNodes
- ChiaSeeders

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

* `/full_node` for ChiaNodes
* `/seeder` for ChiaSeeders

## DNS Hostnames

To configure the ChiaSeeder healthcheck, you need to specify a configuration option in your custom resource file:

```yaml
spec:
  chiaHealthcheck:
    enabled: true
    dnsHostname: "my-seeder.chia.domain.url"
```

This configures the hostname used to check for DNS responses, and corresponds to the `--dns-hostname` flag from chia-healthcheck. The seeder healthcheck will be disabled if not provided.
