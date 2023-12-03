package routes

import (
	"github.com/Troom-Corp/troom/internal/controllers"
	"github.com/Troom-Corp/troom/internal/store"
	"github.com/gofiber/fiber/v2"
)

func GetAuthRoutes(apiRouter fiber.Router, store store.InterfaceStore) {
	controller := controllers.GetAuthControllers(store)

	auth := apiRouter.Group("/auth")
	auth.Add("POST", "/users/sign_up", controller.UserSignUp)
	auth.Add("POST", "/users/sign_in", controller.UserSignIn)
	auth.Add("POST", "/users/validate_credentials", controller.ValidateCredentials)
	auth.Add("POST", "/logout", controller.Logout)
}
