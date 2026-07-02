package router

import (
	"knmp-backend/internal/handler"
	"knmp-backend/internal/middleware"
	"knmp-backend/internal/repository"
	"knmp-backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterSatdikRoutes wires the /satdik group. `protected` sudah menuntut JWT.
// Read (list/detail) terbuka untuk semua admin login; mutasi hanya super_admin.
func RegisterSatdikRoutes(protected fiber.Router, db *gorm.DB) {
	repo := repository.NewSatdikRepository(db)
	svc := service.NewSatdikService(repo)
	h := handler.NewSatdikHandler(svc)

	group := protected.Group("/satdik")
	group.Get("/", h.List)
	group.Get("/:id", h.GetByID)

	// Mutations restricted to super_admin via route-level middleware (a nested
	// group with middleware would leak the role check onto the read routes too).
	superAdmin := middleware.RequireRole("super_admin")
	group.Post("/", superAdmin, h.Create)
	group.Put("/:id", superAdmin, h.Update)
	group.Delete("/:id", superAdmin, h.Delete)
}
