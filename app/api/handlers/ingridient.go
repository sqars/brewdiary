package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sqars/brewdiary/app/models"
	"github.com/sqars/brewdiary/app/utils"
)

// NewIngridientHandler is function constructor for IngridientHandler
func NewIngridientHandler(db *gorm.DB) *IngridientHandler {
	db.DropTableIfExists(&models.Ingridient{})
	db.AutoMigrate(&models.Ingridient{})
	return &IngridientHandler{DB: db}
}

// IngridientHandler is struct with api handlers for ingridient
type IngridientHandler struct {
	DB *gorm.DB
}

// AddIngridient adds ingridient into database
func (i *IngridientHandler) AddIngridient(w http.ResponseWriter, r *http.Request) {
	ingridient := models.Ingridient{}
	err := utils.DecodeJSON(r, &ingridient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err := i.DB.Create(&ingridient).Error; err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetIngridient returns ingridient with specified id
func (i *IngridientHandler) GetIngridient(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	ingridient := models.Ingridient{}
	if err := i.DB.Find(&ingridient, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	utils.ResponseJSON(w, http.StatusOK, ingridient)
}

// GetIngridients returns all ingridients from database
func (i *IngridientHandler) GetIngridients(w http.ResponseWriter, r *http.Request) {
	var ingridients []models.Ingridient
	if err := i.DB.Find(&ingridients).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	utils.ResponseJSON(w, http.StatusOK, ingridients)
}
