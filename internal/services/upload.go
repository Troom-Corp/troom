package services

import (
	"fmt"
	"github.com/Troom-Corp/troom/internal/storage"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
)

type UploadInterface interface {
	ValidData(filename string, filesize int) ValidFile
	DeleteOldAvatar(userid int) error
}

type ValidFile struct {
	Ext    string `json:"isExt"`
	Length string `json:"isLength"`
}

type UploadAvatar struct{}

const FileMaxLength = 5242880

func (u UploadAvatar) ValidData(filename string, filesize int) ValidFile {
	var isFileValid ValidFile
	ext := filepath.Ext(filename)

	if ext != ".jpg" {
		isFileValid.Ext = "Недопустимое расширение файла. Допустимые расширения: .jpg, .jpeg, .png"
	}
	if filesize > FileMaxLength {
		isFileValid.Length = "Максимальный размер файла - 5MB"
	}
	return isFileValid
}

func (u UploadAvatar) DeleteOldAvatar(userid int) error {
	var oldAvatar string
	conn, err := storage.Sql.Open()
	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	err = conn.Get(&oldAvatar, fmt.Sprintf("SELECT avatar FROM public.users WHERE userid = %d", userid))
	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	if oldAvatar != "" {
		err = os.Remove(fmt.Sprintf("./uploads/%s", oldAvatar))
		if err != nil {
			return fiber.NewError(500, "Говно адрес")
		}
	}
	return nil
}
