package main

import (
	"app/internal/application"
	"fmt"
	"os"
)

func main() {
	cfg := &application.ConfigAppDefault{
		ServerAddr: os.Getenv("SERVER_ADDR"),
		DbFile:     "/Users/jdoffo/Desktop/Practica Bootcamp/Bootcamp-GoWeb/Desafio-Cierre/docs/db/tickets.csv",
	}
	app := application.NewApplicationDefault(cfg)

	// - setup
	err := app.SetUp()
	if err != nil {
		fmt.Println(err)
		return
	}

	// - run
	err = app.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
