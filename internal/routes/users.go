package routes

import (
	"github.com/Troom-Corp/troom/internal/controllers"
	"github.com/Troom-Corp/troom/internal/store"
	"github.com/gofiber/fiber/v2"
)

func GetUserRoutes(apiRouter fiber.Router, store store.InterfaceStore) {
	controller := controllers.GetUserControllers(store)

	users := apiRouter.Group("/user")
	users.Add("GET", "/profile", controller.Profile)
	users.Add("GET", "", controller.SearchByQuery)
}
