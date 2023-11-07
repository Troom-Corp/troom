package services

import (
	"fmt"
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/storage"
)

type SignUpInterface interface {
	ValidData() ValidUser
}

// SignUpCredentials первичная регистрация
type SignUpCredentials struct {
	FirstName   string
	SecondName  string
	Login       string
	Email       string
	Password    string
	Gender      string
	DateOfBirth string
	Location    string
	Job         string
}

// ValidUser структура, которая отправляет все невалидные данные клиенту
type ValidUser struct {
	Nick     string `json:"isNick"`
	Email    string `json:"isEmail"`
	Password string `json:"isPassword"`
}

func (s SignUpCredentials) ValidData() ValidUser {
	var isUserValid ValidUser

	var userEmail, userNick string
	conn, err := storage.Sql.Open()

	if err != nil {
		conn.Close()
		return ValidUser{}
	}

	getUserQuery := fmt.Sprintf("SELECT email, nick FROM public.users WHERE email='%s' OR nick='%s'", s.Email, s.Login)
	rows, err := conn.Query(getUserQuery)

	for rows.Next() {
		err = rows.Scan(&userEmail, &userNick)
	}

	if userNick == s.Login || len(s.Login) > 20 {
		isUserValid.Nick = "Невалидный Login"
	}

	if userEmail == s.Email {
		isUserValid.Email = "Пользователь с таким email уже существует"
	}

	if !pkg.ValidPassword(s.Password) {
		isUserValid.Password = "Пароль не соотвествует требованиям"
	}

	conn.Close()
	return isUserValid
}
