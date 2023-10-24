package models

import "gorm.io/gorm"

type Store struct {
	gorm.Model
	AddressId int `json:"address_id"`
}
