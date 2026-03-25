package statusbar

import (
	"strings"

	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/config"
	"github.com/Tesohh/isshues/model"
)

type Model struct {
	App *app.App

	Width   int
	HelpBar help.Model
}

func New(app *app.App) Model {
	return Model{
		App:     app,
		Width:   0,
		HelpBar: help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
	}
	return m, nil
}

func (m Model) bottomRight() string {
	return config.MakeWaterMarkReverse(m.App.Viper) + " "
}

func (m Model) View(currentModel model.Helper) string {
	bottomRight := m.bottomRight()
	bottomRightWidth := lipgloss.Width(bottomRight)

	m.HelpBar.SetWidth(m.Width - bottomRightWidth)

	help := ""

	if currentModel.ShowFullHelp() {
		help = m.HelpBar.FullHelpView(currentModel.FullHelp())
	} else {
		help = m.HelpBar.ShortHelpView(currentModel.ShortHelp())
	}

	helpWidth := lipgloss.Width(help)

	whitespaceWidth := max(m.Width-helpWidth-bottomRightWidth, 0)
	whitespace := strings.Repeat(" ", whitespaceWidth)

	return help + whitespace + bottomRight
}
