package controllers

import (
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/services"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type UserControllers struct {
	UserServices        services.UserInterface
	UserProfileServices services.ProfileInterface
}

func (u UserControllers) GetUserByNick(c *fiber.Ctx) error {
	userNick := c.Params("nick")

	u.UserServices = services.User{Login: userNick}
	user, err := u.UserServices.ReadByLogin()

	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (u UserControllers) SearchUsersByQuery(c *fiber.Ctx) error {
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
	authHeader := c.Get("authorization")
	authToken := strings.SplitN(authHeader, " ", 2)[1]
	userId, _, err := pkg.GetIdentity(authToken)
	if err != nil {
		return fiber.NewError(500, "Ошибка при открытии профиля")
	}
	u.UserServices = services.User{UserId: userId}

	err = u.UserServices.Delete()
	if err != nil {
		return err
	}
	return fiber.NewError(200, "Пользователь успешно удален")
}

func (u UserControllers) Profile(c *fiber.Ctx) error {
	authHeader := c.Get("authorization")
	authToken := strings.SplitN(authHeader, " ", 2)[1]
	userId, _, err := pkg.GetIdentity(authToken)
	if err != nil {
		return fiber.NewError(500, "Ошибка при открытии профиля")
	}

	userProfile, err := services.User{UserId: userId}.UserProfile()
	if err != nil {
		return err
	}
	return c.JSON(userProfile)
}

//func (u UserControllers) UpdateInfo(c *fiber.Ctx) error {
//	var newUserCredentials services.User
//	c.BodyParser(&newUserCredentials)
//
//	authHeader := c.Get("authorization")
//	authToken := strings.SplitN(authHeader, " ", 2)[1]
//	userId, _, err := pkg.GetIdentity(authToken)
//	if err != nil {
//		return fiber.NewError(500, "Ошибка при обновлении данных")
//	}
//
//	newUserCredentials.UserId = userId
//	u.UserProfileServices = newUserCredentials
//	err = u.UserProfileServices.UpdateInfo()
//	if err != nil {
//		return err
//	}
//	return fiber.NewError(200, "Пользователь успешно обновлен")
//}

//func (u UserControllers) UpdateLogin(c *fiber.Ctx) error {
//	login := struct {
//		Login string
//	}{}
//	c.BodyParser(&login)
//
//	authHeader := c.Get("authorization")
//	authToken := strings.SplitN(authHeader, " ", 2)[1]
//	userId, _, err := pkg.GetIdentity(authToken)
//	if err != nil {
//		return fiber.NewError(500, "Ошибка при обновлении данных")
//	}
//
//	u.UserProfileServices = services.User{UserId: userId, Login: login.Login}
//	err = u.UserProfileServices.UpdateLogin()
//	if err != nil {
//		return err
//	}
//	return fiber.NewError(200, "Логин успешно обновлен")
//}
