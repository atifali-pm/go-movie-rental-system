package routes

import (
	"github.com/atifali-pm/go-movie-rental-system/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Post("/categories", controllers.CreateCategory)
	app.Post("/films", controllers.CreateFilm)
	app.Post("/actors", controllers.SaveActor)
	app.Post("/actors-in-film", controllers.SaveActorInFilm)
	app.Get("/films/:id", controllers.FilmDetails)

	app.Post("/customers", controllers.CreateCustomer)

	app.Post("/payments", controllers.MakePayment)
	app.Post("/return-film", controllers.ReturnFilm)

	// app.Get("/films", controllers.FilmLists)
	app.Get("/films", controllers.FilmsList)
	app.Put("/films/:id", controllers.UpdateFilm)
	app.Delete("/films/:id", controllers.DeleteFilm)
}
