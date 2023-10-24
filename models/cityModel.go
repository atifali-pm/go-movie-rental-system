package models

import "gorm.io/gorm"

type City struct {
	gorm.Model
	CountryId uint   `json:"country_id"`
	City      string `json:"city"`
}
