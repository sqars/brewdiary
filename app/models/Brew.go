package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Brew - struct model for brew
type Brew struct {
	ID          int              `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
	Name        string           `json:"name" gorm:"UNIQUE;NOT NULL"`
	Location    string           `json:"location"`
	Comments    string           `json:"comments"`
	Ingridients []BrewIngridient `json:"ingridients"`
}

// Get returns brew from db
func (b *Brew) Get(db *gorm.DB) error {
	err := db.Find(b, b.ID).Error
	return err
}

// Create creates new brew in DB
func (b *Brew) Create(db *gorm.DB) error {
	err := db.Create(b).Error
	return err
}

// Delete removes brew from db
func (b *Brew) Delete(db *gorm.DB) error {
	err := db.Delete(b).Error
	return err
}

// Update updates brew in db
func (b *Brew) Update(db *gorm.DB) error {
	err := db.Model(&b).Updates(Brew{
		Name:        b.Name,
		Location:    b.Location,
		Comments:    b.Comments,
		Ingridients: b.Ingridients,
	}).Error
	return err
}

// GetAll returns all brews from db
func (b *Brew) GetAll(db *gorm.DB) ([]Brew, error) {
	brews := []Brew{}
	err := db.Find(&brews).Error
	if err != nil {
		return nil, err
	}
	return brews, nil
}

// OK validates if brew valeus are correct
func (b *Brew) OK() error {
	if len(b.Name) == 0 {
		return ErrMissingField("name")
	}
	return nil
}
