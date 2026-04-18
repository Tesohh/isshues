package root

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/common"
	"github.com/Tesohh/isshues/model"
	"github.com/Tesohh/isshues/model/projects"
	"github.com/Tesohh/isshues/model/statusbar"
	tint "github.com/lrstanley/bubbletint/v2"
)

type Model struct {
	App   *app.App
	Theme *tint.Tint

	NavStack []model.NavModel
	// ProjectsView         projects.Model
	// IssuesViewSideBySide issues.SideBySideModel

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

func (m Model) testChangeTHeme() tea.Msg {
	time.Sleep(5 * time.Second)
	log.Info("changed theme")
	theme, _ := tint.GetTint("gruvbox_dark")
	return model.ThemeChangedMsg{NewTheme: theme}
}

func (m Model) Active() model.NavModel {
	return m.NavStack[len(m.NavStack)-1]
}

func (m Model) Propagate(msg tea.Msg) []tea.Cmd {
	cmds := []tea.Cmd{}
	for i, nav := range m.NavStack {
		var cmd tea.Cmd
		m.NavStack[i], cmd = nav.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cmds
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.Active().Init(), m.StatusBar.Init()) //, m.testChangeTHeme)
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

		cmds = append(cmds, m.Propagate(msg)...)

		var sbcmd tea.Cmd
		m.StatusBar, sbcmd = m.StatusBar.Update(msg)
		cmds = append(cmds, sbcmd)

	default:
		var cmd tea.Cmd
		m.NavStack[len(m.NavStack)-1], cmd = m.Active().Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	statusbar := m.StatusBar.View(m.Active())

	// statusbar = m.StatusBar.HelpBar.ShortHelpView(m.ProjectsView.ShortHelp())

	v := tea.NewView(m.Active().View() + statusbar)

	v.AltScreen = true
	v.WindowTitle = fmt.Sprintf("isshues %s / %s / %s", common.GetVersion(), m.App.Viper.GetString("company.name"), m.Active().Title())
	v.BackgroundColor = m.Theme.Bg

	return v
}
