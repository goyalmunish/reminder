mkfile_name := $(lastword $(MAKEFILE_LIST))
mkfile_path := $(abspath $(mkfile_name))
repo_path := $(realpath $(dir $(mkfile_path)))
repo_name := $(notdir ${repo_path})

.PHONY: gobuild
gobuild:
	go build -v ./...

.PHONY: run
run:
	go run ./...

.PHONY: lint
lint:
	. ./scripts/go_lint

.PHONY: fmt
fmt:
	. ./scripts/go_fmt

.PHONY: test
test:
	. ./scripts/go_test

.PHONY: coverage
coverage:
	. ./scripts/go_coverage

.PHONY: open
open:
	. ./scripts/open_data_file
