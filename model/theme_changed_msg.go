package model

import tint "github.com/lrstanley/bubbletint/v2"

type ThemeChangedMsg struct {
	NewTheme *tint.Tint
}
