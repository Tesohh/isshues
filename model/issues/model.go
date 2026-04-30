package issues

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/Tesohh/isshues/app"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/Tesohh/isshues/model"
	tint "github.com/lrstanley/bubbletint/v2"
)

type viewData struct {
	issuesMap     map[int64]*issueAndRelations
	users         []db.User  // list of all assignee users from all issues in this view
	labels        []db.Label // list of all labels from all issues in this view
	shallowIssues []db.Issue // list of all issues with a incoming relationship from all issues in this view
}

type Model struct {
	app   *app.App
	theme *tint.Tint

	showFullHelp     bool
	fullScreenHeight int

	userId    int64
	projectId int64

	project  db.Project
	views    []db.View
	viewData map[int64]viewData // if a entry (for a view id) doesnt exist, it means it wasn't loaded
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
	case UpdateProjectMsg: // theoretically, this only happens once.
		m.project = msg.Project
		m.views = msg.Views
	case UpdateViewDataMsg:
		m.viewData[msg.viewID] = msg.viewData
	case model.ThemeChangedMsg:
		m.theme = msg.NewTheme
	}
	return m, nil
}

func (m Model) View() string {
	s := "hello this is the view...\n"
	for _, view := range m.views {
		s += view.Name + "\n"
	}
	return s
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
