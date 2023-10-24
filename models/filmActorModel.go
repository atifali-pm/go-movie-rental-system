package models

import (
	"gorm.io/gorm"
)

type FilmActor struct {
	gorm.Model     // This includes ID, CreatedAt, and UpdatedAt fields
	FilmId     int `json:"film_id"`
	ActorId    int `json:"actor_id"`
}
