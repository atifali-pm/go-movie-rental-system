package routes

import (
	"github.com/atifali-pm/go-movie-rental-system/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Post("/categories", controllers.CreateCategory)
	app.Post("/actors", controllers.SaveActor)
	app.Post("/actors-in-film", controllers.SaveActorInFilm)

	app.Post("/customers", controllers.CreateCustomer)
	app.Get("/customers", controllers.CustomersList)
	app.Get("/customers/:id", controllers.CustomerDetail)
	app.Delete("/customers/:id", controllers.DeleteCustomer)
	app.Put("/customers/:id", controllers.UpdateCustomer)

	app.Post("/payments", controllers.MakePayment)
	app.Post("/return-film", controllers.ReturnFilm)
	app.Get("/payments/customer/:customer_id", controllers.PaymentsListByCustomer)

	// app.Get("/films", controllers.FilmLists)
	app.Post("/films", controllers.CreateFilm)
	app.Get("/films", controllers.FilmsList)
	app.Put("/films/:id", controllers.UpdateFilm)
	app.Delete("/films/:id", controllers.DeleteFilm)
	app.Get("/films/:id", controllers.FilmDetails)

	app.Post("/staff", controllers.CreateStaff)
	app.Get("/staff/:id", controllers.StaffDetail)
	app.Put("/staff/:id", controllers.UpdateStaff)
	app.Delete("/staff/:id", controllers.DeleteStaff)
	app.Get("/staff", controllers.StaffList)

}
