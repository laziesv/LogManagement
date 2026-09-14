package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"

	"logmanagement/backend/internal/model"
)

func parseFilter(c fiber.Ctx) (model.Filter, error) {
	tenant := c.Locals("user").(model.User).Tenant
	if t := c.Query("tenant"); t != "" && t != tenant {
		return model.Filter{}, fiber.NewError(403, "Tenant access denied")
	}
	f := model.Filter{Tenant: tenant, Source: c.Query("source"), Query: c.Query("q"), From: time.Now().UTC().Add(-24 * time.Hour), To: time.Now().UTC(), Limit: 50}
	if len(f.Query) > 200 {
		return f, fiber.NewError(400, "Search is limited to 200 characters")
	}
	for key, dst := range map[string]*time.Time{"from": &f.From, "to": &f.To} {
		if s := c.Query(key); s != "" {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return f, fiber.NewError(400, key+" must be RFC3339")
			}
			*dst = t
		}
	}
	if f.From.After(f.To) || f.To.Sub(f.From) > 31*24*time.Hour {
		return f, fiber.NewError(400, "Choose a time range of up to 31 days")
	}
	for key, dst := range map[string]*int{"limit": &f.Limit, "offset": &f.Offset} {
		if s := c.Query(key); s != "" {
			n, e := strconv.Atoi(s)
			if e != nil || n < 0 {
				return f, fiber.NewError(400, "Invalid pagination")
			}
			*dst = n
		}
	}
	if f.Limit < 1 || f.Limit > 100 || f.Offset > 100000 {
		return f, fiber.NewError(400, "Pagination out of range")
	}
	return f, nil
}
