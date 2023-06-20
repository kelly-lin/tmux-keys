# tmux-keys

`tmux-keys` enables configuring single and chord (a sequence of keys)
tmux keybindings declaratively in a `YAML` file.

## Dependencies

- `Go`

## Install

`go install github.com/kelly-lin/tmux-keys`

## Usage

1. Configure `tmux-keys.yml` (see [configuration](#configuration)) and place it
in the default config directory (`$HOME/.config/tmux-keys/tmux-keys.yml`)
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

You can configure your keybindings by declaring them in a `tmux-keys.yml` file,
see [example config](./docs/tmux-keys.yml).
