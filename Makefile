.PHONY: install
install: ## Install kava-bridge
	go install -mod=readonly ./cmd/kava-bridged

.PHONY: start
start: install
	./contrib/devnet/init-new-chain.sh
	kava-bridged start

.PHONY: lint
lint: ## Run golint
	golint -set_exit_status ./...

.PHONY: golangci-lint
golangci-lint: ## Run golangci-lint
	golangci-lint run

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: build
build: ## Run go build
	go build ./...

.PHONY: test
test: ## Run go test
	go test ./...

.PHONY: cover
cover: ## Run tests with coverage and save to coverage.html
	go test -coverprofile=c.out ./...
	go tool cover -html=c.out -o coverage.html

.PHONY: watch
watch: ## Run tests on file changes
	while sleep 0.5; do find . -type f -name '*.go' | entr -d go test ./...; done

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
