package main

import (
	"github.com/Tesohh/isshues/config"
	"github.com/spf13/viper"
)

func default_groups() []config.DefaultGroup {
	return []config.DefaultGroup{
		{
			Name:        "admins",
			Color:       "#AA89E1",
			Mentionable: false,
			AddCreator:  true,
			Permissions: []string{"write-issues", "read-issues", "edit-project", "delete-project"},
		},
		{
			Name:        "devs",
			Mentionable: false,
			AddCreator:  false,
			Permissions: []string{"write-issues", "read-issues"},
		},
		{
			Name:        "guests",
			Mentionable: false,
			AddCreator:  false,
			Permissions: []string{"read-issues"},
		},
	}
}

func default_config(viper *viper.Viper) {
	viper.SetDefault("ssh.host", "0.0.0.0")
	viper.SetDefault("ssh.port", "2222")

	viper.SetDefault("company.name", "pausetta.org")
	viper.SetDefault("company.logo", "")
	viper.SetDefault("company.logo_bg", "#AA89E1")
	viper.SetDefault("company.logo_fg", "#FFFFFF")
	viper.SetDefault("company.hide", false)

	viper.SetDefault("default_groups", default_groups())
}
