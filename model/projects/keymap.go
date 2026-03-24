package projects

import (
	"charm.land/bubbles/v2/key"
)

// type Keymap struct {
// 	Up key.Binding
// }
//
// func (k Keymap) ShortHelp() []key.Binding {}
// func (k Keymap) FullHelp() []key.Binding  {}

func (m ProjectsView) ShortHelp() []key.Binding {
	if m.creationForm == nil {
		return m.list.ShortHelp()
	} else {
		return m.creationForm.KeyBinds()
	}
}

func (m ProjectsView) FullHelp() [][]key.Binding {
	if m.creationForm == nil {
		return m.list.FullHelp()
	} else {
		return [][]key.Binding{m.creationForm.KeyBinds()}
	}
}
