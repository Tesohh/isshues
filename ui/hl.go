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
	Base    color.Color // Main background
	Surface color.Color // Secondary background (on top of Base)
	Overlay color.Color // Tertiary background (on top of Surface)

	Muted  color.Color // Tertiary foreground
	Subtle color.Color // Secondary foreground
	Text   color.Color // Regular foreground

	Accent   color.Color // Things that must stand out
	Emphasis color.Color // Things that must stand out, but not as much (eg. @your_name in issue component)
	Error    color.Color // Self explanatory
	Warning  color.Color // Self explanatory
	Success  color.Color // Self explanatory

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

	HLKeyStatusTodo      = "status-todo"
	HLKeyStatusProgress  = "status-progress"
	HLKeyStatusDone      = "status-done"
	HLKeyStatusCancelled = "status-cancelled"
)

type HLMap map[string]HLSet

func (m HLMap) GetHL(key HLKey, theme *tint.Tint) color.Color {
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
