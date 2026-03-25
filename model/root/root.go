package root

import (
	tea "charm.land/bubbletea/v2"
	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/model/projects"
	"github.com/Tesohh/isshues/model/statusbar"
)

// TODO: consider just moving the model here.
// Instead of checkign the status and then redirecting the updates there, just send them to the last in the stack
type Status interface {
	// Title() string
}

type ViewingProjects struct{}

type RootModel struct {
	App         *app.App
	StatusStack []Status

	ProjectsView projects.ProjectsView

	StatusBar statusbar.Model

	UserId int64
}

func New(app *app.App, userId int64) RootModel {
	return RootModel{
		UserId:       userId,
		ProjectsView: projects.New(userId, app),
		StatusBar:    statusbar.New(app),
		StatusStack:  []Status{ViewingProjects{}},
	}
}

func (m RootModel) Init() tea.Cmd {
	return tea.Batch(m.ProjectsView.Init(), m.StatusBar.Init())
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "ctrl+c" {
		return m, tea.Quit
	}

	var statusBarCmd tea.Cmd
	m.StatusBar, statusBarCmd = m.StatusBar.Update(msg)
	cmds = append(cmds, statusBarCmd)

	switch msg := msg.(type) {
	default:
		switch m.StatusStack[len(m.StatusStack)-1].(type) {
		case ViewingProjects:
			var cmd tea.Cmd
			m.ProjectsView, cmd = m.ProjectsView.Update(msg)
			cmds = append(cmds, cmd)

		}
	}

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() tea.View {
	statusbar := m.StatusBar.View(m.ProjectsView)

	// statusbar = m.StatusBar.HelpBar.ShortHelpView(m.ProjectsView.ShortHelp())

	v := tea.NewView(m.ProjectsView.View() + statusbar)

	v.AltScreen = true
	return v
}
