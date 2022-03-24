################################################################################
###                             Project Settings                             ###
################################################################################
PROJECT_NAME := kava-bridge
PROJECT_DIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

GO ?= go
PKGS ?= ./...

################################################################################
###                                   Help                                   ###
################################################################################
help: ## Display this help message
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

################################################################################
###                                 Targets                                  ###
################################################################################
.PHONY: install
install: ## Install kava-bridge
	$(GO) install -mod=readonly ./cmd/kava-bridged
	$(GO) install -mod=readonly ./cmd/kava-relayer

.PHONY: start
start: install ## Start kava-bridge chain locally
	./contrib/devnet/init-new-chain.sh
	kava-bridged start

.PHONY: lint
lint: ## Run golint
	golint -set_exit_status $(PKGS)

.PHONY: golangci-lint
golangci-lint: ## Run golangci-lint
	golangci-lint run

.PHONY: vet
vet: ## Run go vet
	$(GO) vet $(PKGS)

.PHONY: build
build: ## Run go build
	$(GO) build $(PKGS)

.PHONY: test
test: ## Run go test
	$(GO) test $(PKGS)

.PHONY: test-integration
test-integration: ## Run go integration tests
	$(GO) test -tags integration ./testing

.PHONY: cover
cover: ## Run tests with coverage and save to coverage.html
	$(GO) test -coverprofile=c.out $(PKGS)
	$(GO) tool cover -html=c.out -o coverage.html

.PHONY: watch
watch: ## Run tests on file changes
	while sleep 0.5; do find . -type f -name '*.go' | entr -d $(GO) test $(PKGS); done

.PHONY: watch-integration
watch-integration: ## Run integration tests on file changes
	while sleep 0.5; do find . -type f -name '*.go' | entr -d $(GO) test -tags integration ./testing; done

.PHONY: clean
clean: ## Clean up build and temporary files
	rm c.out coverage.html

.PHONY: install-devtools
install-devtools: ## Install solc and abigen used by compile-contracts
	cd contract && npm install
	$(GO) install github.com/ethereum/go-ethereum/cmd/abigen@latest

JQ ?= jq
NPM ?= npm
SOLC ?= npx solc
ABIGEN ?= abigen

.PHONY: compile-contracts
compile-contracts: contract/ethermint_json/ERC20MintableBurnable.json relayer/bridge.go ## Compiles contracts and creates ethermint compatible json

contract/artifacts/contracts/ERC20MintableBurnable.sol/ERC20MintableBurnable.json: contract/contracts/ERC20MintableBurnable.sol
	cd contract && $(NPM) run compile

# Ethermint has their own json format for a compiled contract. The following
# converts the abi field to a stringified array, renames bytecode field name to
# bin with the leading `0x` trimmed.
contract/ethermint_json/ERC20MintableBurnable.json: contract/artifacts/contracts/ERC20MintableBurnable.sol/ERC20MintableBurnable.json
	mkdir -p contract/ethermint_json
	$(JQ) '.abi = (.abi | tostring) | {abi, bin: .bytecode[2:] }' < $< > $@

contract/artifacts/contracts_Bridge_sol_Bridge.abi: contract/contracts/Bridge.sol
	cd contract && $(SOLC) --abi contracts/Bridge.sol --base-path . --include-path node_modules/ -o artifacts

contract/artifacts/contracts_Bridge_sol_Bridge.bin: contract/contracts/Bridge.sol
	cd contract && $(SOLC) --optimize --bin contracts/Bridge.sol --base-path . --include-path node_modules/ -o artifacts

relayer/bridge.go: contract/artifacts/contracts_Bridge_sol_Bridge.bin contract/artifacts/contracts_Bridge_sol_Bridge.abi
	$(ABIGEN) --bin $< --abi $(word 2,$^) --pkg relayer --type Bridge --out $@

################################################################################
###                                 Includes                                 ###
################################################################################
include protobuf.mk
