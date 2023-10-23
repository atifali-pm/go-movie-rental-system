package models

import "time"

type Film struct {
	Id              int       `json:"id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	LanguageId      int       `json:"language_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	ReleaseYear     int       `json:"release_year"`
	Rental_Duration int       `json:"rental_duration"`
	RentalRate      float32   `json:"rental_rate"`
	Length          int       `json:"length"`
	ReplacementCost float32   `json:"replacement_cost"`
	Rating          int       `json:"rating"`
	SpecialFeature  string    `json:"special_feature"`
	FullText        string    `json:"full_text"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
