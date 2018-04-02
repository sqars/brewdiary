package models

// Composition - model for junction table for Brew and Ingridients
type Composition struct {
	ID           uint       `json:"id" gorm:"primary_key"`
	Quantity     int        `json:"quantity"`
	BrewID       uint       `json:"-"`
	Ingridient   Ingridient `json:"ingridient"`
	IngridientID uint       `json:"-"`
}
