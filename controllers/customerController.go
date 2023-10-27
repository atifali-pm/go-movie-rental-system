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

type CustomerResponse struct {
	ID        int             `json:"id"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Email     string          `json:"email"`
	Active    bool            `json:"active"`
	Address   AddressResponse `json:"address"`
}

type AddressResponse struct {
	Address    string `json:"address"`
	Address2   string `json:"address2"`
	District   string `json:"district"`
	PostalCode string `json:"postal_code"`
	Phone      string `json:"phone"`
	City       string `json:"city"`
	Country    string `json:"country"`
}

func UpdateCustomer(c *fiber.Ctx) error {

	customerId := c.Params("id")
	var customer models.Customer

	db.DB.Find(&customer, "id=?", customerId)

	if customer.ID <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": true,
			"message": "Customer not found",
		})
	}

	var data CustomerData
	c.BodyParser(&data)

	var country models.Country
	if err := db.DB.Where("country = ?", data.Country).Find(&country).Error; err == nil {
		country.Country = data.Country
		if err := db.DB.Save(&country).Error; err != nil {
			return c.Status(500).SendString("Error while saving the country")
		}
	}

	log.Printf("Country name from DB: %s\n Country Name from Payload: %s", country.Country, data.Country)

	var city models.City
	if err := db.DB.Where("city = ?", data.City).Find(&city).Error; err == nil {
		city.City = data.City
		city.CountryId = country.ID
		if err := db.DB.Save(&city).Error; err != nil {
			return c.Status(500).SendString("Error while saving the city")
		}
	}

	var address models.Address
	db.DB.Where("id = ?", customer.AddressId).Find(&address)
	address.CityId = int(city.ID)
	address.Address = data.Address.Address
	address.Address2 = data.Address.Address2
	address.District = data.Address.District
	address.PostalCode = data.Address.PostalCode
	address.Phone = data.Address.Phone
	if err := db.DB.Save(&address).Error; err != nil {
		return c.Status(500).SendString("Error while saving the address")
	}

	customer.FirstName = data.FirstName
	customer.LastName = data.LastName
	customer.Active = data.Active
	customer.Email = data.Email
	db.DB.Save(&customer)

	return c.Status(201).JSON(fiber.Map{
		"success":  true,
		"messsage": "success",
		"data":     data,
	})
}

func CustomerDetail(c *fiber.Ctx) error {
	customerId := c.Params("id")

	if customerId == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"status":  fiber.StatusNotFound,
			"message": "Customer not found!",
			"error":   map[string]interface{}{},
		})
	}

	var customer models.Customer
	if err := db.DB.First(&customer, customerId).Error; err != nil {
		return c.Status(404).SendString("customer not found")
	}

	var address models.Address
	db.DB.Where("id=?", customer.AddressId).Find(&address)

	var city models.City
	db.DB.Where("id=?", address.CityId).Find(&city)

	var country models.Country
	db.DB.Where("id=?", city.CountryId).Find(&country)

	var customerRes CustomerResponse
	customerRes.ID = int(customer.ID)
	customerRes.FirstName = customer.FirstName
	customerRes.LastName = customer.LastName
	customerRes.Email = customer.Email
	customerRes.Active = customer.Active
	customerRes.Address = AddressResponse{
		Address:    address.Address,
		Address2:   address.Address2,
		District:   address.District,
		PostalCode: address.PostalCode,
		Phone:      address.Phone,
		City:       city.City,
		Country:    country.Country,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"status":  200,
		"message": "success",
		"body":    customerRes,
	})

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
	if err := db.DB.Where("country = ?", data.Country).Find(&country).Error; err == nil {
		country.Country = data.Country
		if err := db.DB.Save(&country).Error; err != nil {
			return c.Status(500).SendString("Error while saving the country")
		}
	}

	var city models.City
	if err := db.DB.Where("city = ?", data.City).Find(&city).Error; err == nil {
		city.City = data.City
		city.CountryId = country.ID
		if err := db.DB.Save(&city).Error; err != nil {
			return c.Status(500).SendString("Error while saving the city")
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
		"data":    data,
	})

}

func CustomersList(c *fiber.Ctx) error {
	var customers []models.Customer

	// Apply pagination and limits
	limit, page := 10, 1 // You can customize these values
	offset := (page - 1) * limit

	// Check if a query parameter for matching films is provided
	query := c.Query("q")
	if query != "" {

		if err := db.DB.Where("first_name LIKE ?", "%"+query+"%").
			Or("last_name LIKE ?", "%"+query+"%").Order("customers.created_at DESC").Limit(limit).Offset(offset).Find(&customers).Error; err != nil {
			return c.Status(500).SendString("Error while fetching customers")
		}

	} else {
		if err := db.DB.Order("customers.created_at DESC").Limit(limit).Offset(offset).Find(&customers).Error; err != nil {
			return c.Status(500).SendString("Error while fetching customers")
		}
	}

	var CustomerResponses []CustomerResponse

	for _, customer := range customers {

		var address models.Address
		db.DB.Where("id=?", customer.AddressId).Find(&address)

		var city models.City
		db.DB.Where("id=?", address.CityId).Find(&city)

		var country models.Country
		db.DB.Where("id=?", city.CountryId).Find(&country)

		CustomerResponses = append(CustomerResponses, CustomerResponse{
			ID:        int(customer.ID),
			FirstName: customer.FirstName,
			LastName:  customer.LastName,
			Email:     customer.Email,
			Active:    customer.Active,
			Address: AddressResponse{
				Address:    address.Address,
				Address2:   address.Address2,
				District:   address.District,
				PostalCode: address.PostalCode,
				Phone:      address.Phone,
				City:       city.City,
				Country:    country.Country,
			},
		})
	}

	totalRecords := len(customers)
	totalPageCount := (totalRecords + limit - 1) / limit

	// Create metadata
	meta := MetaInfo{
		PerPage:      limit,
		TotalPages:   totalPageCount,
		QueryInput:   query,
		TotalRecords: totalRecords,
	}

	response := map[string]interface{}{
		"customers": CustomerResponses,
		"meta":      meta,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"status":  200,
		"message": "success",
		"body":    response,
	})

}

func DeleteCustomer(c *fiber.Ctx) error {
	customerId := c.Params("id")

	var customer models.Customer

	db.DB.First(&customer, customerId)
	if customer.ID <= 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"messgae": "Customer not found!",
		})
	}

	var address models.Address
	db.DB.Where("id=?", customer.AddressId).Delete(&address)

	db.DB.Delete(&customer)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "customer is removed!",
	})

}
