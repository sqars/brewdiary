package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Ingridient - struct model for ingridient
type Ingridient struct {
	ID        int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name" gorm:"UNIQUE;NOT NULL"`
	Comments  string    `json:"comments"`
}

// Get returns ingridient from db
func (i *Ingridient) Get(db *gorm.DB) error {
	err := db.Find(i, i.ID).Error
	return err
}

// Create creates new ingridient in DB
func (i *Ingridient) Create(db *gorm.DB) error {
	err := db.Create(i).Error
	return err
}

// Delete removes ingridient from db
func (i *Ingridient) Delete(db *gorm.DB) error {
	err := db.Delete(i).Error
	return err
}

// Update updates ingridient in db
func (i *Ingridient) Update(db *gorm.DB) error {
	err := db.Model(&i).Updates(Ingridient{
		Name:     i.Name,
		Comments: i.Comments,
	}).Error
	return err
}

// GetAll returns all ingridients from db
func (i *Ingridient) GetAll(db *gorm.DB) ([]Ingridient, error) {
	ingridients := []Ingridient{}
	err := db.Find(&ingridients).Error
	if err != nil {
		return nil, err
	}
	return ingridients, nil
}

// OK validates if ingridient valeus are correct
func (i *Ingridient) OK() error {
	if len(i.Name) == 0 {
		return ErrMissingField("name")
	}
	return nil
}
