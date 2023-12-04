package routes

import (
	"github.com/Troom-Corp/troom/internal/controllers"
	"github.com/Troom-Corp/troom/internal/middleware"
	"github.com/Troom-Corp/troom/internal/store"
	"github.com/gofiber/fiber/v2"
)

func GetAuthRoutes(apiRouter fiber.Router, store store.InterfaceStore) {
	controller := controllers.GetAuthControllers(store)

	auth := apiRouter.Group("/auth")
	auth.Add("GET", "/profile", middleware.JWTMiddleware, controller.Profile)
	auth.Add("POST", "/users/sign_up", controller.SignUp)
	auth.Add("POST", "/users/sign_in", controller.SignIn)
	auth.Add("POST", "/users/validate_credentials", controller.ValidateCredentials)
	auth.Add("POST", "/sign_out", controller.Logout)
}
