package ui

import (
	"charm.land/lipgloss/v2"
	tint "github.com/lrstanley/bubbletint/v2"
)

var HLDefs = HLMap{
	tint.TintRosePine.ID: HLSet{
		Base:            lipgloss.Color("#191724"),
		Surface:         lipgloss.Color("#1f1d2e"),
		Overlay:         lipgloss.Color("#26233a"),
		Muted:           lipgloss.Color("#6e6a86"),
		Subtle:          lipgloss.Color("#908caa"),
		Text:            lipgloss.Color("#e0def4"),
		Accent:          lipgloss.Color("#e40078"), // custom fuchsia - taken from Tesohh/dotfiles/hypr
		Emphasis:        lipgloss.Darken(lipgloss.Color("#c4a7e7"), 0.2),
		Error:           lipgloss.Color("#eb6f92"),
		Warning:         lipgloss.Color("#f6c177"),
		Success:         lipgloss.Color("#31748f"),
		StatusTodo:      lipgloss.Color("#9ccfd8"),
		StatusProgress:  lipgloss.Color("#3174af"), // slightly tuned Pine (bluer)
		StatusDone:      lipgloss.Color("#c4a7e7"),
		StatusCancelled: lipgloss.Color("#eb6f92"),
	},
	tint.TintRosePineMoon.ID: HLSet{
		Base:            lipgloss.Color("#232136"),
		Surface:         lipgloss.Color("#2a273f"),
		Overlay:         lipgloss.Color("#393552"),
		Muted:           lipgloss.Color("#6e6a86"),
		Subtle:          lipgloss.Color("#908caa"),
		Text:            lipgloss.Color("#e0def4"),
		Accent:          lipgloss.Color("#c40058"), // slightly tuned down fuchsia
		Emphasis:        lipgloss.Darken(lipgloss.Color("#c4a7e7"), 0.2),
		Error:           lipgloss.Color("#eb6f92"),
		Warning:         lipgloss.Color("#f6c177"),
		Success:         lipgloss.Color("#3e8fb0"),
		StatusTodo:      lipgloss.Color("#9ccfd8"),
		StatusProgress:  lipgloss.Color("#3e8fd0"), // slightly tuned Pine (bluer)
		StatusDone:      lipgloss.Color("#c4a7e7"),
		StatusCancelled: lipgloss.Color("#eb6f92"),
	},
	tint.TintRosePineDawn.ID: HLSet{
		Base:            lipgloss.Color("#faf4ed"),
		Surface:         lipgloss.Color("#fffaf3"),
		Overlay:         lipgloss.Color("#f2e9e1"),
		Muted:           lipgloss.Color("#9893a5"),
		Subtle:          lipgloss.Color("#797593"),
		Text:            lipgloss.Color("#575279"),
		Accent:          lipgloss.Color("#b40038"), // slightly tuned down fuchsia
		Emphasis:        lipgloss.Darken(lipgloss.Color("#907aa9"), 0.2),
		Error:           lipgloss.Color("#b4637a"),
		Warning:         lipgloss.Color("#ea9d34"),
		Success:         lipgloss.Color("#286983"),
		StatusTodo:      lipgloss.Color("#56949f"),
		StatusProgress:  lipgloss.Color("#2869a3"), // slightly tuned Pine (bluer)
		StatusDone:      lipgloss.Color("#907aa9"),
		StatusCancelled: lipgloss.Color("#b4637a"),
	},
}
