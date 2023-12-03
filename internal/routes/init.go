package routes

import (
	"github.com/Troom-Corp/troom/internal/store"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App, store store.InterfaceStore) {
	api := app.Group("/api")

	GetAuthRoutes(api, store)
	GetUserRoutes(api, store)
}
