package root

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Tesohh/isshues/common"
)

func (m Model) View() tea.View {
	content := m.Active().View()
	statusbar := m.StatusBar.View(m.Active())

	render := lipgloss.JoinVertical(lipgloss.Top, content, statusbar)

	v := tea.NewView(render)

	v.AltScreen = true
	v.WindowTitle = fmt.Sprintf("isshues %s / %s / %s", common.GetVersion(), m.App.Viper.GetString("company.name"), m.Active().Title())
	v.BackgroundColor = m.Theme.Bg

	return v
}
