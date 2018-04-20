package main

import (
	"log"
	"os"

	"github.com/sqars/brewdiary/app"
	"github.com/sqars/brewdiary/config"
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

	fHTTP := openOrCreateFile("logs/http.log")
	defer fHTTP.Close()

	fDB := openOrCreateFile("logs/db.log")
	defer fDB.Close()

	a := app.NewApp(conf, fDB, fHTTP)
	defer a.DB.Close()

	err = a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
