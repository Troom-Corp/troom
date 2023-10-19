package controllers

import (
	"fmt"
	jwt "github.com/Russian-LinkedIn/jwt_token-service"
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/services"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type AuthControllers struct {
	SignInService services.SignInInterface
	SignUpService services.SignUpInterface
}

func (a AuthControllers) SignIn(c *fiber.Ctx) error {
	var credentials services.SignInCredentials
	c.BodyParser(&credentials)
	a.SignInService = credentials

	fmt.Println(credentials)

	authUser, err := a.SignInService.ValidData()
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:  "refresh_token",
		Value: jwt.SignJWT(authUser.UserId),
	})

	return c.JSON(jwt.SignJWT(authUser.UserId))
}

func (a AuthControllers) SignUp(c *fiber.Ctx) error {
	var credentials services.SignUpCredentials
	c.BodyParser(&credentials)
	a.SignUpService = credentials

	hashedPassword, err := pkg.Encode([]byte(credentials.Password))
	if err != nil {
		return fiber.NewError(500, "Ошибка при создании пользователя")
	}
	newUser := services.User{FirstName: credentials.FirstName, SecondName: credentials.SecondName, Email: credentials.Email, Password: string(hashedPassword)}

	err = a.SignUpService.ValidData()
	if err != nil {
		return err
	}
	userId, err := newUser.Create()
	if err != nil {
		return fiber.NewError(500, "Ошибка при создании пользователя")
	}

	c.Cookie(&fiber.Cookie{
		Name:  "refresh_token",
		Value: jwt.SignJWT(userId),
	})

	return c.JSON(jwt.SignJWT(userId))
}

func (a AuthControllers) RefreshToken(c *fiber.Ctx) error {
	token := strings.SplitN(c.Get("authorization"), " ", 2)[1]
	userId, _ := jwt.GetIdentity(token)
	return c.JSON(jwt.SignJWT(userId))
}
