package issues

import (
	tea "charm.land/bubbletea/v2"
	"github.com/Tesohh/isshues/app"
	tint "github.com/lrstanley/bubbletint/v2"
)

type SideBySideModel struct {
	app   *app.App
	theme *tint.Tint

	showFullHelp     bool
	fullScreenHeight int

	userId int64
}

func NewSideBySide(userId int64, app *app.App, theme *tint.Tint) SideBySideModel {
	m := SideBySideModel{
		app:          app,
		theme:        theme,
		showFullHelp: false,
		userId:       userId,
	}

	return m
}

func (m SideBySideModel) Init() tea.Cmd {
	return nil
}

func (m SideBySideModel) Update(msg tea.Msg) (SideBySideModel, tea.Cmd) {
	return m, nil
}

func (m SideBySideModel) View() string {
	return "hello this is the side by side view"
}
