package statusbar

import (
	"strings"

	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/config"
	"github.com/Tesohh/isshues/model"
	tint "github.com/lrstanley/bubbletint/v2"
)

type Model struct {
	App   *app.App
	Theme *tint.Tint

	Width   int
	HelpBar help.Model
}

func New(app *app.App, theme *tint.Tint) Model {
	help := help.New()
	help.Styles.ShortKey = help.Styles.ShortKey.Foreground(lipgloss.Darken(theme.Fg, 0.3))
	help.Styles.ShortDesc = help.Styles.ShortDesc.Foreground(lipgloss.Darken(theme.Fg, 0.5))
	help.Styles.ShortSeparator = help.Styles.ShortSeparator.Foreground(lipgloss.Darken(theme.Fg, 0.5))
	help.Styles.FullKey = help.Styles.ShortKey.Foreground(lipgloss.Darken(theme.Fg, 0.3))
	help.Styles.FullDesc = help.Styles.ShortDesc.Foreground(lipgloss.Darken(theme.Fg, 0.5))
	help.Styles.FullSeparator = help.Styles.FullSeparator.Foreground(lipgloss.Darken(theme.Fg, 0.5))

	return Model{
		App:     app,
		Theme:   theme,
		Width:   0,
		HelpBar: help,
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
