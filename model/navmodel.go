package model

import tea "charm.land/bubbletea/v2"

// SubModel describes models that are controlled by the Root model or other SubModels.
type SubModel[M SubModel[M]] interface {
	Init() tea.Cmd
	Update(tea.Msg) (M, tea.Cmd)
	View() string
}

// Titler describes any model that can provide a title
// (eg. used for showing in the statusbar, window title...)
type Titler interface {
	Title() string
}

// Rehydrater describes any model that can be rehydrated,
// similarly to Init(), called when the model becomes active but already exists.
// for example, when receiving a Refresh message, set a "dirty" bit, and on Rehydrate, query new data if "dirty"
type Rehydrater interface {
	Rehydrate() tea.Cmd
}

// NavModel describes any model that can be placed in the NavStack of the root model.
// It is a submodel and it also must be able to provide Help keys, a title, and be able to be rehydrated
type NavModel interface {
	SubModel[NavModel]
	Helper
	Titler
	Rehydrater
}
