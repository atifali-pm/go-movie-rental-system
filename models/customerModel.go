package models

import "gorm.io/gorm"

type Customer struct {
	gorm.Model

	AddressId int     `json:"address_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Active    bool    `json:"active"`
	Address   Address `json:"address"`
}
