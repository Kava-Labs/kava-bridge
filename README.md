# Kava Bridge

Repository for Kava Ethereum Bridge

### Directories

- `(root)` - Relayer module
- `contracts` - Bridge Ethereum Contract

## Development

```text
> make help

build                          Run go build
clean                          Clean up build and temporary files
compile-contracts              Compiles contracts and creates ethermint compatible json
cover                          Run tests with coverage and save to coverage.html
golangci-lint                  Run golangci-lint
help                           Display this help message
install-devtools               Install solc and abigen used by compile-contracts
install                        Install kava-bridge
lint                           Run golint
proto-all                      Run all protobuf targets
proto-check-breaking           Check for breaking changes against master
proto-format                   Format protobuf files
proto-gen                      Generate protobuf files
proto-lint                     Lint protobuf files
proto-swagger-gen              Generate protobuf swagger
proto-update-deps              Update proto
start-geth                     Start private geth chain locally with the Bridge contract
start                          Start kava-bridge chain locally
test-integration               Run go integration tests
test                           Run go test
vet                            Run go vet
watch-integration              Run integration tests on file changes
watch                          Run tests on file changes
```

### Install

```
make install
```

Installs `kava-bridge` and `kava-relayer` to `$GOPATH/bin`

### Relayer Usage

```bash
kava-relayer network generate-network-secret
kava-relayer network generate-node-key
kava-relayer network show-node-id

# Initialize configuration file. By default this creates the file ~/.kava-relayer/config.yaml
# Add network secret, node key, peer list.
kava-relayer init

# Connect to p2p network
kava-relayer network connect

# Start single signer relayer without P2P network
kava-relayer start 
```

### Test

```
make test
```
Runs all unit tests.

```
make watch
```
Runs all unit tests on file changes

### Development Smart Contracts

#### Geth

`make start-geth` automatically deploys the following smart contracts. These 
addresses are for **local testing only**. They are not mainnet addresses and are
subject to change. These only exist on the local Ethereum node.

* Signer Address (miner): `0x21E360e198Cde35740e88572B59f2CAdE421E6b1`
* Bridge Address `0xb588617416D0B0A3C29618bf8Fb6aC0cAd4Ede7f`
* Bridge Relayer: `0xa2F728F997f62F47D4262a70947F6c36885dF9fa`
* WETH Address: `0x6098c27D41ec6dc280c2200A737D443b0AaA2E8F`
* ERC20 Tokens (mintable by signer, see [`ERC20MintableBurnable.sol`]):
  * `MEOW` Address (18 decimals): `0x8223259205A3E31C54469fCbfc9F7Cf83D515ff6`
  * `USDC` Address (6 decimals): `0x4Fb48E68842bb59f07569c623ACa5826b600F8F7`
* Test User Address: `0x7Bbf300890857b8c241b219C6a489431669b3aFA`
* Multicall Address: `0xB94efB606287D37732Fe871BDdD015c5E7Ab2e76`
* Multicall2 Address: `0x77C3c07d77a5E99Fffdc2635CcdB66b16e3a1Bed`

#### Kava EVM

Kava EVM also has some contracts included in the genesis state. Note that the
`WETH` specific contract does **not** exist on Kava EVM. When it is bridged from
Ethereum, it is a [`ERC20MintableBurnable.sol`] deployed by the bridge module.

* Multicall Address: `0xeA7100edA2f805356291B0E55DaD448599a72C6d`
* Multicall2 Address: `0x62d2f38dAA1153b381c6ed2A48e7f4673303ac9A`
* WKAVA Address `0x70C79B608aBBC502c2F61f38E04190fB407BefCF`

These contracts were first deployed manually to Kava EVM, exported as Kava
JSON state, then copied over to the `init-new-chain.sh` script. The following
is **not** necessary to deploy Multicall and WKAVA, only when adding additional
smart contracts to the genesis state.

```bash
# First comment out the contract additions in init-new-chain.sh to prevent
# address collisions, then start the local chain.
make start

cd contract
npx hardhat run scripts/init_kava_evm.ts --network kava

# Stop chain, then export. Copy over necessary values to init-new-chain.sh
# app_state.auth.accounts and app_state.evm.accounts
kava-bridged export | jq > export.json
```

[`ERC20MintableBurnable.sol`]: ./contract/contracts/ERC20MintableBurnable.sol
