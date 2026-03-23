package config

type DefaultGroup struct {
	Name        string   `mapstructure:"name"`
	Color       string   `mapstructure:"color"`
	Mentionable bool     `mapstructure:"mentionable"`
	Permissions []string `mapstructure:"permissions"`
	AddCreator  bool     `mapstructure:"add_creator"`
}
