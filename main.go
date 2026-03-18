package main

import (
	"github.com/Tesohh/isshues/app"
)

const (
	host = "0.0.0.0"
	port = "2222"
)

func main() {
	app := app.NewApp(host, port)
	app.Start()
}
