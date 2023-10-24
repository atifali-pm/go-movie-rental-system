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

	country := models.Country{
		Country: data.Country,
	}
	db.DB.Create(&country)

	city := models.City{
		CountryId: country.ID,
		City:      data.City,
	}
	db.DB.Create(&city)

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
