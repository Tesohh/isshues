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
	err     error
	HelpBar help.Model
}

func New(app *app.App, theme *tint.Tint) Model {
	help := help.New()

	m := Model{
		App:     app,
		Theme:   theme,
		Width:   0,
		err:     nil,
		HelpBar: help,
	}

	m = m.refreshTheme()
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) refreshTheme() Model {
	m.HelpBar.Styles.ShortKey = m.HelpBar.Styles.ShortKey.Foreground(lipgloss.Darken(m.Theme.Fg, 0.3))
	m.HelpBar.Styles.ShortDesc = m.HelpBar.Styles.ShortDesc.Foreground(lipgloss.Darken(m.Theme.Fg, 0.5))
	m.HelpBar.Styles.ShortSeparator = m.HelpBar.Styles.ShortSeparator.Foreground(lipgloss.Darken(m.Theme.Fg, 0.5))
	m.HelpBar.Styles.FullKey = m.HelpBar.Styles.ShortKey.Foreground(lipgloss.Darken(m.Theme.Fg, 0.3))
	m.HelpBar.Styles.FullDesc = m.HelpBar.Styles.ShortDesc.Foreground(lipgloss.Darken(m.Theme.Fg, 0.5))
	m.HelpBar.Styles.FullSeparator = m.HelpBar.Styles.FullSeparator.Foreground(lipgloss.Darken(m.Theme.Fg, 0.5))
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width

	case model.ThemeChangedMsg:
		m.Theme = msg.NewTheme
		m = m.refreshTheme()

	case model.ErrMsg:
		m.err = msg.Err
		cmd = errResetCmd

	case errResetMsg:
		m.err = nil
	}
	return m, cmd
}

func (m Model) bottomRight() string {
	return config.MakeWaterMarkReverse(m.App.Viper, m.Theme) + " "
}

func (m Model) View(currentModel model.Helper) string {
	bottomRight := m.bottomRight()
	bottomRightWidth := lipgloss.Width(bottomRight)

	m.HelpBar.SetWidth(m.Width - bottomRightWidth)

	var help string

	if m.err == nil {
		if currentModel.ShowFullHelp() {
			help = m.HelpBar.FullHelpView(currentModel.FullHelp())
		} else {
			help = m.HelpBar.ShortHelpView(currentModel.ShortHelp())
		}

	} else {
		help = lipgloss.NewStyle().Foreground(m.Theme.Red).Render(m.err.Error())
	}

	helpWidth := lipgloss.Width(help)
	whitespaceWidth := max(m.Width-helpWidth-bottomRightWidth, 0)
	whitespace := strings.Repeat(" ", whitespaceWidth)

	return help + whitespace + bottomRight
}
