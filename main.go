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

	a, err := app.NewApp(conf)
	defer a.DB.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
