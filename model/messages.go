package model

import (
	"errors"

	tea "charm.land/bubbletea/v2"
	tint "github.com/lrstanley/bubbletint/v2"
)

var (
	ErrInternal = errors.New("internal error")
)

type ThemeChangedMsg struct {
	NewTheme *tint.Tint
}

type ErrMsg struct {
	Err error
}

func ErrInternalMsg() ErrMsg {
	return ErrMsg{Err: ErrInternal}
}

func MakeErrCmd(err error) func() tea.Msg {
	return func() tea.Msg {
		return ErrMsg{err}
	}
}
