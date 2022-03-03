################################################################################
###                             Project Settings                             ###
################################################################################
PROJECT_NAME := kava-bridge
PROJECT_DIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

PKGS ?= ./...

################################################################################
###                                   Help                                   ###
################################################################################
help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

################################################################################
###                                 Targets                                  ###
################################################################################
.PHONY: install
install: ## Install kava-bridge
	go install -mod=readonly ./cmd/kava-bridged

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
	go vet $(PKGS)

.PHONY: build
build: ## Run go build
	go build $(PKGS)

.PHONY: test
test: ## Run go test
	go test $(PKGS)

.PHONY: cover
cover: ## Run tests with coverage and save to coverage.html
	go test -coverprofile=c.out $(PKGS)
	go tool cover -html=c.out -o coverage.html

.PHONY: watch
watch: ## Run tests on file changes
	while sleep 0.5; do find . -type f -name '*.go' | entr -d go test $(PKGS); done

.PHONY: clean
clean: ## Clean up build and temporary files
	rm c.out coverage.html

################################################################################
###                                 Includes                                 ###
################################################################################
include protobuf.mk
