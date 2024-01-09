.DEFAULT_GOAL := help

BIN_NAME := golang-minify
BIN_DIR := ./bin
X_BIN_DIR := $(BIN_DIR)/goxz
VERSION := $$(make -s app-version)

GOBIN ?= $(shell go env GOPATH)/bin

.PHONY: build
build: ## make binary
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BIN_NAME) main.go

.PHONY: x-build
x-build: $(GOBIN)/goxz
	goxz -d $(X_BIN_DIR) -n $(BIN_NAME) .

.PHONY: upload-binary
upload-binary: $(GOBIN)/ghr ## x-build
	ghr "v$(VERSION)" $(X_BIN_DIR)

$(GOBIN)/goxz:
	@go install github.com/Songmu/goxz/cmd/goxz@latest

$(GOBIN)/ghr:
	@go install github.com/tcnksm/ghr@latest


.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
