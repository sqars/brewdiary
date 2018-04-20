package main

import (
	"log"

	"github.com/sqars/brewdiary/app"
	"github.com/sqars/brewdiary/config"
	"github.com/sqars/brewdiary/logger"
	"github.com/sqars/brewdiary/utils"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	f := utils.OpenOrCreateFile("logs/log.log")
	defer f.Close()

	logger.Init(f)
	a := app.NewApp(conf)
	defer a.DB.Close()

	err = a.Run()
	if err != nil {
		logger.Error.Println(err)
	}
}
