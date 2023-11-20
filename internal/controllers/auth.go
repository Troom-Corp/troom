package controllers

import (
	"encoding/json"
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/services"
	"github.com/gofiber/fiber/v2"
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

	signInUser, err := services.ProfileInfo{UserId: userId}.UserProfile()
	if err != nil {
		return err
	}

	return c.JSON(signInUser)
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

	newUserObj, err := newUser.Create()

	if err != nil {
		return err
	}

	return c.JSON(newUserObj)
}
