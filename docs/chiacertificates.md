# ChiaCertificates

NOTE: This resource currently only outputs Chia service certificate-key pairs. But there is no API yet to mount these inside chia-operator spawned Pods.

This resource intakes a Chia CA (certificate authority) Secret that contains at minimum a private CA certificate and key, and generates a new Secret that contains all public and private certificate-key pairs for Chia services.

Because this resource requires a pre-existing CA Secret, it is common to use this in conjunction with a ChiaCA, or a manually created CA Secret.

Example usage:

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaCertificates
metadata:
  name: my-certificates
spec:
  secret: my-ca
```

If applied, this example will create a Secret with all chia cert-key pairs named `my-certificates` from a private certificate authority in a Secret in the same namespace named `my-ca`.
