package router

import (
	"knmp-backend/internal/handler"
	"knmp-backend/internal/repository"
	"knmp-backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterProductRoutes wires repository -> service -> handler and mounts
// the /products route group.
//
// To activate, add this single line inside SetupRoutes in router.go:
//
//	RegisterProductRoutes(api, db)
func RegisterProductRoutes(api fiber.Router, db *gorm.DB) {
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)

	group := api.Group("/products")
	group.Post("/", h.Create)
	group.Get("/", h.List)
	group.Get("/:id", h.GetByID)
	group.Put("/:id", h.Update)
	group.Delete("/:id", h.Delete)
}
