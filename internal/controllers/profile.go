package controllers

import (
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/services"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type ProfileControllers struct {
	ProfileServices services.ProfileInterface
}

func (p ProfileControllers) UpdatePassword(c *fiber.Ctx) error {
	authHeader := c.Get("authorization")
	authToken := strings.SplitN(authHeader, " ", 2)[1]
	userId, _, _ := pkg.GetIdentity(authToken)

	var newPasswordCredentials services.NewPasswordCredentials
	c.BodyParser(&newPasswordCredentials)

	if newPasswordCredentials.ResetPassword(userId) != nil {
		return fiber.NewError(500, "Ошибка")
	}

	return nil
}

func (p ProfileControllers) CheckCode(c *fiber.Ctx) error {
	authHeader := c.Get("authorization")
	authToken := strings.SplitN(authHeader, " ", 2)[1]
	userId, _, _ := pkg.GetIdentity(authToken)

	var newPasswordCredentials services.NewPasswordCredentials
	c.BodyParser(&newPasswordCredentials)

	if err := newPasswordCredentials.CheckCode(userId); err != nil {
		return err
	}

	return fiber.NewError(200, "успешно изменили пароль")
}
