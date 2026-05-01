package tabs

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Tesohh/isshues/model"
	tint "github.com/lrstanley/bubbletint/v2"
)

var (
	tabStyle       = lipgloss.NewStyle().Padding(0, 1)
	activeTabStyle = tabStyle

	tabGapStyle = tabStyle.
			BorderTop(false).
			BorderLeft(false).
			BorderRight(false)
)

var nerdFontNumbers = []rune{'󰲠', '󰲢', '󰲤', '󰲦', '󰲨', '󰲪', '󰲬', '󰲮', '󰲰', '󰿬'}

type Tab struct {
	ID    int64
	Title string
}

type Model struct {
	width    int
	selected int
	tabs     []Tab
	theme    *tint.Tint
}

type UpdateTabsMsg struct {
	Tabs []Tab
}

func New(width int, tabs []Tab, theme *tint.Tint) Model {
	return Model{
		width:    width,
		selected: 0,
		tabs:     tabs,
		theme:    theme,
	}
}

func NewTab(id int64, title string) Tab {
	return Tab{ID: id, Title: title}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case UpdateTabsMsg:
		m.tabs = msg.Tabs
		m.selected = min(m.selected, len(m.tabs)-1)
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case model.ThemeChangedMsg:
		m.theme = msg.NewTheme
	case tea.KeyPressMsg:
		r := msg.Key().Code
		if r >= '0' && r <= '9' {
			num := int(r - '0')
			if num == 0 {
				num = 10
			}

			index := min(num-1, len(m.tabs)-1)
			m.selected = index
		}
	}
	return m, nil
}

func (m Model) View() string {
	out := []string{}

	for i, tab := range m.tabs {
		if i == m.selected {
			out = append(out, activeTabStyle.Background(m.theme.Purple).Foreground(m.theme.Bg).Render(string(nerdFontNumbers[i])+"  "+tab.Title))
		} else {
			out = append(out, tabStyle.Background(lipgloss.Darken(m.theme.Black, 0.5)).Render(string(nerdFontNumbers[i])+"  "+tab.Title))
		}
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, out...)
	gap := tabGapStyle.Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(row)-2)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)

	return row
}

// returns -1 if the selected tab does not exist (often, if there are no tabs)
func (m Model) SelectedID() int64 {
	if m.selected > len(m.tabs)-1 {
		return -1
	}

	return m.tabs[m.selected].ID
}
