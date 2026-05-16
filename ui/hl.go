package ui

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	tint "github.com/lrstanley/bubbletint/v2"
)

// Additional layer over bubbletint for defining themes more precisely.
//
// if a `HLSet` is not found, fallback to the `tint.Tint`
//
// Keys mostly follow rose pine palette terminology: https://rosepinetheme.com/palette/
type HLSet struct {
	Base    color.Color // Main background. Also Chroma background
	Surface color.Color // Secondary background (on top of Base)
	Overlay color.Color // Tertiary background (on top of Surface)

	Muted  color.Color // Tertiary foreground
	Subtle color.Color // Secondary foreground
	Text   color.Color // Regular foreground. Also chroma text

	Accent   color.Color // Things that must stand out
	Emphasis color.Color // Things that must stand out, but not as much (eg. @your_name in issue component)
	Error    color.Color // Self explanatory. Also chroma error
	Warning  color.Color // Self explanatory
	Success  color.Color // Self explanatory

	// The following HLs are for markdown rendering
	Bold   color.Color
	Italic color.Color

	ChromaComment             color.Color
	ChromaCommentPreproc      color.Color
	ChromaKeyword             color.Color
	ChromaKeywordReserved     color.Color
	ChromaKeywordNamespace    color.Color
	ChromaKeywordType         color.Color
	ChromaOperator            color.Color
	ChromaPunctuation         color.Color
	ChromaName                color.Color
	ChromaNameConstant        color.Color
	ChromaNameBuiltin         color.Color
	ChromaNameTag             color.Color
	ChromaNameAttribute       color.Color
	ChromaNameClass           color.Color
	ChromaNameDecorator       color.Color
	ChromaNameFunction        color.Color
	ChromaLiteralNumber       color.Color
	ChromaLiteralString       color.Color
	ChromaLiteralStringEscape color.Color
	ChromaGenericDeleted      color.Color
	ChromaGenericInserted     color.Color
	ChromaGenericSubheading   color.Color

	StatusTodo      color.Color // Self explanatory
	StatusProgress  color.Color // Self explanatory
	StatusDone      color.Color // Self explanatory
	StatusCancelled color.Color // Self explanatory
}

type HLKey string

const (
	HLKeyBase    = "base"
	HLKeySurface = "surface"
	HLKeyOverlay = "overlay"

	HLKeyMuted  = "muted"
	HLKeySubtle = "subtle"
	HLKeyText   = "text"

	HLKeyAccent   = "accent"
	HLKeyEmphasis = "emphasis"
	HLKeyError    = "error"
	HLKeyWarning  = "warning"
	HLKeySuccess  = "success"

	HLKeyBold   = "bold"   // also used for chroma generic strong
	HLKeyItalic = "italic" // also used for chroma generic emphasis

	HLKeyChromaComment             = "chroma-comment"
	HLKeyChromaCommentPreproc      = "chroma-comment-preproc"
	HLKeyChromaKeyword             = "chroma-keyword"
	HLKeyChromaKeywordReserved     = "chroma-keyword-reserved"
	HLKeyChromaKeywordNamespace    = "chroma-keyword-namespace"
	HLKeyChromaKeywordType         = "chroma-keyword-type"
	HLKeyChromaOperator            = "chroma-operator"
	HLKeyChromaPunctuation         = "chroma-punctuation"
	HLKeyChromaName                = "chroma-name"
	HLKeyChromaNameConstant        = "chroma-name-constant"
	HLKeyChromaNameBuiltin         = "chroma-name-builtin"
	HLKeyChromaNameTag             = "chroma-name-tag"
	HLKeyChromaNameAttribute       = "chroma-name-attribute"
	HLKeyChromaNameClass           = "chroma-name-class"
	HLKeyChromaNameDecorator       = "chroma-name-decorator"
	HLKeyChromaNameFunction        = "chroma-name-function"
	HLKeyChromaLiteralNumber       = "chroma-literal-number"
	HLKeyChromaLiteralString       = "chroma-literal-string"
	HLKeyChromaLiteralStringEscape = "chroma-literal-string-escape"
	HLKeyChromaGenericDeleted      = "chroma-generic-deleted"
	HLKeyChromaGenericInserted     = "chroma-generic-inserted"
	HLKeyChromaGenericSubheading   = "chroma-generic-subheading"

	HLKeyStatusTodo      = "status-todo"
	HLKeyStatusProgress  = "status-progress"
	HLKeyStatusDone      = "status-done"
	HLKeyStatusCancelled = "status-cancelled"
)

type HLMap map[string]HLSet

