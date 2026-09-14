package handler

import (
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) Health(c fiber.Ctx) error {
	if err := h.store.Ping(c.Context()); err != nil {
		return fiber.NewError(503, "Database unavailable")
	}
	return c.JSON(fiber.Map{"status": "ok"})
}
