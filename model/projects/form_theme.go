package projects

import (
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	tint "github.com/lrstanley/bubbletint/v2"
)

// TODO: move this somewhere else....
type dynamicFormTheme struct {
	*tint.Tint
}

func (d dynamicFormTheme) Theme(isDark bool) *huh.Styles {
	t := huh.ThemeBase(isDark)

	var (
		base     = d.Bg
		text     = d.Fg
		subtext1 = lipgloss.Darken(d.Fg, 0.6)
		subtext0 = lipgloss.Darken(d.Fg, 0.4)
		// overlay1 = lipgloss.Darken(d.Fg, 0.3)
		overlay0 = lipgloss.Darken(d.Fg, 0.2)
		green    = d.Green
		red      = d.Red
		pink     = d.Purple
		mauve    = d.Blue
		cursor   = lipgloss.Darken(d.Fg, 0.6)
	)

	t.Focused.Base = t.Focused.Base.BorderForeground(subtext1)
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(mauve)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(mauve)
	t.Focused.Directory = t.Focused.Directory.Foreground(mauve)
	t.Focused.Description = t.Focused.Description.Foreground(subtext0)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(pink)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(pink)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(pink)
	t.Focused.Option = t.Focused.Option.Foreground(text)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(pink)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(green)
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(green)
	t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.Foreground(text)
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(text)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(base).Background(pink)
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(text).Background(base)

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(cursor)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(overlay0)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(pink)
	t.Focused.TextInput.Text = t.Focused.TextInput.Text.Foreground(text)

	t.Blurred = t.Focused
	t.Blurred.TextInput.Text = t.Blurred.TextInput.Text.Foreground(text)
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description
	return t
}
