package handler

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"logmanagement/backend/internal/model"
)

func (h *Handler) Alerts(c fiber.Ctx) error {
	a, err := h.store.Alerts(c.Context(), c.Locals("user").(model.User).Tenant)
	if err != nil {
		return err
	}
	return c.JSON(a)
}

func (h *Handler) Acknowledge(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return fiber.NewError(400, "Invalid alert ID")
	}
	ok, err := h.store.Acknowledge(c.Context(), c.Locals("user").(model.User).Tenant, id)
	if err != nil {
		return err
	}
	if !ok {
		return fiber.ErrNotFound
	}
	return c.SendStatus(204)
}

func (h *Handler) GetRule(c fiber.Ctx) error {
	r, err := h.store.GetRule(c.Context(), c.Locals("user").(model.User).Tenant)
	if err != nil {
		return err
	}
	return c.JSON(r)
}

func (h *Handler) SetRule(c fiber.Ctx) error {
	var r model.Rule
	if json.Unmarshal(c.Body(), &r) != nil || r.Threshold < 2 || r.Threshold > 1000 || r.WindowMinutes < 1 || r.WindowMinutes > 60 {
		return fiber.NewError(400, "Threshold must be 2–1000 and window 1–60 minutes")
	}
	if err := h.store.SetRule(c.Context(), c.Locals("user").(model.User).Tenant, r); err != nil {
		return err
	}
	return c.JSON(r)
}
