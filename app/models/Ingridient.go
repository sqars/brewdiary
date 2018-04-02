package models

import (
	"github.com/jinzhu/gorm"
)

const (
	hop   = "yeast"
	yeast = "yeast"
	malt  = "malt"
	other = "other"
)

// Ingridient - struct model for brew ingridient
type Ingridient struct {
	gorm.Model
	Name     string `json:"name" gorm:"UNIQUE;NOT NULL"`
	Type     string `json:"type" gorm:"NOT NULL"`
	Comments string `json:"comments"`
}

// // ValidateType validates type of a ingridient
// func (i *Ingridient) ValidateType() error {

// }
