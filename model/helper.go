package model

import "charm.land/bubbles/v2/key"

type Helper interface {
	ShortHelp() []key.Binding
	FullHelp() [][]key.Binding
	ShowFullHelp() bool
}
