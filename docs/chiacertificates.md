# ChiaCertificates

This resource intakes a Chia CA (certificate authority) Secret that contains at minimum a private CA certificate and key, and generates a new Secret that contains all public and private certificate-key pairs for Chia services.

Because this resource requires a pre-existing CA Secret, it is common to use this in conjunction with a ChiaCA, or a manually created CA Secret.

NOTE: This resource currently only outputs Chia service certificate-key pairs, there is no API to mount these certificates in chia-operator spawned resources.

Example usage:

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaCertificates
metadata:
  name: my-certificates
spec:
  caSecretName: my-ca # required: name of the chia CA Secret to use
  secret: my-certificate-secret # optional: name of the Secret to create (defaults to "chiacertificates")
```

If applied, this example will create a Secret with all chia cert-key pairs named `my-certificate-secret` from a private certificate authority in a Secret in the same namespace named `my-ca`. 

## More Info

This page contains documentation specific to this resource. Please see the [Chia CA](chiaca.md) documentation for information on generating a CA Secret.
