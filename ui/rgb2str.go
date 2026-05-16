package ui

import (
	"fmt"
	"image/color"
)

func Rgb2Str(rgb color.Color) string {
	r, g, b, _ := rgb.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}
