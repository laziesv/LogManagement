package handler

import (
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) Logs(c fiber.Ctx) error {
	f, err := parseFilter(c)
	if err != nil {
		return err
	}
	events, total, err := h.store.Logs(c.Context(), f)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": events, "total": total, "limit": f.Limit, "offset": f.Offset})
}

func (h *Handler) Stats(c fiber.Ctx) error {
	f, err := parseFilter(c)
	if err != nil {
		return err
	}
	s, err := h.store.Stats(c.Context(), f)
	if err != nil {
		return err
	}
	return c.JSON(s)
}
