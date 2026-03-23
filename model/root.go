package model

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/model/projects"
)

// TODO: consider just moving the model here.
// Instead of checkign the status and then redirecting the updates there, just send them to the last in the stack
type Status interface {
	// Title() string
}

type ViewingProjects struct{}

type RootModel struct {
	App          *app.App
	StatusStack  []Status
	ProjectsView projects.ProjectsView
	UserId       int64
}

func NewRoot(app *app.App, userId int64) RootModel {
	return RootModel{
		UserId:       userId,
		ProjectsView: projects.New(userId, app),
		StatusStack:  []Status{ViewingProjects{}},
	}
}

func (m RootModel) Init() tea.Cmd {
	return tea.Batch(m.ProjectsView.Init())
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "ctrl+c" {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	default:
		switch status := m.StatusStack[len(m.StatusStack)-1].(type) {
		case ViewingProjects:
			_ = status
			var cmd tea.Cmd
			m.ProjectsView, cmd = m.ProjectsView.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

var kabashiStyle = lipgloss.NewStyle()

func (m RootModel) View() tea.View {
	v := tea.NewView(kabashiStyle.Render(m.ProjectsView.View()))
	v.AltScreen = true
	return v
}
