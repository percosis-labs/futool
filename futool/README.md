# futool

Assorted dev tools for working with the fury blockchain.

To get started with running a local fury network, check out our docs on [Getting Started](https://docs.fury.io/docs/cosmos/getting-started).

## Installation

```bash
make install
```

## Initialization: futool testnet

Note that the most accurate documentation lives in the CLI itself. It's recommended you read through `futool testnet bootstrap --help`.

Option 1:

The `futool testnet bootstrap` command starts a local Fury blockchain as a
background docker container called `generated-furynode-1`. The bootstrap command
only starts the Fury blockchain and Fury REST server services.

```bash
# Start new testnet
futool testnet bootstrap --fury.configTemplate master
```

The endpoints are exposed to localhost:

* RPC: http://localhost:26657
* REST: http://localhost:1317
* GRPC: http://localhost:9090
* GRPC Websocket: http://localhost:9091
* EVM JSON-RPC: http://localhost:8545
* EVM Websocket: http://localhost:8546

Option 2:

To generate a testnet for fury, binance chain, and a deputy that relays swaps between them:

```bash
# Generate a new futool configuration based off template files
futool testnet gen-config fury binance deputy --fury.configTemplate master

# Pull latest docker images. Docker must be running.
cd ./full_configs/generated && docker-compose pull

# start the testnet
futool testnet up

# When finished with usage, shut down the processes
futool testnet down
```

### Flags

Additional flags can be added when initializing a testnet to add additional
services:

`--ibc`: Run Fury testnet with an additional IBC chain. The IBC chain runs in the container named `ibcnode`. It has primary denom `ufury`.

Example:

```bash
# Run Fury testnet with an additional IBC chain
futool testnet bootstrap --fury.configTemplate master --ibc
```

`--geth`: Run a go-ethereum node alongside the Fury testnet. The geth node is
initialized with the Fury Bridge contract and test ERC20 tokens. The Fury EVM
also includes Multicall contracts deployed. The contract addresses can be found
on the [Percosis-Labs/fury-bridge](https://github.com/Percosis-Labs/fury-bridge#development)
README.

Example:

```bash
# Run the testnet with a geth node in parallel
futool testnet bootstrap --fury.configTemplate master --geth
```

Geth node ports are **not** default, as the Fury EVM will use default JSON-RPC
ports:

Fury EVM RPC Ports:

* HTTP JSON-RPC: `8545`
* WS-RPC port: `8546`

Geth RPC Ports:

* HTTP JSON-RPC: `8555`
* WS-RPC port: `8556`

To connect to the associated Ethereum wallet with Metamask, setup a new network with the following parameters:
* New RPC URL: `http://localhost:8555`
* Chain ID: `88881` (configured from the [genesis](config/templates/geth/initstate/genesis.json#L3))
* Currency Symbol: `ETH`

Finally, connect the mining account by importing the JSON config in [this directory](config/templates/geth/initstate/.geth/keystore)
with [this password](config/templates/geth/initstate/eth-password).

## Automated Chain Upgrade

Futool supports running upgrades on a chain. To do this requires the fury final docker image to have a registered upgrade handler.
The upgrade will start a chain with the docker container tag from `--upgrade-base-image-tag`. Once it reaches height `--upgrade-height`, it halts the chain for an upgrade named `--upgrade-name`. At that point, the container is restated with the desired container: `FURY_TAG` if defined, of if not defined, the default tag for the config template.

**Example**:
Test a chain upgrade from v0.19.2 -> v0.21.0 at height 15.

Using an overridden docker image tag:
```
$ FURY_TAG=v0.21.0 futool testnet bootstrap --upgrade-name v0.21.0 --upgrade-height 15 --upgrade-base-image-tag v0.19.2
```

Using a config template:
```
Test a chain upgrade from v0.19.2 -> v0.21.0:
$ futool testnet bootstrap --fury.configTemplate v0.21 --upgrade-name v0.21.0 --upgrade-height 15 --upgrade-base-image-tag v0.19.2
```

## Usage: futool testnet

REST APIs for both blockchains are exposed on localhost:

- Fury: http://localhost:1317
- Binance Chain: http://localhost:8080

You can also interact with the blockchain using the `fury` command line. In a
new terminal window, set up an alias to `fury` on the dockerized fury node and
use it to send a query.

```bash
# Add an alias to the dockerized fury cli
alias dfury='docker exec -it generated_furynode_1 fury'

# Confirm that the alias has been added
alias fury

# For versions before v0.16.x
alias dfucli='docker exec -it generated_furynode_1 fucli'
```

Note that for some architectures or docker versions, the containers are generated with hyphens (`-`) instead of underscores (`_`).

You can test the set up and alias by executing a sample query:

```bash
dfury status
dfury q cdp params
```

The chain has several accounts that are funded from genesis. A list of the account names can be found [here](config/common/addresses.json).

The binary is pre-configured to have these keys in its keyring so you should be able to use them directly.
```bash
# Example sending funds from `whale` to another account
dfury tx bank send whale [fury-address-to-fund] 1000000ufury --gas-prices 0.001ufury -y

# Check transaction result by tx hash
dfury q tx [tx-hash]
```
### A note about eth accounts

Account keys can be created with two different algorithms in Fury: `secp256k1` and `eth_secp256k1`.
Which algorithm is used is dictate by the presence of the `--eth` flag on key creation.

Eth accounts can be exported for use in ethereum wallets like Metamask. A list of of the pre-funded eth accounts can be found [here](config/generate/genesis/auth.accounts/eth-accounts.json).
Notable, `whale2` is an eth account. These keys can be easily imported into a wallet via their private keys:
```bash
# DANGEROUS EXPORT OF PRIVATE KEY BELOW! BE CAREFUL WITH YOUR PRIVATE KEYS FOR MAINNET ACCOUNTS.
dfury keys unsafe-export-eth-key whale2
```
The above will output the hex-encoded ethereum private key that can be directly imported to Metamask or another EVM-supporting wallet.

You can always import or generate new eth accounts as well:
```bash
# generate new account
dfury keys add new-eth-account --eth

# recover an eth account from a mnemonic
dfury keys add new-eth-account2 --eth --recover
eth flag specified: using coin-type 60 and signing algorithm eth_secp256k1
> Enter your bip39 mnemonic
# enter your mnemonic here

# import an eth account from a hex-encoded ethereum private key
fury keys unsafe-import-eth-key new-eth-account3 [priv-key]
```

### ERC20 token

The master template includes a pre-deployed ERC20 token with the name "USDT". The token is configured to be converted to an sdk coin of the denom `erc20/tether/usdt`.

Token Address: `0xeA7100edA2f805356291B0E55DaD448599a72C6d`
Funded Account: `whale2` - `0x03db6b11F47d074a532b9eb8a98aB7AdA5845087` (1000 USDC)

## Shut down: futool testnet

When you're done make sure to shut down the futool testnet. Always shut down the futool testnets before pulling the latest image from docker, otherwise you may experience errors.

```bash
futool testnet down
```

# Updating fury genesis

When new versions of fury are released, they often involve changes to genesis.
The fury `master` template includes a genesis.json that is generated from a pure state:
* Ensure the desired version of `fury` is in your path as `fury`
* Run `make generate-fury-genesis`
* The script will create a genesis with desired accounts & validator
* Updates to the genesis should be made in [`update-fury-genesis.sh`](./config/generate/genesis/generate-fury-genesis.sh)
