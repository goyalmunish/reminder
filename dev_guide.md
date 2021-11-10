<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Developer Guide](#developer-guide)
    - [How to run in development mode?](#how-to-run-in-development-mode)
    - [Run Tests](#run-tests)
    - [Format Files](#format-files)
    - [Debugging Program](#debugging-program)
    - [Build Docker Image](#build-docker-image)
    - [Additional Conventions](#additional-conventions)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Developer Guide

## How to run in development mode?

Here's how you can run in development mode:

```sh
cd reminder/

# way 1
go run cmd/reminder/main.go

# way 2
cd cmd/reminder
go run .

# way 3
. ./scripts/go_run
```

## Run Tests

You can make use of [`go_test`](./scripts/go_test) to run test suite:

```sh
cd reminder/

# run tests while supressing printing to console
. ./scripts/go_test

# run tests without supressing printing to console
CONSOLE_PRINT=true . ./scripts/go_test
```

## Format Files

You can make use of [`go_fmt`](./scripts/go_fmt) to auto-format all `.go` files:

```sh
cd reminder/

. ./scripts/go_fmt
```

## Debugging Program

Here are some hints for debugging:

```sh
cd reminder/
cd cmd/reminder

~/bin/dlv debug

> break main.main
> continue
> <next>, <list>, <continue>
> call <function>
> print <expression>
```

## Build Docker Image

_Make use of [`build_image.sh`](./scripts/build_image.sh) to build and push (requires admin rights) Docker image:_

```sh
cd reminder/
. ./scripts/build_image.sh
```

## Additional Conventions

Here are some additional conventions followed:

- Functions (not methods) are prefixed with `F` or `f`.
