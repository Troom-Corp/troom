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
	Nick       string
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
	var userEmail, userNick string
	conn, err := storage.Sql.Open()

	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе")
	}

	getUserQuery := fmt.Sprintf("SELECT email, nick FROM public.users WHERE email='%s' OR nick='%s'", s.Email, s.Nick)
	rows, err := conn.Query(getUserQuery)

	for rows.Next() {
		err = rows.Scan(&userEmail, &userNick)
	}

	if userNick == s.Nick {
		conn.Close()
		return fiber.NewError(409, "Пользователь с таким nick уже существует")
	}

	if userEmail == s.Email {
		conn.Close()
		return fiber.NewError(409, "Пользователь с таким email уже существует")
	}

	if len(userNick) > 20 {
		return fiber.NewError(409, "Nick должен соотвествовать требованиям")
	}

	if !PasswordValidator(s.Password) {
		conn.Close()
		return fiber.NewError(409, "Пароль должен соотвествовать требованиям")
	}

	conn.Close()
	return nil
}
