package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
)

func Security(origin string) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		c.SetContext(ctx)
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("Cache-Control", "no-store")
		if c.Method() != "GET" && c.Method() != "HEAD" && c.Method() != "OPTIONS" {
			origin := c.Get("Origin")
			if origin != "" && origin != origin {
				return fiber.NewError(403, "Origin not allowed")
			}
			if c.Cookies("session") != "" && origin == "" {
				return fiber.NewError(403, "Origin required for cookie-authenticated writes")
			}
		}
		return c.Next()
	}
}
