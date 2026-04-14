package shorthand

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	t.Run("+feat !crit +frontend add Nuke", func(t *testing.T) {
		input := "+feat !crit +frontend add Nuke"
		want := parserCaptures{
			Text:         "add Nuke",
			Mentions:     nil,
			Labels:       []string{"feat", "frontend"},
			Priorities:   []string{"crit"},
			Dependencies: nil,
			Descriptions: nil,
		}
		got := Parse(input)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %#v != want: %#v", got, want)
		}
	})

	t.Run("+feat +gfx add Nuke 3D graphics @lallos >1 !low", func(t *testing.T) {
		input := "+feat +gfx add Nuke 3D graphics @lallos >1 !low"
		want := parserCaptures{
			Text:         "add Nuke 3D graphics",
			Mentions:     []string{"lallos"},
			Labels:       []string{"feat", "gfx"},
			Priorities:   []string{"low"},
			Dependencies: []string{"1"},
			Descriptions: nil,
		}
		got := Parse(input)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %#v != want: %#v", got, want)
		}
	})

	t.Run("+idea consider adding nuke 4D graphics @quantum-team >2", func(t *testing.T) {
		input := "+idea consider adding nuke 4D graphics @quantum-team >2"
		want := parserCaptures{
			Text:         "consider adding nuke 4D graphics",
			Mentions:     []string{"quantum-team"},
			Labels:       []string{"idea"},
			Priorities:   nil,
			Dependencies: []string{"2"},
		}
		got := Parse(input)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %#v != want: %#v", got, want)
		}
	})

	t.Run("+bug fix exploding phone bug @nobody", func(t *testing.T) {
		input := "+bug fix exploding phone bug @nobody"
		want := parserCaptures{
			Text:         "fix exploding phone bug",
			Mentions:     []string{"nobody"},
			Labels:       []string{"bug"},
			Priorities:   nil,
			Dependencies: nil,
			Descriptions: nil,
		}
		got := Parse(input)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %#v != want: %#v", got, want)
		}
	})

	t.Run("+bug fix exploding phone bug \"as a user, i would like it if my phone didnt explode.\"", func(t *testing.T) {
		input := "+bug fix exploding phone bug \"as a user, i would like it if my phone didnt explode.\""
		want := parserCaptures{
			Text:         "fix exploding phone bug",
			Mentions:     nil,
			Labels:       []string{"bug"},
			Priorities:   nil,
			Dependencies: nil,
			Descriptions: []string{"as a user, i would like it if my phone didnt explode."},
		}
		got := Parse(input)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %#v != want: %#v", got, want)
		}
	})
}
