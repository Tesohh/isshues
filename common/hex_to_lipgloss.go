package common

import (
	"errors"
	"strconv"

	"charm.land/lipgloss/v2"
)

var (
	HexLen7Err     = errors.New("hex color must be 7 characters")
	HexWrongFormat = errors.New("hex format is invalid")
)

func HexToLipgloss(hex string) (lipgloss.RGBColor, error) {
	color := lipgloss.RGBColor{}

	if len(hex) != 7 {
		return color, HexLen7Err
	}

	if hex[0] != '#' {
		return color, HexWrongFormat
	}

	r, err := strconv.ParseInt(hex[1:3], 16, 9)
	if err != nil {
		return color, err
	}
	g, err := strconv.ParseInt(hex[3:5], 16, 9)
	if err != nil {
		return color, err
	}
	b, err := strconv.ParseInt(hex[5:7], 16, 9)
	if err != nil {
		return color, err
	}

	color.R = uint8(r)
	color.G = uint8(g)
	color.B = uint8(b)

	return color, nil
}
