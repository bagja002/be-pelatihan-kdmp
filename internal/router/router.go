package router

import (
	"knmp-backend/internal/middleware"
	"knmp-backend/internal/storage"
	"knmp-backend/pkg/token"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes registers every route group. After scaffolding a new entity
// with the generator, register it here with a single line — either on the
// public `api` group or the authenticated `protected` group.
func SetupRoutes(app *fiber.App, db *gorm.DB, tm *token.Manager, store *storage.Store) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api/v1")

	// Authentication endpoints (public).
	RegisterAuthRoutes(api, db, tm)

	// Peserta public flow: verify NIK -> short-lived token, then self get/update.
	// Registered on the public group BEFORE the protected /peserta/:id routes so
	// that the literal /peserta/me is matched for peserta tokens.
	RegisterPesertaAuthRoutes(api, db, tm)
	RegisterPesertaSelfRoutes(api, db, tm)

	// Pelatih SDM: registrasi mandiri publik (link terbuka).
	RegisterPelatihPublicRoutes(api, db, store)

	// Protected group: every route below requires a valid JWT access token.
	protected := api.Group("", middleware.RequireAuth(tm))

	// ────────────────────────────────────────────────────────────
	// Register entity routes below (one line per entity).
	// Use `protected` to require auth, or `api` for public access.
	RegisterSatdikRoutes(protected, db)
	RegisterUserRoutes(protected, db)
	RegisterPesertaImportRoutes(protected, db)
	RegisterPesertaAdminRoutes(protected, db)
	RegisterPelatihAdminRoutes(protected, db, store)
	// ────────────────────────────────────────────────────────────
}
