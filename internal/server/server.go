package server

import (
	"github.com/Troom-Corp/troom/internal/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var userControllers = controllers.UserControllers{}
var postControllers = controllers.PostsControllers{}
var commentControllers = controllers.CommentControllers{}
var authControllers = controllers.AuthControllers{}

func Start() {
	app := fiber.New()
	api := app.Group("/api")

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// users group
	users := api.Group("/users")
	users.Get("/", userControllers.GetUser)
	users.Delete("/", userControllers.DeleteUser)
	users.Patch("/", userControllers.PatchUser)

	posts := api.Group("/posts")
	posts.Get("/", postControllers.GetPost)
	posts.Delete("/", postControllers.DeletePost)
	posts.Patch("/", postControllers.PatchPost)

	comments := api.Group("/comments")
	comments.Post("/", commentControllers.CreateComment)
	comments.Get("/", commentControllers.GetComment)
	comments.Delete("/", commentControllers.DeleteComment)
	comments.Patch("/", commentControllers.PatchComment)

	auth := api.Group("/auth")
	auth.Post("/sign_in", authControllers.SignIn)
	auth.Post("/sign_up", authControllers.SignUp)
	auth.Post("/refresh_token", authControllers.RefreshToken)

	app.Listen(":5000")

}
