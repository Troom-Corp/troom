package services

import (
	"fmt"
	"github.com/Troom-Corp/troom/internal/storage"
	"github.com/gofiber/fiber/v2"
	"regexp"
	"unicode"
)

type SignUpInterface interface {
	ValidData() error
}

type SignUpCredentials struct {
	FirstName  string
	SecondName string
	Email      string
	Password   string
}

func PasswordValidator(password string) bool {
	var containNums bool = false

	for k := range password {
		if unicode.IsDigit(rune(password[k])) {
			containNums = true
			break
		}
	}

	containUpper, _ := regexp.Match(`[A-Z]`, []byte(password))
	containSymbols, _ := regexp.Match(`[!@#$%^&*_-]`, []byte(password))

	if (len(password) > 8 && len(password) < 20) && containNums && containUpper && containSymbols {
		return true
	}

	return false
}

func (s SignUpCredentials) ValidData() error {
	var userEmail string
	conn, err := storage.Sql.Open()

	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе")
	}

	getUserQuery := fmt.Sprintf("SELECT email FROM public.users WHERE email='%s'", s.Email)
	rows, err := conn.Query(getUserQuery)

	for rows.Next() {
		err = rows.Scan(&userEmail)
	}

	if userEmail != "" {
		conn.Close()
		return fiber.NewError(409, "Пользователь с таким email уже существует")
	}

	if !PasswordValidator(s.Password) {
		conn.Close()
		return fiber.NewError(409, "Пароль должен соотвествовать требованиям")
	}

	conn.Close()
	return nil
}
