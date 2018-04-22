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
	db.AutoMigrate(&models.Brew{}, &models.Composition{}, &models.Fermentation{})
	db.Set("gorm:auto_preload", true).Find(&models.Brew{})
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
	fermentation := models.Fermentation{}
	b.DB.Model(&brew).Preload("Ingridient").Find(&ingridients)
	b.DB.Model(&brew).Related(&fermentation)
	brew.Fermentation = fermentation
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
		fermentation := models.Fermentation{}
		tx.Where("brew_id = ?", brews[i].ID).Preload("Ingridient").Find(&ingridients)
		tx.Model(&brews[i]).Related(&fermentation)
		brews[i].Fermentation = fermentation
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
	if err := tx.Model(&brew).Updates(brewPUT).Error; err != nil {
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
	reqData := ingridientReq{}
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
	if err := tx.Model(&brew).Association("Ingridients").Append(composition).Error; err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	ingridients := []models.Composition{}
	tx.First(&brew, id).Preload("Ingridient").Find(&ingridients)
	brew.Ingridients = ingridients
	responseJSON(
		w,
		http.StatusOK,
		brew,
	)
	tx.Commit()
}

// AddFermentation adds fermentation specific information to brew
func (b *BrewHandler) AddFermentation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	fermentation := models.Fermentation{}
	err = decodeJSON(r, &fermentation)
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
	if err = tx.Model(&brew).Association("Fermentation").Append(fermentation).Error; err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	ingridients := []models.Composition{}
	tx.Model(&brew).Preload("Ingridient").Find(&ingridients)
	brew.Fermentation = fermentation
	brew.Ingridients = ingridients
	responseJSON(
		w,
		http.StatusOK,
		brew,
	)
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
