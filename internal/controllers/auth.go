package controllers

import (
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/services"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

type AuthControllers struct {
	SignInService services.SignInInterface
	SignUpService services.SignUpInterface
}

func (a AuthControllers) SignIn(c *fiber.Ctx) error {
	var credentials services.SignInCredentials
	err := c.BodyParser(&credentials)

	if err != nil {
		return fiber.NewError(500, "Неизвестная ошибка")
	}

	a.SignInService = credentials
	userId, err := a.SignInService.ValidData()

	if err != nil {
		return err
	}

	accessToken, _ := pkg.CreateAccessToken(userId)
	refreshToken, _ := pkg.CreateRefreshToken(userId)

	c.Cookie(&fiber.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
	})

	return c.JSON(accessToken)
}

func (a AuthControllers) SignUp(c *fiber.Ctx) error {
	var credentials services.SignUpCredentials
	err := c.BodyParser(&credentials)

	if err != nil {

		return fiber.NewError(500, "Ошибка при создании пользователя")
	}

	hashedPassword, err := pkg.Encode([]byte(credentials.Password))
	if err != nil {
		return fiber.NewError(500, "Ошибка при создании пользователя")
	}

	a.SignUpService = credentials
	err = a.SignUpService.ValidData()
	if err != nil {
		return err
	}

	newUser := services.User{
		Nick:        credentials.Nick,
		FirstName:   credentials.FirstName,
		SecondName:  credentials.SecondName,
		Email:       credentials.Email,
		Password:    string(hashedPassword),
		Gender:      credentials.Gender,
		DateOfBirth: credentials.DateOfBirth,
		Location:    credentials.Location,
		Job:         credentials.Job,
	}

	userId, err := newUser.Create()

	if err != nil {
		return fiber.NewError(500, "Ошибка при создании пользователя")
	}

	accessToken, _ := pkg.CreateAccessToken(userId)
	refreshToken, _ := pkg.CreateRefreshToken(userId)

	c.Cookie(&fiber.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
	})

	return c.JSON(accessToken)
}

func (a AuthControllers) RefreshToken(c *fiber.Ctx) error {
	authHeader := c.Get("authorization")
	if authHeader == "Bearer" {
		return fiber.NewError(401, "Access токена нет")
	}
	headerToken := strings.SplitN(authHeader, " ", 2)[1]
	userRefreshToken := c.Cookies("refresh_token")
	if userRefreshToken == "" {
		return fiber.NewError(401, "Refresh токена нет")
	}
	accessUserId, _, _ := pkg.GetIdentity(headerToken)
	refreshUserId, expTime, _ := pkg.GetIdentity(userRefreshToken)
	if expTime < time.Now().Unix() {
		c.ClearCookie("refresh_token")
	}

	if accessUserId != refreshUserId {
		return fiber.NewError(401, "Вы пытаетесь обновить чужой токен")
	}

	newAccessToken, _ := pkg.CreateAccessToken(refreshUserId)
	return c.JSON(newAccessToken)
}
