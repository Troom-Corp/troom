package controllers

import (
	"github.com/Troom-Corp/troom/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type UserControllers struct {
	UserServices services.UserInterface
}

func (u UserControllers) UserId(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(500, "Неизвестная ошибка")
	}

	u.UserServices = services.User{UserId: userId}
	user, err := u.UserServices.ReadById()

	if err != nil {
		return fiber.NewError(500, "Неизвестная ошибка")
	}

	return c.JSON(user)
}

func (u UserControllers) PatchUser(c *fiber.Ctx) error {
	var user services.User
	c.BodyParser(&user)
	u.UserServices = user
	err := u.UserServices.Update()
	if err != nil {
		return fiber.NewError(500, "Ошибка при обновлении данных")
	}
	return fiber.NewError(200, "Данные успешно обновлены")
}

func (u UserControllers) AllUsers(c *fiber.Ctx) error {
	// получаем все queries
	searchQuery := c.Query("search_query")

	if searchQuery != "" {
		u.UserServices = services.User{}
		resultUsers, err := u.UserServices.SearchByQuery(searchQuery)
		if err != nil {
			return fiber.NewError(500, "Ошибка при поиске пользователей")
		}
		return c.JSON(resultUsers)
	}

	u.UserServices = services.User{}
	allUsers, err := u.UserServices.ReadAll()
	if err != nil {
		return fiber.NewError(404, "Ошибка при получении пользователей")
	}
	return c.JSON(allUsers)

}

func (u UserControllers) DeleteUser(c *fiber.Ctx) error {
	userId, _ := strconv.Atoi(c.Query("user_id"))
	u.UserServices = services.User{UserId: userId}

	err := u.UserServices.Delete()
	if err != nil {
		return fiber.NewError(500, "Ошибка при удалении пользователя")
	}
	return fiber.NewError(200, "Пользователь успешно удален")
}
