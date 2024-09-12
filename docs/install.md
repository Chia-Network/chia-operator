# Installation

There are two parts to this Operator. The CRDs (ChiaCA, ChiaFarmer, ChiaNode, etc.) and the actual operator manager Deployment and related objects. You can install these components in two methods, either by cloning the repository and generating the manifests yourself with kustomize, or with `kubectl apply` on the generated manifests on all releases, starting with release `0.2.1`.

## Using the release manifests

Install the CRDs:

```bash
kubectl apply -f https://github.com/Chia-Network/chia-operator/releases/latest/download/crd.yaml
```

Install the controller manager:

```bash
kubectl apply -f https://github.com/Chia-Network/chia-operator/releases/latest/download/manager.yaml
```

## Using kustomize

Clone this repository (and change to its directory):

```bash
git clone https://github.com/Chia-Network/chia-operator.git
cd chia-operator
```

You're currently on the main branch which defaults to the latest versions of this operator, chia, and all sidecars (chia-exporter, chia-healthcheck, etc.) You can switch to the latest release tag for a more stable experience:

```bash
git checkout $(git describe --tags `git rev-list --tags --max-count=1`)
```

Install the CRDs:

```bash
make install
```

Deploy the operator:

```bash
make deploy
```
