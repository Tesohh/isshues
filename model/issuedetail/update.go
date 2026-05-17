package issuedetail

import (
	tea "charm.land/bubbletea/v2"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.descriptionViewport, cmd = m.descriptionViewport.Update(msg)
	cmds = append(cmds, cmd)

	m.tabs, cmd = m.tabs.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m = m.SetSize(msg.Width, msg.Height) // TODO: pading
	}

	return m, tea.Batch(cmds...)
}
