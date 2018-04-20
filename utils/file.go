package utils

import (
	"log"
	"os"
)

// OpenOrCreateFile takes path as arguments and returns opened file,
// if here is no file, function creates it
func OpenOrCreateFile(path string) *os.File {
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
