package root

import (
	tea "charm.land/bubbletea/v2"
	"github.com/Tesohh/isshues/model"
	"github.com/Tesohh/isshues/model/issues"
	"github.com/Tesohh/isshues/model/projects"
)

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

		cmds = append(cmds, m.PropagateNav(msg)...)

		var sbcmd tea.Cmd
		m.StatusBar, sbcmd = m.StatusBar.Update(msg)
		cmds = append(cmds, sbcmd)

	case projects.SwitchToProjectMsg:
		issuesModel := issues.New(m.UserId, msg.ProjectId, m.App, m.Theme)
		cmds = append(cmds, issuesModel.Init())
		m.NavStack = append(m.NavStack, issuesModel)

	default:
		var cmd tea.Cmd
		m.NavStack[len(m.NavStack)-1], cmd = m.Active().Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
