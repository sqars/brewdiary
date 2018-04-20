package logger

import (
	"io"
	"log"
)

var (
	// Info is logger for INFO logs
	Info *log.Logger
	// Warning is logger for WARNING logs
	Warning *log.Logger
	// Error is logger for ERROR logs
	Error *log.Logger
	// DB is logger for DB logs
	DB *log.Logger
)

// Init initialize logger variables with proper prefix
func Init(
	logFile io.Writer,
) {

	Info = log.New(logFile,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(logFile,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(logFile,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	DB = log.New(logFile,
		"DATABASe: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
