package router

import (
	"knmp-backend/internal/handler"
	"knmp-backend/internal/middleware"
	"knmp-backend/internal/repository"
	"knmp-backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterUserRoutes mounts the /users group (super_admin only).
func RegisterUserRoutes(protected fiber.Router, db *gorm.DB) {
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	group := protected.Group("/users", middleware.RequireRole("super_admin"))
	group.Post("/", h.Create)
	group.Get("/", h.List)
	group.Get("/:id", h.GetByID)
	group.Put("/:id", h.Update)
	group.Delete("/:id", h.Delete)
}
