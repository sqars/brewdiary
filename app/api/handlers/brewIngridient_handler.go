package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sqars/brewdiary/app/models"
)

// NewBrewIngridientHandler is function constructor for BrewIngridient Handler
func NewBrewIngridientHandler(db *gorm.DB) *BrewIngridientHandler {
	db.AutoMigrate(&models.BrewIngridient{})
	return &BrewIngridientHandler{DB: db}
}

// BrewIngridientHandler is struct with api handlers for brew ingridient
type BrewIngridientHandler struct {
	DB *gorm.DB
}

// AddBrewIngridient adds ingridient to brew in database
func (b *BrewIngridientHandler) AddBrewIngridient(w http.ResponseWriter, r *http.Request) {
	tx := b.DB.Begin()
	brewID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	// get brew from db
	brew := models.Brew{ID: brewID}
	err = brew.Get(b.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	// decode brewIngridient payload
	brewIngridient := models.BrewIngridient{}
	err = decodeJSON(r, &brewIngridient, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	brewIngridient.BrewID = brewID
	// check if ingridient specified in payload exists
	ingridient := models.Ingridient{ID: brewIngridient.IngridientID}
	err = ingridient.Get(b.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// create brew ingridient
	brewIngridient.Ingridient = ingridient
	err = brewIngridient.Create(b.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// update brew with new ingridient
	brew.Ingridients = append(brew.Ingridients, brewIngridient)
	err = brew.Update(b.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tx.Commit()
	responseJSON(w, http.StatusCreated, brewIngridient)
}

// GetBrewIngridient returns brew ingridient brew from database
func (b *BrewIngridientHandler) GetBrewIngridient(w http.ResponseWriter, r *http.Request) {
	tx := b.DB.Begin()
	brewID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	iID, err := strconv.Atoi(mux.Vars(r)["iid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	bi := models.BrewIngridient{
		BrewID: brewID,
		ID:     iID,
	}
	err = bi.Get(b.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	tx.Commit()
	responseJSON(w, http.StatusOK, bi)
}

// DeleteBrewIngridient remove brewIngridient from db for brew
func (b *BrewIngridientHandler) DeleteBrewIngridient(w http.ResponseWriter, r *http.Request) {
	tx := b.DB.Begin()
	brewID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	iID, err := strconv.Atoi(mux.Vars(r)["iid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	bi := models.BrewIngridient{
		BrewID: brewID,
		ID:     iID,
	}
	err = bi.Delete(b.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	tx.Commit()
	responseJSON(w, http.StatusOK, bi)
}

// UpdateBrewIngridient remove brewIngridient from db for brew
func (b *BrewIngridientHandler) UpdateBrewIngridient(w http.ResponseWriter, r *http.Request) {
	tx := b.DB.Begin()
	brewID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	iID, err := strconv.Atoi(mux.Vars(r)["iid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	// decode brewIngridient payload
	bi := models.BrewIngridient{}
	err = decodeJSON(r, &bi, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bi.BrewID = brewID
	bi.ID = iID
	err = bi.Update(b.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = bi.Get(b.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	tx.Commit()
	responseJSON(w, http.StatusOK, bi)
}
