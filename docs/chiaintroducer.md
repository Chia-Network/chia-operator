# ChiaIntroducer

Specifying a ChiaIntroducer will create a kubernetes Deployment and some Services for a Chia introducer.

The majority of people do not need to run an introducer. Introducers in Chia serve the purpose of introducing full_nodes in a network to other full_node peers on that network.

Here's a ChiaIntroducer example custom resource (CR):

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaIntroducer
metadata:
  name: my-introducer
spec:
```

## Certificate Authority

If you have your own Certificate Authority to pass to initialize chia from:

```yaml
spec:
  chia:
    caSecretName: chiaca-secret
```

[See the chiaca documentation](chiaca.md#manually-create-a-ca-secret) for information on creating a certificate authority Secret for chia.

## More Info

This page contains documentation specific to this resource. Please see the rest of the documentation for information on more available configurations.

* [Generic options for all chia-operator resources.](all.md)
* [chia-exporter configuration](chia-exporter.md)
* [Services and networking](services-networking.md)
* [Storage](storage.md)
