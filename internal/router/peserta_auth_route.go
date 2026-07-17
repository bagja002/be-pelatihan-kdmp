package router

import (
	"knmp-backend/internal/handler"
	"knmp-backend/internal/repository"
	"knmp-backend/internal/service"
	"knmp-backend/pkg/token"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterPesertaAuthRoutes mounts /peserta-auth (public).
// Rate limit sengaja dilepas: di belakang reverse proxy semua klien terlihat
// satu IP, sehingga limiter per-IP memblokir massal (429) saat ramai.
func RegisterPesertaAuthRoutes(api fiber.Router, db *gorm.DB, tm *token.Manager) {
	satdikRepo := repository.NewSatdikRepository(db)
	pesertaRepo := repository.NewPesertaRepository(db)
	svc := service.NewPesertaAuthService(satdikRepo, pesertaRepo, tm)
	h := handler.NewPesertaAuthHandler(svc)

	group := api.Group("/peserta-auth")
	group.Post("/verify", h.Verify)
}
