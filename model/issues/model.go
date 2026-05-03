package issues

import (
	"fmt"
	"slices"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"github.com/Tesohh/isshues/app"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/Tesohh/isshues/model"
	"github.com/Tesohh/isshues/model/tabs"
	tint "github.com/lrstanley/bubbletint/v2"
)

type Model struct {
	app   *app.App
	theme *tint.Tint

	showFullHelp     bool
	fullScreenWidth  int
	fullScreenHeight int

	userId    int64
	projectId int64

	tabs tabs.Model

	project    db.Project
	views      []db.View          // TODO consider switching this to a map
	viewData   map[int64]viewData // if a entry (for a view id) doesnt exist, it means it wasn't loaded
	viewModels map[int64]Panels   // TEMP for now only panels, figure something out for tables in future (interface?)
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
		viewModels:   make(map[int64]Panels),
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

func (m Model) PropagateAllModels(msg tea.Msg) []tea.Cmd {
	cmds := []tea.Cmd{}
	for k, model := range m.viewModels {
		var cmd tea.Cmd
		m.viewModels[k], cmd = model.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cmds
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
		log.Info("successfully loaded view", "viewID", msg.viewID)

		// Create (or recreate) the actual view
		index := slices.IndexFunc(m.views, func(v db.View) bool {
			return v.ID == msg.viewID
		})
		if index == -1 {
			log.Warn("updated tab not found somehow")
			break
		}

		view := m.views[index]
		if view.Style == db.ViewStylePanels {
			var cmd tea.Cmd
			m.viewModels[view.ID], cmd = NewPanels(m.userId, m.app, m.project).
				SetSize(m.fullScreenWidth, m.fullScreenHeight-1).
				SetTheme(m.theme).
				SetViewAndData(view, msg.viewData)
			cmds = append(cmds, cmd)
		} else {
			cmds = append(cmds, model.MakeErrCmd(fmt.Errorf("style %s unsupported", view.Style)))
		}

	case tabs.SwitchTabMsg:
		id := m.tabs.SelectedID()
		if _, ok := m.viewData[id]; !ok {
			cmds = append(cmds, m.MakeLoadIssuesForSelectedViewCmd())
		}

	case tea.KeyPressMsg:
		id := m.tabs.SelectedID()

		var cmd tea.Cmd
		if model, ok := m.viewModels[id]; ok {
			m.viewModels[id], cmd = model.Update(msg)
		}
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		m.fullScreenWidth = msg.Width
		m.fullScreenHeight = msg.Height
		cmds = append(cmds, m.PropagateAllModels(msg)...)

	case model.ThemeChangedMsg:
		m.theme = msg.NewTheme
		cmds = append(cmds, m.PropagateAllModels(msg)...)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	style := lipgloss.NewStyle().
		Height(m.fullScreenHeight - 1).
		MaxHeight(m.fullScreenHeight - 1).
		Width(m.fullScreenWidth).
		MaxWidth(m.fullScreenWidth)

	viewRender := ""
	if viewModel, ok := m.viewModels[m.tabs.SelectedID()]; ok {
		viewRender = viewModel.View()
	}

	render := style.Render(lipgloss.JoinVertical(lipgloss.Top,
		m.tabs.View(),
		viewRender,
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
