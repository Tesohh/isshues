package model

import (
	tea "charm.land/bubbletea/v2"
	tint "github.com/lrstanley/bubbletint/v2"
)

type ThemeChangedMsg struct {
	NewTheme *tint.Tint
}

type ErrMsg struct {
	Err error
}

func MakeErrCmd(err error) func() tea.Msg {
	return func() tea.Msg {
		return ErrMsg{err}
	}
}
