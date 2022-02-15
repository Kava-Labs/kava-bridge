.PHONY: install
install:
	go install

.PHONY: lint
lint:
	golint -set_exit_status ./...

.PHONY: golangci-lint
golangci-lint:
	golangci-lint run

.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test ./...

.PHONY: cover
cover:
	go test -coverprofile=c.out ./...
	go tool cover -html=c.out -o coverage.html

.PHONY: watch
watch:
	while sleep 0.5; do find . -type f -name '*.go' | entr -d go test ./...; done
