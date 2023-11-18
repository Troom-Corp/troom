package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/services"
	"github.com/Troom-Corp/troom/internal/storage"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

type AuthControllers struct {
	SignInService services.SignInInterface
	SignUpService services.SignUpInterface
}

func (a AuthControllers) UserSignIn(c *fiber.Ctx) error {
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
		Name:    "refresh_token",
		Value:   refreshToken,
		Expires: time.Now().Add(time.Minute),
	})

	return c.JSON(accessToken)
}

func (a AuthControllers) UserSignUp(c *fiber.Ctx) error {
	var credentials services.SignUpCredentials
	err := c.BodyParser(&credentials)

	if err != nil {
		return fiber.NewError(500, "Ошибка при чтении JSON")
	}
	a.SignUpService = credentials
	isUserValid := a.SignUpService.ValidData()

	if isUserValid.Login != "" || isUserValid.Email != "" || isUserValid.Password != "" {
		isUserValidString, _ := json.Marshal(isUserValid)
		return fiber.NewError(409, string(isUserValidString))
	}

	hashedPassword, err := pkg.Encode([]byte(credentials.Password))
	if err != nil {
		return fiber.NewError(500, "Ошибка при хешировании пароля")
	}

	newUser := services.User{
		Login:       credentials.Login,
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
	refreshUserId, _, _ := pkg.GetIdentity(userRefreshToken)

	if accessUserId != refreshUserId {
		return fiber.NewError(401, "Вы пытаетесь обновить чужой токен")
	}

	newAccessToken, _ := pkg.CreateAccessToken(refreshUserId)
	return c.JSON(newAccessToken)
}

func (a AuthControllers) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:  "refresh_token",
		Value: "",
		Path:  "/",
	})
	return fiber.NewError(200, "Вы успешно вышли из аккаунта")
}

func (a AuthControllers) TestSingIn(c *fiber.Ctx) error {
	var authUser services.User
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
	conn, err := storage.Sql.Open()

	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}
	defer conn.Close()
	getUserQuery := fmt.Sprintf("SELECT * FROM public.users WHERE email = '%s' OR nick = '%s'", credentials.Login, credentials.Login)
	conn.Get(&authUser, getUserQuery)

	accessToken, _ := pkg.CreateAccessToken(userId)
	refreshToken, _ := pkg.CreateRefreshToken(userId)

	c.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Value:   refreshToken,
		Expires: time.Now().Add(time.Minute),
	})
	c.Cookie(&fiber.Cookie{
		Name:    "access_token",
		Value:   accessToken,
		Expires: time.Now().Add(time.Minute),
	})

	return c.JSON(authUser)
}
