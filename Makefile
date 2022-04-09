mkfile_name := $(lastword $(MAKEFILE_LIST))
mkfile_path := $(abspath $(mkfile_name))
repo_path := $(realpath $(dir $(mkfile_path)))
repo_name := $(notdir ${repo_path})

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
