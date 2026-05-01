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

	case projects.SwitchToProjectMsg:
		issuesModel := issues.New(m.UserId, msg.ProjectId, m.App, m.Theme)

		// also send a resize update so that the model instantly has the right size
		newIssuesModel, _ := issuesModel.Update(tea.WindowSizeMsg{Height: m.Height, Width: m.Width})
		issuesModel = newIssuesModel.(issues.Model)

		cmds = append(cmds, issuesModel.Init())
		m.NavStack = append(m.NavStack, issuesModel)

	case tea.WindowSizeMsg:
		cmds = append(cmds, m.PropagateNav(msg)...)
		m.Height = msg.Height
		m.Width = msg.Width

	default:
		var cmd tea.Cmd
		m.NavStack[len(m.NavStack)-1], cmd = m.Active().Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
