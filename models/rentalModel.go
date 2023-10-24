package models

import (
	"time"

	"gorm.io/gorm"
)

type Rental struct {
	gorm.Model
	StaffId     int       `json:"staff_id"`
	CustomerId  int       `json:"customer_id"`
	InventoryId int       `json:"inventory_id"`
	RentalDate  time.Time `json:"rental_date"`
	RenturnDate time.Time `json:"return_date"`
}
