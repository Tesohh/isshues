package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"charm.land/log/v2"
	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/cli"
	"github.com/fsnotify/fsnotify"
	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
	viperlib "github.com/spf13/viper"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("cannot connect to database!", err)
	}
	viper := viperlib.New()
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$XDG_CONFIG_HOME/isshues")
	viper.AddConfigPath("$HOME/.config/isshues")
	viper.AddConfigPath(".")

	default_config(viper)

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
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()

	app := app.NewApp(conn, viper, cli.RootCmd)
	app.Start()
}
