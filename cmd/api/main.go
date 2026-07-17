package main

import (
	"log"
	"time"

	"knmp-backend/internal/config"
	"knmp-backend/internal/database"
	"knmp-backend/internal/router"
	"knmp-backend/internal/storage"
	"knmp-backend/pkg/crypto"
	"knmp-backend/pkg/response"
	"knmp-backend/pkg/token"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// Install the field-level encryption key before any DB access.
	if err := crypto.SetKey(cfg.EncryptionKey); err != nil {
		log.Fatalf("crypto init: %v", err)
	}

	db := database.Connect(cfg)
	if err := database.Migrate(db); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	tm := token.NewManager(cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL, cfg.PesertaTokenTTL, "knmp-backend")

	app := fiber.New(fiber.Config{
		AppName:   "knmp-backend",
		BodyLimit: cfg.BodyLimit,
		// StreamRequestBody: body tidak ditampung penuh di RAM — multipart
		// besar (bahan ajar ratusan MB) di-spill fasthttp ke temp dir OS.
		StreamRequestBody:     true,
		ReadTimeout:           cfg.ReadTimeout,
		WriteTimeout:          cfg.WriteTimeout,
		IdleTimeout:           60 * time.Second,
		DisableStartupMessage: true,
		ErrorHandler:          response.ErrorHandler(cfg.IsProduction()),
	})

	// ── Security & observability middleware (order matters) ──
	app.Use(requestid.New())
	app.Use(recover.New())

	// Security headers: HSTS, CSP, X-Frame-Options, nosniff, referrer policy…
	app.Use(helmet.New(helmet.Config{
		XSSProtection:         "0",
		ContentSecurityPolicy: "default-src 'self'; frame-ancestors 'none'",
		ReferrerPolicy:        "no-referrer",
		HSTSMaxAge:            31536000,
		HSTSPreloadEnabled:    true,
	}))

	// Strict CORS: explicit origin allow-list, never a wildcard with creds.
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip} ${status} ${method} ${path} rid=${locals:requestid} ${latency}\n",
	}))

	// Global rate limiter (per IP). Nonaktif bila RATE_LIMIT_MAX <= 0 —
	// di belakang reverse proxy semua klien terlihat 1 IP sehingga limiter
	// per-IP justru memblokir massal; aktifkan lagi hanya bila proxy sudah
	// meneruskan IP asli (X-Forwarded-For) dan Fiber dikonfigurasi membacanya.
	if cfg.RateLimitMax > 0 {
		app.Use(limiter.New(limiter.Config{
			Max:        cfg.RateLimitMax,
			Expiration: cfg.RateLimitWindow,
			LimitReached: func(c *fiber.Ctx) error {
				return response.TooManyRequests(c, "rate limit exceeded")
			},
		}))
	}

	store := storage.New(cfg.UploadDir, cfg.MaxUploadBytes)
	// Store terpisah untuk bahan ajar: root sama, batas ukuran lebih besar.
	bahanStore := storage.New(cfg.UploadDir, cfg.MaxBahanAjarBytes)
	router.SetupRoutes(app, db, tm, store, bahanStore)

	addr := ":" + cfg.AppPort
	log.Printf("server (%s) listening on %s", cfg.AppEnv, addr)
	log.Fatal(app.Listen(addr))
}
