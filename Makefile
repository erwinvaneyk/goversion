.DEFAULT_GOAL := help

VERSION := "v0.1.0-SNAPSHOT"

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build goversion.
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
	go test ./...

.PHONY: verify
verify: ## Run all code analysis tools and linters.
	# Check if codebase is formatted.
	@echo "gofmt -l ."
	@bash -c "[ -z $$(gofmt -l .) ] && echo 'OK' || (echo 'ERROR: files are not formatted:' && gofmt -l . && false)"
	# Run static checks on codebase.
	go vet ./...

.PHONY: format
format: ## Run all formatters on the codebase.
	go fmt ./...

.PHONY: release
release: clean generate verify test ## Build and release goversion, publishing the artifacts on Github and Dockerhub.
	goreleaser release --rm-dist --skip-publish