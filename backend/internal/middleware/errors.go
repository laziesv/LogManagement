package middleware

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	code := 500
	message := "Internal server error"
	var fe *fiber.Error
	if errors.As(err, &fe) {
		code = fe.Code
		message = fe.Message
	}
	if code >= 500 {
		log.Printf("request failed: %v", err)
	}
	return c.Status(code).JSON(fiber.Map{"error": message})
}
