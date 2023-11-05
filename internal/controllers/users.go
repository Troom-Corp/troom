package controllers

import (
	"github.com/Troom-Corp/troom/internal/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type UserControllers struct {
	UserServices services.UserInterface
}

func (u UserControllers) GetUserByNick(c *fiber.Ctx) error {
	userNick := c.Params("nick")

	u.UserServices = services.User{Nick: userNick}
	user, err := u.UserServices.ReadByNick()

	if err != nil {
		return err
	}

	return c.JSON(user)
}

//func (u UserControllers) PatchUser(c *fiber.Ctx) error {
//	var user services.User
//	c.BodyParser(&user)
//	u.UserServices = user
//	err := u.UserServices.Update()
//	if err != nil {
//		return err
//	}
//	return fiber.NewError(200, "Пользователь успешно изменен")
//}

func (u UserControllers) GetAllUsers(c *fiber.Ctx) error {
	// получаем все queries
	searchQuery := c.Query("search_query")

	if searchQuery != "" {
		u.UserServices = services.User{}
		resultUsers, err := u.UserServices.SearchByQuery(searchQuery)
		if err != nil {
			return fiber.NewError(500, "Ошибочка")
		}
		return c.JSON(resultUsers)
	}

	u.UserServices = services.User{}
	allUsers, err := u.UserServices.ReadAll()
	if err != nil {
		return fiber.NewError(404, "Ошибочка")
	}
	return c.JSON(allUsers)
}

func (u UserControllers) DeleteUser(c *fiber.Ctx) error {
	userId, _ := strconv.Atoi(c.Query("user_id"))
	u.UserServices = services.User{UserId: userId}

	err := u.UserServices.Delete()
	if err != nil {
		return err
	}
	return fiber.NewError(200, "Пользователь успешно удален")
}
