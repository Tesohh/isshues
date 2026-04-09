package config

type DefaultLabel struct {
	Name   string `mapstructure:"name"`
	Color  string `mapstructure:"color"`
	Symbol string `mapstructure:"symbol"`
}
