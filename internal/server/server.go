package server

import (
	"github.com/Troom-Corp/troom/internal/controllers"
	"github.com/Troom-Corp/troom/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var userControllers = controllers.UserControllers{}
var postControllers = controllers.PostsControllers{}
var commentControllers = controllers.CommentControllers{}
var authControllers = controllers.AuthControllers{}
var companyControllers = controllers.CompanyControllers{}
var vacanciesControllers = controllers.VacancyControllers{}
var profileControllers = controllers.ProfileControllers{}

func Start() {
	app := fiber.New()
	api := app.Group("/api")

	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	api.Get("/profile", middleware.Middleware, userControllers.Profile)

	// users group
	users := api.Group("/users")
	users.Get("/", userControllers.SearchUsersByQuery)
	users.Get("/@:nick", userControllers.GetUserByNick)
	users.Delete("/", middleware.Middleware, userControllers.DeleteUser)
	//users.Patch("/", middleware.Middleware, userControllers.UpdateInfo)
	//users.Patch("/update_login", middleware.Middleware, userControllers.UpdateLogin)
	users.Patch("/reset_password", middleware.Middleware, profileControllers.UpdatePassword)
	users.Patch("/reset_password/check_code", middleware.Middleware, profileControllers.CheckCode)

	// posts group
	posts := api.Group("/posts")
	posts.Get("/", postControllers.AllPost)
	posts.Get("/:id", postControllers.PostId)
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
	auth.Post("/logout", authControllers.Logout)
	auth.Post("/refresh_token", authControllers.RefreshToken)

	// companies group
	company := api.Group("/companies")
	company.Get("/", companyControllers.AllCompanies)
	company.Get("/:id", companyControllers.CompanyId)
	company.Delete("/", companyControllers.DeleteCompany)
	company.Patch("/", companyControllers.PatchCompany)

	// vacancies group
	vacancies := api.Group("/vacancies")
	vacancies.Get("/", vacanciesControllers.AllVacancies)
	vacancies.Get("/:id", vacanciesControllers.VacancyId)
	vacancies.Patch("/", vacanciesControllers.PatchVacancy)
	vacancies.Delete("/", vacanciesControllers.DeleteVacancy)

	app.Listen(":5000")
}
