package root

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/Tesohh/isshues/common"
)

func (m Model) View() tea.View {
	statusbar := m.StatusBar.View(m.Active())

	v := tea.NewView(m.Active().View() + statusbar)

	v.AltScreen = true
	v.WindowTitle = fmt.Sprintf("isshues %s / %s / %s", common.GetVersion(), m.App.Viper.GetString("company.name"), m.Active().Title())
	v.BackgroundColor = m.Theme.Bg

	return v
}
