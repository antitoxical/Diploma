package auth

import (
	"strings"

	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware validates the Authorization: Bearer <token> header.
// On success it stores *Claims in c.Locals("admin").
func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodOptions {
			return c.Next()
		}

		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(
				web.ErrorResp("error_missing_token", "Authorization header with Bearer token is required"),
			)
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := ParseToken(tokenStr, secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				web.ErrorResp("error_invalid_token", "Invalid or expired token"),
			)
		}

		c.Locals("admin", claims)
		return c.Next()
	}
}
