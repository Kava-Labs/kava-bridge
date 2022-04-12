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

`make start-geth` automatically deploys the following smart contracts. These 
addresses are for **local testing only**. They are not mainnet addresses and are
subject to change.

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

[`ERC20MintableBurnable.sol`]: ./contract/contracts/ERC20MintableBurnable.sol
