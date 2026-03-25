package projects

import (
	"errors"

	"charm.land/bubbles/v2/key"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	tint "github.com/lrstanley/bubbletint/v2"
)

func prefixValidator(s string) error {
	if len(s) != 4 {
		return errors.New("prefix must be 4 characters long")
	}
	return nil
}

func MakeForm(theme *tint.Tint) *huh.Form {
	keymap := huh.NewDefaultKeyMap()
	keymap.Quit = key.NewBinding(key.WithKeys("esc"), key.WithHelp("quit", "quits the form"))

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Key("title").Title("Project title"),
			huh.NewInput().Key("prefix").Title("Project prefix").Validate(prefixValidator),
			huh.NewConfirm().Key("confirm").Affirmative("Create").Negative(""),
		),
	).WithKeyMap(keymap).WithShowHelp(false).WithTheme(dynamicFormTheme{theme})
}

var formStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

// TODO: move this somewhere else....
type dynamicFormTheme struct {
	*tint.Tint
}

// TODO: Redo this from scratch...
func (t dynamicFormTheme) Theme(isDark bool) *huh.Styles {
	charm := huh.ThemeBase(isDark)
	var (
		normalFg = t.Fg
		border   = lipgloss.Darken(t.Fg, 0.5)
		indigo   = t.Purple
		fuchsia  = t.Blue
		green    = t.Green
		red      = t.Red
	)

	charm.Focused.Base = charm.Focused.Base.BorderForeground(border)
	charm.Focused.Card = charm.Focused.Base
	charm.Focused.Title = charm.Focused.Title.Foreground(indigo).Bold(true)
	charm.Focused.NoteTitle = charm.Focused.NoteTitle.Foreground(indigo).Bold(true).MarginBottom(1)
	charm.Focused.Directory = charm.Focused.Directory.Foreground(indigo)
	charm.Focused.Description = charm.Focused.Description.Foreground(lipgloss.Color("243"))
	charm.Focused.ErrorIndicator = charm.Focused.ErrorIndicator.Foreground(red)
	charm.Focused.ErrorMessage = charm.Focused.ErrorMessage.Foreground(red)
	charm.Focused.SelectSelector = charm.Focused.SelectSelector.Foreground(border)
	charm.Focused.NextIndicator = charm.Focused.NextIndicator.Foreground(fuchsia)
	charm.Focused.PrevIndicator = charm.Focused.PrevIndicator.Foreground(fuchsia)
	charm.Focused.Option = charm.Focused.Option.Foreground(normalFg)
	charm.Focused.MultiSelectSelector = charm.Focused.MultiSelectSelector.Foreground(fuchsia)
	charm.Focused.SelectedOption = charm.Focused.SelectedOption.Foreground(green)
	charm.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.Color("#02A877")).SetString("✓ ")
	charm.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.Color("243")).SetString("• ")
	charm.Focused.UnselectedOption = charm.Focused.UnselectedOption.Foreground(normalFg)
	charm.Focused.FocusedButton = charm.Focused.FocusedButton.Foreground(normalFg).Background(fuchsia)
	charm.Focused.Next = charm.Focused.FocusedButton
	charm.Focused.BlurredButton = charm.Focused.BlurredButton.Foreground(normalFg).Background(lipgloss.Color("252"))

	charm.Focused.TextInput.Cursor = charm.Focused.TextInput.Cursor.Foreground(lipgloss.Darken(t.Fg, 0.2))
	charm.Focused.TextInput.Placeholder = charm.Focused.TextInput.Placeholder.Foreground(lipgloss.Color("238"))
	charm.Focused.TextInput.Prompt = charm.Focused.TextInput.Prompt.Foreground(fuchsia)

	charm.Blurred = charm.Focused
	charm.Blurred.Base = charm.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	charm.Blurred.Card = charm.Blurred.Base
	charm.Blurred.NextIndicator = lipgloss.NewStyle().Foreground(border)
	charm.Blurred.PrevIndicator = lipgloss.NewStyle().Foreground(border)

	charm.Group.Title = charm.Focused.Title
	charm.Group.Description = charm.Focused.Description
	return charm
}
