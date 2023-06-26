package generate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const PREFIX_TABLE_NAME = "prefix"

type Table struct {
	Name     string `yaml:"table_name"`
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
		keys, err := splitKeys(binding.Keys)
		if err != nil {
			return nil, err
		}

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

    tableName = table.Name
	}
	return result, nil
}

// Split keys into separate strings based on tmux binding syntax.
// The following should be single keys: x (case sensitive), C-x, M-x, [digit], symbols and
// quotation marks, F[1-13], special keys (Up, Down, Left, Right, Space) - case
// insensitive.
func splitKeys(keys string) ([]string, error) {
	if keys == "" {
		return []string{}, errors.New("keys is empty")
	}

	result := []string{}

	current := ""
	items := strings.Split(keys, " ")
	for _, item := range items {

		current += item

		if err := validateKey(current); err != nil {
			return []string{}, err
		}

		if current == `'` {
			current = `"'"`
		}

		if current == `"` {
			current = `'"'`
		}

		result = append(result, current)
		current = ""
	}

	return result, nil
}

func validateKey(key string) error {
	isFnKey := func(key string) bool {
		key = strings.ToLower(key)
		return key[0] == 'f' && len(key) > 1
	}

	isValidFnKey := func(key string) bool {
		asInt, err := strconv.Atoi(key[1:])
		if err != nil {
			return false
		}

		if asInt > 12 || asInt == 0 {
			return false
		}

		return true
	}

	isSpecialKey := func(key string) bool {
		key = strings.ToLower(key)
		specialKeys := map[string]bool{
			"up":       true,
			"down":     true,
			"left":     true,
			"right":    true,
			"bspace":   true,
			"btab":     true,
			"dc":       true,
			"delete":   true,
			"end":      true,
			"enter":    true,
			"escape":   true,
			"home":     true,
			"ic":       true,
			"insert":   true,
			"npage":    true,
			"pagedown": true,
			"pgdn":     true,
			"ppage":    true,
			"pageup":   true,
			"pgup":     true,
			"space":    true,
			"tab":      true,
		}
		_, ok := specialKeys[key]
		return ok
	}

	isModifierKey := func(key string) bool {
    if len(key) == 2 {
      if key[0] == '^' {
        return true
      }
    }

		key = strings.ToLower(key)
    if len(key) == 3 {
      if key[1] == '-' && (key[0] == 'c' || key[0] == 'm' || key[0] == 's') {
        return true
      }
    }

		return false
	}

	if isFnKey(key) {
		if !isValidFnKey(key) {
			return fmt.Errorf("%q is not a supported key function key, supported function keys are F1 to F12", key)
		}
		return nil
	}

	if isModifierKey(key) {
		return nil
	}

	if !isSpecialKey(key) && len(key) > 1 {
		return fmt.Errorf("%q is not a supported special key, see tmux manual for reference", key)
	}

	return nil
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
