package controllers

import (
	"log"

	db "github.com/atifali-pm/go-movie-rental-system/config"
	"github.com/atifali-pm/go-movie-rental-system/models"
	"github.com/gofiber/fiber/v2"
)

type FilmResponse struct {
	ID              uint              `json:"id"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	ReleaseYear     int               `json:"release_year"`
	Rental_Duration int               `json:"rental_duration"`
	RentalRate      float32           `json:"rental_rate"`
	Length          int               `json:"length"`
	ReplacementCost float32           `json:"replacement_cost"`
	Rating          int               `json:"rating"`
	SpecialFeature  string            `json:"special_feature"`
	FullText        string            `json:"full_text"`
	Language        LanguageResponse  `json:"language"`
	Actors          []ActorResponse   `json:"actors"`
	Categories      []CategoryReponse `json:"categories"`
}

type LanguageResponse struct {
	Name string `json:"name"`
}

type CategoryReponse struct {
	Name string `json:"name"`
}

type ActorResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type MetaInfo struct {
	PerPage      int
	TotalPages   int
	QueryInput   string
	TotalRecords int
}

func CreateFilm(c *fiber.Ctx) error {
	// Parse the incoming JSON payload
	var film models.Film

	if err := c.BodyParser(&film); err != nil {
		return c.Status(400).SendString("Invalid request body")
	}

	// Check if a film with the same title already exists
	var existingFilm models.Film
	if err := db.DB.Where("title = ?", film.Title).First(&existingFilm).Error; err == nil {
		return c.Status(409).SendString("Film with the same title already exists")
	}

	// Create and save the language
	var language models.Language
	if err := db.DB.Where("name = ?", film.Language.Name).First(&language).Error; err != nil {
		// Language doesn't exist, save it
		if err := db.DB.Save(&language).Error; err != nil {
			return c.Status(500).SendString("Error while saving the language")
		}
	} else {
		// Language already exists, use the existing one
		film.Language = language
	}

	var newFilmData models.Film
	newFilmData.Title = film.Title
	newFilmData.Description = film.Description
	newFilmData.ReleaseYear = film.ReleaseYear
	newFilmData.RentalRate = film.RentalRate
	newFilmData.Rating = film.Rating
	newFilmData.Length = film.Length
	newFilmData.ReplacementCost = film.ReplacementCost
	newFilmData.Rating = film.Rating
	// newFilmData.Language = film.Language
	// newFilmData.Actors = film.Actors
	// newFilmData.Categories = film.Categories
	newFilmData.SpecialFeature = film.SpecialFeature
	newFilmData.FullText = film.FullText
	newFilmData.Language = film.Language
	newFilmData.LanguageId = film.LanguageId

	// Create and save the film
	if err := db.DB.Create(&newFilmData).Error; err != nil {
		return c.Status(500).SendString("Error while creating the film")
	}

	// Iterate through the actors and save them if they don't exist
	for i := range film.Actors {
		actor := &film.Actors[i]
		var existingActor models.Actor
		var filmActor models.FilmActor
		if err := db.DB.Where("first_name = ? AND last_name = ?", actor.FirstName, actor.LastName).First(&existingActor).Error; err != nil {
			// Actor doesn't exist, save it
			log.Printf("Creating new actors: %s", actor.FirstName)
			if err := db.DB.Create(&actor).Error; err != nil {
				return c.Status(500).SendString("Error while saving actors")
			}
			filmActor.ActorId = int(actor.ID)
		} else {
			filmActor.ActorId = int(existingActor.ID)
		}
		filmActor.FilmId = int(newFilmData.ID)
		db.DB.Save(&filmActor)
	}

	// Iterate through the categories and save them if they don't exist
	for i := range film.Categories {
		category := &film.Categories[i]
		var existingCategory models.Category
		var filmCategory models.FilmCategory
		if err := db.DB.Where("name = ?", category.Name).First(&existingCategory).Error; err != nil {
			// Category doesn't exist, save it
			log.Printf("Creating new categories: %s", category.Name)
			if err := db.DB.Create(&category).Error; err != nil {
				return c.Status(500).SendString("Error while saving categories")
			}
			filmCategory.CategoryId = int(category.ID)
		} else {
			filmCategory.CategoryId = int(existingCategory.ID)
		}
		filmCategory.FilmId = int(newFilmData.ID)
		db.DB.Save(&filmCategory)
	}

	// Return a success response
	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"status":  201,
		"message": "success",
		"body":    film,
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
	filmResponse.ID = uint(film.ID)
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
		cateogoryReponse := CategoryReponse{
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

func FilmsList(c *fiber.Ctx) error {
	var films []models.Film

	// Apply pagination and limits
	limit, page := 10, 1 // You can customize these values
	offset := (page - 1) * limit
	db.DB.Limit(limit).Offset(offset)

	// Check if a query parameter for matching films is provided
	query := c.Query("q")
	if query != "" {
		// // Subquery for actors
		// actorQuery := db.DB.Model(&films).
		//     Joins("JOIN film_actors ON films.id = film_actors.film_id").
		//     Joins("JOIN actors ON film_actors.actor_id = actors.id").
		//     Where("actors.name LIKE ?", "%"+query+"%").
		//     Select("films.id")

		// // Subquery for categories
		// categoryQuery := db.DB.Model(&Film{}).
		//     Joins("JOIN film_categories ON films.id = film_categories.film_id").
		//     Joins("JOIN categories ON film_categories.category_id = categories.id").
		//     Where("categories.name LIKE ?", "%"+query+"%").
		//     Select("films.id")

		// // Subquery for languages
		// languageQuery := db.DB.Model(&Film{}).
		//     Joins("JOIN languages ON films.language_id = languages.id").
		//     Where("languages.name LIKE ?", "%"+query+"%").
		//     Select("films.id")

		// db = db.Where("films.name LIKE ?", "%"+query+"%").
		//     Or("films.id IN (?)", actorQuery).
		//     Or("films.id IN (?)", categoryQuery).
		//     Or("films.id IN (?)", languageQuery)

		// db.DB = db.DB.
		// 	Joins("LEFT JOIN languages ON films.language_id = languages.id").
		// 	Joins("LEFT JOIN film_actors AS fa ON films.id = fa.film_id").
		// 	Joins("LEFT JOIN actors ON fa.actor_id = actors.id").
		// 	Joins("LEFT JOIN film_categories AS fc ON films.id = fc.film_id").
		// 	Joins("LEFT JOIN categories ON fc.category_id = categories.id").
		// 	Where("films.title LIKE ?", "%"+query+"%").
		// 	Or("actors.first_name LIKE ?", "%"+query+"%").
		// 	Or("categories.name LIKE ?", "%"+query+"%").
		// 	Or("languages.name LIKE ?", "%"+query+"%")
		db.DB = db.DB.Where("title LIKE ?", "%"+query+"%")
	}

	if err := db.DB.Preload("Actors").Preload("Categories").Order("films.created_at DESC").Find(&films).Error; err != nil {
		return c.Status(500).SendString("Error while fetching films")
	}

	var filmResponses []FilmResponse
	for _, film := range films {
		var actorsResponse []ActorResponse
		for _, actor := range film.Actors {
			actorsResponse = append(actorsResponse, ActorResponse{
				FirstName: actor.FirstName,
				LastName:  actor.LastName,
			})
		}

		var categoriesResponse []CategoryReponse
		for _, category := range film.Categories {
			categoriesResponse = append(categoriesResponse, CategoryReponse{
				Name: category.Name,
			})
		}

		var language models.Language
		db.DB.Where("id=?", film.LanguageId).Find(&language)

		filmResponses = append(filmResponses, FilmResponse{
			ID:              uint(film.ID),
			Title:           film.Title,
			Description:     film.Description,
			ReleaseYear:     film.ReleaseYear,
			Rental_Duration: film.Rental_Duration,
			RentalRate:      film.RentalRate,
			Length:          film.Length,
			ReplacementCost: film.ReplacementCost,
			Rating:          film.Rating,
			SpecialFeature:  film.SpecialFeature,
			FullText:        film.FullText,
			Actors:          actorsResponse,
			Categories:      categoriesResponse,
			Language: LanguageResponse{
				Name: language.Name,
			},
		})
	}

	// Calculate total pages based on total records and per page limit
	totalRecords := len(films)
	totalPageCount := (totalRecords + limit - 1) / limit

	// Create metadata
	meta := MetaInfo{
		PerPage:      limit,
		TotalPages:   totalPageCount,
		QueryInput:   query,
		TotalRecords: totalRecords,
	}

	response := map[string]interface{}{
		"films": filmResponses,
		"meta":  meta,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"status":  200,
		"message": "success",
		"body":    response,
	})
}

func DeleteFilm(c *fiber.Ctx) error {
	filmId := c.Params("id")

	var film models.Film

	db.DB.First(&film, filmId)
	if film.ID <= 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"messgae": "film not found!",
		})
	}

	var filmCategory models.FilmCategory
	db.DB.Where("film_id=?", filmId).Delete(&filmCategory)

	var filmActor models.FilmActor
	db.DB.Where("film_id=?", filmId).Delete(&filmActor)

	db.DB.Delete(&film)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "film is removed!",
	})
}

func UpdateFilm(c *fiber.Ctx) error {
	filmId := c.Params("id")

	var updateFilmData models.Film
	if err := db.DB.Where("id=?", filmId).Find(&updateFilmData).Error; err != nil {
		return c.Status(404).SendString("Film not found!")
	}

	var film models.Film
	if err := c.BodyParser(&film); err != nil {
		return c.Status(400).SendString("Invalid request body")
	}

	// Create and save the language
	var language models.Language
	if err := db.DB.Where("name = ?", film.Language.Name).First(&language).Error; err != nil {
		// Language doesn't exist, save it
		if err := db.DB.Create(&language).Error; err != nil {
			return c.Status(500).SendString("Error while saving the language")
		}
	} else {
		// Language already exists, use the existing one
		film.Language = language
	}

	// Iterate through the actors and save them if they don't exist
	for i := range film.Actors {
		actor := &film.Actors[i]
		var existingActor models.Actor
		if err := db.DB.Where("first_name = ? AND last_name = ?", actor.FirstName, actor.LastName).First(&existingActor).Error; err != nil {
			// Actor doesn't exist, save it
			if err := db.DB.Create(&actor).Error; err != nil {
				return c.Status(500).SendString("Error while saving actors")
			}
		}
	}

	// Iterate through the categories and save them if they don't exist
	for i := range film.Categories {
		category := &film.Categories[i]
		var existingCategory models.Category
		if err := db.DB.Where("name = ?", category.Name).First(&existingCategory).Error; err != nil {
			// Category doesn't exist, save it
			if err := db.DB.Create(&category).Error; err != nil {
				return c.Status(500).SendString("Error while saving categories")
			}
		}
	}

	updateFilmData.Title = film.Title
	updateFilmData.Description = film.Description
	updateFilmData.ReleaseYear = film.ReleaseYear
	updateFilmData.RentalRate = film.RentalRate
	updateFilmData.Rating = film.Rating
	updateFilmData.Length = film.Length
	updateFilmData.ReplacementCost = film.ReplacementCost
	updateFilmData.Rating = film.Rating
	updateFilmData.Language = film.Language
	// updateFilmData.Actors = film.Actors
	// updateFilmData.Categories = film.Categories
	updateFilmData.SpecialFeature = film.SpecialFeature
	updateFilmData.FullText = film.FullText

	// Save the film with updated actor and category associations
	if err := db.DB.Save(&updateFilmData).Error; err != nil {
		return c.Status(500).SendString("Error while saving the film with actors and categories")
	}
	updateFilmData.Actors = film.Actors
	updateFilmData.Categories = film.Categories
	return c.Status(201).JSON(updateFilmData)
}
