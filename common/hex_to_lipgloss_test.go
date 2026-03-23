package common

import (
	"testing"

	"charm.land/lipgloss/v2"
)

var hexlgtests = []struct {
	in  string
	out lipgloss.RGBColor
	err error
}{
	{"#", lipgloss.RGBColor{}, HexLen7Err},
	{"#FF0000", lipgloss.RGBColor{R: 255, G: 0, B: 0}, nil},
}

func TestHexToLipgloss(t *testing.T) {
	for _, tt := range hexlgtests {
		t.Run(tt.in, func(t *testing.T) {
			color, err := HexToLipgloss(tt.in)

			if err != tt.err {
				t.Errorf("wrong error! got %v, want %v", err, tt.err)
			} else if color != tt.out {
				t.Errorf("wrong color! got %d %d %d, want %d %d %d", color.R, color.G, color.B, tt.out.R, tt.out.G, tt.out.B)
				// t.Errorf("wrong color! got #%x%x%x, want #%x%x%x", color.R, color.G, color.B, tt.out.R, tt.out.G, tt.out.B)
			}
		})
	}
}
