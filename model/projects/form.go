package projects

import (
	"errors"

	"charm.land/bubbles/v2/key"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

func prefixValidator(s string) error {
	if len(s) != 4 {
		return errors.New("prefix must be 4 characters long")
	}
	return nil
}

var formStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

func MakeForm() *huh.Form {
	keymap := huh.NewDefaultKeyMap()
	keymap.Quit = key.NewBinding(key.WithKeys("esc"), key.WithHelp("quit", "quits the form"))

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Key("title").Title("Project title"),
			huh.NewInput().Key("prefix").Title("Project prefix").Validate(prefixValidator),
			huh.NewConfirm().Key("confirm").Affirmative("Create").Negative(""),
		),
	).WithKeyMap(keymap)
}
