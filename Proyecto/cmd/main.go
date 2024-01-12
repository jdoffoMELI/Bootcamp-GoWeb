package main

import (
	"os"
	"proyecto/internal/application"
)

func main() {
	/* Set environment variables */
	os.Setenv("TOKEN", "123456") // Token to access data modification operations

	/* Run the application */
	app := application.NewApplicationDefault("localhost:8080")
	app.Run()
}
