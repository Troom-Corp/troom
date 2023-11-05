package services

import (
	"fmt"
	"github.com/Troom-Corp/troom/internal/storage"
	"github.com/gofiber/fiber/v2"
	"regexp"
)

type SignUpInterface interface {
	ValidData() error
	ValidPassword() error
}

// первичная регистрация
type SignUpCredentials struct {
	FirstName   string
	SecondName  string
	Nick        string
	Email       string
	Password    string
	Gender      string
	DateOfBirth string
	Location    string
	Job         string
}

func (s SignUpCredentials) ValidPassword() error {
	containNums, _ := regexp.Match(`[0123456789]`, []byte(s.Password))
	containUpper, _ := regexp.Match(`[A-Z]`, []byte(s.Password))
	containSymbols, _ := regexp.Match(`[!@#$%^&*_-]`, []byte(s.Password))
	if (len(s.Password) > 8 && len(s.Password) < 20) && containNums && containUpper && containSymbols {
		return fiber.NewError(200, "Пароль соотвествует требованиям")
	}

	return fiber.NewError(409, "Пароль не соотвествует требованиям")
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

	conn.Close()
	return nil
}
