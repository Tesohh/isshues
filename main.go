package main

import (
	"context"
	"log"
	"os"

	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/cli"
	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
)

const (
	host = "0.0.0.0"
	port = "2222"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalln("cannot connect to database!", err)
	}

	app := app.NewApp(host, port, conn, cli.RootCmd)
	app.Start()
}
