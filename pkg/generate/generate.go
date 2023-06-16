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
			tableNameWithKey := concatTableNameWithKey(tableName, key)
			cmd := createTableSwitchCmd(tableNameWithKey)
			if idx == len(keys)-1 {
				cmd = binding.Cmd
			}

			result = append(result, createBindCmd(tableName, key, cmd))

			tableName = tableNameWithKey
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

// Tmux command used to switch to the target key-table.
func createTableSwitchCmd(tableName string) string {
	return "switch-client -T" + tableName
}

func concatTableNameWithKey(tableName, key string) string {
	return tableName + "_" + key
}

// Tmux command to set keybinds to a key-table.
func createBindCmd(tableName, key, cmd string) string {
	return fmt.Sprintf("tmux bind-key -T%s %s %s", tableName, key, cmd)
}
