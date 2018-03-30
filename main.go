package main

import (
	"log"

	"github.com/sqars/brewdiary/app"
	"github.com/sqars/brewdiary/config"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	a := app.NewApp(conf)
	defer a.DB.Close()

	err = a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
