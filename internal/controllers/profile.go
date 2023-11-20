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
}

func (p ProfileControllers) GetResetPasswordLink(c *fiber.Ctx) error {
	var oldPassword services.NewPasswordCredentials
	c.BodyParser(&oldPassword)

	authHeader := c.Get("authorization")
	authToken := strings.SplitN(authHeader, " ", 2)[1]
	userId, _, _ := pkg.GetIdentity(authToken)

	err := oldPassword.GetResetPasswordLink(userId)
	return err
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

func (p ProfileControllers) GetResetEmailLink(c *fiber.Ctx) error {
	var oldEmail services.NewEmailCredentials
	c.BodyParser(&oldEmail)

	authHeader := c.Get("authorization")
	authToken := strings.SplitN(authHeader, " ", 2)[1]
	userId, _, _ := pkg.GetIdentity(authToken)

	err := oldEmail.GetResetEmailLink(userId)
	return err
}

func (p ProfileControllers) ResetEmailByLink(c *fiber.Ctx) error {
	uuidCode := c.Params("uuid")

	var newEmail services.NewEmailCredentials
	c.BodyParser(&newEmail)

	authHeader := c.Get("authorization")
	authToken := strings.SplitN(authHeader, " ", 2)[1]
	userId, _, _ := pkg.GetIdentity(authToken)

	rds := storage.Redis.Open()
	conn, err := storage.Sql.Open()
	if err != nil {
		return err
	}
	defer conn.Close()
	defer rds.Close()

	tokenUserId := rds.Get(context.Background(), strconv.Itoa(userId))

	if tokenUserId.Val() != uuidCode {
		return fiber.NewError(404, "Кажется вы пытаетесь подключиться к невалидной ссылке")
	}

	checkForDuplicateEmailQuery := fmt.Sprintf("SELECT email FROM public.users WHERE email = '%s'", newEmail)
	_, err = conn.Query(checkForDuplicateEmailQuery)
	if err != nil {
		return fiber.NewError(409, "Такая почта уже существует")
	}

	err = newEmail.SetNewEmail(userId)
	if err != nil {
		return fiber.NewError(200, "Ошибка при обновлении почты")
	}

	rds.Del(context.Background(), strconv.Itoa(userId))
	return fiber.NewError(200, "Почта успешно обновлена")
}

func (p ProfileControllers) UpdateLogin(c *fiber.Ctx) error {
	var newLoginCredentials services.NewLoginCredentials
	c.BodyParser(&newLoginCredentials)
	authHeader := c.Get("authorization")
	authToken := strings.SplitN(authHeader, " ", 2)[1]
	userId, _, _ := pkg.GetIdentity(authToken)

	err := newLoginCredentials.SetNewLogin(userId)
	return err
}

func (p ProfileControllers) UpdateInfo(c *fiber.Ctx) error {
	var newProfileInfo services.ProfileInfo
	c.BodyParser(&newProfileInfo)
	authHeader := c.Get("authorization")
	authToken := strings.SplitN(authHeader, " ", 2)[1]
	userId, _, _ := pkg.GetIdentity(authToken)
	newProfileInfo.UserId = userId
	err := newProfileInfo.UpdateInfo()
	return err
}

func (p ProfileControllers) Profile(c *fiber.Ctx) error {
	authHeader := c.Get("authorization")
	authToken := strings.SplitN(authHeader, " ", 2)[1]
	userId, _, err := pkg.GetIdentity(authToken)
	if err != nil {
		return fiber.NewError(500, "Ошибка при открытии профиля")
	}

	userProfile, err := services.ProfileInfo{UserId: userId}.UserProfile()
	if err != nil {
		return err
	}
	return c.JSON(userProfile)
}
