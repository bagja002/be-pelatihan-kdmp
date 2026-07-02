package router

import (
	"time"

	"knmp-backend/internal/handler"
	"knmp-backend/internal/middleware"
	"knmp-backend/internal/repository"
	"knmp-backend/internal/service"
	"knmp-backend/pkg/response"
	"knmp-backend/pkg/token"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"gorm.io/gorm"
)

// RegisterAuthRoutes mounts the /auth group. Credential endpoints get a
// strict rate limit as brute-force protection, and /auth/me is protected by
// the JWT auth middleware.
func RegisterAuthRoutes(api fiber.Router, db *gorm.DB, tm *token.Manager) {
	repo := repository.NewUserRepository(db)
	svc := service.NewAuthService(repo, tm)
	h := handler.NewAuthHandler(svc)

	// Stricter limiter for credential endpoints (anti brute-force).
	strict := limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return response.TooManyRequests(c, "too many attempts, please try again later")
		},
	})

	auth := api.Group("/auth")
	auth.Post("/login", strict, h.Login)
	auth.Post("/refresh", strict, h.Refresh)
	auth.Get("/me", middleware.RequireAuth(tm), h.Me)
}
