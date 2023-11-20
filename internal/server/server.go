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
var uploadControllers = controllers.UploadControllers{}

func Start() {
	app := fiber.New()
	api := app.Group("/api")
	app.Static("/uploads", "./uploads")

	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// users group
	users := api.Group("/users")
	users.Get("/", userControllers.SearchUsersByQuery)
	users.Get("/@:nick", userControllers.GetUserByNick)
	users.Delete("/", middleware.Middleware, userControllers.DeleteUser)

	// upload group
	upload := api.Group("/uploads")
	upload.Patch("/set_avatar", uploadControllers.SetAvatar)

	// profile group
	profile := api.Group("/profile")
	profile.Get("/", middleware.Middleware, profileControllers.Profile)
	profile.Patch("/reset_password/", middleware.Middleware, profileControllers.GetResetPasswordLink)
	profile.Patch("/reset_password/:uuid", middleware.Middleware, profileControllers.ResetPasswordByLink)
	profile.Patch("/reset_email", middleware.Middleware, profileControllers.GetResetEmailLink)
	profile.Patch("/reset_email/:uuid", middleware.Middleware, profileControllers.ResetEmailByLink)
	profile.Patch("/update_login", middleware.Middleware, profileControllers.UpdateLogin)
	profile.Patch("/update_info", middleware.Middleware, profileControllers.UpdateInfo)

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
	auth.Post("/users/sign_in", authControllers.UserSignIn)
	auth.Post("/users/sign_up", authControllers.UserSignUp)
	auth.Post("/logout", authControllers.Logout)
	auth.Post("/refresh_token", authControllers.RefreshToken)
	auth.Post("/testauth", authControllers.TestSingIn)

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
