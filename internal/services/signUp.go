package services

import (
	"fmt"
	"net/mail"
	"regexp"

	"github.com/Troom-Corp/troom/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type SignUpInterface interface {
	ValidPassword() error
	EmailNotInBase() error
	ValidEmail() error
	ValidNick() error
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

const ContainNums = `[0123456789]`
const ContainUpper = `[A-Z]`
const ContainLovver = `[a-z]`
const ContainSymbols = `[!@#$%^&*_-]`

func (s SignUpCredentials) ValidPassword() error {
	containNums, _ := regexp.Match(ContainNums, []byte(s.Password))
	containUpper, _ := regexp.Match(ContainUpper, []byte(s.Password))
	containSymbols, _ := regexp.Match(ContainSymbols, []byte(s.Password))
	containLover, _ := regexp.Match(ContainLovver, []byte(s.Password))
	if (len(s.Password) > 8 && len(s.Password) < 20) && containNums && containUpper && containSymbols && containLover {
		return fiber.NewError(200, "Пароль соотвествует требованиям")
	}

	return fiber.NewError(409, "Пароль не соотвествует требованиям")
}

func (s SignUpCredentials) EmailNotInBase() error {
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

	if userEmail == s.Email {
		conn.Close()
		return fiber.NewError(409, "Пользователь с таким email уже существует")
	}

	conn.Close()
	return nil
}

func (s SignUpCredentials) ValidEmail() error {
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
	conn.Close()

	_, err = mail.ParseAddress(s.Email)

	if err == nil && userEmail != s.Email {
		return fiber.NewError(200, "Почта соотвествует требованиям")
	} else if userEmail == s.Email {
		return fiber.NewError(409, "Почта уже есть в базе")
	}
	return fiber.NewError(409, "Почта не соотвествует требованиям")
}

func (s SignUpCredentials) ValidNick() error {
	var userNick string
	conn, err := storage.Sql.Open()

	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе")
	}

	getUserQuery := fmt.Sprintf("SELECT nick FROM public.users WHERE nick='%s'", s.Nick)
	rows, err := conn.Query(getUserQuery)

	for rows.Next() {
		err = rows.Scan(&userNick)
	}
	conn.Close()

	containNums, _ := regexp.Match(ContainNums, []byte(s.Nick))
	containUpper, _ := regexp.Match(ContainUpper, []byte(s.Nick))
	containSymbols, _ := regexp.Match(ContainSymbols, []byte(s.Nick))
	containLover, _ := regexp.Match(ContainLovver, []byte(s.Nick))

	if len(s.Nick) < 20 && containNums && containUpper && containSymbols && containLover && (userNick != s.Nick) {
		return fiber.NewError(200, "Ник соотвествует требованиям")
	} else if userNick == s.Nick {
		return fiber.NewError(409, "Ник уже есть в базе")
	}
	return fiber.NewError(409, "Ник не соотвествует требованиям")
}
