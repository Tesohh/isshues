package main

const (
	host = "localhost"
	port = "23234"
)

func main() {
	app := NewApp()
	app.Start()
}
