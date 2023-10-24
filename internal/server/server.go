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
var vacanciesControllers = controllers.VacancyControllers{}

func Start() {
	app := fiber.New()
	api := app.Group("/api")

	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// users group
	users := api.Group("/users")
	users.Get("/", userControllers.AllUsers)
	users.Get("/:id", userControllers.UserId)
	users.Delete("/", userControllers.DeleteUser)
	users.Patch("/", userControllers.PatchUser)

	// posts group
	posts := api.Group("/posts")
	posts.Get("/", postControllers.AllPost)
	users.Get("/:id", postControllers.PostId)
	posts.Delete("/", postControllers.DeletePost)
	posts.Patch("/", postControllers.PatchPost)

	// comments group
	comments := api.Group("/comments")
	comments.Post("/", commentControllers.CreateComment)
	comments.Get("/:post_id", commentControllers.CommentByPostId)
	comments.Delete("/", commentControllers.DeleteComment)
	comments.Patch("/", commentControllers.PatchComment)

	// authorization group
	auth := api.Group("/auth")
	auth.Post("/sign_in", authControllers.SignIn)
	auth.Post("/sign_up", authControllers.SignUp)
	auth.Post("/refresh_token", authControllers.RefreshToken)

	// vacancies group
	vacancies := api.Group("/vacancies")
	vacancies.Get("/", vacanciesControllers.AllVacancies)
	vacancies.Get("/:id", vacanciesControllers.VacancyId)
	vacancies.Patch("/", vacanciesControllers.PatchVacancy)
	vacancies.Delete("/", vacanciesControllers.DeleteVacancy)

	app.Listen(":5000")
}
