package issues

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/Tesohh/isshues/app"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/Tesohh/isshues/model"
	tint "github.com/lrstanley/bubbletint/v2"
)

type Model struct {
	app   *app.App
	theme *tint.Tint

	showFullHelp     bool
	fullScreenHeight int

	userId    int64
	projectId int64

	project db.Project
}

func New(userId int64, projectId int64, app *app.App, theme *tint.Tint) Model {
	m := Model{
		app:          app,
		theme:        theme,
		showFullHelp: false,
		userId:       userId,
		projectId:    projectId,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.LoadProjectCmd)
}

func (m Model) Title() string {
	if m.project.Title == "" {
		return "Issues"
	}
	return m.project.Title
}

func (m Model) Rehydrate() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (model.NavModel, tea.Cmd) {
	switch msg := msg.(type) {
	case UpdateProjectMsg:
		m.project = msg.Project
	case model.ThemeChangedMsg:
		m.theme = msg.NewTheme
	}
	return m, nil
}

func (m Model) View() string {
	return "hello this is the view..."
}

// TODO:
func (m Model) ShortHelp() []key.Binding {
	return []key.Binding{}
}
func (m Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
func (m Model) ShowFullHelp() bool {
	return false
}
