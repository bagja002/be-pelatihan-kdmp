package router

import (
	"knmp-backend/internal/middleware"
	"knmp-backend/pkg/token"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes registers every route group. After scaffolding a new entity
// with the generator, register it here with a single line — either on the
// public `api` group or the authenticated `protected` group.
func SetupRoutes(app *fiber.App, db *gorm.DB, tm *token.Manager) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api/v1")

	// Authentication endpoints (public).
	RegisterAuthRoutes(api, db, tm)

	// Protected group: every route below requires a valid JWT access token.
	protected := api.Group("", middleware.RequireAuth(tm))

	// ────────────────────────────────────────────────────────────
	// Register entity routes below (one line per entity).
	// Use `protected` to require auth, or `api` for public access.
	RegisterProductRoutes(protected, db)
	// ────────────────────────────────────────────────────────────
}
