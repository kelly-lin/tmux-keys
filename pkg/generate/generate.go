package generate

import "fmt"

const PREFIX_TABLE_NAME = "prefix"

type Binding struct {
	TableName string
	Keys      string
	Cmd       string
}

// Generates the tmux commands to set keybinds described by the keytable config.
func Generate(bindings []Binding) ([]string, error) {
	result := []string{}

	for _, binding := range bindings {
		keys := splitKeys(binding.Keys)
		tableName := binding.TableName
		for idx, key := range keys {
			// If we are on the last key we want to bind to the command, otherwise bind
			// to the key-table.
			cmd := "switch-client -T" + concatTableNameKey(tableName, key)
			if idx == len(keys)-1 {
				cmd = binding.Cmd
			}

			result = append(result, createBindCmd(tableName, key, cmd))

			tableName = concatTableNameKey(tableName, key)
		}
	}
	return result, nil
}

// Split keys into separate strings based on tmux binding syntax.
func splitKeys(keys string) []string {
	result := []string{}
	for _, key := range keys {
		result = append(result, string(key))
	}
	return result
}

func createBindCmd(tableName, key, cmd string) string {
	return fmt.Sprintf("tmux bind-key -T%s %s %s", tableName, key, cmd)
}

func concatTableNameKey(tableName, key string) string {
	return tableName + "_" + key
}
