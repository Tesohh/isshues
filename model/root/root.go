package root

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/model"
	"github.com/Tesohh/isshues/model/projects"
	"github.com/Tesohh/isshues/model/statusbar"
	tint "github.com/lrstanley/bubbletint/v2"
)

// TODO: consider just moving the model here.
// Instead of checkign the status and then redirecting the updates there, just send them to the last in the stack
type Status interface {
	// Title() string
}

type ViewingProjects struct{}

type Model struct {
	App         *app.App
	StatusStack []Status
	Theme       *tint.Tint

	ProjectsView projects.Model

	StatusBar statusbar.Model

	UserId int64
}

func New(app *app.App, userId int64, theme *tint.Tint) Model {
	return Model{
		UserId:       userId,
		ProjectsView: projects.New(userId, app, theme),
		StatusBar:    statusbar.New(app, theme),
		StatusStack:  []Status{ViewingProjects{}},
		Theme:        theme,
	}
}

func (m Model) testChangeTHeme() tea.Msg {
	time.Sleep(5 * time.Second)
	log.Info("changed theme")
	theme, _ := tint.GetTint("gruvbox_dark")
	return model.ThemeChangedMsg{NewTheme: theme}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.ProjectsView.Init(), m.StatusBar.Init()) //, m.testChangeTHeme)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	if key, ok := msg.(tea.KeyPressMsg); ok && key.String() == "ctrl+c" {
		return m, tea.Quit
	}

	var statusBarCmd tea.Cmd
	m.StatusBar, statusBarCmd = m.StatusBar.Update(msg)
	cmds = append(cmds, statusBarCmd)

	switch msg := msg.(type) {
	case model.ThemeChangedMsg:
		m.Theme = msg.NewTheme

		var pcmd tea.Cmd
		m.ProjectsView, pcmd = m.ProjectsView.Update(msg)
		cmds = append(cmds, pcmd)

		var sbcmd tea.Cmd
		m.StatusBar, sbcmd = m.StatusBar.Update(msg)
		cmds = append(cmds, sbcmd)

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

func (m Model) View() tea.View {
	statusbar := m.StatusBar.View(m.ProjectsView)

	// statusbar = m.StatusBar.HelpBar.ShortHelpView(m.ProjectsView.ShortHelp())

	v := tea.NewView(m.ProjectsView.View() + statusbar)

	v.AltScreen = true
	v.BackgroundColor = m.Theme.Bg

	return v
}
