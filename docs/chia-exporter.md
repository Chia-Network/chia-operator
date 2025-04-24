# Chia Exporter

[Chia Exporter](https://github.com/Chia-Network/chia-exporter) is an optional component to all Chia resources installed by this operator. It will run as a sidecar container in the Pod to Chia components.

## Enable/Disable

The chia-exporter sidecar will be enabled by default. But you can explicitly enable or disable it with the following:

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

## Supplemental Configuration

There are some niche configuration options for chia-exporter that the majority of people will not need. It is recommended to leave these alone unless you know what you're doing.

In any Chia custom resource yaml file you can set the following:

```yaml
spec:
  chiaExporter:
    configSecretName: chia-exporter-config
```

Where `chia-exporter-config` is the name of a Kubernetes Secret in the same namespace as the supplied chia resource. An example Secret definition would look like the following:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: chia-exporter-config
stringData:
  CHIA_EXPORTER_MAXMIND_COUNTRY_DB_PATH: "/chia-data/country.db"
  CHIA_EXPORTER_MAXMIND_ASN_DB_PATH: "/chia-data/asn.db"
  CHIA_EXPORTER_MYSQL_HOST: "10.0.0.10"
  CHIA_EXPORTER_MYSQL_PASSWORD: "mypassword"
  CHIA_EXPORTER_MYSQL_DB_NAME: "mydbname"
```

This just sets a few non-default options in the environment variables of a chia-exporter sidecar container.

## Specify the version of chia-exporter

Operator releases tend to pin to the current latest version of chia-exporter (at the time the release was published) but if you'd like to manage the version of chia-exporter yourself, there's a field to do so:

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaNode
metadata:
  name: my-seeder
spec:
  chiaExporter:
    enabled: true
    image: ghcr.io/chia-network/chia-exporter:0.14.3
```

The example shows a ChiaNode (full_node) resource on 0.2.1 of chia, but this field is also available on other resources.
