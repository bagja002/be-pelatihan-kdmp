package router

import (
	"time"

	"knmp-backend/internal/handler"
	"knmp-backend/internal/repository"
	"knmp-backend/internal/service"
	"knmp-backend/pkg/response"
	"knmp-backend/pkg/token"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"gorm.io/gorm"
)

// RegisterPesertaAuthRoutes mounts /peserta-auth (public, strict rate limit).
func RegisterPesertaAuthRoutes(api fiber.Router, db *gorm.DB, tm *token.Manager) {
	satdikRepo := repository.NewSatdikRepository(db)
	pesertaRepo := repository.NewPesertaRepository(db)
	svc := service.NewPesertaAuthService(satdikRepo, pesertaRepo, tm)
	h := handler.NewPesertaAuthHandler(svc)

	strict := limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return response.TooManyRequests(c, "terlalu banyak percobaan, coba lagi nanti")
		},
	})

	group := api.Group("/peserta-auth")
	group.Post("/verify", strict, h.Verify)
}
