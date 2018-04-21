package logger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
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
	// HTTP logs http traffic
	HTTP *HTTPLogger
)

// HTTPLogger struct embeding log.Logger
type HTTPLogger struct {
	logger *log.Logger
}

func newHTTPLogger(l *log.Logger) *HTTPLogger {
	return &HTTPLogger{logger: l}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// LogRequest logs request info
func (l *HTTPLogger) logRequest(r http.Request) {
	l.logger.Println(
		fmt.Sprintf("--> Request: %s, %s, %s", r.Method, r.URL, r.Proto),
	)
}

func (l *HTTPLogger) logResponse(code int) {
	l.logger.Println(
		fmt.Sprintf("<-- Response: %s, %s", strconv.Itoa(code), http.StatusText(code)),
	)
}

// LogTraffic logs incoming request info and outcoming response ifno
func LogTraffic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := loggingResponseWriter{w, http.StatusOK}
		HTTP.logRequest(*r)
		h.ServeHTTP(&lrw, r)
		HTTP.logResponse(lrw.status)
	})
}

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
		"DATABASE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	HTTP = newHTTPLogger(
		log.New(logFile,
			"",
			log.Ldate|log.Ltime),
	)
}
