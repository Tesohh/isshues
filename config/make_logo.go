package config

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/Tesohh/isshues/common"
	tint "github.com/lrstanley/bubbletint/v2"
	"github.com/spf13/viper"
)

func MakeLogo(viper *viper.Viper, theme *tint.Tint) string {
	bg := common.KeyToColor(theme, viper.GetString("company.logo_bg"))
	fg := common.KeyToColor(theme, viper.GetString("company.logo_fg"))

	style := lipgloss.NewStyle().Background(bg).Foreground(fg).PaddingLeft(1).PaddingRight(2) // kabashi
	logo := style.Render(viper.GetString("company.logo"))

	return logo
}

func MakeWaterMark(viper *viper.Viper, theme *tint.Tint) string {
	logo := MakeLogo(viper, theme)
	name := lipgloss.NewStyle().Foreground(theme.Fg).Render(viper.GetString("company.name"))

	return fmt.Sprintf("%s %s", logo, name)
}

func MakeWaterMarkReverse(viper *viper.Viper, theme *tint.Tint) string {
	logo := MakeLogo(viper, theme)
	name := lipgloss.NewStyle().Foreground(theme.Fg).Render(viper.GetString("company.name"))

	return fmt.Sprintf("%s %s", name, logo)
}
