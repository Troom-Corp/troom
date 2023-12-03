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
	// Open the database connection
	store := store.NewStore()
	store.Open()
	defer store.Close()

	// init the fiber app instance
	app := fiber.New()

	// cors settings
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool { return true },
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// init the all routes
	routes.InitRoutes(app, store)

	// set the static route for images
	app.Static("/uploads", "./uploads")

	// listen the http requests
	app.Listen(":5000")
}
