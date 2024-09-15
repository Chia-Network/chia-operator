# ChiaNetwork

ChiaNetwork resources contain the configuration for a blockchain network (network overrides, network name, network port, etc.) When you deploy a ChiaNetwork resource, you can use it in the spec for your other Chia-deploying resources (ChiaNode, etc.) to contain all of your network-related configuration in one place.

NOTE: This API uses a new tool for setting complex values in your configuration file. Being a new tool, it may be appropriate to back up your chia configuration file before attaching a ChiaNetwork resource to any of your other resources.

The most common use-case for ChiaNetwork resources will be for standing up non-default networks. Most people farming mainnet or the default testnet probably won't need to use this.

Here's an example ChiaNetwork custom resource (CR):

```yaml
apiVersion: k8s.chia.net/v1
kind: ChiaNetwork
metadata:
  name: testnetz
spec:
  constants:
    MIN_PLOT_SIZE: 18
    GENESIS_CHALLENGE: "ccd5bb71183532bff220ba46c268991a3ff07eb358e8255a65c30a2dce0e5fbb"
    GENESIS_PRE_FARM_POOL_PUZZLE_HASH: "d23da14695a188ae5708dd152263c4db883eb27edeb936178d4d988b8f3ce5fc"
    GENESIS_PRE_FARM_FARMER_PUZZLE_HASH: "3d8765d3a597ec1d99663f6c9816d915b9f68613ac94009884c4addaefcce6af"
  config:
    address_prefix: txch
    default_full_node_port: 58444
  networkName: testnetz
  networkPort: 58444
  introducerAddress: intro.testnetz.example.com
  dnsIntroducerAddress: dnsintro.testnetz.example.com
```

- `networkName` is the name of the network. The name of the ChiaNetwork resource will be used if this is left unspecified.
- `networkPort` is the full_node port to use on this network.
- `introducerAddress` is the address to an introducer on this network.
- `dnsIntroducerAddress` is the address to a DNS introducer (seeder) on this network.
- `constants` are the network constants to be defined in the chia config for this network underneath network_overrides.
- `config` is the config to be defined in the chia config for this network underneath network_overrides.

## Usage

On a Chia-deploying resource (ChiaNode, ChiaFarmer, etc.) you can specify a ChiaNetwork resource to use for configuration like so:

```yaml
spec:
  chia:
    chiaNetwork: "testnetz"
```

testnetz is the name of the ChiaNetwork deployed in the same kubernetes namespace.

## Precedence

Several of these configuration options are also available on Chia-deploying resources (ChiaNode, ChiaFarmer, etc.) If specified on the ChiaNetwork, the ChiaNetwork resource's fields will take precedence.
