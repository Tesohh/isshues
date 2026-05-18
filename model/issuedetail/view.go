package issuedetail

func (m Model) View() string {
	header := m.MakeHeader()
	tabs := m.tabs.View()
	description := m.descriptionViewport.View()
	return header + "\n" + tabs + "\n" + description
}
