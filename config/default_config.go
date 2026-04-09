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

func defaultLabels() []DefaultLabel {
	return []DefaultLabel{
		{Name: "feat", Color: "blue", Symbol: "󰇈"},
		{Name: "fix", Color: "red", Symbol: ""},
		{Name: "chore", Color: "green", Symbol: "󰃢"},
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

	viper.SetDefault("priorities.crit.value", 10)
	viper.SetDefault("priorities.crit.color", "red")
	viper.SetDefault("priorities.high.value", 5)
	viper.SetDefault("priorities.high.color", "yellow")
	viper.SetDefault("priorities.med.value", 3)
	viper.SetDefault("priorities.med.color", "cyan")
	viper.SetDefault("priorities.default.value", 1)
	viper.SetDefault("priorities.low.value", 0)
	viper.SetDefault("priorities.low.color", "black")

	viper.SetDefault("default_groups", default_groups())
	viper.SetDefault("default_labels", defaultLabels())
}
