package root

import (
	tea "charm.land/bubbletea/v2"
	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/model"
	"github.com/Tesohh/isshues/model/projects"
	"github.com/Tesohh/isshues/model/statusbar"
	tint "github.com/lrstanley/bubbletint/v2"
)

type Model struct {
	App   *app.App
	Theme *tint.Tint

	Width  int
	Height int

	NavStack  []model.NavModel
	StatusBar statusbar.Model

	UserId int64
}

func New(app *app.App, userId int64, theme *tint.Tint) Model {
	return Model{
		UserId:    userId,
		StatusBar: statusbar.New(app, theme),
		NavStack:  []model.NavModel{projects.New(userId, app, theme)},
		Theme:     theme,
	}
}

func (m Model) Active() model.NavModel {
	return m.NavStack[len(m.NavStack)-1]
}

func (m Model) PropagateNav(msg tea.Msg) []tea.Cmd {
	cmds := []tea.Cmd{}
	for i, nav := range m.NavStack {
		var cmd tea.Cmd
		m.NavStack[i], cmd = nav.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cmds
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.Active().Init(), m.StatusBar.Init())
}
