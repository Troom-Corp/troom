package controllers

import (
	"context"
	"fmt"
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/services"
	"github.com/Troom-Corp/troom/internal/storage"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

type ProfileControllers struct {
	ProfileServices services.ProfileInterface
}

func (p ProfileControllers) GetResetLink(c *fiber.Ctx) error {
	var oldPassword services.NewPasswordCredentials
	c.BodyParser(&oldPassword)

	authHeader := c.Get("authorization")
	authToken := strings.SplitN(authHeader, " ", 2)[1]
	userId, _, _ := pkg.GetIdentity(authToken)

	if oldPassword.GetResetLink(userId) != nil {
		return fiber.NewError(500, "Ошибка")
	}

	return fiber.NewError(200, "Ссылка была отправлена на вашу почту")
}

func (p ProfileControllers) ResetPasswordByLink(c *fiber.Ctx) error {
	uuidCode := c.Params("uuid")

	fmt.Println(uuidCode)

	var newPassword services.NewPasswordCredentials
	c.BodyParser(&newPassword)

	authHeader := c.Get("authorization")
	authToken := strings.SplitN(authHeader, " ", 2)[1]
	userId, _, _ := pkg.GetIdentity(authToken)

	rds := storage.Redis.Open()
	tokenUserId := rds.Get(context.Background(), strconv.Itoa(userId))

	if tokenUserId.Val() != uuidCode {
		rds.Close()
		return fiber.NewError(404, "Кажется вы пытаетесь подключиться к невалидной ссылке")
	}

	if !pkg.ValidPassword(newPassword.NewPassword) {
		rds.Close()
		return fiber.NewError(409, "Пароль не соответствует требованиям")
	}

	if newPassword.SetNewPassword(userId) != nil {
		rds.Close()
		return fiber.NewError(500, "Ошибка при сбрасывании пароля")
	}

	rds.Del(context.Background(), strconv.Itoa(userId))
	return fiber.NewError(200, "Пароль успешно обновлен")
}
