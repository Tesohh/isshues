package statusbar

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type errResetMsg struct{}

func errResetCmd() tea.Msg {
	time.Sleep(5 * time.Second)
	return errResetMsg{}
}
