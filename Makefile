PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin
$(shell [ -f bin ] || mkdir -p $(PROJECT_BIN))
PATH := $(PROJECT_BIN):$(PATH)

GOLANGCI_LINT = $(PROJECT_BIN)/golangci-lint

default: help

.PHONY: help
help: ## Show help
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@egrep '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

.PHONY: .install-linter
.install-linter: ## Download and copy binary of golangci-lint to ./bin/ dir
	[ -f $(PROJECT_BIN)/golangci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v1.49.0

.PHONY: lint
lint: .install-linter  ## Run golangci-lint
	$(GOLANGCI_LINT) run ./... --config=./.golangci.yml

.PHONY: swag
swag:  ## Generate swagger documentation using swag
	go run github.com/swaggo/swag/cmd/swag init -g internal/usecase/controller/http/v1/router.go

generate: swag  ## Generate mocks and swagger documentation using go generate and swag
	go generate ./...
