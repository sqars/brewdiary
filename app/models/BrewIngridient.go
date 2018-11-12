package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// BrewIngridient - struct model for brew ingridient
type BrewIngridient struct {
	ID           int        `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	Quantity     int        `json:"quantity" gorm:"NOT_NULL"`
	BrewID       int        `json:"brewId"`
	Ingridient   Ingridient `json:"ingridient" gorm:"foreignkey:IngridientID"`
	IngridientID int        `json:"ingridientId"`
}

// Get returns brew ingridient from db
func (b *BrewIngridient) Get(db *gorm.DB) error {
	err := db.Find(b, b.ID).Error
	return err
}

// Create creates new brew ingridient in DB
func (b *BrewIngridient) Create(db *gorm.DB) error {
	err := db.Create(b).Error
	return err
}

// Delete removes brew ingridient from db
func (b *BrewIngridient) Delete(db *gorm.DB) error {
	err := db.Delete(b).Error
	return err
}

// Update updates brew ingridient in db
func (b *BrewIngridient) Update(db *gorm.DB) error {
	err := db.Model(&b).Updates(BrewIngridient{
		Quantity: b.Quantity,
	}).Error
	return err
}

// GetAll returns all brew ingridient from db
func (b *BrewIngridient) GetAll(db *gorm.DB) ([]BrewIngridient, error) {
	ingridients := []BrewIngridient{}
	err := db.Find(&ingridients).Error
	if err != nil {
		return nil, err
	}
	return ingridients, nil
}

// OK validates if brew ingridient valeus are correct
func (b *BrewIngridient) OK() error {
	if b.Quantity == 0 {
		return ErrMissingField("quantity")
	}
	if b.IngridientID == 0 {
		return ErrMissingField("ingridientId")
	}
	return nil
}
