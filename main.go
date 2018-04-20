package main

import (
	"log"
	"os"

	"github.com/sqars/brewdiary/app"
	"github.com/sqars/brewdiary/config"
	"github.com/sqars/brewdiary/logger"
)

func openOrCreateFile(path string) *os.File {
	if _, err := os.Stat(path); err == nil {
		f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			log.Fatal(err)
		}
		return f
	}
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	f := openOrCreateFile("logs/http.log")
	defer f.Close()

	logger.Init(f)
	a := app.NewApp(conf)
	defer a.DB.Close()

	err = a.Run()
	if err != nil {
		logger.Error.Println(err)
	}
}
