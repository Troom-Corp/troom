package services

import "path/filepath"

type UploadInterface interface {
	ValidData(filename string, filesize int) ValidFile
}

type ValidFile struct {
	Ext    string `json:"isExt"`
	Lenght string `json:"isLenght"`
}

type UploadCredentials struct{}

const FileMaxLeignht = 5242880

func (u UploadCredentials) ValidData(filename string, filesize int) ValidFile {
	var isFileValid ValidFile
	ext := filepath.Ext(filename)
	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
		isFileValid.Ext = "Недопустимое расширение файла. Допустимые расширения: .jpg, .jpeg, .png"
	}
	if filesize > FileMaxLeignht {
		isFileValid.Lenght = "Максимальный размер файла - 5MB"
	}
	return isFileValid
}
