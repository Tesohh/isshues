package issuedetail

func (m Model) View() string {
	tabs := m.tabs.View()
	description := m.descriptionViewport.View()
	return tabs + "\n" + description
}
