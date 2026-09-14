package router

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"

	"logmanagement/backend/internal/handler"
	"logmanagement/backend/internal/middleware"
	"logmanagement/backend/internal/repository"
)

// register is the route table: endpoints, middleware order, and named handlers.
func register(server *fiber.App, h *handler.Handler, store repository.Store) {
	auth := middleware.Authenticate(store)
	admin := middleware.RequireAdmin
	ingestAuth := middleware.IngestAuth(store)
	server.Get("/api/health", h.Health)
	server.Post("/api/auth/login", limiter.New(limiter.Config{Max: 10, Expiration: time.Minute}), h.Login)
	server.Get("/api/auth/me", auth, h.Me)
	server.Post("/api/auth/logout", auth, h.Logout)
	server.Post("/api/ingest", ingestAuth, h.Ingest)
	// Assignment-compatible alias for the same ingest handler.
	server.Post("/ingest", ingestAuth, h.Ingest)
	server.Get("/api/logs", auth, h.Logs)
	server.Get("/api/stats", auth, h.Stats)
	server.Get("/api/alerts", auth, h.Alerts)
	server.Post("/api/alerts/:id/acknowledge", auth, admin, h.Acknowledge)
	server.Get("/api/rule", auth, h.GetRule)
	server.Put("/api/rule", auth, admin, h.SetRule)
}
