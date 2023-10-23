package routes

import (
	"github.com/atifali-pm/go-movie-rental-system/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Post("/categories", controllers.CreateCategory)
	app.Post("/films", controllers.StoreFilm)
}
