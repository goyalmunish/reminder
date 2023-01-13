## Other Ways of Running

### Easily run the tool via Docker

_This is the easiest way to get going, if you have [Docker](https://docs.docker.com/get-docker/) installed. Just download the [`reminder` image](https://hub.docker.com/r/goyalmunish/reminder/tags) by issuing the following commands:_

_**Using Script:**_

Make sure first to clone the repo and `cd` into it.

```sh
# pull latest reminder image, make sure ~/reminder directory exists, and run the tool
. ./scripts/run_via_docker.sh

# run the tool (just run, without pulling image and other initialization steps)
. ./scripts/run_via_docker.sh fast
```

_**Directly using `docker` command:**_

```sh
# pull the image (or get the latest image)
docker pull goyalmunish/reminder

# make sure the directory for the data file exists on the host machine
mkdir -p ~/reminder

# spin up the container, with data file shared from the host machine
docker run -it --rm --name reminder -v ~/reminder:/root/reminder goyalmunish/reminder
```

_For subsequent runs, better add the below alias to `~/.bashrc` ( or `~/.zshrc`, etc), so that you can invoke the tool, just by typing `reminder` (or any other alias that you prefer):_

```sh
# define the alias
alias reminder='docker run -it --rm --name reminder -v ~/reminder:/root/reminder goyalmunish/reminder'
```

_Then, run the tool using `reminder` command._

### Non-Docker Setup

Check for available installers on [**releases**](https://github.com/goyalmunish/reminder/releases) page. Otherwise,

#### Install `go`

On Mac, you can just install it with `brew` as:

```sh
brew install go@1.18
```

_For other platforms, check [official `go` download and install guide](https://go.dev/dl/)._

Check installed version:

```sh
go version
```

#### Install the tool (optional)

Clone the repo:

```sh
git clone git@github.com:goyalmunish/reminder.git
```

If this results in Permission issues, such as `git@github.com: Permission denied (publickey).`, then either you [Setup Git](https://git-scm.com/book/en/v2/Getting-Started-First-Time-Git-Setup) or just use `git clone https://github.com/goyalmunish/reminder.git` instead.

Install the tool as:

```sh
# cd into the local copy of the repo
cd reminder

# install the tool
go install ./cmd/reminder

# move the binary to /usr/local/bin/
mv ${GOPATH}/bin/reminder /usr/local/bin/reminder
```

#### Run the tool

If you have installed the tool, and your `go/bin` path is alreay in `PATH`, then you can just run it as:

```sh
reminder
```

Otherwise, you can just run it as (without installing, directly from clone of the repo):

```sh
# cd into the local copy of the repo
cd reminder

# running the tool using `make`
make run

# or as
go run ./cmd/reminder
```