func (m HLMap) Get(key HLKey, theme *tint.Tint) color.Color {
	set, ok := m[theme.ID]
	if !ok {
		return highlightFallback(key, theme)
	}

	switch key {
	case HLKeyBase:
		return set.Base
	case HLKeySurface:
		return set.Surface
	case HLKeyOverlay:
		return set.Overlay
	case HLKeyMuted:
		return set.Muted
	case HLKeySubtle:
		return set.Subtle
	case HLKeyText:
		return set.Text
	case HLKeyAccent:
		return set.Accent
	case HLKeyEmphasis:
		return set.Emphasis
	case HLKeyError:
		return set.Error
	case HLKeyWarning:
		return set.Warning
	case HLKeySuccess:
		return set.Success
	case HLKeyChromaComment:
		return set.ChromaComment
	case HLKeyChromaCommentPreproc:
		return set.ChromaCommentPreproc
	case HLKeyChromaKeyword:
		return set.ChromaKeyword
	case HLKeyChromaKeywordReserved:
		return set.ChromaKeywordReserved
	case HLKeyChromaKeywordNamespace:
		return set.ChromaKeywordNamespace
	case HLKeyChromaKeywordType:
		return set.ChromaKeywordType
	case HLKeyChromaOperator:
		return set.ChromaOperator
	case HLKeyChromaPunctuation:
		return set.ChromaPunctuation
	case HLKeyChromaName:
		return set.ChromaName
	case HLKeyChromaNameConstant:
		return set.ChromaNameConstant
	case HLKeyChromaNameBuiltin:
		return set.ChromaNameBuiltin
	case HLKeyChromaNameTag:
		return set.ChromaNameTag
	case HLKeyChromaNameAttribute:
		return set.ChromaNameAttribute
	case HLKeyChromaNameClass:
		return set.ChromaNameClass
	case HLKeyChromaNameDecorator:
		return set.ChromaNameDecorator
	case HLKeyChromaNameFunction:
		return set.ChromaNameFunction
	case HLKeyChromaLiteralNumber:
		return set.ChromaLiteralNumber
	case HLKeyChromaLiteralString:
		return set.ChromaLiteralString
	case HLKeyChromaLiteralStringEscape:
		return set.ChromaLiteralStringEscape
	case HLKeyChromaGenericDeleted:
		return set.ChromaGenericDeleted
	case HLKeyChromaGenericInserted:
		return set.ChromaGenericInserted
	case HLKeyBold:
		return set.Bold
	case HLKeyItalic:
		return set.Italic
	case HLKeyChromaGenericSubheading:
		return set.ChromaGenericSubheading
	case HLKeyStatusTodo:
		return set.StatusTodo
	case HLKeyStatusProgress:
		return set.StatusProgress
	case HLKeyStatusDone:
		return set.StatusDone
	case HLKeyStatusCancelled:
		return set.StatusCancelled
	default:
		log.Warn("hlkey doesn't exist", "key", key)
		return lipgloss.Color("#FF0000")
	}
}

func highlightFallback(key HLKey, theme *tint.Tint) color.Color {
	var fgmodifier, bgmodifier func(color.Color, float64) color.Color
	if theme.Dark {
		fgmodifier = lipgloss.Darken
		bgmodifier = lipgloss.Lighten
	} else {
		fgmodifier = lipgloss.Lighten
		bgmodifier = lipgloss.Darken
	}

	switch key {
	case HLKeyBase:
		return theme.Bg
	case HLKeySurface:
		return bgmodifier(theme.Bg, 0.2)
	case HLKeyOverlay:
		return bgmodifier(theme.Bg, 0.3)
	case HLKeyMuted:
		return fgmodifier(theme.Fg, 0.4)
	case HLKeySubtle:
		return fgmodifier(theme.Fg, 0.3)
	case HLKeyText:
		return theme.Fg
	case HLKeyAccent:
		return theme.Purple
	case HLKeyEmphasis:
		return fgmodifier(theme.Purple, 0.2)
	case HLKeyError:
		return theme.Red
	case HLKeyWarning:
		return theme.Yellow
	case HLKeySuccess:
		return theme.Green
	case HLKeyBold:
		return theme.Red
	case HLKeyItalic:
		return fgmodifier(theme.Purple, 0.2)
	case HLKeyChromaComment:
		return fgmodifier(theme.Fg, 0.4)
	case HLKeyChromaCommentPreproc:
		return theme.Cyan
	case HLKeyChromaKeyword:
		return theme.Cyan
	case HLKeyChromaKeywordReserved:
		return theme.Cyan
	case HLKeyChromaKeywordNamespace:
		return theme.Cyan
	case HLKeyChromaKeywordType:
		return theme.Blue
	case HLKeyChromaOperator:
		return theme.Cyan
	case HLKeyChromaPunctuation:
		return theme.Fg
	case HLKeyChromaName:
		return theme.Blue
	case HLKeyChromaNameConstant:
		return theme.Purple
	case HLKeyChromaNameBuiltin:
		return theme.Blue
	case HLKeyChromaNameTag:
		return theme.Cyan
	case HLKeyChromaNameAttribute:
		return theme.Green
	case HLKeyChromaNameClass:
		return theme.Blue
	case HLKeyChromaNameDecorator:
		return theme.Green
	case HLKeyChromaNameFunction:
		return theme.Green
	case HLKeyChromaLiteralNumber:
		return theme.Fg
	case HLKeyChromaLiteralString:
		return theme.Yellow
	case HLKeyChromaLiteralStringEscape:
		return theme.Blue
	case HLKeyChromaGenericDeleted:
		return theme.Red
	case HLKeyChromaGenericInserted:
		return theme.Green
	case HLKeyChromaGenericSubheading:
		return theme.Purple
	case HLKeyStatusTodo:
		return theme.Green
	case HLKeyStatusProgress:
		return theme.Blue
	case HLKeyStatusDone:
		return theme.Purple
	case HLKeyStatusCancelled:
		return theme.Red
	default:
		log.Warn("hlkey doesn't exist (fallbacked)", "key", key)
		return lipgloss.Color("#4A412A")
	}
}
