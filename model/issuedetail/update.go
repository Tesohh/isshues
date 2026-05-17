package issuedetail

import (
	tea "charm.land/bubbletea/v2"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.descriptionViewport, cmd = m.descriptionViewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m = m.SetSize(msg.Width, msg.Height) // TODO: pading
	}

	return m, cmd
}
