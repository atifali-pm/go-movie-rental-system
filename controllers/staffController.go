package controllers

import (
	"log"

	db "github.com/atifali-pm/go-movie-rental-system/config"
	"github.com/atifali-pm/go-movie-rental-system/models"
	"github.com/gofiber/fiber/v2"
)

type StaffData struct {
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

type PaymentResponse struct {
	Amount      float32                 `json:"amount"`
	PaymentDate string                  `json:"payment_date"`
	Customer    PaymentCustomerResponse `json:"customer"`
}

type PaymentCustomerResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func StaffList(c *fiber.Ctx) error {
	var staffList []models.Staff

	// Apply pagination and limits
	limit, page := 10, 1 // You can customize these values
	offset := (page - 1) * limit

	// Check if a query parameter for matching films is provided
	query := c.Query("q")
	if query != "" {

		if err := db.DB.Where("first_name LIKE ?", "%"+query+"%").
			Or("last_name LIKE ?", "%"+query+"%").Order("staffs.created_at DESC").Limit(limit).Offset(offset).Find(&staffList).Error; err != nil {
			return c.Status(500).SendString("Error while fetching records")
		}

	} else {
		if err := db.DB.Order("staffs.created_at DESC").Limit(limit).Offset(offset).Find(&staffList).Error; err != nil {
			return c.Status(500).SendString("Error while fetching records")
		}
	}

	var StaffResponses []StaffData

	for _, staff := range staffList {

		StaffResponses = append(StaffResponses, StaffData{
			FirstName:  staff.FirstName,
			LastName:   staff.LastName,
			Email:      staff.Email,
			Active:     staff.Active,
			AddressId:  staff.AddressId,
			StoreId:    staff.StoreId,
			Username:   staff.Username,
			Password:   staff.Password,
			PictureURL: staff.PictureURL,
		})
	}

	totalRecords := len(staffList)
	totalPageCount := (totalRecords + limit - 1) / limit

	// Create metadata
	meta := MetaInfo{
		PerPage:      limit,
		TotalPages:   totalPageCount,
		QueryInput:   query,
		TotalRecords: totalRecords,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"status":  200,
		"message": "success",
		"data":    StaffResponses,
		"meta":    meta,
	})

}

func CreateStaff(c *fiber.Ctx) error {
	var data StaffData

	err := c.BodyParser(&data)
	if err != nil {
		log.Fatalf("Error in the payload, staff not created %v", err)
	}

	if data.FirstName == "" || data.LastName == "" || data.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Fields are required",
		})
	}

	var existingEmail models.Staff
	if err := db.DB.Where("email = ?", data.Email).First(&existingEmail).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email already exists!",
		})
	}

	staff := models.Staff{
		AddressId:  data.AddressId,
		FirstName:  data.FirstName,
		LastName:   data.LastName,
		Active:     data.Active,
		Email:      data.Email,
		Username:   data.Username,
		Password:   data.Password,
		PictureURL: data.PictureURL,
		StoreId:    data.StoreId,
	}
	db.DB.Create(&staff)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    data,
	})

}

func StaffDetail(c *fiber.Ctx) error {
	staffId := c.Params("id")

	if staffId == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"status":  fiber.StatusNotFound,
			"message": "Staff not found!",
			"error":   map[string]interface{}{},
		})
	}

	var staff models.Staff
	if err := db.DB.First(&staff, staffId).Error; err != nil {
		return c.Status(404).SendString("Staff not found")
	}

	var payments []models.Payment
	db.DB.Where("staff_id=?", staff.ID).Find(&payments)

	var PaymentResponses []PaymentResponse
	for _, payment := range payments {

		var customer models.Customer
		db.DB.Where("id=?", payment.CustomerId).Find(&customer)

		PaymentResponses = append(PaymentResponses, PaymentResponse{
			Amount:      payment.Amount,
			PaymentDate: payment.PaymentDate,
			Customer: PaymentCustomerResponse{
				FirstName: customer.FirstName,
				LastName:  customer.LastName,
				Email:     customer.Email,
			},
		})
	}

	response := map[string]interface{}{
		"payments": PaymentResponses,
		"staff":    staff,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"status":  200,
		"message": "success",
		"body":    response,
	})

}

func UpdateStaff(c *fiber.Ctx) error {
	staffId := c.Params("id")
	var staff models.Staff

	db.DB.Find(&staff, "id=?", staffId)

	if staff.ID <= 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": true,
			"message": "Staff not found",
		})
	}

	var data StaffData
	c.BodyParser(&data)

	staff.FirstName = data.FirstName
	staff.LastName = data.LastName
	staff.Active = data.Active
	staff.Email = data.Email
	staff.Username = data.Username
	staff.Password = data.Password
	staff.AddressId = data.AddressId
	staff.PictureURL = data.PictureURL
	db.DB.Save(&staff)

	return c.Status(201).JSON(fiber.Map{
		"success":  true,
		"messsage": "success",
		"data":     data,
	})

}

func DeleteStaff(c *fiber.Ctx) error {
	staffId := c.Params("id")

	var staff models.Staff

	db.DB.First(&staff, staffId)
	if staff.ID <= 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"messgae": "Staff not found!",
		})
	}

	db.DB.Delete(&staff)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Staff is removed!",
	})
}
