package models

import "time"

type Actor struct {
	Id        int       `json:"id" gorm:"type:INT(10) UNSINGED NOT NULL AUTO_INCREMENT;primaryKey"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
