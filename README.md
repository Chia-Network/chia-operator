# chia-operator

A Kubernetes operator for managing all of your Chia services in your favorite container orchestrator!

Easily run Chia components in Kubernetes by applying simple manifests. A whole farm can be run with each component isolated in its own pod, with a chia-exporter sidecar for remote monitoring, and chia-healthcheck for intelligent healthchecking (for supported chia services.)

## Quickstart

### Install

Install the latest version of chia-operator:

```bash
kubectl apply --server-side -f https://github.com/Chia-Network/chia-operator/releases/latest/download/crd.yaml
kubectl apply -f https://github.com/Chia-Network/chia-operator/releases/latest/download/manager.yaml
```

The operator Deployment will be installed in the `chia-operator-system` namespace.

### Prometheus metrics (Optional)

If you have the Prometheus Operator installed in your cluster and would like to use the bundled ServiceMonitor to scrape chia-operator metrics:

```bash
kubectl apply -f https://github.com/Chia-Network/chia-operator/releases/latest/download/monitor.yaml
```

The ServiceMonitor will be installed in the `chia-operator-system` namespace.

### Install Chia Services

The operator should be running in your cluster now and ready to go! Take a look at the [documentation](docs/README.md) and get to installing some Chia resources. If you're a farmer, see the [Start a Farm](docs/start-a-farm.md) guide, or view these individually:

* [ChiaCA](docs/chiaca.md) (required so your chia services can all talk to each other!)
* [Node](docs/chianode.md)
* [Farmer](docs/chiafarmer.md)
* [Harvester](docs/chiaharvester.md)
* [Wallet](docs/chiawallet.md)

Other Chia services are also available:
* [Crawler](docs/chiacrawler.md)
* [DataLayer](docs/chiadatalayer.md)
* [Introducer](docs/chiaintroducer.md)
* [Seeder](docs/chiaseeder.md)
* [Timelord](docs/chiatimelord.md)

For more information on specific configurations:

* [Generic options for chia resources](docs/all.md)
* [chia-exporter configuration](docs/chia-exporter.md)
* [chia-healthcheck configuration](docs/chia-healthcheck.md)
* [Services and networking](docs/services-networking.md)
* [Storage](docs/storage.md)
