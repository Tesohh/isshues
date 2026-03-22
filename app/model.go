package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Tesohh/isshues/projects"
)

type Status interface{}

type ViewingProjects struct{}

type model struct {
	app          *App
	statusStack  []Status
	projectsView projects.ProjectsView
	id           string
	userId       int64
}

func initialModel(app *App, userId int64) model {
	return model{
		userId:       userId,
		projectsView: projects.New(userId, app),
		statusStack:  []Status{ViewingProjects{}},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.projectsView.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "ctrl+c" {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	default:
		switch status := m.statusStack[len(m.statusStack)-1].(type) {
		case ViewingProjects:
			_ = status
			var cmd tea.Cmd
			m.projectsView, cmd = m.projectsView.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

var kabashiStyle = lipgloss.NewStyle()

func (m model) View() tea.View {
	v := tea.NewView(kabashiStyle.Render(m.projectsView.View()))
	v.AltScreen = true
	return v
}
