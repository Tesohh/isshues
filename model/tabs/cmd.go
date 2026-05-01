package tabs

import tea "charm.land/bubbletea/v2"

type SwitchTabMsg struct{}

func SwitchTabCmd() tea.Msg {
	return SwitchTabMsg{}
}
