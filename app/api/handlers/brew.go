package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/sqars/brewdiary/app/models"
)

// NewBrewHandler is function constructor for Brew Handler
func NewBrewHandler(db *gorm.DB) *BrewHandler {
	db.AutoMigrate(&models.Brew{})
	return &BrewHandler{DB: db}
}

// BrewHandler is struct with api handlers for brew
type BrewHandler struct {
	DB *gorm.DB
}

// AddBrew adds Brew into database
func (b *BrewHandler) AddBrew(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	brew := models.Brew{}
	err := decoder.Decode(&brew)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err := b.DB.Create(&brew).Error; err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
