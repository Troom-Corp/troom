package services

import (
	"context"
	"fmt"
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type SignInInterface interface {
	ValidData() (User, error)
}

type SignInCredentials struct {
	Email    string
	Password string
}

func (s SignInCredentials) ValidData() (User, error) {
	var user User
	conn := storage.SqlInterface.New()

	getUserQuery := fmt.Sprintf("SELECT * FROM public.users WHERE email='%s'", s.Email)
	conn.QueryRow(context.Background(), getUserQuery).Scan(&user.UserId, &user.FirstName, &user.SecondName, &user.Email, &user.Password, &user.Photo, &user.Bio, &user.Phone, &user.Links, &user.Followers, &user.Subscribers)

	if pkg.Decode([]byte(user.Password), []byte(s.Password)) != nil {
		storage.SqlInterface.Close(conn)
		return User{}, fiber.NewError(401, "Неправильные данные пользователя")
	}

	storage.SqlInterface.Close(conn)
	return user, nil
}
