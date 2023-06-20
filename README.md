# tmux-keys

`tmux-keys` enables configuring single and chord (a sequence of keys)
tmux keybindings declaratively in a `YAML` file.

## Dependencies

- [`go`](https://go.dev/)
- [`tmux`](https://github.com/tmux/tmux/wiki)

## Install

`go install github.com/kelly-lin/tmux-keys`

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
