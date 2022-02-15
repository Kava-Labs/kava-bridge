# Kava Bridge

Repository for Kava Ethereum Bridge

### Directories

- `(root)` - Relayer module
- `contracts` - Bridge Ethereum Contract

## Development

```
> make help

build                         Run go build
cover                         Run tests with coverage and save to coverage.html
golangci-lint                 Run golangci-lint
help                          Display this help message
install                       Install kava-bridge
test                          Run go test
watch                         Run tests on file changes
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

