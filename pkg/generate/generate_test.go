package generate

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTableBindingCmds(t *testing.T) {
	for desc, tc := range map[string]struct {
		Tables []Table
		Want   []string
	}{
		"1 layer binding": {
			Tables: []Table{
				{
					Name: PREFIX_TABLE_NAME,
					Bindings: []Binding{
						{Keys: "f", Cmd: "new-window"},
					},
				},
			},
			Want: []string{"tmux bind-key -Tprefix f new-window"},
		},
		"2 layer binding": {
			Tables: []Table{
				{
					Name: PREFIX_TABLE_NAME,
					Bindings: []Binding{
						{Keys: "f d", Cmd: "new-window"},
					},
				},
			},
			Want: []string{
				"tmux bind-key -Tprefix f switch-client -Tprefix_f",
				"tmux bind-key -Tprefix_f d new-window",
			},
		},
		"3 layer binding": {
			Tables: []Table{
				{
					Name: PREFIX_TABLE_NAME,
					Bindings: []Binding{
						{Keys: "f d c", Cmd: "new-window"},
					},
				},
			},
			Want: []string{
				"tmux bind-key -Tprefix f switch-client -Tprefix_f",
				"tmux bind-key -Tprefix_f d switch-client -Tprefix_f_d",
				"tmux bind-key -Tprefix_f_d c new-window",
			},
		},
		"bindings in different tables with the same starting key should have bindings in unique tables": {
			Tables: []Table{
				{
					Name: PREFIX_TABLE_NAME,
					Bindings: []Binding{
						{Keys: "f d", Cmd: "list-sessions"},
					},
				},
				{
					Name: "some-other-table-name",
					Bindings: []Binding{
						{Keys: "f d", Cmd: "new-window"},
					},
				},
			},
			Want: []string{
				"tmux bind-key -Tprefix f switch-client -Tprefix_f",
				"tmux bind-key -Tprefix_f d list-sessions",
				"tmux bind-key -Tsome-other-table-name f switch-client -Tsome-other-table-name_f",
				"tmux bind-key -Tsome-other-table-name_f d new-window",
			},
		},
		"multiple tables": {
			Tables: []Table{
				{
					Name:     PREFIX_TABLE_NAME,
					Bindings: []Binding{{Keys: "f", Cmd: "new-window"}},
				},
				{
					Name: "some-other-table-name",
					Bindings: []Binding{
						{Keys: "f d", Cmd: "new-window"},
					},
				},
			},
			Want: []string{
				"tmux bind-key -Tprefix f new-window",
				"tmux bind-key -Tsome-other-table-name f switch-client -Tsome-other-table-name_f",
				"tmux bind-key -Tsome-other-table-name_f d new-window",
			},
		},
	} {
		t.Run(desc, func(t *testing.T) {
			got := []string{}
			for _, table := range tc.Tables {
				cmds, _ := createTableBindingCmds(table)
				got = append(got, cmds...)
			}

			assert.Equal(t, tc.Want, got)
		})
	}
}

func TestSplitKeys(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		type Want struct {
			Keys []string
		}

		type TC struct {
			Keys string
			Want Want
		}

		for _, tc := range []TC{
			{
				Keys: "a",
				Want: Want{Keys: []string{"a"}},
			},
			{
				Keys: "A",
				Want: Want{Keys: []string{"A"}},
			},
			{
				Keys: "F",
				Want: Want{Keys: []string{"F"}},
			},
			{
				Keys: "a b",
				Want: Want{Keys: []string{"a", "b"}},
			},
			{
				Keys: "C-x",
				Want: Want{Keys: []string{"C-x"}},
			},
			{
				Keys: "M-x",
				Want: Want{Keys: []string{"M-x"}},
			},
			{
				Keys: "S-x",
				Want: Want{Keys: []string{"S-x"}},
			},
			{
				Keys: "Space",
				Want: Want{Keys: []string{"Space"}},
			},
			{
				Keys: "Space Right",
				Want: Want{Keys: []string{"Space", "Right"}},
			},
			{
				Keys: "'",
				Want: Want{Keys: []string{`"'"`}},
			},
			{
				Keys: `"`,
				Want: Want{Keys: []string{`'"'`}},
			},
			{
				Keys: "^C",
				Want: Want{Keys: []string{"^C"}},
			},
			{
				Keys: "F1",
				Want: Want{Keys: []string{"F1"}},
			},
			{
				Keys: "F12",
				Want: Want{Keys: []string{"F12"}},
			},
		} {
			keys, _ := splitKeys(tc.Keys)
			assert.Equal(t, tc.Want.Keys, keys)
		}
	})

	t.Run("error", func(t *testing.T) {
		type Want struct {
			Error error
		}

		type TC struct {
			Keys string
			Want Want
		}

		for _, tc := range []TC{
			{
				Keys: "",
				Want: Want{Error: errors.New("keys is empty")},
			},
			{
				Keys: "F0",
				Want: Want{Error: errors.New(`"F0" is not a supported key function key, supported function keys are F1 to F12`)},
			},
			{
				Keys: "F13",
				Want: Want{Error: errors.New(`"F13" is not a supported key function key, supported function keys are F1 to F12`)},
			},
			{
				Keys: "FF",
				Want: Want{Error: errors.New(`"FF" is not a supported key function key, supported function keys are F1 to F12`)},
			},
			{
				Keys: "FF",
				Want: Want{Error: errors.New(`"FF" is not a supported key function key, supported function keys are F1 to F12`)},
			},
			{
				Keys: "spacee",
				Want: Want{Error: errors.New(`"spacee" is not a supported special key, see tmux manual for reference`)},
			},
			{
				Keys: "meta",
				Want: Want{Error: errors.New(`"meta" is not a supported special key, see tmux manual for reference`)},
			},
			{
				Keys: "shift",
				Want: Want{Error: errors.New(`"shift" is not a supported special key, see tmux manual for reference`)},
			},
			{
				Keys: "upp",
				Want: Want{Error: errors.New(`"upp" is not a supported special key, see tmux manual for reference`)},
			},
		} {
			_, got := splitKeys(tc.Keys)
			assert.Equal(t, tc.Want.Error, got)
		}
	})
}

