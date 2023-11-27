package server

import (
	"github.com/Troom-Corp/troom/internal/controllers"
	"github.com/Troom-Corp/troom/internal/routes"
	"github.com/Troom-Corp/troom/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var authControllers = controllers.AuthControllers{}

func Start() {
	app := fiber.New()

	store := store.NewStore()
	store.Open()
	defer store.Close()

	routes.InitRoutes(app, store)

	app.Static("/uploads", "./uploads")

	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Listen(":5000")
}
