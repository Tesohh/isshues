package common

import (
	"image/color"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	tint "github.com/lrstanley/bubbletint/v2"
)

func KeyToColor(theme *tint.Tint, key string) color.Color {
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
	case "":
		return theme.Fg
	default:
		return tint.FromHex("#FF0000")
	}
}

func NullableKeyToColor(theme *tint.Tint, defaultColor color.Color, key pgtype.Text) color.Color {
	if !key.Valid || (key.Valid && key.String == "") {
		return defaultColor
	}

	return KeyToColor(theme, key.String)
}
