package middleware

import (
	"github.com/gofiber/fiber/v3"

	"logmanagement/backend/internal/model"
	"logmanagement/backend/internal/repository"
)

func Authenticate(store repository.Store) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies("session")
		if token == "" {
			return fiber.ErrUnauthorized
		}
		u, err := store.Session(c.Context(), token)
		if err != nil {
			return fiber.ErrUnauthorized
		}
		c.Locals("user", u)
		return c.Next()
	}
}
func RequireAdmin(c fiber.Ctx) error {
	if c.Locals("user").(model.User).Role != "admin" {
		return fiber.ErrForbidden
	}
	return c.Next()
}
func IngestAuth(store repository.Store) fiber.Handler {
	return func(c fiber.Ctx) error {
		if key := c.Get("X-API-Key"); key != "" {
			tenant, err := store.APIKeyTenant(c.Context(), key)
			if err != nil {
				return fiber.ErrUnauthorized
			}
			c.Locals("tenant", tenant)
			return c.Next()
		}
		u, err := store.Session(c.Context(), c.Cookies("session"))
		if err != nil {
			return fiber.ErrUnauthorized
		}
		if u.Role != "admin" {
			return fiber.ErrForbidden
		}
		c.Locals("tenant", u.Tenant)
		return c.Next()
	}
}
