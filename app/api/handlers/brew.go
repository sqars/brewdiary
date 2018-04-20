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
	db.DropTableIfExists(&models.Brew{}, &models.Composition{})
	db.AutoMigrate(&models.Brew{}, &models.Composition{})
	return &BrewHandler{DB: db}
}

// BrewHandler is struct with api handlers for brew
type BrewHandler struct {
	DB *gorm.DB
}

// AddBrew adds Brew into database
func (b *BrewHandler) AddBrew(w http.ResponseWriter, r *http.Request) {
	brew := models.Brew{}
	err := decodeJSON(r, &brew)
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

// GetBrew returns brew for specified id
func (b *BrewHandler) GetBrew(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	brew := models.Brew{}
	if err := b.DB.Find(&brew, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	ingridients := []models.Composition{}
	b.DB.Where("brew_id = ?", id).Preload("Ingridient").Find(&ingridients)
	brew.Ingridients = ingridients
	responseJSON(w, http.StatusOK, brew)
}

// GetBrews returns all brews from database
func (b *BrewHandler) GetBrews(w http.ResponseWriter, r *http.Request) {
	tx := b.DB.Begin()
	var brews []models.Brew
	if err := tx.Find(&brews).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	for i := range brews {
		ingridients := []models.Composition{}
		tx.Where("brew_id = ?", brews[i].ID).Preload("Ingridient").Find(&ingridients)
		brews[i].Ingridients = ingridients
	}
	responseJSON(w, http.StatusOK, brews)
}

// UpdateBrew updates brew with specified id
func (b *BrewHandler) UpdateBrew(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	brewPUT := models.Brew{}
	err = decodeJSON(r, &brewPUT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	tx := b.DB.Begin()
	brew := models.Brew{}
	if err := tx.First(&brew, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	brew.Name = brewPUT.Name
	brew.Comments = brewPUT.Comments
	brew.Location = brewPUT.Location
	if err := tx.Save(&brew).Error; err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	responseJSON(w, http.StatusOK, brew)
	tx.Commit()
}

// DeleteBrew deletes brew with specified id
func (b *BrewHandler) DeleteBrew(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	b.DB.Where("id = ?", id).Delete(&models.Brew{})
	if err := b.DB.Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// AddIngridient adds Ingridient into brew Ingridients
func (b *BrewHandler) AddIngridient(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	reqData := ingridientRequest{}
	err = decodeJSON(r, &reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	tx := b.DB.Begin()
	ingridient := models.Ingridient{}
	if err := tx.Find(&ingridient, reqData.IngridientID).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	brew := models.Brew{}
	if err := tx.First(&brew, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	composition := models.Composition{
		Quantity:   reqData.Quantity,
		Ingridient: ingridient,
	}
	ingridients := []models.Composition{}
	tx.Where("brew_id = ?", id).Preload("Ingridient").Find(&ingridients)
	brew.Ingridients = append(ingridients, composition)
	if err := tx.Save(&brew).Error; err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	responseJSON(w, http.StatusOK, tx.Model(&brew).Related(&models.Composition{}).Value)
	tx.Commit()
}

// DeleteIngridient removes ingridient from brew
func (b *BrewHandler) DeleteIngridient(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	cid, err := strconv.Atoi(mux.Vars(r)["cid"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	b.DB.Where("id = ? AND brew_id = ?", cid, id).Delete(&models.Composition{})
	if err = b.DB.Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
