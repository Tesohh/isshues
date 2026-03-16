package app

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type (
	errMsg  error
	chatMsg struct {
		id   string
		text string
	}
)

type model struct {
	app         *App
	viewport    viewport.Model
	messages    []string
	id          string
	textarea    *textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	taStyles := ta.Styles()
	taStyles.Focused.CursorLine = lipgloss.NewStyle()
	ta.SetStyles(taStyles)

	ta.ShowLineNumbers = false

	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    &ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	ti, tiCmd := m.textarea.Update(msg)
	m.textarea = &ti
	vp, vpCmd := m.viewport.Update(msg)
	m.viewport = vp

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Key().Mod {
		case tea.ModCtrl: // We're only interested in ctrl+<key>
			switch msg.Key().Code {
			case 'c':
				return m, tea.Quit
			}
		}
		switch msg.Key().Code {
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.app.Broadcast(chatMsg{
				id:   m.id,
				text: m.textarea.Value(),
			})
			m.textarea.Reset()
		}

	case chatMsg:
		m.messages = append(m.messages, m.senderStyle.Render(msg.id)+": "+msg.text)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() tea.View {
	v := tea.NewView(fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n")
	return v
}
