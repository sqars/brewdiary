package models

import (
	"time"
)

// Brew - struct model for brew
type Brew struct {
	ID             uint          `json:"id" gorm:"primary_key"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
	Name           string        `json:"name" gorm:"UNIQUE;NOT NULL"`
	Location       string        `json:"location"`
	Comments       string        `json:"comments"`
	Ingridients    []Composition `json:"ingridients"`
	Fermentation   Fermentation  `json:"fermentation"`
	FermentationID uint          `json:"-"`
}
