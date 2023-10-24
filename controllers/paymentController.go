package controllers

import (
	"log"
	"time"

	db "github.com/atifali-pm/go-movie-rental-system/config"
	"github.com/atifali-pm/go-movie-rental-system/models"
	"github.com/gofiber/fiber/v2"
)

type PaymentData struct {
	Rental      RentalData `json:"rental"`
	CustomerId  int        `json:"customer_id"`
	StaffId     int        `json:"staff_id"`
	Amount      float32    `json:"amount"`
	PaymentDate string     `json:"payment_date"`
}

type RentalData struct {
	CustomerId  int    `json:"customer_id"`
	StaffId     int    `json:"staff_id"`
	InventoryId int    `json:"inventory_id"`
	RentalDate  string `json:"rental_date"`
	ReturnDate  string `json:"return_date"`
}

func MakePayment(c *fiber.Ctx) error {
	var data PaymentData

	err := c.BodyParser(&data)
	if err != nil {
		log.Fatalf("Product error in post request %v", err)
	}

	if data.Amount <= 0 || data.PaymentDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": true,
			"message": "Fields are required",
		})
	}

	rental := models.Rental{
		StaffId:     data.Rental.StaffId,
		CustomerId:  data.Rental.CustomerId,
		InventoryId: data.Rental.InventoryId,
		RentalDate:  data.Rental.RentalDate,
		RenturnDate: data.Rental.ReturnDate,
	}
	db.DB.Create(&rental)

	payment := models.Payment{
		RentalId:    int(rental.ID),
		CustomerId:  rental.CustomerId,
		StaffId:     rental.StaffId,
		Amount:      data.Amount,
		PaymentDate: time.Now().UTC(),
	}
	db.DB.Create(&payment)

	return c.Status(200).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"success": true,
		"message": "Payment done!",
		"data":    data,
	})
}
