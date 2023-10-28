package middleware

import (
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

func Middleware(c *fiber.Ctx) error {
	authHeader := c.Get("authorization")
	if authHeader == "Bearer" {
		return fiber.NewError(401, "Сначала нужно зарегистрироваться")
	}
	headerToken := strings.SplitN(authHeader, " ", 2)[1]
	_, expiredTime, err := pkg.GetIdentity(headerToken)

	if err != nil {
		return fiber.NewError(401, "Невалидный токен")
	}

	if expiredTime < time.Now().Unix() {
		return fiber.NewError(401, "Токен просрочен")
	}

	return c.Next()
}
