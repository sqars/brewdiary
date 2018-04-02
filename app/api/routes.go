package api

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/sqars/brewdiary/app/api/handlers"
)

// Init initialize api and register app routes
func Init(db *gorm.DB) *mux.Router {
	brewHandler := handlers.NewBrewHandler(db)
	ingridientHandler := handlers.NewIngridientHandler(db)
	r := mux.NewRouter()

	// brew specific routes
	r.HandleFunc("/brew", brewHandler.GetBrews).Methods("GET")
	r.HandleFunc("/brew/{id:[0-9]+}", brewHandler.GetBrew).Methods("GET")
	r.HandleFunc("/brew", brewHandler.AddBrew).Methods("POST")
	r.HandleFunc("/brew/{id:[0-9]+}", brewHandler.DeleteBrew).Methods("DELETE")
	r.HandleFunc("/brew/{id:[0-9]+}", brewHandler.UpdateBrew).Methods("PUT")
	r.HandleFunc("/brew/{id:[0-9]+}/ingridients", brewHandler.AddIngridient).Methods("PUT")
	r.HandleFunc("/brew/{id:[0-9]+}/ingridients/{cid:[0-9]+}", brewHandler.DeleteIngridient).Methods("DELETE")

	// ingridient specific routes
	r.HandleFunc("/ingridient", ingridientHandler.GetIngridients).Methods("GET")
	r.HandleFunc("/ingridient/{id:[0-9]+}", ingridientHandler.GetIngridient).Methods("GET")
	r.HandleFunc("/ingridient", ingridientHandler.AddIngridient).Methods("POST")
	r.HandleFunc("/ingridient/{id:[0-9]+}", ingridientHandler.UpdateIngridient).Methods("PUT")
	r.HandleFunc("/ingridient/{id:[0-9]+}", ingridientHandler.DeleteIngridient).Methods("DELETE")

	return r
}
