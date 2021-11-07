<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Reminder](#reminder)
    - [How to run?](#how-to-run)
        - [Running via Docker](#running-via-docker)
        - [Non-Docker Setup](#non-docker-setup)
            - [Install `go`](#install-go)
            - [Install `reminder` command](#install-reminder-command)
            - [Run `reminder`](#run-reminder)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Reminder

## How to run?

### Running via Docker

This is the easiest way to get going, if you have [Docker](https://docs.docker.com/get-docker/) installed. Just issue the following commands:

```sh
# pull the image
docker pull goyalmunish/reminder

# make sure the directory for data file existing on host machine
mkdir -p ~/reminder

# spin up the container, with data file shared from host machine
docker run -it -v ~/reminder:/root/reminder goyalmunish/reminder
```

### Non-Docker Setup

#### Install `go`

On Mac, you can just install it with `brew` as:

```sh
brew install golang
```

For other platforms, check [official `go` download and install guide](https://golang.org/dl/).

Otherwise, you can also use one of the [Golang Offical Images](https://hub.docker.com/_/golang) to run command from a Docker container. For example,

```sh
GOLANG_IMAGE=golang:1.17.2-alpine3.14
GOLANG_VERSION=1.17

# run the image
docker pull ${GOLANG_IMAGE}
docker run -it -d --privileged --name golang${GOLANG_VERSION} ${GOLANG_IMAGE}

# exec into the container
docker exec -it golang${GOLANG_VERSION} /bin/sh
```

If `git` and `ssh` are not available (for instance case of fresh `alpine` image, from above), install them as:

```sh
apk add git
apk add openssh
```

Check installed version:

```sh
go version
```

#### Install `reminder` command

Clone the repo:

```sh
git clone git@github.com:goyalmunish/reminder.git
```

If this results in Permission issues, such as `git@github.com: Permission denied (publickey).`, then either you [Setup Git](https://git-scm.com/book/en/v2/Getting-Started-First-Time-Git-Setup) or just use `git clone https://github.com/goyalmunish/reminder.git` instead.

Install the command as:

```sh
cd reminder
go install cmd/reminder/main.go
cp ${GOPATH}/bin/main ~/reminder
```

#### Run `reminder`

If your `go/bin` path is alreay in `PATH`, then you can just run the command as:

```sh
reminder
```

You may like to check the installation path:

```sh
which reminder
```
