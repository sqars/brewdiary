package api

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/sqars/brewdiary/app/api/handlers"
)

// Init initialize api and register app routes
func Init(db *gorm.DB) *mux.Router {
	brewHandler := handlers.NewBrewHandler(db)
	brewIngridientHandler := handlers.NewBrewIngridientHandler(db)
	ingridientHandler := handlers.NewIngridientHandler(db)
	r := mux.NewRouter()

	// brew specific routes
	r.HandleFunc("/brew", brewHandler.GetBrews).Methods("GET")
	r.HandleFunc("/brew/{id:[0-9]+}", brewHandler.GetBrew).Methods("GET")
	r.HandleFunc("/brew", brewHandler.AddBrew).Methods("POST")
	r.HandleFunc("/brew/{id:[0-9]+}", brewHandler.DeleteBrew).Methods("DELETE")
	r.HandleFunc("/brew/{id:[0-9]+}", brewHandler.UpdateBrew).Methods("PATCH")
	// brew ingridient routes
	r.HandleFunc("/brew/{id:[0-9]+}/ingr", brewIngridientHandler.AddBrewIngridient).Methods("POST")
	r.HandleFunc("/brew/{id:[0-9]+}/ingr/{iid:[0-9]+}", brewIngridientHandler.GetBrewIngridient).Methods("GET")
	r.HandleFunc("/brew/{id:[0-9]+}/ingr/{iid:[0-9]+}", brewIngridientHandler.DeleteBrewIngridient).Methods("DELETE")
	r.HandleFunc("/brew/{id:[0-9]+}/ingr/{iid:[0-9]+}", brewIngridientHandler.UpdateBrewIngridient).Methods("PATCH")

	// ingridient specific routes
	r.HandleFunc("/ingr", ingridientHandler.GetIngridients).Methods("GET")
	r.HandleFunc("/ingr/{id:[0-9]+}", ingridientHandler.GetIngridient).Methods("GET")
	r.HandleFunc("/ingr", ingridientHandler.AddIngridient).Methods("POST")
	r.HandleFunc("/ingr/{id:[0-9]+}", ingridientHandler.DeleteIngridient).Methods("DELETE")
	r.HandleFunc("/ingr/{id:[0-9]+}", ingridientHandler.UpdateIngridient).Methods("PATCH")

	return r
}
