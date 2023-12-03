package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Troom-Corp/troom/internal/models"
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/store"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"strconv"
)

type UserControllers struct {
	UserServices store.InterfaceUser
}

func (u UserControllers) Profile(c *fiber.Ctx) error {
	jwt := c.Cookies("token")
	ID, err := pkg.GetIdentity(jwt)
	if err != nil {
		return fiber.NewError(500, "Ошибка при получении профиля")
	}

	userProfile, err := u.UserServices.FindOne("userid", ID)

	if err != nil {
		return fiber.NewError(500, "Ошибка при получении профиля")
	}

	return c.JSON(userProfile)
}

func (u UserControllers) SearchByQuery(c *fiber.Ctx) error {
	queries := c.Queries()
	limit, _ := strconv.Atoi(queries["limit"])
	page, _ := strconv.Atoi(queries["page"])

	queryUser, err := u.UserServices.QuerySearch(queries["q"], limit, page)
	if err != nil {
		return fiber.NewError(500, "Ошибка при поиске пользователя")
	}

	return c.JSON(queryUser)
}

func (u UserControllers) SetAvatar(c *fiber.Ctx) error {
	jwt := c.Cookies("token")
	ID, err := pkg.GetIdentity(jwt)

	avatar, err := c.FormFile("avatar")

	if err != nil {
		return fiber.NewError(500, "Ошибка при загрузке аватара")
	}

	avatarInstance := models.Avatar{
		Image: avatar,
	}

	isValidAvatar := avatarInstance.Validate()

	if isValidAvatar.Extension != "" || isValidAvatar.Size != "" {
		isValidJSON, _ := json.Marshal(isValidAvatar)
		return fiber.NewError(400, string(isValidJSON))
	}

	imageUUID := pkg.GenerateUUID()
	imageExt := filepath.Ext(avatar.Filename)

	oldAvatar, err := u.UserServices.UploadAvatar(ID, imageUUID+imageExt)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(500, "Ошибка при загрузке фотографии")
	}

	if oldAvatar != "" {
		os.Remove(fmt.Sprintf("./uploads/%s", oldAvatar))
	}

	err = c.SaveFile(avatar, fmt.Sprintf("./uploads/%s", imageUUID+imageExt))
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(500, "Ошибка при загрузке фотографии")
	}

	return fiber.NewError(200, "Аватар загружен")
}

func (u UserControllers) SetLayout(c *fiber.Ctx) error {
	jwt := c.Cookies("token")
	ID, _ := pkg.GetIdentity(jwt)

	layout, err := c.FormFile("layout")

	if err != nil {
		return fiber.NewError(500, "Ошибка при загрузке layout")
	}

	layoutInstance := models.Avatar{
		Image: layout,
	}

	isValidAvatar := layoutInstance.Validate()

	if isValidAvatar.Extension != "" || isValidAvatar.Size != "" {
		isValidJSON, _ := json.Marshal(isValidAvatar)
		return fiber.NewError(400, string(isValidJSON))
	}

	imageUUID := pkg.GenerateUUID()
	imageExt := filepath.Ext(layout.Filename)

	oldLayout, err := u.UserServices.UploadLayout(ID, imageUUID+imageExt)
	if err != nil {
		return fiber.NewError(500, "Ошибка при загрузке фотографии")
	}
	if oldLayout != "" {
		os.Remove(fmt.Sprintf("./uploads/%s", oldLayout))
	}

	err = c.SaveFile(layout, fmt.Sprintf("./uploads/%s", imageUUID+imageExt))
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(500, "Ошибка при загрузке фотографии")
	}

	return fiber.NewError(200, "Layout загружен")
}

func (u UserControllers) DeleteAvatar(c *fiber.Ctx) error {
	jwt := c.Cookies("token")
	ID, _ := pkg.GetIdentity(jwt)

	deletedAvatar, err := u.UserServices.DeleteAvatar(ID)
	if err != nil {
		return fiber.NewError(500, "Ошибка при удалении аватара")
	}

	if deletedAvatar != "" {
		os.Remove(fmt.Sprintf("./uploads/%s", deletedAvatar))
	}
	return fiber.NewError(200, "Аватар успешно удален")
}

func (u UserControllers) DeleteLayout(c *fiber.Ctx) error {
	jwt := c.Cookies("token")
	ID, _ := pkg.GetIdentity(jwt)

	deletedLayout, err := u.UserServices.DeleteLayout(ID)
	if err != nil {
		return fiber.NewError(500, "Ошибка при удалении layout")
	}
	if deletedLayout != "" {
		os.Remove(fmt.Sprintf("./uploads/%s", deletedLayout))
	}

	return fiber.NewError(200, "Layout успешно удален")
}

func GetUserControllers(store store.InterfaceStore) UserControllers {
	return UserControllers{
		UserServices: store.Users(),
	}
}
