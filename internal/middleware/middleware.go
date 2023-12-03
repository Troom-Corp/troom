package middleware

import (
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware checks is the JWT valid and exist in general
func JWTMiddleware(c *fiber.Ctx) error {
	jwt := c.Cookies("token")
	if jwt == "" {
		return fiber.NewError(409, "Unauthorized")
	}

	_, err := pkg.GetIdentity(jwt)
	if err != nil {
		return fiber.NewError(409, "Unauthorized")
	}

	return c.Next()
}
