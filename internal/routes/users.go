package routes

import (
	"github.com/Troom-Corp/troom/internal/controllers"
	"github.com/Troom-Corp/troom/internal/middleware"
	"github.com/Troom-Corp/troom/internal/store"
	"github.com/gofiber/fiber/v2"
)

// GetUserRoutes collect users controllers
func GetUserRoutes(apiRouter fiber.Router, store store.InterfaceStore) {
	controller := controllers.GetUserControllers(store)

	users := apiRouter.Group("/users", middleware.JWTMiddleware)
	users.Add("GET", "/profile", controller.Profile)
	users.Add("GET", "", controller.SearchByQuery)
	users.Add("POST", "/avatar", controller.SetAvatar)
	users.Add("DELETE", "/avatar", controller.DeleteAvatar)
	users.Add("POST", "/layout", controller.SetLayout)
	users.Add("DELETE", "/layout", controller.DeleteLayout)
}
