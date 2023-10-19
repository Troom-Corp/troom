package middleware

import (
	jwt "github.com/Russian-LinkedIn/jwt_token-service"
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
	_, expiredTime := jwt.GetIdentity(headerToken)

	if expiredTime < time.Now().Unix() {
		return c.JSON("Token has expired")
	}

	return c.Next()
}
