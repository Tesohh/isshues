package shorthand

import (
	"errors"
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
			t.Errorf("got: %v != want: %v", got, want)
		}
	})

	t.Run("invalid capture with name path", func(t *testing.T) {
		input := "+feat !OOPS +frontend add Nuke"
		captures := Parse(input)

		wantv := 1
		wantw := WarningInternalErrorDefaulting
		got, warning := parsePriorityWithViper(captures.Priorities, viper)

		if got != wantv {
			t.Errorf("got: %v != want: %v", got, wantv)
		}

		if !errors.Is(warning, WarningInvalidPriority) {
			t.Errorf("got: %v != want: %v", warning, wantw)
		}
	})

	t.Run("valid capture with integer path", func(t *testing.T) {
		input := "+feat !999 +frontend add Nuke"
		captures := Parse(input)

		want := 999
		got, _ := parsePriorityWithViper(captures.Priorities, viper)

		if got != want {
			t.Errorf("got: %v != want: %v", got, want)
		}
	})
}
