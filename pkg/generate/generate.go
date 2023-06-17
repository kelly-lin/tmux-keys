package generate

import "fmt"

const PREFIX_TABLE_NAME = "prefix"

type Table struct {
	Name     string
	Bindings []Binding
}

type Binding struct {
	Keys string
	Cmd  string
}

func Generate(tables []Table) ([]string, error) {
	result := []string{}

	for _, table := range tables {
		cmds, err := createTableBindingCmds(table)
		if err != nil {
			return nil, err
		}
		result = append(result, cmds...)
	}

	return result, nil
}

// Generates the tmux commands to set keybinds described by the keytable config.
func createTableBindingCmds(table Table) ([]string, error) {
	result := []string{}

	tableName := table.Name
	for _, binding := range table.Bindings {
		keys := splitKeys(binding.Keys)
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
