package controllers

import (
	"log"

	db "github.com/atifali-pm/go-movie-rental-system/config"
	"github.com/atifali-pm/go-movie-rental-system/models"
	"github.com/gofiber/fiber/v2"
)

func StoreFilm(c *fiber.Ctx) error {
	body := struct {
		Title           string  `json:"title"`
		LanguageId      int     `json:"language_id"`
		Description     string  `json:"description"`
		ReleaseYear     int     `json:"release_year"`
		RentalDuration  int     `json:"rental_duration"`
		RentalRate      float32 `json:"rental_rate"`
		Length          int     `json:"length"`
		ReplacementCost float32 `json:"replacement_cost"`
		Rating          int     `json:"rating"`
		SpecialFeature  string  `json:"special_feature"`
		FullText        string  `json:"full_text"`
		CategoryId      int     `json:"category_id"`
	}{}

	err := c.BodyParser(&body)

	if err != nil {
		log.Fatalf("Film storing error %v", err)
	}

	if body.CategoryId == 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": true,
			"message": "Required fields can not be empty",
		})
	}

	if body.Title == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": true,
			"message": "Required fields can not be empty",
		})
	}

	film := models.Film{
		Title:           body.Title,
		LanguageId:      body.LanguageId,
		Description:     body.Description,
		ReleaseYear:     body.ReleaseYear,
		Rental_Duration: body.RentalDuration,
		RentalRate:      body.RentalRate,
		Length:          body.Length,
		ReplacementCost: body.ReplacementCost,
		Rating:          body.Rating,
		SpecialFeature:  body.SpecialFeature,
		FullText:        body.FullText,
	}

	db.DB.Create(&film)

	filmCategory := models.FilmCategory{
		FilmId: film.Id,
	}

	db.DB.Create(&filmCategory)

	inventory := models.Inventory{
		FilmId: film.Id,
	}

	db.DB.Create(&inventory)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    film,
	})
}
