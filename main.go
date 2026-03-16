package main

import (
	"github.com/Tesohh/isshues/app"
)

const (
	host = "localhost"
	port = "23234"
)

func main() {
	app := app.NewApp(host, port)
	app.Start()
}
