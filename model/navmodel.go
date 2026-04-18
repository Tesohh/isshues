package model

import tea "charm.land/bubbletea/v2"

type NavModel interface {
	Init() tea.Cmd
	Update(tea.Msg) (NavModel, tea.Cmd)
	View() string

	// This model shall provide help keys
	Helper

	// Title used for statusbar, window title...
	Title() string

	// similar to Init(), called when the model becomes active but already exists.
	// for example, when receiving a Refresh message, set a "dirty" bit, and on Rehydrate, query new data if "dirty"
	Rehydrate() tea.Cmd
}
