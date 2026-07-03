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

	// SECURITY: enforce an operator role on every CRUD route. Without this a
	// peserta-role token (role="peserta", uid=peserta.ID) passes RequireAuth
	// and reaches the scope logic, which does users.FindByID(uid). Because user
	// and peserta ids are independent auto-increment sequences that overlap, a
	// peserta could inherit a User's scope — a peserta whose id collides with
	// the super_admin (id 1) would gain full cross-satdik access. The check is
	// attached PER ROUTE (not on the group) so it does not leak onto the
	// separately-registered public /peserta/me self routes.
	adminOnly := middleware.RequireRole("super_admin", "admin")
	group.Get("/", adminOnly, h.List)
	group.Post("/", adminOnly, h.Create)
	group.Get("/:id", adminOnly, h.GetByID)
	group.Put("/:id", adminOnly, h.Update)
	group.Delete("/:id", adminOnly, h.Delete)
}

// RegisterPesertaSelfRoutes: peserta mengambil & meng-update datanya sendiri.
// `api` publik; token peserta divalidasi RequireAuth.
func RegisterPesertaSelfRoutes(api fiber.Router, db *gorm.DB, tm *token.Manager) {
	h := newPesertaHandler(db)
	group := api.Group("/peserta", middleware.RequireAuth(tm))
	group.Get("/me", h.GetSelf)
	group.Put("/me", h.UpdateSelf)
}

// RegisterPesertaImportRoutes: unggah xlsx (super_admin only).
func RegisterPesertaImportRoutes(protected fiber.Router, db *gorm.DB) {
	pesertaRepo := repository.NewPesertaRepository(db)
	satdikRepo := repository.NewSatdikRepository(db)
	svc := service.NewPesertaImportService(pesertaRepo, satdikRepo)
	h := handler.NewPesertaImportHandler(svc)

	// Route-level role check only — a group with RequireRole on the /peserta
	// prefix would leak the super_admin requirement onto the admin CRUD routes.
	protected.Post("/peserta/import", middleware.RequireRole("super_admin"), h.Import)
}
