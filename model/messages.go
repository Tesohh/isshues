package model

import (
	"errors"

	tea "charm.land/bubbletea/v2"
	tint "github.com/lrstanley/bubbletint/v2"
)

var (
	InternalErr = errors.New("internal error")
)

type ThemeChangedMsg struct {
	NewTheme *tint.Tint
}

type ErrMsg struct {
	Err error
}

func InternalErrMsg() ErrMsg {
	return ErrMsg{Err: InternalErr}
}

func MakeErrCmd(err error) func() tea.Msg {
	return func() tea.Msg {
		return ErrMsg{err}
	}
}
