.PHONY: all
all: simple simple-manual generated generated-manual

.PHONY: simple
simple:
	@echo "[simple]"
	go run $(shell goversion ldflags --version v1.0.1) ./simple

.PHONY: simple-manual
simple-manual:
	@echo "[simple-manual]"
	go run -ldflags ' \
		-X "github.com/erwinvaneyk/go-version.version=v1.0.0" \
		-X "github.com/erwinvaneyk/go-version.gitCommit=$(shell git rev-parse HEAD)" \
		-X "github.com/erwinvaneyk/go-version.buildDate=$(shell date)"' \
	    ./simple

.PHONY: generated
generated:
	@echo "[generated]"
	go run $(shell goversion ldflags --pkg main --version v1.0.2) ./generated

.PHONY: generated-manual
generated-manual:
	@echo "[generated-manual]"
	go run -ldflags ' \
    		-X "main.version=v1.0.3" \
    		-X "main.gitCommit=$(shell git rev-parse HEAD)" \
    		-X "main.buildDate=$(shell date)"' \
    		./generated
