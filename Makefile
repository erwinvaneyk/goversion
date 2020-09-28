.DEFAULT_GOAL := help

VERSION := $(shell git describe --abbrev=0 2>/dev/null || echo "v0.0.0-SNAPSHOT")
PATH  := $(PWD)/bin:$(PATH)
SHELL := env PATH=$(PATH) /bin/bash

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build goversion.
	@shell which goversion &>/dev/null || go build -o ./bin/goversion ./cmd/goversion
	go build $(shell goversion ldflags --version ${VERSION}) -o ./bin/goversion ./cmd/goversion

.PHONY: clean
clean: ## Clean up all build and release-related resources.
	rm -rf ./bin
	rm -rf ./dist

.PHONY: install
install: build ## Build and install goversion.
	cp -f ./bin/goversion ${GOBIN}

.PHONY: generate
generate: ## Run all code generators
	go generate ./...

.PHONY: test
test: ## Run all unit tests.
	go test -v ./...

.PHONY: verify
verify: ## Run all static analysis checks.
	# Check if codebase is formatted.
	@which goimports > /dev/null || ! echo 'goimports not found'
	@bash -c "[ -z \"$(goimports -l cmd pkg)\" ] && echo 'OK' || (echo 'ERROR: files are not formatted:' && goimports -l cmd pkg && echo -e \"\nRun 'make format' or manually fix the formatting issues.\n\" && false)"
	# Run static checks on codebase.
	go vet ./cmd/... ./pkg/...

.PHONY: format
format: ## Run all formatters on the codebase.
	# Format the Go codebase.
	goimports -w cmd pkg

	# Format the go.mod file.
	go mod tidy

.PHONY: release
release: clean generate verify test ## Build and release goversion, publishing the artifacts on Github and Dockerhub.
	GOVERSION_LDFLAGS=$(shell goversion ldflags --print-ldflag=false --version=${VERSION}) goreleaser release --rm-dist --skip-publish --snapshot --debug