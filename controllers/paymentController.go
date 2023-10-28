package controllers

import (
	"log"
	"time"

	db "github.com/atifali-pm/go-movie-rental-system/config"
	"github.com/atifali-pm/go-movie-rental-system/models"
	"github.com/gofiber/fiber/v2"
)

type PaymentData struct {
	Rental      RentalData       `json:"rental"`
	Customer    CustomerResponse `json:"customer"`
	Staff       StaffData        `json:"staff"`
	Amount      float32          `json:"amount"`
	PaymentDate string           `json:"payment_date"`
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

	if data.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": true,
			"message": "Fields are required",
		})
	}

	if data.Rental.InventoryId <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": true,
			"message": "Fields are required",
		})
	}

	// Check if a film exists in inventory
	var existingInventory models.Inventory
	db.DB.Where("id = ?", data.Rental.InventoryId).First(&existingInventory)

	if existingInventory.Copies <= 1 {
		return c.Status(409).SendString("All copies are rented!")
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
		PaymentDate: data.PaymentDate,
	}
	db.DB.Create(&payment)

	existingInventory.Copies -= 1
	db.DB.Save(&existingInventory)

	return c.Status(200).JSON(fiber.Map{
		"status":    fiber.StatusOK,
		"success":   true,
		"message":   "Payment done!",
		"data":      data,
		"Inventory": existingInventory,
	})
}

func ReturnFilm(c *fiber.Ctx) error {
	body := struct {
		StaffId     int `json:"staff_id"`
		CustomerId  int `json:"customer_id"`
		InventoryId int `json:"inventory_id"`
	}{}

	err := c.BodyParser(&body)

	if err != nil {
		log.Fatalf("File can not be returned %v", err)
	}

	if body.CustomerId <= 0 || body.StaffId <= 0 || body.InventoryId <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Missing required fields",
		})
	}

	var inventory models.Inventory
	db.DB.Where("id = ?", body.InventoryId).First(&inventory)
	if inventory.Copies >= 5 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "All rented copied returned, no more returns needed!",
		})
	}

	var rental models.Rental
	db.DB.Where("customer_id = ?", body.CustomerId).Where("staff_id = ?", body.StaffId).Where("inventory_id = ?", body.InventoryId).First(&rental)

	rental.RenturnDate = time.Now().UTC().String()
	db.DB.Create(&rental)

	inventory.Copies += 1
	db.DB.Save(&inventory)

	return c.Status(200).JSON(fiber.Map{
		"success":   true,
		"messsage":  "success",
		"data":      rental,
		"inventory": inventory,
	})

}

func PaymentsListByCustomer(c *fiber.Ctx) error {

	var paymentList []models.Payment

	limit, page := 10, 1
	offset := (page - 1) * limit
	db.DB.Limit(limit).Offset(offset)
	customerId := c.Params("customer_id")
	if customerId == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"messgae": "Wrong payload!",
		})
	}

	if err := db.DB.Where("payments.customer_id", customerId).Order("payments.created_at DESC").Find(&paymentList).Error; err != nil {
		return c.Status(500).SendString("Error while fetching films")
	}

	var PaymentResponses []PaymentData

	for _, payment := range paymentList {

		var customer models.Customer
		db.DB.Where("id=?", payment.CustomerId).Find(&customer)

		var rental models.Rental
		db.DB.Where("id=?", payment.RentalId).Find(&rental)

		var staff models.Staff
		db.DB.Where("id=?", payment.StaffId).Find(&staff)

		PaymentResponses = append(PaymentResponses, PaymentData{
			Rental: RentalData{
				InventoryId: rental.InventoryId,
				RentalDate:  rental.RentalDate,
				ReturnDate:  rental.RenturnDate,
			},
			Customer: CustomerResponse{
				FirstName: customer.FirstName,
				LastName:  customer.LastName,
				Email:     customer.Email,
				Active:    customer.Active,
			},
			Staff: StaffData{
				FirstName:  staff.FirstName,
				LastName:   staff.LastName,
				Email:      staff.Email,
				Active:     staff.Active,
				Username:   staff.Username,
				PictureURL: staff.PictureURL,
			},
			Amount:      payment.Amount,
			PaymentDate: payment.PaymentDate,
		})

	}

	return c.Status(200).JSON(fiber.Map{
		"success":  true,
		"messsage": "success",
		"data":     PaymentResponses,
	})
}
