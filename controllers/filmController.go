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
	PerPage    int
	TotalPages int
	QueryInput string
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
		LanguageId:      uint(body.LanguageId),
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

	if err := db.DB.Preload("Actors").Preload("Categories").Preload("Language").Find(&films).Error; err != nil {
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

		log.Printf("film id %v", film.Language.Name)

		filmResponses = append(filmResponses, FilmResponse{
			ID:              uint(film.Id),
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
				Name: film.Language.Name,
			},
		})
	}

	// Calculate total pages based on total records and per page limit
	totalRecords := len(films)
	totalPageCount := (totalRecords + limit - 1) / limit

	// Create metadata
	meta := MetaInfo{
		PerPage:    limit,
		TotalPages: totalPageCount,
		QueryInput: query,
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
