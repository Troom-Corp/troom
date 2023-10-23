package services

import (
	"context"
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
	var user User
	conn := storage.SqlInterface.New()

	getUserQuery := fmt.Sprintf("SELECT * FROM public.users WHERE email='%s'", s.Email)
	conn.QueryRow(context.Background(), getUserQuery).Scan(&user.UserId, &user.FirstName, &user.SecondName, &user.Email, &user.Password, &user.Photo, &user.Bio, &user.Phone, &user.Links, &user.Followers, &user.Subscribers)

	if user.Email != "" {
		storage.SqlInterface.Close(conn)
		return fiber.NewError(409, "Пользователь с таким email уже существует")
	}

	if !PasswordValidator(s.Password) {
		storage.SqlInterface.Close(conn)
		return fiber.NewError(409, "Пароль должен соотвествовать требованиям")
	}

	storage.SqlInterface.Close(conn)
	return nil
}
