A simple directory scoped version manager for commands. `xvm` lets you set different paths for the same command in different directories.

## Install

```shell
go install github.com/ajinasokan/xvm@latest

# add xvm overrides dir to path

export PATH=~/.xvm:$PATH
```

## Get started

```shell
# Enable xvm for a command

$ xvm enable flutter /Users/myname/flutter/bin/flutter

symlink /Users/myname/go/bin/xvm->flutter created at /Users/myname/.xvm/flutter

$ flutter --version

Flutter 2.0.6 • channel stable

# Change it for a directory

$ cd myproject

$ xvm set flutter /Users/myname/flutterlatest/bin/flutter

saved config .xvm.conf

$ flutter --version

Flutter 2.2.0 • channel stable
```

## How it works

- When enabled for a command xvm writes a symlink inside `~/.xvm` that simply points to `xvm` binary. Since `~/.xvm` has precedence in the PATH from now on the command will always execute `xvm`.
- If the current directory has `.xvm.conf` with an override for the above enabled command then executing the command will launch the override.
- If current directory doesn't have `.xvm.conf` but one of its parent up the tree have it then `xvm` will use that override.
- If there are no overrides in the directory tree then `xvm` will launch the default path given at the time of enabling command.

## All commands

```
$ xvm enable <command> <default path>

Enable xvm for a command. If no local configuration overrides the command, xvm will use default path for execution.

$ xvm disable <command>

Removes xvm for command.

$ xvm get <command>

Checks for the command overrides in the directory tree. Lists the path that will be used.

$ xvm set <command> <path>

Overrides the path of the command in the current directory. Writes a .xvm.conf file if it doesn't exist.

$ xvm unset <command>

Removes the command from the overrides in the current directory.
```