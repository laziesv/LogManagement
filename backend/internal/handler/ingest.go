package handler

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"

	"logmanagement/backend/internal/model"
	"logmanagement/backend/internal/normalize"
)

func (h *Handler) Ingest(c fiber.Ctx) error {
	body := c.Body()
	if strings.HasPrefix(c.Get("Content-Type"), "multipart/form-data") {
		file, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(400, "Choose a JSON file")
		}
		f, err := file.Open()
		if err != nil {
			return err
		}
		defer f.Close()
		body, err = io.ReadAll(io.LimitReader(f, 2*1024*1024+1))
		if err != nil {
			return err
		}
		if len(body) > 2*1024*1024 {
			return fiber.ErrRequestEntityTooLarge
		}
	}
	batch, err := normalize.DecodeBatch(body)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}
	tenant := c.Locals("tenant").(string)
	events := make([]model.Event, 0, len(batch))
	now := time.Now().UTC()
	for i, raw := range batch {
		e, err := normalize.Normalize(raw, tenant, now)
		if err != nil {
			return fiber.NewError(400, fmt.Sprintf("Log %d: %s", i+1, err))
		}
		events = append(events, e)
	}
	if err = h.store.Insert(c.Context(), events); err != nil {
		return err
	}
	return c.Status(201).JSON(fiber.Map{"accepted": len(events), "tenant": tenant})
}
