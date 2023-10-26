package models

import (
	"gorm.io/gorm"
)

type Actor struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Films     []Film `gorm:"many2many:film_actors"`
}
