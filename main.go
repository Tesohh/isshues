package main

import (
	"errors"
	"os"

	"charm.land/log/v2"
	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/cli"
	db_complex "github.com/Tesohh/isshues/db/complex"
	"github.com/fsnotify/fsnotify"
	_ "github.com/joho/godotenv/autoload"
	tint "github.com/lrstanley/bubbletint/v2"
	viperlib "github.com/spf13/viper"
)

func main() {
	pool, err := db_complex.GetDbPool(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("cannot connect to database!", err)
	}

	viper := viperlib.New()
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$XDG_CONFIG_HOME/isshues")
	viper.AddConfigPath("$HOME/.config/isshues")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./.ignored/config") // for development purposes

	var fileLookupError viperlib.ConfigFileNotFoundError
	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &fileLookupError) {
			log.Warn("no config file found, using default settings", "err", err)
		} else {
			log.Fatal("no config file found, using default settings", "err", err)
		}
	}

	log.Info("Welcome", "company.name", viper.Get("company.name"))

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Config file changed:", "filename", e.Name)
	})
	viper.WatchConfig()

	tint.NewDefaultRegistry()

	app := app.NewApp(pool, viper, cli.RootCmd)
	app.Start()
}
