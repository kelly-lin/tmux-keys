package generate

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
						{Keys: "fd", Cmd: "new-window"},
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
						{Keys: "fdc", Cmd: "new-window"},
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
						{Keys: "fd", Cmd: "list-sessions"},
					},
				},
				{
					Name: "some-other-table-name",
					Bindings: []Binding{
						{Keys: "fd", Cmd: "new-window"},
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
						{Keys: "fd", Cmd: "new-window"},
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
	type TC struct {
		Keys string
		Want []string
	}

	for desc, tc := range []TC{
		{
			Keys: "",
			Want: []string{},
		},
		{
			Keys: "a",
			Want: []string{"a"},
		},
		{
			Keys: "ab",
			Want: []string{"a", "b"},
		},
	} {
		assert.Equal(t, tc.Want, splitKeys(tc.Keys), desc)
	}
}
