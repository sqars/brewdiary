package api

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/sqars/brewdiary/app/api/handlers"
)

// Init initialize api and register app routes
func Init(db *gorm.DB) *mux.Router {
	brewHandler := handlers.NewBrewHandler(db)
	r := mux.NewRouter()
	s := r.PathPrefix("/brew").Subrouter()

	s.Methods("POST").HandlerFunc(brewHandler.AddBrew)

	return r
}
