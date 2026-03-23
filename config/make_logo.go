package config

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"github.com/Tesohh/isshues/common"
	"github.com/spf13/viper"
)

func MakeLogo(viper *viper.Viper) string {
	bg, err := common.HexToLipgloss(viper.GetString("company.logo_bg"))
	if err != nil {
		log.Warn("company.logo_bg hex parse error", "err", err)
		return "ERROR"
	}

	fg, err := common.HexToLipgloss(viper.GetString("company.logo_fg"))
	if err != nil {
		log.Warn("company.logo_fg hex parse error", "err", err)
		return "ERROR"
	}

	style := lipgloss.NewStyle().Background(bg).Foreground(fg).PaddingLeft(1).PaddingRight(2) // kabashi
	logo := style.Render(viper.GetString("company.logo"))

	return logo
}

func MakeWaterMark(viper *viper.Viper) string {
	logo := MakeLogo(viper)
	return fmt.Sprintf("%s %s", logo, viper.GetString("company.name"))
}
