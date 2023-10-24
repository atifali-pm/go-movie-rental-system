package models

import (
	"gorm.io/gorm"
)

type Rental struct {
	gorm.Model
	StaffId     int    `json:"staff_id"`
	CustomerId  int    `json:"customer_id"`
	InventoryId int    `json:"inventory_id"`
	RentalDate  string `json:"rental_date"`
	RenturnDate string `json:"return_date"`
}
