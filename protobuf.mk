ifndef PROJECT_NAME
$(error PROJECT_NAME must be defined)
endif

ifndef PROJECT_DIR
$(error PROJECT_DIR must be defined)
endif

DOCKER ?= docker
DOCKER_FLAGS := -v $(PROJECT_DIR):/workspace --workdir /workspace
DOCKER_BUF := $(DOCKER) run --rm $(DOCKER_FLAGS) bufbuild/buf
RSYNC ?= rsync

protoVer := v0.3
protoImageName := tendermintdev/sdk-proto-gen:$(protoVer)
containerProtoGen := $(PROJECT_NAME)-proto-gen-$(protoVer)
containerProtoGenSwagger := $(PROJECT_NAME)-proto-gen-swagger-$(protoVer)
containerProtoFmt := $(PROJECT_NAME)-proto-fmt-$(protoVer)

define docker_run
	if $(DOCKER) ps -a --format '{{.Names}}' | grep -Eq "^${1}$$"; then $(DOCKER) start -a $(1); else $(DOCKER) run --name $(1) $(DOCKER_FLAGS) $(protoImageName) $(2); fi
endef

define sync_proto_deps
	$(RSYNC) -r --chmod 644 --include "*.proto" --include='*/' --exclude='*' --prune-empty-dirs $(1) $(2)
endef

.PHONY: proto-all
proto-all: proto-gen proto-format proto-lint proto-swagger-gen ## Run all protobuf targets

.PHONY: proto-gen
proto-gen: ## Generate protobuf files
	@echo "Generating Protobuf files"
	@$(call docker_run,$(containerProtoGen),sh ./scripts/protocgen.sh)

.PHONY: proto-swagger-gen
proto-swagger-gen: ## Generate protobuf swagger
	@echo "Generating Protobuf Swagger"
	@$(call docker_run,$(containerProtoGenSwagger),sh ./scripts/protoc-swagger-gen.sh)

.PHONY: proto-format
proto-format: ## Format protobuf files
	@echo "Formatting Protobuf files"
	@$(call docker_run,$(containerProtoFmt),find ./ -not -path "./third_party/*" -name *.proto -exec clang-format -style=file -i {} \;)

.PHONY: proto-lint
proto-lint: ## Lint protobuf files
	@$(DOCKER_BUF) lint --error-format=json --exclude-path ./third_party/proto/cosmos_proto,./third_party/proto/gogoproto

PROTO_CHECK_REF ?= .git\#branch=main
.PHONY: proto-check-breaking
proto-check-breaking: ## Check for breaking changes against master
	@$(DOCKER_BUF) breaking --against $(PROTO_CHECK_REF)

COSMOS_PROTO_PATH := $(shell go list -m -f '{{.Dir}}' github.com/cosmos/cosmos-proto)
GOGO_PROTO_PATH := $(shell go list -m -f '{{.Dir}}' github.com/gogo/protobuf)
.PHONY: proto-update-deps
proto-update-deps: ## Update proto
	mkdir -p third_party/proto
	$(call sync_proto_deps,$(COSMOS_PROTO_PATH)/proto third_party)
	$(call sync_proto_deps,$(GOGO_PROTO_PATH)/gogoproto third_party/proto)
