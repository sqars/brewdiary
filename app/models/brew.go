package models

import (
	"github.com/jinzhu/gorm"
)

// Brew - struct model for brew
type Brew struct {
	gorm.Model
	Name     string `json:"name" gorm:"UNIQUE;NOT NULL"`
	Num      int    `json:"num" gorm:"AUTO_INCREMENT;UNIQUE;NOT NULL"`
	Comments string `json:"comments"`
}
