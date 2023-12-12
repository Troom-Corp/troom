package models

import (
	"mime/multipart"
	"path/filepath"
)

type IsValidImage struct {
	Extension string
	Size      string
}

type Avatar struct {
	Image *multipart.FileHeader
}

const FileMaxSize = 6291456

func (a Avatar) Validate() IsValidImage {
	var isImage IsValidImage

	extension := filepath.Ext(a.Image.Filename)
	size := a.Image.Size

	if extension != ".png" && extension != ".jpg" && extension != ".jpeg" {
		isImage.Extension = "Невалидное расширение"
	}

	if size > FileMaxSize {
		isImage.Size = "Слишком большой размер файла"
	}

	return isImage
}
