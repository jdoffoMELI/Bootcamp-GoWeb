package main

import "proyecto/internal/application"

func main() {
	app := application.NewApplicationDefault("localhost:8080")
	app.Run()
}
