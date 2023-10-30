package services

import (
	"fmt"
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type SignInInterface interface {
	ValidData() (int, error)
}

type SignInCredentials struct {
	Login    string // email or nick
	Password string
}

func (s SignInCredentials) ValidData() (int, error) {
	var userId int
	var userPassword string

	conn, err := storage.Sql.Open()
	if err != nil {
		return userId, err
	}

	getUserQuery := fmt.Sprintf("SELECT userid, password FROM public.users WHERE email='%s' OR nick='%s'", s.Login, s.Login)

	rows, err := conn.Query(getUserQuery)
	for rows.Next() {
		rows.Scan(&userId, &userPassword)
	}

	if pkg.Decode([]byte(userPassword), []byte(s.Password)) != nil {
		conn.Close()
		return 0, fiber.NewError(401, "Неправильные данные пользователя")
	}

	conn.Close()
	return userId, nil
}
