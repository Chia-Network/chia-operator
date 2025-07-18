# Troubleshooting

This troubleshooting guide lists potential issues by operator resource type.

## ChiaWallet

### I changed my mnemonic Secret, how can I make my ChiaWallet use the new mnemonic?

You can start a rollout of the wallet's Deployment resource, and when the new wallet Pod starts, it will recognize the new mnemonic. This also applies to other resources that use mnemonic keys like ChiaDataLayer. In a future Chia Operator release, this operator may watch for changes to the mnemonic Secret and trigger the rollout for you.

Trigger a rollout of the Deployment like so (change the Deployment namespace and name to match your ChiaWallet's Deployment):

```bash
kubectl rollout restart -n ${namespace} deployment/${name}
```
