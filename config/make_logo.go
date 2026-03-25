package config

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/spf13/viper"
)

func MakeLogo(viper *viper.Viper) string {
	bg := lipgloss.Color(viper.GetString("company.logo_bg"))
	fg := lipgloss.Color(viper.GetString("company.logo_fg"))

	style := lipgloss.NewStyle().Background(bg).Foreground(fg).PaddingLeft(1).PaddingRight(2) // kabashi
	logo := style.Render(viper.GetString("company.logo"))

	return logo
}

func MakeWaterMark(viper *viper.Viper) string {
	logo := MakeLogo(viper)
	return fmt.Sprintf("%s %s", logo, viper.GetString("company.name"))
}

func MakeWaterMarkReverse(viper *viper.Viper) string {
	logo := MakeLogo(viper)
	return fmt.Sprintf("%s %s", viper.GetString("company.name"), logo)
}
