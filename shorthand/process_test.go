package shorthand

import (
	"errors"
	"testing"

	"github.com/spf13/viper"
)

func TestParsePriorityWithViper(t *testing.T) {
	viper := viper.New()

	viper.SetDefault("priorities.crit.value", 10)
	viper.SetDefault("priorities.high.value", 5)
	viper.SetDefault("priorities.med.value", 3)
	viper.SetDefault("priorities.default.value", 1)
	viper.SetDefault("priorities.low.value", 0)

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
