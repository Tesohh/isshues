package projects

import (
	"context"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/Tesohh/isshues/app"
	db "github.com/Tesohh/isshues/db/generated"
)

type ProjectsView struct {
	list list.Model
	app  *app.App

	creationForm *huh.Form

	userId int64
}

func New(userId int64, app *app.App) ProjectsView {
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
	var cmd tea.Cmd
	var formIsNew bool

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	case RefreshProjectsMsg:
		cmd = m.FetchProjectsCmd
	case UpdateProjectsMsg:
		m.list.SetItems(itemsFromProjects(msg.Projects))
	case tea.KeyPressMsg:
		if msg.String() == "+" && m.creationForm == nil {
			ctx := context.Background()
			hasPermission, _ := m.app.DB.UserHasGlobalPermission(ctx, db.UserHasGlobalPermissionParams{
				UserID:             m.userId,
				GlobalPermissionID: "create-projects",
			})
			// TODO: handle error

			if !hasPermission {
				break
			}

			m.creationForm = MakeForm()
			cmd = m.creationForm.Init()
			formIsNew = true
		}
	}

	var listCmd, formCmd tea.Cmd

	if m.creationForm != nil && !formIsNew {
		var form huh.Model
		form, formCmd = m.creationForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.creationForm = f
		}

		switch m.creationForm.State {
		case huh.StateCompleted:
			title := m.creationForm.GetString("title")
			prefix := strings.ToUpper(m.creationForm.GetString("prefix"))

			formCmd = m.MakeCreateProjectCmd(title, prefix)

			m.creationForm = nil
		case huh.StateAborted:
			m.creationForm = nil
		}
	} else {
		m.list, listCmd = m.list.Update(msg)
	}

	return m, tea.Batch(cmd, listCmd, formCmd)

}
func (m ProjectsView) View() string {
	layers := []*lipgloss.Layer{
		lipgloss.NewLayer(m.list.View()),
	}

	if m.creationForm != nil {
		vw := m.list.Width()
		vh := m.list.Height()

		view := formStyle.Render(m.creationForm.View())
		fw := lipgloss.Width(view)
		fh := lipgloss.Height(view)

		layers = append(layers, lipgloss.NewLayer(view).X((vw-fw)/2).Y((vh-fh)/2).Z(1))
	}

	return lipgloss.NewCompositor(layers...).Render()
}
