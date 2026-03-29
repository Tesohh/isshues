package shorthand

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	t.Run("+feat !crit +frontend add Nuke", func(t *testing.T) {
		input := "+feat !crit +frontend add Nuke"
		want := parserCaptures{
			Raws:         []string{"add", "Nuke"},
			Mentions:     []string{},
			Labels:       []string{"feat", "frontend"},
			Priorities:   []string{"crit"},
			Dependencies: []string{},
		}
		got := Parse(input)
		if reflect.DeepEqual(got, want) {
			t.Errorf("got: %v != want: %v", got, want)
		}
	})

	t.Run("+feat +gfx add Nuke 3D graphics @lallos >1 !low", func(t *testing.T) {
		input := "+feat +gfx add Nuke 3D graphics @lallos >1 !low"
		want := parserCaptures{
			Raws:         []string{"add", "Nuke", "3D", "graphics"},
			Mentions:     []string{"lallos"},
			Labels:       []string{"feat", "gfx"},
			Priorities:   []string{"low"},
			Dependencies: []string{"1"},
		}
		got := Parse(input)
		if reflect.DeepEqual(got, want) {
			t.Errorf("got: %v != want: %v", got, want)
		}
	})

	t.Run("+idea consider adding nuke 4D graphics @quantum-team >2", func(t *testing.T) {
		input := "+idea consider adding nuke 4D graphics @quantum-team >2"
		want := parserCaptures{
			Raws:         []string{"consider", "adding", "nuke", "4D", "graphics"},
			Mentions:     []string{"quantum-team"},
			Labels:       []string{"idea"},
			Priorities:   []string{},
			Dependencies: []string{"2"},
		}
		got := Parse(input)
		if reflect.DeepEqual(got, want) {
			t.Errorf("got: %v != want: %v", got, want)
		}
	})

	t.Run("+bug fix exploding phone bug @nobody", func(t *testing.T) {
		input := "+bug fix exploding phone bug @nobody"
		want := parserCaptures{
			Raws:         []string{"fix", "exploding", "phone", "bug"},
			Mentions:     []string{"nobody"},
			Labels:       []string{"bug"},
			Priorities:   []string{},
			Dependencies: []string{},
		}
		got := Parse(input)
		if reflect.DeepEqual(got, want) {
			t.Errorf("got: %v != want: %v", got, want)
		}
	})
}
