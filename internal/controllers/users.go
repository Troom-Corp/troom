package controllers

import (
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/store"
	"github.com/gofiber/fiber/v2"
)

type UserControllers struct {
	UserServices store.InterfaceUser
}

func (u *UserControllers) Profile(c *fiber.Ctx) error {
	jwt := c.Cookies("next-auth.session-token")
	ID, err := pkg.GetIdentity(jwt)
	if err != nil {
		return fiber.NewError(500, "Ошибка при получении профиля")
	}

	userProfile, err := u.UserServices.FindByID(ID)

	if err != nil {
		return fiber.NewError(500, "Ошибка при получении профиля")
	}

	return c.JSON(userProfile)
}

func (u *UserControllers) SearchByQuery(c *fiber.Ctx) error {
	queryParam := c.Query("login")

	queryUser, err := u.UserServices.FindByQuery(queryParam)
	if err != nil {
		return fiber.NewError(500, "Ошибка при поиске пользователя")
	}

	return c.JSON(queryUser)
}

func GetUserControllers(store store.InterfaceStore) UserControllers {
	return UserControllers{
		UserServices: store.Users(),
	}
}
