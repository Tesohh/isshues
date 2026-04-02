package shorthand

import (
	"testing"

	"github.com/Tesohh/isshues/config"
	"github.com/spf13/viper"
)

func TestParsePriorityWithViper(t *testing.T) {
	viper := viper.New()
	config.ApplyDefaultConfig(viper)

	t.Run("valid default path", func(t *testing.T) {
		input := "+feat +frontend add Nuke"
		captures := Parse(input)

		want := 1
		got, _ := parsePriorityWithViper(captures.Priorities, viper)
		if got != want {
			t.Errorf("got: %v != want: %v", got, want)
		}
	})

	t.Run("valid capture with name path", func(t *testing.T) {
		input := "+feat !crit +frontend add Nuke"
		captures := Parse(input)

		want := 10
		got, _ := parsePriorityWithViper(captures.Priorities, viper)

		if got != want {
			t.Errorf("got: %v != want: %v , %#v", got, want, captures)
		}
	})
}
