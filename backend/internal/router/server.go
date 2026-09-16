package router

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"logmanagement/backend/internal/handler"
	"logmanagement/backend/internal/middleware"
	"logmanagement/backend/internal/repository"
)

type Config struct {
	Origin        string
	SecureCookies bool
}

func NewServer(store repository.Store, cfg Config) *fiber.App {
	server := fiber.New(fiber.Config{
		AppName: "LogManagement", BodyLimit: 2 * 1024 * 1024,
		ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second,
		ErrorHandler: middleware.ErrorHandler,
	})
	server.Use(recover.New())
	server.Use(middleware.AccessLog)
	server.Use(middleware.Security(cfg.Origin))
	register(server, handler.New(store, cfg.SecureCookies), store)
	return server
}
