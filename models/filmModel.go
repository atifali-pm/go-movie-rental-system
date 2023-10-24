package models

import (
	"gorm.io/gorm"
)

type Film struct {
	Id              int     `json:"id"`
	gorm.Model              // This includes ID, CreatedAt, and UpdatedAt fields
	LanguageId      int     `json:"language_id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	ReleaseYear     int     `json:"release_year"`
	Rental_Duration int     `json:"rental_duration"`
	RentalRate      float32 `json:"rental_rate"`
	Length          int     `json:"length"`
	ReplacementCost float32 `json:"replacement_cost"`
	Rating          int     `json:"rating"`
	SpecialFeature  string  `json:"special_feature"`
	FullText        string  `json:"full_text"`

	Actors     []Actor    `gorm:"many2many:film_actors;"`
	Categories []Category `gorm:"many2many:film_categories"`
}
