package main

import (
	"github.com/kelly-lin/tmux-keys/pkg/generate"
	"gopkg.in/yaml.v3"

	"fmt"
	"os"
)

const (
	GENERATE_CMD = "generate"
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
	case GENERATE_CMD:
		if len(os.Args) < 3 {
			fmt.Println("please provide config file")
			os.Exit(1)
		}

		configFilePath := os.Args[2]
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
			fmt.Println(cmd)
		}

	default:
		fmt.Printf("tmux-keys: command not found: %s\nsee \"tmux-keys help\" for usage\n", cmd)
		os.Exit(1)
	}
}

func print_usage() {
	println(`tmux-keys

Usage: 
        tmux-keys [flags] [command]

Commands:
        generate        generate keybindings
`)
}
