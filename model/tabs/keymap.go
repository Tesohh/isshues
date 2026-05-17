package tabs

import "charm.land/bubbles/v2/key"

type Keymap struct {
	HelpKey *key.Binding // if not nil, will use this as the help key
	Tab1    key.Binding
	Tab2    key.Binding
	Tab3    key.Binding
	Tab4    key.Binding
	Tab5    key.Binding
	Tab6    key.Binding
	Tab7    key.Binding
	Tab8    key.Binding
	Tab9    key.Binding
	Tab10   key.Binding
}

func (k Keymap) Help() []key.Binding {
	if k.HelpKey != nil {
		return []key.Binding{*k.HelpKey}
	} else {
		return []key.Binding{k.Tab1, k.Tab2, k.Tab3, k.Tab4, k.Tab5, k.Tab6, k.Tab7, k.Tab8, k.Tab9, k.Tab10}
	}
}

var DefaultKeymap = Keymap{
	HelpKey: new(key.NewBinding(key.WithHelp("ctrl+(1-0)", "switch tab"))),
	Tab1:    key.NewBinding(key.WithKeys("ctrl+1")),
	Tab2:    key.NewBinding(key.WithKeys("ctrl+2")),
	Tab3:    key.NewBinding(key.WithKeys("ctrl+3")),
	Tab4:    key.NewBinding(key.WithKeys("ctrl+4")),
	Tab5:    key.NewBinding(key.WithKeys("ctrl+5")),
	Tab6:    key.NewBinding(key.WithKeys("ctrl+6")),
	Tab7:    key.NewBinding(key.WithKeys("ctrl+7")),
	Tab8:    key.NewBinding(key.WithKeys("ctrl+8")),
	Tab9:    key.NewBinding(key.WithKeys("ctrl+9")),
	Tab10:   key.NewBinding(key.WithKeys("ctrl+0")),
}
