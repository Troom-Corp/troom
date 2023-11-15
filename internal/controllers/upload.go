package controllers

import (
	"encoding/json"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FileControllers struct {
	FileInterface FileInterface
}

func (f FileControllers) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")

	if err != nil {
		return fiber.NewError(500, "Неизвестная ошибка")
	}

	isFileValid := f.ValidData(file.Filename, int(file.Size))

	if isFileValid.Ext != "" || isFileValid.Lenght != "" {
		isFileValidString, _ := json.Marshal(isFileValid)
		return fiber.NewError(409, string(isFileValidString))
	}

	newFileName := uuid.New().String() + filepath.Ext(file.Filename)
	err = c.SaveFile(file, "./uploads/"+newFileName)

	if err != nil {
		return fiber.NewError(500, "Ошибка при загрузке файла")
	}

	return fiber.NewError(200, "Вы успешно загрузили файл")
}
