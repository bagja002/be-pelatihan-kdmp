package router

import (
	"knmp-backend/internal/handler"
	"knmp-backend/internal/middleware"
	"knmp-backend/internal/repository"
	"knmp-backend/internal/service"
	"knmp-backend/internal/storage"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func newPelatihHandler(db *gorm.DB, store *storage.Store) *handler.PelatihHandler {
	repo := repository.NewPelatihRepository(db)
	svc := service.NewPelatihService(repo, store)
	return handler.NewPelatihHandler(svc, store)
}

// RegisterPelatihPublicRoutes: registrasi mandiri (publik).
// Rate limit sengaja dilepas: di belakang reverse proxy semua klien terlihat
// satu IP, sehingga limiter per-IP memblokir massal (429) saat ramai.
func RegisterPelatihPublicRoutes(api fiber.Router, db *gorm.DB, store *storage.Store) {
	h := newPelatihHandler(db, store)

	group := api.Group("/pelatih")
	group.Post("/register", h.Register)
	group.Post("/lookup", h.Lookup)  // cari data by NIP (edit mandiri)
	group.Put("/self", h.UpdateSelf) // perbarui data sendiri by NIP
}

// RegisterPelatihAdminRoutes: kelola pelatih untuk admin/super_admin.
func RegisterPelatihAdminRoutes(protected fiber.Router, db *gorm.DB, store *storage.Store) {
	h := newPelatihHandler(db, store)
	adminOnly := middleware.RequireRole("super_admin", "admin")

	group := protected.Group("/pelatih")
	group.Get("/", adminOnly, h.List)
	group.Get("/export", adminOnly, h.Export)
	group.Get("/sertifikat/:sertifikatID/berkas", adminOnly, h.DownloadSertifikat)
	group.Get("/:id", adminOnly, h.GetByID)
	group.Get("/:id/cv", adminOnly, h.DownloadCV)
	group.Put("/:id", adminOnly, h.AdminUpdate)
	group.Delete("/:id", adminOnly, h.Delete)
}
