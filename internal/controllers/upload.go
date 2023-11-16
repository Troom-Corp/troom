package controllers

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UploadControllers struct {
	UploadInterface services.UploadInterface
}

func (u UploadControllers) GetPhoto(c *fiber.Ctx) error {
	filename := c.Params("filename")
	return c.SendFile("./uploads/" + filename)
}

func (u UploadControllers) SetAvatar(c *fiber.Ctx) error {
	var newProfileInfo services.ProfileInfo
	file, err := c.FormFile("file")

	if err != nil {
		return fiber.NewError(500, "Неизвестная ошибка")
	}

	isFileValid := u.UploadInterface.ValidData(file.Filename, int(file.Size))

	if isFileValid.Ext != "" || isFileValid.Lenght != "" {
		isFileValidString, _ := json.Marshal(isFileValid)
		return fiber.NewError(409, string(isFileValidString))
	}

	authHeader := c.Get("authorization")
	authToken := strings.SplitN(authHeader, " ", 2)[1]
	userId, _, _ := pkg.GetIdentity(authToken)
	user, _ := services.User{UserId: userId}.ReadByLogin()

	newFileName := uuid.New().String() + filepath.Ext(file.Filename)
	user.Avatar = newFileName

	userData, err := json.Marshal(user)
	if err != nil {
		return fiber.NewError(500, "Неизвестная ошибка")
	}

	err = json.Unmarshal(userData, &newProfileInfo)
	if err != nil {
		return fiber.NewError(500, "Неизвестная ошибка")
	}

	err = newProfileInfo.UpdateInfo(userId)

	if err != nil {
		return fiber.NewError(500, "Ошибка при загрузке аватарки у пользователя")
	}

	err = c.SaveFile(file, "./uploads/"+newFileName)

	if err != nil {
		return fiber.NewError(500, "Ошибка при загрузке файла")
	}

	return fiber.NewError(200, "Вы успешно загрузили файл")
}
