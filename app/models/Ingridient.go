package models

import (
	"errors"
	"time"
)

const (
	hop   = "hop"
	yeast = "yeast"
	malt  = "malt"
	other = "other"
)

// Ingridient - struct model for brew ingridient
type Ingridient struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name" gorm:"UNIQUE;NOT NULL"`
	Type      string    `json:"type" gorm:"NOT NULL"`
	Comments  string    `json:"comments"`
}

// ValidateType validates type of a ingridient
func (i *Ingridient) ValidateType() error {
	types := []string{hop, yeast, malt, other}
	for _, t := range types {
		if t == i.Type {
			return nil
		}
	}
	return errors.New("Invalid ingridient type")
}
