package config

import (
	"github.com/spf13/viper"
)

func default_groups() []DefaultGroup {
	return []DefaultGroup{
		{
			Name:        "admins",
			Color:       "red",
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

func ApplyDefaultConfig(viper *viper.Viper) {
	viper.SetDefault("ssh.host", "0.0.0.0")
	viper.SetDefault("ssh.port", "2222")

	viper.SetDefault("company.name", "pausetta.org")
	viper.SetDefault("company.logo", "")
	viper.SetDefault("company.logo_bg", "purple")
	viper.SetDefault("company.logo_fg", "bg")
	viper.SetDefault("company.hide", false)

	viper.SetDefault("default_groups", default_groups())
}
