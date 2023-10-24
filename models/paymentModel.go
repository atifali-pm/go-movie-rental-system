package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model

	RentalId    int       `json:"rental_id"`
	CustomerId  int       `json:"customer_id"`
	StaffId     int       `json:"staff_id"`
	Amount      float32   `json:"amount"`
	PaymentDate time.Time `json:"payment_time"`
}
