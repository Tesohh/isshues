package projects

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	db "github.com/Tesohh/isshues/db/generated"
)

type item struct {
	title, desc string
}

func itemFromProject(p db.Project) item {
	return item{
		title: fmt.Sprintf("#%s %s", p.Prefix, p.Title),
		desc:  "TODO!",
	}
}
func itemsFromProjects(ps []db.Project) []list.Item {
	items := []list.Item{}
	for _, p := range ps {
		items = append(items, itemFromProject(p))
	}
	return items
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type apper interface {
	GetDB() *db.Queries
}

type ProjectsView struct {
	list list.Model
	app  apper

	userId int64
}

func New(userId int64, app apper) ProjectsView {
	m := ProjectsView{
		list:   list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		app:    app,
		userId: userId,
	}

	m.list.Title = "Projects"

	return m
}

func (m ProjectsView) Init() tea.Cmd {
	return m.FetchProjectsCmd
}

func (m ProjectsView) Update(msg tea.Msg) (ProjectsView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	case UpdateProjectsMsg:
		m.list.SetItems(itemsFromProjects(msg.Projects))

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd

}
func (m ProjectsView) View() string {
	v := m.list.View()
	return v
}
