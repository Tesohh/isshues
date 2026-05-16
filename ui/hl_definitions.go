package ui

import (
	"charm.land/lipgloss/v2"
	tint "github.com/lrstanley/bubbletint/v2"
)

//nolint:dupl
var HLDefs = HLMap{
	tint.TintRosePine.ID: HLSet{
		Base:    lipgloss.Color("#191724"),
		Surface: lipgloss.Color("#1f1d2e"),
		Overlay: lipgloss.Color("#26233a"),

		Muted:  lipgloss.Color("#6e6a86"),
		Subtle: lipgloss.Color("#908caa"),
		Text:   lipgloss.Color("#e0def4"),

		Accent:   lipgloss.Color("#e40078"),
		Emphasis: lipgloss.Darken(lipgloss.Color("#c4a7e7"), 0.2),
		Error:    lipgloss.Color("#eb6f92"),
		Warning:  lipgloss.Color("#f6c177"),
		Success:  lipgloss.Color("#31748f"),

		Bold:   lipgloss.Color("#eb6f92"),                       // Love
		Italic: lipgloss.Darken(lipgloss.Color("#c4a7e7"), 0.2), // Emphasis

		ChromaComment:             lipgloss.Color("#6e6a86"), // Muted
		ChromaCommentPreproc:      lipgloss.Color("#9ccfd8"), // Foam
		ChromaKeyword:             lipgloss.Color("#31748f"), // Pine
		ChromaKeywordReserved:     lipgloss.Color("#31748f"), // Pine
		ChromaKeywordNamespace:    lipgloss.Color("#31748f"), // Pine
		ChromaKeywordType:         lipgloss.Color("#9ccfd8"), // Foam
		ChromaOperator:            lipgloss.Color("#908caa"), // Subtle
		ChromaPunctuation:         lipgloss.Color("#908caa"), // Subtle
		ChromaName:                lipgloss.Color("#9ccfd8"), // Foam
		ChromaNameConstant:        lipgloss.Color("#c4a7e7"), // Iris
		ChromaNameBuiltin:         lipgloss.Color("#eb6f92"), // Love
		ChromaNameTag:             lipgloss.Color("#9ccfd8"), // Foam
		ChromaNameAttribute:       lipgloss.Color("#ebbcba"), // Rose
		ChromaNameClass:           lipgloss.Color("#9ccfd8"), // Foam
		ChromaNameDecorator:       lipgloss.Color("#c4a7e7"), // Rose
		ChromaNameFunction:        lipgloss.Color("#c4a7e7"), // Rose
		ChromaLiteralNumber:       lipgloss.Color("#ebbcba"), // Rose
		ChromaLiteralString:       lipgloss.Color("#f6c177"), // Gold
		ChromaLiteralStringEscape: lipgloss.Color("#9ccfd8"), // Foam
		ChromaGenericDeleted:      lipgloss.Color("#eb6f92"), // Love
		ChromaGenericInserted:     lipgloss.Color("#31748f"), // Pine
		ChromaGenericSubheading:   lipgloss.Color("#c4a7e7"), // Iris

		StatusTodo:      lipgloss.Color("#9ccfd8"),
		StatusProgress:  lipgloss.Color("#3174af"),
		StatusDone:      lipgloss.Color("#c4a7e7"),
		StatusCancelled: lipgloss.Color("#eb6f92"),
	},
	tint.TintRosePineMoon.ID: HLSet{
		Base:    lipgloss.Color("#232136"),
		Surface: lipgloss.Color("#2a273f"),
		Overlay: lipgloss.Color("#393552"),

		Muted:  lipgloss.Color("#6e6a86"),
		Subtle: lipgloss.Color("#908caa"),
		Text:   lipgloss.Color("#e0def4"),

		Accent:   lipgloss.Color("#c40058"),
		Emphasis: lipgloss.Darken(lipgloss.Color("#c4a7e7"), 0.2),
		Error:    lipgloss.Color("#eb6f92"),
		Warning:  lipgloss.Color("#f6c177"),
		Success:  lipgloss.Color("#3e8fb0"),

		Bold:   lipgloss.Color("#eb6f92"),                       // Love
		Italic: lipgloss.Darken(lipgloss.Color("#c4a7e7"), 0.2), // Emphasis

		ChromaComment:             lipgloss.Color("#6e6a86"), // Muted
		ChromaCommentPreproc:      lipgloss.Color("#9ccfd8"), // Foam
		ChromaKeyword:             lipgloss.Color("#3e8fb0"), // Pine
		ChromaKeywordReserved:     lipgloss.Color("#3e8fb0"), // Pine
		ChromaKeywordNamespace:    lipgloss.Color("#3e8fb0"), // Pine
		ChromaKeywordType:         lipgloss.Color("#9ccfd8"), // Foam
		ChromaOperator:            lipgloss.Color("#908caa"), // Subtle
		ChromaPunctuation:         lipgloss.Color("#908caa"), // Subtle
		ChromaName:                lipgloss.Color("#9ccfd8"), // Foam
		ChromaNameConstant:        lipgloss.Color("#c4a7e7"), // Iris
		ChromaNameBuiltin:         lipgloss.Color("#eb6f92"), // Love
		ChromaNameTag:             lipgloss.Color("#9ccfd8"), // Foam
		ChromaNameAttribute:       lipgloss.Color("#ea9a97"), // Rose
		ChromaNameClass:           lipgloss.Color("#9ccfd8"), // Foam
		ChromaNameDecorator:       lipgloss.Color("#ea9a97"), // Rose
		ChromaNameFunction:        lipgloss.Color("#ea9a97"), // Rose
		ChromaLiteralNumber:       lipgloss.Color("#ea9a97"), // Rose
		ChromaLiteralString:       lipgloss.Color("#f6c177"), // Gold
		ChromaLiteralStringEscape: lipgloss.Color("#9ccfd8"), // Foam
		ChromaGenericDeleted:      lipgloss.Color("#eb6f92"), // Love
		ChromaGenericInserted:     lipgloss.Color("#3e8fb0"), // Pine
		ChromaGenericSubheading:   lipgloss.Color("#c4a7e7"), // Iris

		StatusTodo:      lipgloss.Color("#9ccfd8"),
		StatusProgress:  lipgloss.Color("#3e8fd0"),
		StatusDone:      lipgloss.Color("#c4a7e7"),
		StatusCancelled: lipgloss.Color("#eb6f92"),
	},
	tint.TintRosePineDawn.ID: HLSet{
		Base:    lipgloss.Color("#faf4ed"),
		Surface: lipgloss.Color("#fffaf3"),
		Overlay: lipgloss.Color("#f2e9e1"),

		Muted:  lipgloss.Color("#9893a5"),
		Subtle: lipgloss.Color("#797593"),
		Text:   lipgloss.Color("#575279"),

		Accent:   lipgloss.Color("#b40038"),
		Emphasis: lipgloss.Darken(lipgloss.Color("#907aa9"), 0.2),
		Error:    lipgloss.Color("#b4637a"),
		Warning:  lipgloss.Color("#ea9d34"),
		Success:  lipgloss.Color("#286983"),

		Bold:   lipgloss.Color("#b4637a"),                       // Love
		Italic: lipgloss.Darken(lipgloss.Color("#907aa9"), 0.2), // Emphasis

		ChromaComment:             lipgloss.Color("#9893a5"), // Muted
		ChromaCommentPreproc:      lipgloss.Color("#56949f"), // Foam
		ChromaKeyword:             lipgloss.Color("#286983"), // Pine
		ChromaKeywordReserved:     lipgloss.Color("#286983"), // Pine
		ChromaKeywordNamespace:    lipgloss.Color("#286983"), // Pine
		ChromaKeywordType:         lipgloss.Color("#56949f"), // Foam
		ChromaOperator:            lipgloss.Color("#797593"), // Subtle
		ChromaPunctuation:         lipgloss.Color("#797593"), // Subtle
		ChromaName:                lipgloss.Color("#56949f"), // Foam
		ChromaNameConstant:        lipgloss.Color("#907aa9"), // Iris
		ChromaNameBuiltin:         lipgloss.Color("#b4637a"), // Love
		ChromaNameTag:             lipgloss.Color("#56949f"), // Foam
		ChromaNameAttribute:       lipgloss.Color("#d7827e"), // Rose
		ChromaNameClass:           lipgloss.Color("#56949f"), // Foam
		ChromaNameDecorator:       lipgloss.Color("#d7827e"), // Rose
		ChromaNameFunction:        lipgloss.Color("#d7827e"), // Rose
		ChromaLiteralNumber:       lipgloss.Color("#d7827e"), // Rose
		ChromaLiteralString:       lipgloss.Color("#ea9d34"), // Gold
		ChromaLiteralStringEscape: lipgloss.Color("#56949f"), // Foam
		ChromaGenericDeleted:      lipgloss.Color("#b4637a"), // Love
		ChromaGenericInserted:     lipgloss.Color("#286983"), // Pine
		ChromaGenericSubheading:   lipgloss.Color("#907aa9"), // Iris

		StatusTodo:      lipgloss.Color("#56949f"),
		StatusProgress:  lipgloss.Color("#2869a3"),
		StatusDone:      lipgloss.Color("#907aa9"),
		StatusCancelled: lipgloss.Color("#b4637a"),
	},
}
