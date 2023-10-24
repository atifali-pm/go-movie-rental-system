package models

import (
	"gorm.io/gorm"
)

type FilmCategory struct {
	gorm.Model     // This includes ID, CreatedAt, and UpdatedAt fields
	FilmId     int `json:"film_id"`
	CategoryId int `json:"category_id"`
}
