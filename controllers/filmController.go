package controllers

import (
	"log"

	db "github.com/atifali-pm/go-movie-rental-system/config"
	"github.com/atifali-pm/go-movie-rental-system/models"
	"github.com/gofiber/fiber/v2"
)

type FilmResponse struct {
	ID              uint               `json:"id"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	ReleaseYear     int                `json:"release_year"`
	Rental_Duration int                `json:"rental_duration"`
	RentalRate      float32            `json:"rental_rate"`
	Length          int                `json:"length"`
	ReplacementCost float32            `json:"replacement_cost"`
	Rating          int                `json:"rating"`
	SpecialFeature  string             `json:"special_feature"`
	FullText        string             `json:"full_text"`
	Language        LanguageResponse   `json:"language"`
	Actors          []ActorResponse    `json:"actors"`
	Categories      []CateogoryReponse `json:"categories"`
}

type LanguageResponse struct {
	Name string `json:"name"`
}

type CateogoryReponse struct {
	Name string `json:"name"`
}

type ActorResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

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

func FilmDetails(c *fiber.Ctx) error {
	filmId := c.Params("id")

	if filmId == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"status":  fiber.StatusNotFound,
			"message": "Film not found!",
			"error":   map[string]interface{}{},
		})
	}

	var film models.Film
	if err := db.DB.Preload("Actors").Preload("Categories").First(&film, filmId).Error; err != nil {
		return c.Status(404).SendString("Film not found")
	}

	var language models.Language
	db.DB.Where("id=?", film.LanguageId).Find(&language)

	var filmResponse FilmResponse
	filmResponse.ID = uint(film.Id)
	filmResponse.Title = film.Title
	filmResponse.Description = film.Description
	filmResponse.ReleaseYear = film.ReleaseYear
	filmResponse.Rental_Duration = film.Rental_Duration
	filmResponse.RentalRate = film.RentalRate
	filmResponse.Length = film.Length
	filmResponse.ReplacementCost = film.ReplacementCost
	filmResponse.Rating = film.Rating
	filmResponse.SpecialFeature = film.SpecialFeature
	filmResponse.FullText = film.FullText
	filmResponse.Language = LanguageResponse{
		Name: language.Name,
	}

	for _, actor := range film.Actors {
		actorResponse := ActorResponse{
			FirstName: actor.FirstName,
			LastName:  actor.LastName,
		}
		filmResponse.Actors = append(filmResponse.Actors, actorResponse)
	}

	for _, category := range film.Categories {
		cateogoryReponse := CateogoryReponse{
			Name: category.Name,
		}
		filmResponse.Categories = append(filmResponse.Categories, cateogoryReponse)
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"status":  200,
		"message": "success",
		"body":    filmResponse,
	})
}
