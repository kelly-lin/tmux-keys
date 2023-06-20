# tmux-keys

`tmux-keys` enables configuring single and chord (a sequence of keys)
tmux keybindings declaratively in a `YAML` file. The configuration below will
open a new window by pressing the keys: `prefix` -> `n` -> `w`.

```yaml
key_binds:
  - table_name: prefix # this binds the keys to the prefix table
    bindings: # begin binding declaration
      - keys: n w # the key sequence: prefix -> n -> w
        cmd: new-window # open a new window (this can be any tmux command)
```

## Dependencies

- [`go`](https://go.dev/)
- [`tmux`](https://github.com/tmux/tmux/wiki)

## Install

`go install github.com/kelly-lin/tmux-keys/cmd/tmux-keys`

## Usage

1. Configure `tmux-keys.yml` (see [configuration](#configuration)) and place it
in the default config directory `$HOME/.config/tmux-keys/tmux-keys.yml`
or any other directory that you choose.
2. Generate and source the keybindings in `tmux.conf` by adding in the following
command:

```shell
run-shell 'command -v tmux-keys &>/dev/null && tmux-keys bind'
```

if you have located your `tmux-keys.yml` somewhere else, then provide the filepath
as an argument to the bind command:

```shell
run-shell 'command -v tmux-keys &>/dev/null && tmux-keys bind path/to/tmux-keys.yml'.
```

## Configuration

You configure your keybindings by declaring them in a `tmux-keys.yml` file,
see below example configuration.

```yaml
key_binds:
  - table_name: prefix # this binds the keys to the prefix table
    bindings: # begin binding declaration
      - keys: n w # the key sequence: prefix -> n -> w
        cmd: new-window # open a new window (this can be any tmux command)
      - keys: s l
        cmd: switch-client -l
      - keys: z
        cmd: switch-client -Tsome_cool_table # if you want to link a keybind to another table (declared below)
                                             # this is how you would do it.
  - table_name: some_cool_table # this is accessed by the binding above for prefix -> z
    bindings:
      - keys: s v # this means that this binding can be accessed by prefix -> z -> s -> v
        cmd: split-window
```

### Configuration Structure

The root key in the configuration file is `key_binds` which contains a list of
`table`.

Each `table` has a name `table_name` and a list of `binding` stored in
the `bindings` key. The `table_name` is the `tmux` `key-table` (see `man tmux`).
Most of the time, you will be wanting to be binding to the default `tmux`
`key-table`'s, `prefix` and `root`. If you want to bind to a custom table, you
can declare them in the config.

Each `binding` has the fields `keys` and `cmd`. `keys` takes a space separated
string which each component describes each key to be pressed to activate the
binding and `cmd` is any `tmux` command.
