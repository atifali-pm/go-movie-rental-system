package controllers

import (
	"log"

	db "github.com/atifali-pm/go-movie-rental-system/config"
	"github.com/atifali-pm/go-movie-rental-system/models"
	"github.com/gofiber/fiber/v2"
)

type CustomerData struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Active    bool    `json:"active"`
	Address   Address `json:"address"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
}

type Address struct {
	Address    string `json:"address"`
	Address2   string `json:"address2"`
	District   string `json:"district"`
	PostalCode string `json:"postal_code"`
	Phone      string `json:"phone"`
}

func CreateCustomer(c *fiber.Ctx) error {
	var data CustomerData

	err := c.BodyParser(&data)
	if err != nil {
		log.Fatalf("Customer not created %v", err)
	}

	if data.FirstName == "" || data.LastName == "" || data.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Fields are required",
		})
	}

	var existingEmail models.Customer
	if err := db.DB.Where("email = ?", data.Email).First(&existingEmail).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email already exists!",
		})
	}

	var country models.Country
	if err := db.DB.Where("country = ?", data.Country).Find(&country).Error; err != nil {
		if err := db.DB.Save(&country).Error; err != nil {
			return c.Status(500).SendString("Error while saving the language")
		}
	}

	var city models.City
	if err := db.DB.Where("city = ?", data.City).Find(&city).Error; err != nil {
		city.City = data.City
		city.CountryId = country.ID
		if err := db.DB.Save(&city).Error; err != nil {
			return c.Status(500).SendString("Error while saving the language")
		}
	}

	address := models.Address{
		CityId:     int(city.ID),
		Address:    data.Address.Address,
		Address2:   data.Address.Address2,
		District:   data.Address.District,
		PostalCode: data.Address.PostalCode,
		Phone:      data.Address.Phone,
	}
	db.DB.Create(&address)

	customer := models.Customer{
		AddressId: int(address.ID),
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Active:    data.Active,
		Email:     data.Email,
	}
	db.DB.Create(&customer)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    customer,
	})

}
