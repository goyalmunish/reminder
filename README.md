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
            - [Install the tool (optional)](#install-the-tool-optional)
            - [Run the tool](#run-the-tool)
    - [Features/Issues to be worked upon](#featuresissues-to-be-worked-upon)
    - [Contributing towards development](#contributing-towards-development)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Reminder

## Yet Another Reminder Tool/App. Why?

_This is a simple **terminal-based fully-interactive reminder tool**. It is not for everybody, but for the folks with a software engineering background who need the **speed** of the command-line and have to manage many day-to-day to-do items (like meeting minutes, reminders about ad-hoc tasks, etc.) in an organized way._

_Apart from being **fully-interactive** and **terminal-based**, the other major features it comes with are:_

- Easy to categorize **tasks** (referred as **"notes"**) with **tags** (🏷 ).
- List and **manage tasks** (change status, update text, add comments, update due-date, etc) of a given tag.
- Each task can be associated with **multiple tags**, and so show up under all of its tags.
- A given task:
    - can be **updated** (📝) with its text, and also can be enhanced with time-stamped **comments** (💬); so that you can track how and when the progress happened
    - can be marked **done** (✅) or **pending** (⏰); marking it as "done" makes it disappear (soft-delete)
    - can be associated with **due-date** (📅); tasks with upcoming deadlines automatically show up under the **"Approaching Due Date"** option under **Main Menu**
    - can be set as "main" or non-main (incidental); tasks marked as "main", show up under dedicated view **Main Notes**
- **Full-text search** (🔎) among all tasks.
- **Tag-groups** for grouping tags, for managing priority-levels (⬆️ ⬇️) or workflow-stages. For example, a task (note) can be part of only one tag out of tags (for example, `priority-low`, `priority-medium`, and `priority-high` ) part of same tag-group.
- Provides you with **"Register Basic Tags"** functionality to seed basic tags which have special meaning to the workflow.
- All of your **data** (📋) remains with **only you**; so, any of your sensitive information burried inside any of your tasks, doesn't leave your machine.
- The **data** remains in a human-readable and usable format. This is useful when you require to edit your file manually.
- Allows your to **Look Ahead** whole year in advance.
- Easily take **time-stamped backups** (💾).
- Provides you a way to easily add/remove tags to any of the existing tasks.

Nothing is hidden (except your data)! The tool is Open Source. You are welcome to use, recommend features, raise bugs, and enhance it further.

## How to Use?

The [Screencast of Basic Features](./assets/videos/screencast_basic_features.mov) can provide you with gist of how the tool looks like and its basic functionality (but there is a lot more that you can do with it).

Once you invoke the tool (for example, by using the [alias **`reminder`**)](#easily-run-the-tool-via-docker-recommended), you are presented with its **Main Menu**. Use **Up-Arrow** and **Down-Arrow** keys to navigate up and down:

<p align="center">
  <img src="./assets/images/screen_home_list_stuff.png" width="100%">
</p>

You may like to take a good look at the options of the **Main Menu** (as shown above), as we'll be talking about it through out this guide.

Also note that you can choose an option by pressing **Enter** key, and use **Ctrl-c** to jump from any nested level to **one level up** (towards the the **Main Menu**).

In [`reminder`](https://github.com/goyalmunish/reminder), the **tags** are the main method of categorizing tasks. When you first time start the app, the tag list is empty, but you can use the **"Register Basic Tags"** option (the option pointed out by the selection-cursor in the above figure) to register basic tags (as listed in the figure below). Then, you can use the **"List Stuff"** option to list them out (as shown below figure).

As we'll see, the **"List Stuff"** option is the most frequently used option in this list. It lets you add tags, add tasks (also referred to as "notes") under those tags, update those tasks; so almost 90% of use-cases.

<p align="center">
  <img src="./assets/images/screen_basic_tags_02.png" width="100%">
</p>

Now, from within the **"List Stuff"** option, you can add a new tag using **"Add Tag"** (as shown at the bottom of the above figure) or choose an existing tag to add a **task** to it. For example, the following figure shows state of the UI when you select a tag (such as **"priority-urgent"**) to add a new task under it:

<p align="center">
  <img src="./assets/images/screen_add_note_01.png" width="100%">
</p>

On selecting a tag (navigating to the tag and hitting **Enter** key), all of its tasks show up as a list of selectable items. You can then **navigate to a given task** and hit **Enter** key to bring up a **menu to update the task** (it lets you change its text, add comments, mark it as pending, mark it as done, add due-date, change its existing tag(s)). The following figures shows you how this menu looks like:

Note: The **"Approaching Due Date"** shows you tasks that require your immediate attention. In general, tasks with a **due-date** in upcoming `7` days start showing up under this option (and remain there until they are marked done). The tags **"repeat-monthly"** and **"repeat-annually"** are special; tasks tagged with them also show up under the **"Approaching Due Date"** option close to their due-dates in their respective monthly and annual frequencies. These rules are also listed under **"Approaching Due Date"** option for a reference.

<p align="center">
  <img src="./assets/images/screen_home_approaching_due_date.png" width="100%">
</p>

With time, you will add more tags and hundreds of tasks under them. These **stats** will show up **on top of your main-menu screen** (as shown in previous image):

Here, the status states that there are currently 21 tags, a total of 164 tasks, and out of them 80 tasks are in the **"pending"** state. The tasks marked as **"done"** disappear (but not deleted, and will still show up under **"Done Notes"** and in Search results).

The **"Search Notes"** option lets you perform a **full-text search** (with each task's status, text and its comments) through all tasks. You can use `[done]` as search text to filter only tasks which are done, similarly use `[pending]` for tasks which are pending.

<p align="center">
  <img src="./assets/images/screen_home_search.png" width="100%">
</p>

The **result list** updates as you add or delete characters in the **search field** (without hitting Enter-key):

<p align="center">
  <img src="./assets/images/screen_search_list.png" width="100%">
</p>

You can navigate to a search entry (a task) and hit **Enter** key to bring up the **menu to update the task** (similar to how we updated tasks under a tag).

Additionally, from the **Main Menu**:

- use the **"Exit"** option to exit the tool. You can come back it to later from where you left off (that is, with your data intact)
- use the **"Create Backup"** option to create manual time-stamped backup of your data file (on host machine)

## How to Run?

### Easily run the tool via Docker (recommended)

_This is the easiest way to get going, if you have [Docker](https://docs.docker.com/get-docker/) installed. Just download the [`reminder` image](https://hub.docker.com/r/goyalmunish/reminder/tags) by issuing the following commands:_

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

_Run the tool:_

```sh
# run the tool
reminder
```

### Non-Docker Setup

#### Install `go`

On Mac, you can just install it with `brew` as:

```sh
brew install go@1.18
```

For other platforms, check [official `go` download and install guide](https://golang.org/dl/).

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

## Features/Issues to be worked upon

Check [**Issues**](https://github.com/goyalmunish/reminder/issues) to track bugs and request for new features.

## Contributing towards development

Check [Development Guide](./dev_guide.md).
