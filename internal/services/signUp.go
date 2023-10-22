package services

import (
	"context"
	"fmt"
	"regexp"
	"unicode"

	"github.com/Troom-Corp/troom/internal"
	"github.com/gofiber/fiber/v2"
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
	getUserQuery := fmt.Sprintf("SELECT * FROM public.users WHERE email='%s'", s.Email)
	internal.Store().QueryRow(context.Background(), getUserQuery).Scan(
		&user.UserId,
		&user.FirstName,
		&user.SecondName,
		&user.Email,
		&user.Password,
		&user.Photo,
		&user.Bio,
		&user.Phone,
		&user.Links,
		&user.Followers,
		&user.Subscribers)

	if user.Email != "" {
		return fiber.NewError(409, "Пользователь с таким email уже существует")
	}
	if !PasswordValidator(s.Password) {
		return fiber.NewError(409, "Пароль должен соотвествовать требованиям")
	}

	return nil
}
