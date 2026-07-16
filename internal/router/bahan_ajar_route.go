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

func newBahanAjarHandler(db *gorm.DB, store *storage.Store) *handler.BahanAjarHandler {
	repo := repository.NewBahanAjarRepository(db)
	svc := service.NewBahanAjarService(repo, store)
	return handler.NewBahanAjarHandler(svc, store)
}

// RegisterBahanAjarPublicRoutes: daftar & unduh bahan ajar (publik, tanpa login).
func RegisterBahanAjarPublicRoutes(api fiber.Router, db *gorm.DB, store *storage.Store) {
	h := newBahanAjarHandler(db, store)
	group := api.Group("/bahan-ajar")
	group.Get("/", h.List)
	group.Get("/:id/berkas/:jenis", h.Download)
}

// RegisterBahanAjarAdminRoutes: kelola bahan ajar untuk admin/super_admin.
func RegisterBahanAjarAdminRoutes(protected fiber.Router, db *gorm.DB, store *storage.Store) {
	h := newBahanAjarHandler(db, store)
	adminOnly := middleware.RequireRole("super_admin", "admin")

	group := protected.Group("/bahan-ajar")
	group.Post("/kategori", adminOnly, h.CreateKategori)
	group.Put("/kategori/:id", adminOnly, h.UpdateKategori)
	group.Delete("/kategori/:id", adminOnly, h.DeleteKategori)
	group.Post("/", adminOnly, h.CreateItem)
	group.Put("/:id", adminOnly, h.UpdateItem)
	group.Delete("/:id", adminOnly, h.DeleteItem)
}
