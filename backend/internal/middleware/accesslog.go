package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
)

func AccessLog(c fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	status := c.Response().StatusCode()
	if err != nil {
		if fe, ok := err.(*fiber.Error); ok {
			status = fe.Code
		} else {
			status = 500
		}
	}
	log.Printf("%s %s -> %d (%s)", c.Method(), c.OriginalURL(), status, time.Since(start).Round(time.Millisecond))
	return err
}
