<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Reminder](#reminder)
    - [Yet Another Reminder Tool/App. Why?](#yet-another-reminder-toolapp-why)
    - [How to run?](#how-to-run)
        - [Easily run the tool via Docker (recommended)](#easily-run-the-tool-via-docker-recommended)
        - [Non-Docker Setup](#non-docker-setup)
            - [Install `go`](#install-go)
            - [Install the tool](#install-the-tool)
            - [Run the tool](#run-the-tool)
    - [Contributing towards development](#contributing-towards-development)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Reminder

## Yet Another Reminder Tool/App. Why?

_This is a **terminal-based fully-interactive reminder tool**. It is not for everybody, but for the folks who spend most of their time on command-line (terminal) on Linux/Unix/macOS._

_Apart from being **fully-interactive** and **terminal-based**, the other major benefits it comes with are:_

- Easy to categorize notes/tasks with **tags** (🏷 ).
- **Tag-groups** for manage priority levels (⬆️ ⬇️).
- Each note/task can be **tagged with multiple keywords**.
- Notes/tasks can be **updated** (📝) and also can be enhanced with **comments** (💬).
- Notes/tasks can be marked **done** (✅) or **pending** (⏰).
- Notes/tasks can be associated with **due date** (📅). Notes/tasks with upcoming deadlines automatically show up under the **current** tag.
- **Full-text search** (🔎) among **pending** notes/tasks.
- All of your **data** (📋) remains with **only you**.
- The **data** remains in human readable and usable format. This is useful in case you choose to move away.
- Easily take **time-stamped backups** (💾).
- Easy to update tags of existing notes/tasks.
- Nothing is hidden (except your data)! The tool is Open Source. You are welcome to use, recommend features, raise bugs, and enhance it further.

## How to run?

### Easily run the tool via Docker (recommended)

_This is the easiest way to get going, if you have [Docker](https://docs.docker.com/get-docker/) installed. Just issue the following commands:_

```sh
# pull the image
docker pull goyalmunish/reminder

# make sure the directory for data file existing on host machine
mkdir -p ~/reminder

# spin up the container, with data file shared from host machine
docker run -it -v ~/reminder:/root/reminder goyalmunish/reminder
```

_For subsequent runs, better add below alias to `~/.bashrc` ( or `~/.zshrc`, etc), so that you can invoke the tool, just by typing `reminder`:_

```sh
alias reminder='docker run -it -v ~/reminder:/root/reminder goyalmunish/reminder'
```

_Run the tool:_

```sh
reminder
```

### Non-Docker Setup

#### Install `go`

On Mac, you can just install it with `brew` as:

```sh
brew install golang
```

For other platforms, check [official `go` download and install guide](https://golang.org/dl/).

Otherwise, you can also use one of the [Golang Offical Images](https://hub.docker.com/_/golang) to run tool from a Docker container. For example,

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

#### Install the tool

Clone the repo:

```sh
git clone git@github.com:goyalmunish/reminder.git
```

If this results in Permission issues, such as `git@github.com: Permission denied (publickey).`, then either you [Setup Git](https://git-scm.com/book/en/v2/Getting-Started-First-Time-Git-Setup) or just use `git clone https://github.com/goyalmunish/reminder.git` instead.

Install the tool as:

```sh
cd reminder
go install cmd/reminder/main.go
mv ${GOPATH}/bin/main ${GOPATH}/bin/reminder
```

#### Run the tool

If your `go/bin` path is alreay in `PATH`, then you can just run the tool as:

```sh
reminder
```

## Contributing towards development

Check [Development Guide](./dev_guide.md).
