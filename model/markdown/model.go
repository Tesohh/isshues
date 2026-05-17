package markdown

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"charm.land/glamour/v2/ansi"
	"github.com/Tesohh/isshues/ui"
	tint "github.com/lrstanley/bubbletint/v2"
)

type Model struct {
	renderer *glamour.TermRenderer
	content  string
	output   string
}

func New() Model {
	return Model{}
}

func (m Model) SetTheme(theme *tint.Tint) Model {
	hl := func(key ui.HLKey) *string {
		return new(ui.Rgb2Str(ui.HLDefs.Get(key, theme)))
	}

	m.renderer, _ = glamour.NewTermRenderer(
		glamour.WithStyles(
			ansi.StyleConfig{
				Document: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{
						BlockPrefix: "\n",
						BlockSuffix: "\n",
						Color:       hl(ui.HLKeyText),
					},
					Margin: new(uint(2)),
				},
				BlockQuote: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{},
					Indent:         new(uint(1)),
					IndentToken:    new("│ "),
				},
				List: ansi.StyleList{
					StyleBlock: ansi.StyleBlock{
						StylePrimitive: ansi.StylePrimitive{
							Color: hl(ui.HLKeyText),
						},
					},
					LevelIndent: 2,
				},
				Heading: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{
						BlockSuffix: "\n",
						Color:       hl(ui.HLKeyAccent),
						Bold:        new(true),
					},
				},
				H1: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{
						Prefix: "# ",
						Bold:   new(true),
					},
				},
				H2: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{
						Prefix: "## ",
					},
				},
				H3: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{
						Prefix: "### ",
					},
				},
				H4: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{
						Prefix: "#### ",
					},
				},
				H5: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{
						Prefix: "##### ",
					},
				},
				H6: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{
						Prefix: "###### ",
					},
				},
				Strikethrough: ansi.StylePrimitive{
					CrossedOut: new(true),
				},
				Emph: ansi.StylePrimitive{
					Color:  hl(ui.HLKeyItalic),
					Italic: new(true),
				},
				Strong: ansi.StylePrimitive{
					Color: hl(ui.HLKeyBold),
					Bold:  new(true),
				},
				HorizontalRule: ansi.StylePrimitive{
					Color:  hl(ui.HLKeyMuted),
					Format: "\n--------\n",
				},
				Item: ansi.StylePrimitive{
					BlockPrefix: "• ",
				},
				Enumeration: ansi.StylePrimitive{
					BlockPrefix: ". ",
					Color:       hl(ui.HLKeyEmphasis),
				},
				Task: ansi.StyleTask{
					StylePrimitive: ansi.StylePrimitive{},
					Ticked:         "[✓] ",
					Unticked:       "[ ] ",
				},
				Link: ansi.StylePrimitive{
					Color:     hl(ui.HLKeyEmphasis),
					Underline: new(true),
				},
				LinkText: ansi.StylePrimitive{
					Color: hl(ui.HLKeyAccent),
				},
				Image: ansi.StylePrimitive{
					Color:     hl(ui.HLKeyEmphasis),
					Underline: new(true),
				},
				ImageText: ansi.StylePrimitive{
					Color:  hl(ui.HLKeyAccent),
					Format: "Image: {{.text}} →",
				},
				Code: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{
						Color: hl(ui.HLKeySuccess),
					},
				},
				CodeBlock: ansi.StyleCodeBlock{
					StyleBlock: ansi.StyleBlock{
						StylePrimitive: ansi.StylePrimitive{
							BackgroundColor: hl(ui.HLKeySurface),
						},
						Margin: new(uint(2)),
					},
					Chroma: &ansi.Chroma{
						Text: ansi.StylePrimitive{
							Color: hl(ui.HLKeyText),
						},
						Error: ansi.StylePrimitive{
							Color:           hl(ui.HLKeyText),
							BackgroundColor: hl(ui.HLKeyError),
						},
						Comment: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaComment),
						},
						CommentPreproc: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaCommentPreproc),
						},
						Keyword: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaKeyword),
						},
						KeywordReserved: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaKeywordReserved),
						},
						KeywordNamespace: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaKeywordNamespace),
						},
						KeywordType: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaKeywordType),
						},
						Operator: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaOperator),
						},
						Punctuation: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaPunctuation),
						},
						Name: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaName),
						},
						NameConstant: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaNameConstant),
						},
						NameBuiltin: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaNameBuiltin),
						},
						NameTag: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaNameTag),
						},
						NameAttribute: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaNameAttribute),
						},
						NameClass: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaNameClass),
						},
						NameDecorator: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaNameDecorator),
						},
						NameFunction: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaNameFunction),
						},
						LiteralNumber: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaLiteralNumber),
						},
						LiteralString: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaLiteralString),
						},
						LiteralStringEscape: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaLiteralStringEscape),
						},
						GenericDeleted: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaGenericDeleted),
						},
						GenericEmph: ansi.StylePrimitive{
							Color:  hl(ui.HLKeyItalic),
							Italic: new(true),
						},
						GenericInserted: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaGenericInserted),
						},
						GenericStrong: ansi.StylePrimitive{
							Color: hl(ui.HLKeyBold),
							Bold:  new(true),
						},
						GenericSubheading: ansi.StylePrimitive{
							Color: hl(ui.HLKeyChromaGenericSubheading),
						},
						Background: ansi.StylePrimitive{
							BackgroundColor: hl(ui.HLKeySurface),
						},
					},
				},
				Table: ansi.StyleTable{
					StyleBlock: ansi.StyleBlock{
						StylePrimitive: ansi.StylePrimitive{},
					},
				},
				DefinitionDescription: ansi.StylePrimitive{
					BlockPrefix: "\n🠶 ",
				},
			}),
		glamour.WithChromaFormatter("terminal16m"),
	)
	return m
}

func (m Model) SetWidth(width int) Model {
	_ = glamour.WithWordWrap(width)(m.renderer)
	return m
}

// as an optimization, markdown is only rendered once when setting the content, not when calling View
func (m Model) SetContent(content string) Model {
	m.content = content
	out, err := m.renderer.Render(m.content)
	m.output = out
	if err != nil {
		m.output += fmt.Sprintf("\n\nerror: %s", err)
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(_ tea.Msg) (Model, tea.Cmd) {
	return m, nil
}
func (m Model) View() string {
	return m.output
}

func (m Model) ShortHelp() []key.Binding {
	return []key.Binding{}
}
func (m Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
func (m Model) ShowFullHelp() bool {
	return false
}
