package routes

import (
	"github.com/Troom-Corp/troom/internal/store"
	"github.com/gofiber/fiber/v2"
)

// InitRoutes collect all routes in itself
func InitRoutes(app *fiber.App, store store.InterfaceStore) {
	// the main group of routes for API
	api := app.Group("/api")

	// connect routes by groups
	GetAuthRoutes(api, store)
	GetUserRoutes(api, store)
}
