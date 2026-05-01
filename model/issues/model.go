package issues

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Tesohh/isshues/app"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/Tesohh/isshues/model"
	"github.com/Tesohh/isshues/model/tabs"
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
	fullScreenWidth  int
	fullScreenHeight int

	userId    int64
	projectId int64

	tabs tabs.Model

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
		tabs:         tabs.New(0, []tabs.Tab{}, theme),
		viewData:     make(map[int64]viewData),
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return m.LoadProjectCmd
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
	cmds := []tea.Cmd{}

	var tabCmd tea.Cmd
	m.tabs, tabCmd = m.tabs.Update(msg)
	cmds = append(cmds, tabCmd)

	switch msg := msg.(type) {

	case UpdateProjectMsg: // theoretically, this only happens once.
		m.project = msg.Project
		m.views = msg.Views

		tabList := make([]tabs.Tab, 0, len(msg.Views))
		for _, view := range msg.Views {
			tabList = append(tabList, tabs.NewTab(view.ID, view.Name))
		}

		var cmd tea.Cmd
		m.tabs, cmd = m.tabs.Update(tabs.UpdateTabsMsg{Tabs: tabList})
		cmds = append(cmds, cmd)

		cmds = append(cmds, m.MakeLoadIssuesForSelectedViewCmd())

	case UpdateViewDataMsg:
		m.viewData[msg.viewID] = msg.viewData
		fmt.Printf("%#v\n", m.viewData[msg.viewID])
	case tea.WindowSizeMsg:
		m.fullScreenWidth = msg.Width
		m.fullScreenHeight = msg.Height

	case model.ThemeChangedMsg:
		m.theme = msg.NewTheme
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := lipgloss.NewStyle().
		Height(m.fullScreenHeight - 1).
		MaxHeight(m.fullScreenHeight - 1).
		Width(m.fullScreenWidth).
		MaxWidth(m.fullScreenWidth)

	render := s.Render(lipgloss.JoinVertical(lipgloss.Top,
		m.tabs.View(),
		"alskjdflkjasfdl",
	))

	return render
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
