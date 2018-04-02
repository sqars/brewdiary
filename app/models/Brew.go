package models

import (
	"github.com/jinzhu/gorm"
)

// Brew - struct model for brew
type Brew struct {
	gorm.Model
	Name        string `json:"name" gorm:"UNIQUE;NOT NULL"`
	Num         int    `json:"num" gorm:"AUTO_INCREMENT;UNIQUE;NOT NULL"`
	Location    string `json:"location"`
	Comments    string `json:"comments"`
	Ingridients []Composition
}

// Composition - model for junction table for Brew and Ingridients
type Composition struct {
	ID           uint `gorm:"primary_key"`
	Quantity     int  `json:"quantity"`
	BrewID       uint
	Ingridient   Ingridient
	IngridientID uint
}
