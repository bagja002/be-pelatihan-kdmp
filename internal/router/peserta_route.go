package router

import (
	"knmp-backend/internal/handler"
	"knmp-backend/internal/middleware"
	"knmp-backend/internal/repository"
	"knmp-backend/internal/service"
	"knmp-backend/pkg/token"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// newPesertaHandler membangun handler peserta (dipakai admin & peserta routes).
func newPesertaHandler(db *gorm.DB) *handler.PesertaHandler {
	repo := repository.NewPesertaRepository(db)
	users := repository.NewUserRepository(db)
	svc := service.NewPesertaService(repo, users)
	return handler.NewPesertaHandler(svc)
}

// RegisterPesertaAdminRoutes: CRUD peserta untuk admin/super_admin (ter-scope).
func RegisterPesertaAdminRoutes(protected fiber.Router, db *gorm.DB) {
	h := newPesertaHandler(db)
	group := protected.Group("/peserta")
	group.Get("/", h.List)
	group.Post("/", h.Create)
	group.Get("/:id", h.GetByID)
	group.Put("/:id", h.Update)
	group.Delete("/:id", h.Delete)
}

// RegisterPesertaSelfRoutes: peserta mengambil & meng-update datanya sendiri.
// `api` publik; token peserta divalidasi RequireAuth.
func RegisterPesertaSelfRoutes(api fiber.Router, db *gorm.DB, tm *token.Manager) {
	h := newPesertaHandler(db)
	group := api.Group("/peserta", middleware.RequireAuth(tm))
	group.Get("/me", h.GetSelf)
	group.Put("/me", h.UpdateSelf)
}
