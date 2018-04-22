package models

import "time"

// Fermentation struct model for fermentation
// one to one relationship with Brew
type Fermentation struct {
	FermentationID       int
	ID                   uint      `json:"id" gorm:"primary_key"`
	WortVolume           int       `json:"WortVolume"`
	TurbulantDateStart   time.Time `json:"turbulantDateStart"`
	YeastApplicationTemp int       `json:"yeastApplicationTemp"`
	InitialDensity       int       `json:"initialDensity"`
	QuietDateStart       time.Time `json:"quietDateStart"`
	QuietDateEnd         time.Time `json:"quietDateEnd"`
	QuietInitialDensity  int       `json:"quietInitialDensity"`
	QuietInitialTemp     int       `json:"quietInitialTemp"`
	FinalDensity         int       `json:"finalDensity"`
	Brew                 Brew
}
