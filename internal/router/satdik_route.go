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

	admin := group.Group("", middleware.RequireRole("super_admin"))
	admin.Post("/", h.Create)
	admin.Put("/:id", h.Update)
	admin.Delete("/:id", h.Delete)
}
