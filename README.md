# hist - Improved shell history

> :warning: **WORK IN PROGRESS**

`hist` aims to be an improvement to the standard history command on Linux
and Mac based systems (currently only zsh is supported).

Command history is stored in an `sqlite3` database (located in the user's
XDG DataDir).

Certain commands can be ignored by defining exclude patterns in the config file
(located in the user's XSD ConfigDir).

An interactive search mode is bound to the `ctrl-c` command. It used a fuzzy
search mode to select a command from your command history.

## Config

By default, config is stored in `XDG_CONFIG_HOME/hist/config.yml`.

Currently only specifying path's to be excluded is supported.

```yaml
---
patterns:
  excludes:
    - ^ls .*
    - ^ll.*
    - ^cd .*
    - ^rm .*
```

This excluded all commands starting with `ls`, `ll`, `cd` and `rm` from being added to history.
When adding new exclude rules, the current config can be tidied up by using `hist tidy`.

## Usage

````shell
Usage: hist [--config CONFIG] <command> [<args>]

Options:
  --config CONFIG [default: ~/.config/hist/config.yml]
  --help, -h             display this help and exit

Commands:
  delete                 Delete commands from history
  get                    Get a command by it's index'
  import                 Import commands from a legacy history file
  list                   List commands
  record                 Record a new command
  search                 Start the interactive fuzzy selection mode
  stats                  Show some statistics
  tidy                   Apply exclude patterns to clean up the hist database
````

### Delete

Delete one or more entries from the database.

```shell
Usage: hist delete [--id ID] [--updated-before UPDATED-BEFORE] [--pattern PATTERN]

Options:
  --id ID, -i ID
  --updated-before UPDATED-BEFORE, -u UPDATED-BEFORE
  --pattern PATTERN, -p PATTERN
                         Delete all records matching the pattern
```

### Get

Get a command from the database by index. This is used when traversing through
history using `arrow up` and `arrow down` keys. 

```shell
Usage: hist get [--index INDEX]

Options:
  --index INDEX
```

### Import

```shell
Usage: hist import [PATH]

Positional arguments:
  PATH
```
### List

```shell
Usage: hist list [--by-count] [--reverse] [--no-count] [--no-last-update] [--with-id] [--limit LIMIT]

Options:
  --by-count
  --reverse
  --no-count
  --no-last-update
  --with-id
  --limit LIMIT, -l LIMIT [default: -1]
```
### Record

```shell
Usage: hist record [COMMAND]

Positional arguments:
  COMMAND
```

### Search

```shell
Usage: hist search
```

### Stats

```shell
Usage: hist stats
```

### Tidy

```shell
Usage: hist tidy
```

## Building and installing

Currently only building the binary yourself is supported. 

```shell
git clone https://github.com/herzrasen/hist.git
cd hist
make && sudo make install
```

Also make sure that you source the shell script in your
`.zshrc` 