func TestValidateKey(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		for _, key := range []string{
			"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
			"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
			"C-a", "C-b", "C-c", "C-d", "C-e", "C-f", "C-g", "C-h", "C-i", "C-j", "C-k", "C-l", "C-m", "C-n", "C-o", "C-p", "C-q", "C-r", "C-s", "C-t", "C-u", "C-v", "C-w", "C-x", "C-y", "C-z",
			"C-A", "C-B", "C-C", "C-D", "C-E", "C-F", "C-G", "C-H", "C-I", "C-J", "C-K", "C-L", "C-M", "C-N", "C-O", "C-P", "C-Q", "C-R", "C-S", "C-T", "C-U", "C-V", "C-W", "C-X", "C-Y", "C-Z",
			"c-a", "c-b", "c-c", "c-d", "c-e", "c-f", "c-g", "c-h", "c-i", "c-j", "c-k", "c-l", "c-m", "c-n", "c-o", "c-p", "c-q", "c-r", "c-s", "c-t", "c-u", "c-v", "c-w", "c-x", "c-y", "c-z",
			"S-a", "S-b", "S-c", "S-d", "S-e", "S-f", "S-g", "S-h", "S-i", "S-j", "S-k", "S-l", "S-m", "S-n", "S-o", "S-p", "S-q", "S-r", "S-s", "S-t", "S-u", "S-v", "S-w", "S-x", "S-y", "S-z",
			"S-A", "S-B", "S-C", "S-D", "S-E", "S-F", "S-G", "S-H", "S-I", "S-J", "S-K", "S-L", "S-M", "S-N", "S-O", "S-P", "S-Q", "S-R", "S-S", "S-T", "S-U", "S-V", "S-W", "S-X", "S-Y", "S-Z",
			"s-a", "s-b", "s-c", "s-d", "s-e", "s-f", "s-g", "s-h", "s-i", "s-j", "s-k", "s-l", "s-m", "s-n", "s-o", "s-p", "s-q", "s-r", "s-s", "s-t", "s-u", "s-v", "s-w", "s-x", "s-y", "s-z",
			"M-a", "M-b", "M-c", "M-d", "M-e", "M-f", "M-g", "M-h", "M-i", "M-j", "M-k", "M-l", "M-m", "M-n", "M-o", "M-p", "M-q", "M-r", "M-s", "M-t", "M-u", "M-v", "M-w", "M-x", "M-y", "M-z",
			"M-A", "M-B", "M-C", "M-D", "M-E", "M-F", "M-G", "M-H", "M-I", "M-J", "M-K", "M-L", "M-M", "M-N", "M-O", "M-P", "M-Q", "M-R", "M-S", "M-T", "M-U", "M-V", "M-W", "M-X", "M-Y", "M-Z",
			"m-a", "m-b", "m-c", "m-d", "m-e", "m-f", "m-g", "m-h", "m-i", "m-j", "m-k", "m-l", "m-m", "m-n", "m-o", "m-p", "m-q", "m-r", "m-s", "m-t", "m-u", "m-v", "m-w", "m-x", "m-y", "m-z",
			"Up", "Down", "Left", "Right", "BSpace", "BTab", "DC", "Delete", "End", "Enter", "Escape", "Home", "IC", "Insert", "NPage", "PageDown", "PgDn", "PPage", "PageUp", "PgUp", "Space", "Tab",
			"UP", "DOWN", "LEFT", "RIGHT", "BSPACE", "BTAB", "DC", "DELETE", "END", "ENTER", "ESCAPE", "HOME", "IC", "INSERT", "NPAGE", "PAGEDOWN", "PGDN", "PPAGE", "PAGEUP", "PGUP", "SPACE", "TAB",
			"up", "down", "left", "right", "bspace", "btab", "dc", "delete", "end", "enter", "escape", "home", "ic", "insert", "npage", "pagedown", "pgdn", "ppage", "pageup", "pgup", "space", "tab",
			"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10", "F11", "F12",
			"f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11", "f12",
		} {
			assert.NoError(t, validateKey(key))
		}
	})

	t.Run("error", func(t *testing.T) {
		for _, key := range []string{
			"aa",
			"c-aa",
			"m-aa",
			"s-aa",
      "-a",
      " a",
      "  a",
      " a ",
      "$a",
		} {
			assert.Error(t, validateKey(key), "key = " + key)
		}
	})
}
