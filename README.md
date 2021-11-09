<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Reminder](#reminder)
    - [Yet Another Reminder Tool/App. Why?](#yet-another-reminder-toolapp-why)
    - [How to Use?](#how-to-use)
    - [How to Run?](#how-to-run)
        - [Easily run the tool via Docker (recommended)](#easily-run-the-tool-via-docker-recommended)
        - [Non-Docker Setup](#non-docker-setup)
            - [Install `go`](#install-go)
            - [Install the tool](#install-the-tool)
            - [Run the tool](#run-the-tool)
    - [Features/Issues to be worked upon](#featuresissues-to-be-worked-upon)
    - [Contributing towards development](#contributing-towards-development)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Reminder

## Yet Another Reminder Tool/App. Why?

_This is a **terminal-based fully-interactive reminder tool**. It is not for everybody, but for the folks who spend most of their time on command-line (terminal) on Linux/Unix/macOS._

_Apart from being **fully-interactive** and **terminal-based**, the other major benefits it comes with are:_

- Easy to categorize notes/tasks with **tags** (🏷 ).
- **Tag-groups** for managing priority levels (⬆️ ⬇️).
- Each note/task can be **tagged with multiple keywords**.
- Notes/tasks can be **updated** (📝) and also can be enhanced with **comments** (💬).
- Notes/tasks can be marked **done** (✅) or **pending** (⏰).
- Notes/tasks can be associated with **due date** (📅). Notes/tasks with upcoming deadlines automatically show up under the **"Pending Notes"** option.
- **Full-text search** (🔎) among **pending** notes/tasks.
- All of your **data** (📋) remains with **only you**.
- The **data** remains in human readable and usable format. This is useful in case you choose to move away.
- Easily take **time-stamped backups** (💾).
- Easy to update tags of existing notes/tasks.

Nothing is hidden (except your data)! The tool is Open Source. You are welcome to use, recommend features, raise bugs, and enhance it further.

## How to Use?

Once you invoke the tool (for example, by using the previously created alias **`reminder`**), you are presented with its **Main Menu**. Use **Up-Arrow** and **Down-Arrow** keys to navigate up and down:

<p align="center">
  <img src="./assets/images/screen_basic_tags_01.png" width="100%">
</p>

Note: You may like to take a good look at the options of the **Main Menu** (as shown above). You can always come back to it by pressing **Ctrl-c**.

Note: Choose an option by pressing **Enter** key, and use **Ctrl-c** to jump from any nested level to the **Main Menu**.

The **tags** are the main method of categorizing notes. In the beginning, the tag list is empty, but you can use the **"Register Basic Tags"** option (the option pointed out with the cursor in above figure) to register basic tags. Then, you can use the **"List Stuff"** option to list them out (as shown below):

Note: As we'll see, the **"List Stuff"** option is the most important option in this list. It lets you add tags, add notes, update notes; so almost 90% of use-cases. The **Main Menu** also lists options such as **"Add Note"** and **"Add Tag"** but you'll rarely have to use them directly.

<p align="center">
  <img src="./assets/images/screen_basic_tags_02.png" width="100%">
</p>

You can add a new tag using **"Add Tag"** or choose an existing tag to add a **note** to it. For example, the following adds a new note to the **"priority-urgent"** tag:

<p align="center">
  <img src="./assets/images/screen_add_note_01.png" width="100%">
</p>

All notes of a selected tag show up under it. From list of notes under a tag, you can **navigate to a given note** and hit **Enter** key to bring up a **menu to update it** (change its text, add comments, mark it as pending, mark it as done, add due date, change its existing tag(s)). This is how this menu looks like:

Note: Notes with a **due date** in upcoming `7` days start showing up under **"Pending Notes"** option (until they are marked done).

<p align="center">
  <img src="./assets/images/screen_search_04.png" width="100%">
</p>

With time, you will add more tags and hundreds of notes under them. This **status** will show up **on top of your main-menu screen** (as shown below):

<p align="center">
  <img src="./assets/images/screen_home.png" width="100%">
</p>

The above status states that there are currently 22 tags, a total of 330 notes, and out of them 204 notes are in the **"pending"** state. The notes marked as **"done"** become **invisible** (but not deleted) and also don't show up in search results.

The **"Search Notes"** option lets you perform a **full-text search** (with each note's text and its comments) through all notes with the **"pending"** state.

<p align="center">
  <img src="./assets/images/screen_search_01.png" width="100%">
</p>

The **result list** updates as you add or delete characters in the **search field** (without hitting Enter-key):

<p align="center">
  <img src="./assets/images/screen_search_03.png" width="100%">
</p>

The figures such as 1618211581 that you see in the above results are timestamps of the comments added to the corresponding notes.

You can navigate to a search entry (a note) and hit **Enter** key to bring up the **menu to update it** (similar to how we updated notes under a tag):

<p align="center">
  <img src="./assets/images/screen_search_04.png" width="100%">
</p>

Additionally, from the **Main Menu**:

- Use the **"Exit"** option to exit the tool. You can come back it to later from where you left off (that is, with your data intact).
- Use the **"Create Backup"** option to create time-stamped backup of your data file (on host machine).
- Use the **"Done Note"** option to display done notes (which are otherwise invisible under the **"List Stuff"** and **"Search Notes"** options).

## How to Run?

### Easily run the tool via Docker (recommended)

_This is the easiest way to get going, if you have [Docker](https://docs.docker.com/get-docker/) installed. Just issue the following commands:_

```sh
# pull the image
docker pull goyalmunish/reminder

# make sure the directory for the data file exists on the host machine
mkdir -p ~/reminder

# spin up the container, with data file shared from the host machine
docker run -it -v ~/reminder:/root/reminder goyalmunish/reminder
```

_For subsequent runs, better add the below alias to `~/.bashrc` ( or `~/.zshrc`, etc), so that you can invoke the tool, just by typing `reminder`:_

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

## Features/Issues to be worked upon

Check [**Issues**](https://github.com/goyalmunish/reminder/issues) to track bugs and request for new features.

## Contributing towards development

Check [Development Guide](./dev_guide.md).
