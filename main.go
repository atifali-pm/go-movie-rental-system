package main

import (
	"fmt"

	db "github.com/atifali-pm/go-movie-rental-system/config"
	"github.com/atifali-pm/go-movie-rental-system/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	fmt.Println("hello")

	db.Connect()

	app := fiber.New()
	app.Use(cors.New())

	routes.Setup(app)
	app.Listen(":3000")
}
