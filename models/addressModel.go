package models

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	CityId     int    `json:"city_id"`
	Address    string `json:"address"`
	Address2   string `json:"address_2"`
	District   string `json:"district"`
	PostalCode string `json:"postal_code"`
	Phone      string `json:"phone"`
}
