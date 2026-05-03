package issues

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Tesohh/isshues/app"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/Tesohh/isshues/model"
	tint "github.com/lrstanley/bubbletint/v2"
)

const PanelsListXScalingFactor = 0.3

type Panels struct {
	app   *app.App
	theme *tint.Tint

	height, width int

	userId   int64
	project  db.Project
	view     db.View
	viewData viewData // if a entry (for a view id) doesnt exist, it means it wasn't loaded

	list list.Model
}

func NewPanels(
	userId int64, app *app.App, project db.Project,
) Panels {
	m := Panels{
		app:     app,
		userId:  userId,
		project: project,
		list:    list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
	}

	m.list.SetShowTitle(false)
	m.list.SetShowHelp(false)

	return m
}

func (m Panels) SetSize(width int, height int) Panels {
	m.width = width
	m.height = height

	m.list.SetSize(int(float64(width)*PanelsListXScalingFactor), height)
	return m
}

func (m Panels) SetTheme(theme *tint.Tint) Panels {
	m.theme = theme
	m.list.SetDelegate(itemDelegate{theme})
	return m
}

func (m Panels) SetViewAndData(view db.View, viewData viewData) (Panels, tea.Cmd) {
	m.view = view
	m.viewData = viewData

	// set list items
	items := itemsFromViewData(m.app, m.theme, &viewData)
	cmd := m.list.SetItems(items)

	return m, cmd
}

func (m Panels) Init() tea.Cmd {
	return nil
}

func (m Panels) Update(msg tea.Msg) (Panels, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m = m.SetSize(msg.Width, msg.Height)

	case tea.KeyPressMsg:
		m.list, cmd = m.list.Update(msg)

	case model.ThemeChangedMsg:
		m = m.SetTheme(msg.NewTheme)
	}
	return m, cmd
}

func (m Panels) View() string {
	listRender := m.list.View()
	return lipgloss.JoinHorizontal(lipgloss.Right, listRender, "GAGA")
}

func (panels *Panels) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (panels *Panels) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

func (panels *Panels) ShowFullHelp() bool {
	return false
}
