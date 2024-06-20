# Chia Exporter

[Chia Exporter](https://github.com/Chia-Network/chia-exporter) is an optional component to all Chia resources installed by this operator. It will run as a sidecar container in the Pod to Chia components.

## Enable/Disable

If any options in the `spec.chiaExporter` section of the configuration is specified, the chia-exporter sidecar will be enabled by default. But you can explicitly enable or disable it with the following:

```yaml
spec:
  chiaExporter:
    enabled: true/false
```

## Add labels/annotations

You may want to add some labels to your chia-exporter Service that get added as labels to your Prometheus metrics.

```yaml
spec:
  chiaExporter:
    service:
      labels:
        network: mainnet
        component: full_node
```

You can do the same thing with annotations.

```yaml
spec:
  chiaExporter:
    service:
      annotations:
        hello: world
```