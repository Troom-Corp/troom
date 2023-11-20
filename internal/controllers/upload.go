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

func (u UploadControllers) SetAvatar(c *fiber.Ctx) error {
	file, err := c.FormFile("avatar")

	if err != nil {
		return fiber.NewError(500, "Ошибка при загрузке аватара")
	}

	u.UploadInterface = services.UploadAvatar{}
	isFileValid := u.UploadInterface.ValidData(file.Filename, int(file.Size))

	if isFileValid.Ext != "" || isFileValid.Length != "" {
		isFileValidString, _ := json.Marshal(isFileValid)
		return fiber.NewError(409, string(isFileValidString))
	}
	newFileName := uuid.New().String() + filepath.Ext(file.Filename)

	authHeader := c.Get("authorization")
	authToken := strings.SplitN(authHeader, " ", 2)[1]
	userId, _, _ := pkg.GetIdentity(authToken)

	err = c.SaveFile(file, "./uploads/"+newFileName)
	if err != nil {
		return fiber.NewError(500, "Ошибка при загрузке файла")
	}

	err = u.UploadInterface.DeleteOldAvatar(userId)
	if err != nil {
		return err
	}
	newProfileInfo := services.ProfileInfo{UserId: userId, Avatar: newFileName}
	err = newProfileInfo.UpdateInfo()

	if err != nil {
		return fiber.NewError(500, "Ошибка при загрузке аватарки у пользователя")
	}

	return fiber.NewError(200, "Вы успешно загрузили файл")
}
