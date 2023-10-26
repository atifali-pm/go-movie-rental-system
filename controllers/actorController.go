package controllers

import (
	"log"

	db "github.com/atifali-pm/go-movie-rental-system/config"
	"github.com/atifali-pm/go-movie-rental-system/models"
	"github.com/gofiber/fiber/v2"
)

func SaveActor(c *fiber.Ctx) error {

	body := struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}{}

	err := c.BodyParser(&body)
	if err != nil {
		log.Fatalf("Actor error while saving %v", err)
	}

	if body.FirstName == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Actor first name is required",
			"error":   map[string]interface{}{},
		})
	}

	if body.LastName == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Actor last name is required",
			"error":   map[string]interface{}{},
		})

	}

	actor := models.Actor{
		FirstName: body.FirstName,
		LastName:  body.LastName,
	}

	db.DB.Create(&actor)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    actor,
	})
}

func SaveActorInFilm(c *fiber.Ctx) error {
	body := struct {
		FilmId  int `json:"film_id"`
		ActorId int `json:"actor_id"`
	}{}

	log.Println(body)
	err := c.BodyParser(&body)
	if err != nil {
		log.Fatalf("Actor error while saving %v", err)
	}

	if body.FilmId <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Film Id is required",
			"error":   map[string]interface{}{},
		})
	}

	if body.ActorId <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Actor Id is required",
			"error":   map[string]interface{}{},
		})
	}

	filmActor := models.FilmActor{
		FilmId:  body.FilmId,
		ActorId: body.ActorId,
	}

	db.DB.Create(&filmActor)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    filmActor,
	})

}
