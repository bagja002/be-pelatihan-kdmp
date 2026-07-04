package handler

import (
	"errors"
	"strconv"

	"knmp-backend/internal/dto"
	"knmp-backend/internal/repository"
	"knmp-backend/internal/service"
	"knmp-backend/internal/storage"
	"knmp-backend/pkg/response"
	"knmp-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type PelatihHandler struct {
	service service.PelatihService
	store   *storage.Store
}

func NewPelatihHandler(s service.PelatihService, store *storage.Store) *PelatihHandler {
	return &PelatihHandler{service: s, store: store}
}

// Register — endpoint publik multipart/form-data untuk registrasi mandiri.
// Field teks: namaLengkap, nip, pendidikan, unitKerja, jabatan.
// Berkas: cv (satu), sertifikat[] (banyak) + sertifikatNama[] (paralel indeks).
func (h *PelatihHandler) Register(c *fiber.Ctx) error {
	req := dto.RegisterPelatihRequest{
		NamaLengkap: c.FormValue("namaLengkap"),
		NIP:         c.FormValue("nip"),
		Pendidikan:  c.FormValue("pendidikan"),
		UnitKerja:   c.FormValue("unitKerja"),
		Jabatan:     c.FormValue("jabatan"),
	}
	if errs := validator.Validate(&req); errs != nil {
		return response.ValidationError(c, errs)
	}

	in := service.RegisterPelatihInput{
		NamaLengkap: req.NamaLengkap,
		NIP:         req.NIP,
		Pendidikan:  req.Pendidikan,
		UnitKerja:   req.UnitKerja,
		Jabatan:     req.Jabatan,
	}

	// CV opsional.
	if cv, err := c.FormFile("cv"); err == nil && cv != nil {
		in.CV = cv
	}

	// Sertifikat: pasangkan nama[i] dengan berkas[i].
	if form, err := c.MultipartForm(); err == nil && form != nil {
		names := form.Value["sertifikatNama[]"]
		files := form.File["sertifikat[]"]
		for i, f := range files {
			nama := ""
			if i < len(names) {
				nama = names[i]
			}
			if nama == "" {
				continue // lewati baris tanpa nama
			}
			in.Sertifikat = append(in.Sertifikat, service.CertUpload{Nama: nama, File: f})
		}
	}

	p, err := h.service.Register(in)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNIPExists):
			return response.Conflict(c, "NIP sudah terdaftar")
		case errors.Is(err, storage.ErrFileType):
			return response.BadRequest(c, "tipe berkas tidak diizinkan (hanya PDF/JPG/PNG)")
		case errors.Is(err, storage.ErrFileTooLarge):
			return response.BadRequest(c, "ukuran berkas melebihi batas")
		default:
			return response.InternalError(c, err)
		}
	}
	return response.Created(c, "pendaftaran berhasil", fiber.Map{"id": p.ID})
}

func (h *PelatihHandler) List(c *fiber.Ctx) error {
	data, err := h.service.List()
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, "pelatih list", data)
}

func (h *PelatihHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	data, err := h.service.Get(uint(id))
	if err != nil {
		return response.NotFound(c, "pelatih not found")
	}
	return response.OK(c, "pelatih detail", data)
}

func (h *PelatihHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	if err := h.service.Delete(uint(id)); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, "pelatih dihapus", nil)
}

func (h *PelatihHandler) DownloadCV(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	p, err := h.service.Get(uint(id))
	if err != nil {
		return response.NotFound(c, "pelatih not found")
	}
	if p.CV == "" {
		return response.NotFound(c, "CV tidak ada")
	}
	return c.Download(h.store.Path(p.CV))
}

func (h *PelatihHandler) DownloadSertifikat(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("sertifikatID"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	s, err := h.service.Sertifikat(uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrSertifikatNotFound) {
			return response.NotFound(c, "sertifikat not found")
		}
		return response.InternalError(c, err)
	}
	if s.Berkas == "" {
		return response.NotFound(c, "berkas sertifikat tidak ada")
	}
	return c.Download(h.store.Path(s.Berkas))
}
