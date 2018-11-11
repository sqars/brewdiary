package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"

	"github.com/sqars/brewdiary/app/models"
)

// NewBrewHandler is function constructor for Brew Handler
func NewBrewHandler(db *gorm.DB) *BrewHandler {
	// db.DropTableIfExists(&models.Brew{}, &models.Fermentation{}, &models.Composition{})
	db.AutoMigrate(&models.Brew{})
	return &BrewHandler{DB: db}
}

// BrewHandler is struct with api handlers for brew
type BrewHandler struct {
	DB *gorm.DB
}

// AddBrew adds Brew into database
func (b *BrewHandler) AddBrew(w http.ResponseWriter, r *http.Request) {
	tx := b.DB.Begin()
	brew := models.Brew{}
	err := decodeJSON(r, &brew)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	err = brew.Create(b.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	tx.Commit()
	responseJSON(w, http.StatusCreated, brew)
}

// GetBrew returns brew for specified id
func (b *BrewHandler) GetBrew(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	brew := models.Brew{ID: id}
	err = brew.Get(b.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	responseJSON(w, http.StatusOK, brew)
}

// DeleteBrew hadler for removing brew
func (b *BrewHandler) DeleteBrew(w http.ResponseWriter, r *http.Request) {
	tx := b.DB.Begin()
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	brew := models.Brew{ID: id}
	err = brew.Delete(b.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	tx.Commit()
	w.WriteHeader(http.StatusOK)
}

// UpdateBrew updates brew in db
func (b *BrewHandler) UpdateBrew(w http.ResponseWriter, r *http.Request) {
	tx := b.DB.Begin()
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	brew := models.Brew{}
	err = decodeJSON(r, &brew)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	brew.ID = id
	err = brew.Update(b.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = brew.Get(b.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	tx.Commit()
	responseJSON(w, http.StatusOK, brew)
}

// GetBrews returns all brews from database
func (b *BrewHandler) GetBrews(w http.ResponseWriter, r *http.Request) {
	tx := b.DB.Begin()
	brew := models.Brew{}
	brews, err := brew.GetAll(b.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tx.Commit()
	responseJSON(w, http.StatusOK, brews)
}
