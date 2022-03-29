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
Installs `kava-bridge` to `$GOPATH/bin`

### Test

```
make test
```
Runs all unit tests.

```
make watch
```
Runs all unit tests on file changes

