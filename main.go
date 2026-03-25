package main

import (
	"context"
	"errors"
	"os"

	"charm.land/log/v2"
	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/cli"
	"github.com/Tesohh/isshues/config"
	"github.com/fsnotify/fsnotify"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	tint "github.com/lrstanley/bubbletint/v2"
	viperlib "github.com/spf13/viper"
)

func main() {
	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
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

	config.ApplyDefaultConfig(viper)

	err = viper.ReadInConfig()

	var fileLookupError viperlib.ConfigFileNotFoundError
	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &fileLookupError) {
			log.Warn("no config file found, using default settings", "err", err)
		} else {
			log.Fatal("no config file found, using default settings", "err", err)
		}
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Config file changed:", "filename", e.Name)
	})
	viper.WatchConfig()

	tint.NewDefaultRegistry()

	app := app.NewApp(conn, viper, cli.RootCmd)
	app.Start()
}
