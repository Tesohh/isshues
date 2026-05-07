package issues

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Tesohh/isshues/ui"
	tint "github.com/lrstanley/bubbletint/v2"
)

type itemDelegate struct {
	theme *tint.Tint
}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 1 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := "  " + i.Title()

	selectedStyle := lipgloss.NewStyle()
	chevronStyle := lipgloss.NewStyle().Background(ui.HLDefs.Get(ui.HLKeyAccent, d.theme)).Foreground(ui.HLDefs.Get(ui.HLKeyBase, d.theme))

	if index == m.Index() {
		str = strings.Replace(str, " ", chevronStyle.Render(">"), 1)
		str = selectedStyle.Render(str)
	}

	_, _ = fmt.Fprint(w, str+"\n"+i.Description())
}
