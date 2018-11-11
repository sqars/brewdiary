package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Brew - struct model for brew
type Brew struct {
	ID        int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name" gorm:"UNIQUE;NOT NULL"`
	Location  string    `json:"location"`
	Comments  string    `json:"comments"`
}

// GetBrew returns brew from db
func (b *Brew) GetBrew(db *gorm.DB) error {
	err := db.Find(b, b.ID).Error
	return err
}

// CreateBrew creates new brew in DB
func (b *Brew) CreateBrew(db *gorm.DB) error {
	err := db.Create(b).Error
	return err
}

// DeleteBrew removes brew from db
func (b *Brew) DeleteBrew(db *gorm.DB) error {
	err := db.Delete(b).Error
	return err
}

// UpdateBrew updates brew in db
func (b *Brew) UpdateBrew(db *gorm.DB) error {
	err := db.Model(&b).Updates(Brew{
		Name:     b.Name,
		Location: b.Location,
		Comments: b.Comments,
	}).Error
	return err
}

// GetBrews returns all brews from db
func (b *Brew) GetBrews(db *gorm.DB) ([]Brew, error) {
	brews := []Brew{}
	err := db.Find(&brews).Error
	if err != nil {
		return nil, err
	}
	return brews, nil
}
