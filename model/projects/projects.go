package projects

import (
	"context"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/Tesohh/isshues/app"
	db "github.com/Tesohh/isshues/db/generated"
	tint "github.com/lrstanley/bubbletint/v2"
)

type ProjectsView struct {
	app   *app.App
	theme *tint.Tint

	list         list.Model
	creationForm *huh.Form

	showFullHelp     bool
	fullScreenHeight int

	userId int64
}

func New(userId int64, app *app.App, theme *tint.Tint) ProjectsView {
	m := ProjectsView{
		list:         list.New([]list.Item{}, itemDelegate{theme}, 0, 0),
		app:          app,
		theme:        theme,
		showFullHelp: false,
		userId:       userId,
	}

	m.list.Title = "Projects"
	m.list.SetShowHelp(false)
	m.list.Styles.Title = m.list.Styles.Title.Background(m.theme.Purple)
	// TODO: changing the list item styles theoretically needs the delegate
	return m
}

func (m ProjectsView) Init() tea.Cmd {
	return tea.Batch(m.FetchProjectsCmd, m.HasCreatePermissionCmd)
}

func (m ProjectsView) Update(msg tea.Msg) (ProjectsView, tea.Cmd) {
	var cmd tea.Cmd
	var formIsNew bool

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		m.fullScreenHeight = msg.Height

	case RefreshProjectsMsg:
		cmd = m.FetchProjectsCmd

	case UpdateProjectsMsg:
		m.list.SetItems(itemsFromProjects(msg.Projects))

	case InitHasCreatePermissionMsg:
		m.list.AdditionalShortHelpKeys = func() []key.Binding {
			return []key.Binding{key.NewBinding(key.WithKeys("+"), key.WithHelp("+", "create project"))}
		}

	case tea.KeyPressMsg:
		if msg.String() == "?" && m.creationForm == nil {
			m.showFullHelp = !m.showFullHelp
		} else if msg.String() == "+" && m.creationForm == nil {
			ctx := context.Background()
			// TODO: do this in a tea.Cmd
			hasPermission, _ := m.app.DB.UserHasGlobalPermission(ctx, db.UserHasGlobalPermissionParams{
				UserID:             m.userId,
				GlobalPermissionID: "create-projects",
			})
			// TODO: handle error

			if !hasPermission {
				break
			}

			m.creationForm = MakeForm(m.theme)
			cmd = m.creationForm.Init()
			formIsNew = true
		}
	}

	// set the correct height on the list
	if m.ShowFullHelp() {
		m.list.SetHeight(m.fullScreenHeight - len(m.FullHelp()))
	} else {
		m.list.SetHeight(m.fullScreenHeight)
	}

	var listCmd, formCmd tea.Cmd

	// Handle form statuses
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

		// TODO: clean this up
		view := formStyle.BorderForeground(lipgloss.Darken(m.theme.Fg, 0.5)).Render(m.creationForm.View())
		fw := lipgloss.Width(view)
		fh := lipgloss.Height(view)

		layers = append(layers, lipgloss.NewLayer(view).X((vw-fw)/2).Y((vh-fh)/2).Z(1))
	}

	return lipgloss.NewCompositor(layers...).Render()
}
