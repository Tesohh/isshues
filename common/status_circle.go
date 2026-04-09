package common

import (
	"charm.land/lipgloss/v2"
	db "github.com/Tesohh/isshues/db/generated"
	tint "github.com/lrstanley/bubbletint/v2"
)

func MakeStatusCircle(theme *tint.Tint, status db.Status) string {
	var color *tint.Color

	switch status {
	case db.StatusTodo:
		color = theme.Green
	case db.StatusProgress:
		color = theme.Blue
	case db.StatusDone:
		color = theme.Purple
	case db.StatusCancelled:
		color = theme.Red
	default:
		color = theme.Fg
	}

	return lipgloss.NewStyle().Foreground(color).Render("◉")
}
