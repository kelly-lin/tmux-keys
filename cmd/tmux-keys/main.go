package main

import (
	"os/exec"
	"path"

	"github.com/kelly-lin/tmux-keys/pkg/generate"
	"gopkg.in/yaml.v3"

	"fmt"
	"os"
)

const (
	BIND_CMD = "bind"
	HELP_CMD = "help"
)

type Config struct {
	KeyBinds []generate.Table `yaml:"key_binds"`
}

func main() {
	if len(os.Args) == 1 {
		print_usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case BIND_CMD:
		configFilePath := path.Join(os.Getenv("HOME"), ".config/tmux-keys/tmux-keys.yml")
		if len(os.Args) > 2 {
			configFilePath = os.Args[2]
		}

		contents, err := os.ReadFile(configFilePath)
		if err != nil {
			fmt.Println("could not read config file")
			os.Exit(1)
		}

		var config Config
		if err := yaml.Unmarshal(contents, &config); err != nil {
			fmt.Printf("could not parse config: %s\n", err)
			os.Exit(1)
		}

		cmds, err := generate.Generate(config.KeyBinds)
		if err != nil {
			fmt.Printf("could not generate keybinds: %s", err)
			os.Exit(1)
		}

		for _, cmd := range cmds {
			err := exec.Command("/bin/sh", "-c", cmd).Run()
			if err != nil {
				fmt.Printf("error while executing command: %s\n", err)
			}
		}

	case HELP_CMD:
		print_usage()

	default:
		fmt.Printf("tmux-keys: command not found: %s\nsee \"tmux-keys help\" for usage\n", cmd)
		os.Exit(1)
	}
}

func print_usage() {
	fmt.Println(`tmux-keys

Usage: 
        tmux-keys <command> [arguments]

Commands:
        bind [config]     set keybindings declared in 'config' if provided, otherwise it will 
                          default to $HOME/.config/tmux-keys/tmux-keys.yml`)
}
