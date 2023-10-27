package models

import "gorm.io/gorm"

type Staff struct {
	gorm.Model
	AddressId  int    `json:"address_id"`
	StoreId    int    `json:"store_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Active     bool   `json:"active"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PictureURL string `json:"picture_url"`
}
