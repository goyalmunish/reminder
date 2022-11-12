<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Developer Guide](#developer-guide)
    - [How to run in development mode?](#how-to-run-in-development-mode)
    - [Run Tests](#run-tests)
    - [Format Files](#format-files)
    - [Linting Code](#linting-code)
    - [Build and Push the Docker Image](#build-and-push-the-docker-image)
    - [Integrating with Google Calendar API](#integrating-with-google-calendar-api)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Developer Guide

## How to run in development mode?

Here's how you can run in development mode:

```sh
cd reminder/

# way 1
go run ./cmd/reminder

# way 2
cd cmd/reminder
go run .

# way 3
make -s run
```

## Run Tests

You can make use of [`go_test`](./scripts/go_test) to run test suite:

```sh
cd reminder/

# run tests while supressing printing to console
. ./scripts/go_test

# run tests without supressing printing to console
CONSOLE_PRINT=true . ./scripts/go_test

# or, using make (versbose)
CONSOLE_PRINT=true make test
# or, using make (don't echo commands, and print statements)
make -s test
```

## Format Files

You can make use of [`go_fmt`](./scripts/go_fmt) to auto-format all `.go` files:

```sh
cd reminder/

. ./scripts/go_fmt

# or using make
make -s fmt
```

## Linting Code

```sh
make -s lint
```

## Build and Push the Docker Image

_Make use of [`build_image.sh`](./scripts/build_image.sh) to build and push (requires admin rights) [Docker image](https://hub.docker.com/r/goyalmunish/reminder/tags):_

```sh
# cd into repo
cd reminder/

# example, setting version
VERSION=v1.0.0

# building images and pushing them
. ./scripts/build_image.sh ${VERSION}
```

## Integrating with Google Calendar API

Refer:

- [Go Quickstart](https://developers.google.com/calendar/api/quickstart/go)
