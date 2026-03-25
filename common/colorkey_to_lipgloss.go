package common

import (
	"strings"

	tint "github.com/lrstanley/bubbletint/v2"
)

func KeyToColor(theme *tint.Tint, key string) *tint.Color {
	switch strings.TrimSpace(strings.ToLower(key)) {
	case "fg":
		return theme.Fg
	case "bg":
		return theme.Bg
	case "bright_black":
		return theme.BrightBlack
	case "bright_blue":
		return theme.BrightBlue
	case "bright_cyan":
		return theme.BrightCyan
	case "bright_green":
		return theme.BrightGreen
	case "bright_purple":
		return theme.BrightPurple
	case "bright_red":
		return theme.BrightRed
	case "bright_white":
		return theme.BrightWhite
	case "bright_yellow":
		return theme.BrightYellow
	case "black":
		return theme.Black
	case "blue":
		return theme.Blue
	case "cyan":
		return theme.Cyan
	case "green":
		return theme.Green
	case "purple":
		return theme.Purple
	case "red":
		return theme.Red
	case "white":
		return theme.White
	case "yellow":
		return theme.Yellow
	default:
		return tint.FromHex("#FF0000")
	}
}
