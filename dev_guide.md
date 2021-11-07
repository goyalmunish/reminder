# Reminder Developer Guide

## Install Go:

```
brew install go@1.17

# add to your .zshrc or .bashrc
export PATH="/usr/local/opt/go@1.17/bin:$PATH"

# pin go version
brew pin go@1.17
```

~Get latest `.pkg` file from https://golang.org/doc/install (for mac) and install it by double clicking.~

Now, relaunch the terminal.

## Install Libraries

```sh
GO111MODULE=on go get golang.org/x/tools/gopls@latest

go get github.com/manifoldco/promptui

go get github.com/go-delve/delve/cmd/dlv
sudo /usr/sbin/DevToolsSecurity -enable
```

## Setup NVim

Note: I have not got this working so far, completely.

Don't run `coa ml_with_p374`.

```nvim
:GoInstallBinaries
```

## How to run in development mode?

```sh
cd ~/MG/cst
cd programs/reminder/

# way 1
go run cmd/reminder/main.go

# way 2
cd cmd/reminder
go run .

# way 3
. ./scripts/go_run
```

## Run Tests

```sh
cd ~/MG/cst
cd programs/reminder/

# run tests while supressing printing to console
. ./scripts/go_test
# run tests without supressing printing to console
CONSOLE_PRINT=true . ./scripts/go_test
```

## Format Files

```sh
cd ~/MG/cst
cd programs/reminder/

. ./scripts/go_fmt
```

## Debugging Program

```sh
cd ~/MG/cst
cd programs/reminder/

cd cmd/reminder

~/bin/dlv debug

> break main.main
> continue
> <next>, <list>, <continue>
> call <function>
> print <expression>
```
