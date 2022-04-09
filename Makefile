mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
repo_name := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

.PHONY: build
build:
	echo "Building reminder..." && \
	go build -v -o ./bin/ ./cmd/${repo_name} && \
	echo "done."

.PHONY: run
run:
	go run ./cmd/${repo_name}

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt:
	. ./scripts/go_fmt

.PHONY: test
test:
	. ./scripts/go_test
