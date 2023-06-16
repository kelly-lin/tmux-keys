package generate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerate(t *testing.T) {
	for desc, tc := range map[string]struct {
		Bindings []Binding
		Want     []string
	}{
		"1 layer binding": {
			Bindings: []Binding{
				{TableName: PREFIX_TABLE_NAME, Keys: "f", Cmd: "new-window"},
			},
			Want: []string{"tmux bind-key -Tprefix f new-window"},
		},
		"2 layer binding": {
			Bindings: []Binding{
				{TableName: PREFIX_TABLE_NAME, Keys: "fd", Cmd: "new-window"},
			},
			Want: []string{
				"tmux bind-key -Tprefix f switch-client -Tprefix_f",
				"tmux bind-key -Tprefix_f d new-window",
			},
		},
		"3 layer binding": {
			Bindings: []Binding{
				{TableName: PREFIX_TABLE_NAME, Keys: "fdc", Cmd: "new-window"},
			},
			Want: []string{
				"tmux bind-key -Tprefix f switch-client -Tprefix_f",
				"tmux bind-key -Tprefix_f d switch-client -Tprefix_f_d",
				"tmux bind-key -Tprefix_f_d c new-window",
			},
		},
		"bindings in different tables with the same starting key": {
			Bindings: []Binding{
				{TableName: PREFIX_TABLE_NAME, Keys: "fd", Cmd: "list-sessions"},
				{TableName: "some-other-table-name", Keys: "fd", Cmd: "new-window"},
			},
			Want: []string{
				"tmux bind-key -Tprefix f switch-client -Tprefix_f",
				"tmux bind-key -Tprefix_f d list-sessions",
				"tmux bind-key -Tsome-other-table-name f switch-client -Tsome-other-table-name_f",
				"tmux bind-key -Tsome-other-table-name_f d new-window",
			},
		},
		"multiple tables": {
			Bindings: []Binding{
				{TableName: PREFIX_TABLE_NAME, Keys: "f", Cmd: "new-window"},
				{TableName: "some-other-table-name", Keys: "fd", Cmd: "new-window"},
			},
			Want: []string{
				"tmux bind-key -Tprefix f new-window",
				"tmux bind-key -Tsome-other-table-name f switch-client -Tsome-other-table-name_f",
				"tmux bind-key -Tsome-other-table-name_f d new-window",
			},
		},
	} {
		t.Run(desc, func(t *testing.T) {
			cmds, _ := Generate(tc.Bindings)
			assert.Equal(t, tc.Want, cmds)
		})
	}
}
