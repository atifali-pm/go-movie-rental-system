package models

import (
	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	FilmId int `json:"film_id"`
}
