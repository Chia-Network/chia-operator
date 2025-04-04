# Quickstart

## Install

Install the latest version of chia-operator:

```bash
kubectl apply --server-side -f https://github.com/Chia-Network/chia-operator/releases/latest/download/crd.yaml
kubectl apply -f https://github.com/Chia-Network/chia-operator/releases/latest/download/manager.yaml
```

## Prometheus metrics (Optional)

If you have the Prometheus Operator installed in your cluster and would like to use the bundled ServiceMonitor to scrape chia-operator metrics:

```bash
kubectl apply -f https://github.com/Chia-Network/chia-operator/releases/latest/download/monitor.yaml
```

This ServiceMonitor is installed in the same namespace as chia-operator, so it will only work in Prometheus Operator configurations that can load ServiceMonitors from any namespace, and also don't have any ServiceMonitorSelectors set.

## Install Chia Services

The operator should be running in your cluster now and ready to go! Get to installing some Chia resources. If you're a farmer, see these guides:

* [ChiaCA](chiaca.md)
* [Node](chianode.md)
* [Farmer](chiafarmer.md)
* [Harvester](chiaharvester.md)
* [Wallet](chiawallet.md)

For more information on specific configurations:

* [Generic options for all chia-operator resources](all.md)
* [chia-exporter configuration](chia-exporter.md)
* [chia-healthcheck configuration](chia-healthcheck.md)
* [Services and networking](services-networking.md)
* [Storage](storage.md)
