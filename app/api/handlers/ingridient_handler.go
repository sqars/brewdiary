package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"

	"github.com/sqars/brewdiary/app/models"
)

// NewIngridientHandler is function constructor for Ingridient Handler
func NewIngridientHandler(db *gorm.DB) *IngridientHandler {
	db.AutoMigrate(&models.Ingridient{})
	return &IngridientHandler{DB: db}
}

// IngridientHandler is struct with api handlers for ingridient
type IngridientHandler struct {
	DB *gorm.DB
}

// AddIngridient adds Ingridient into database
func (i *IngridientHandler) AddIngridient(w http.ResponseWriter, r *http.Request) {
	tx := i.DB.Begin()
	ingridient := models.Ingridient{}
	err := decodeJSON(r, &ingridient, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	err = ingridient.Create(i.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	tx.Commit()
	responseJSON(w, http.StatusCreated, ingridient)
}

// GetIngridient returns ingridient for specified id
func (i *IngridientHandler) GetIngridient(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	ingridient := models.Ingridient{ID: id}
	err = ingridient.Get(i.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	responseJSON(w, http.StatusOK, ingridient)
}

// DeleteIngridient hadler for removing ingridient
func (i *IngridientHandler) DeleteIngridient(w http.ResponseWriter, r *http.Request) {
	tx := i.DB.Begin()
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	ingridient := models.Ingridient{ID: id}
	err = ingridient.Delete(i.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	tx.Commit()
	w.WriteHeader(http.StatusOK)
}

// UpdateIngridient updates ingridient in db
func (i *IngridientHandler) UpdateIngridient(w http.ResponseWriter, r *http.Request) {
	tx := i.DB.Begin()
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	ingridient := models.Ingridient{}
	err = decodeJSON(r, &ingridient, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	ingridient.ID = id
	err = ingridient.Update(i.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = ingridient.Get(i.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	tx.Commit()
	responseJSON(w, http.StatusOK, ingridient)
}

// GetIngridients returns all ingridients from database
func (i *IngridientHandler) GetIngridients(w http.ResponseWriter, r *http.Request) {
	tx := i.DB.Begin()
	ingridient := models.Ingridient{}
	ingridients, err := ingridient.GetAll(i.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tx.Commit()
	responseJSON(w, http.StatusOK, ingridients)
}
