package handler

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"knmp-backend/internal/dto"
	"knmp-backend/internal/repository"
	"knmp-backend/internal/service"
	"knmp-backend/internal/storage"
	"knmp-backend/pkg/response"
	"knmp-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type BahanAjarHandler struct {
	service service.BahanAjarService
	store   *storage.Store
}

func NewBahanAjarHandler(s service.BahanAjarService, store *storage.Store) *BahanAjarHandler {
	return &BahanAjarHandler{service: s, store: store}
}

// List — publik: semua kategori beserta itemnya (tanpa path berkas).
func (h *BahanAjarHandler) List(c *fiber.Ctx) error {
	kategori, err := h.service.ListKategori()
	if err != nil {
		return response.InternalError(c, err)
	}
	out := make([]dto.BahanAjarKategoriResponse, 0, len(kategori))
	for _, k := range kategori {
		out = append(out, dto.ToBahanAjarKategoriResponse(k))
	}
	return response.OK(c, "daftar bahan ajar", out)
}

// Download — publik: unduh berkas pdf|ppt milik satu item.
func (h *BahanAjarHandler) Download(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "id tidak valid")
	}
	item, err := h.service.Item(uint(id))
	if err != nil {
		return h.itemError(c, err)
	}
	var rel string
	switch c.Params("jenis") {
	case "pdf":
		rel = item.FilePdf
	case "ppt":
		rel = item.FilePpt
	default:
		return response.BadRequest(c, "jenis berkas harus pdf atau ppt")
	}
	if rel == "" {
		return response.NotFound(c, "berkas belum diunggah")
	}
	return c.Download(h.store.Path(rel), fmt.Sprintf("%s%s", item.Judul, filepath.Ext(rel)))
}

// CreateKategori — admin: tambah kategori (body JSON).
func (h *BahanAjarHandler) CreateKategori(c *fiber.Ctx) error {
	var req dto.BahanAjarKategoriRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "permintaan tidak valid")
	}
	if errs := validator.Validate(&req); errs != nil {
		return response.ValidationError(c, errs)
	}
	k, err := h.service.CreateKategori(req.Nama, req.Urutan)
	if err != nil {
		return h.kategoriError(c, err)
	}
	return response.Created(c, "kategori dibuat", k)
}

// UpdateKategori — admin: ubah nama/urutan kategori.
func (h *BahanAjarHandler) UpdateKategori(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "id tidak valid")
	}
	var req dto.BahanAjarKategoriRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "permintaan tidak valid")
	}
	if errs := validator.Validate(&req); errs != nil {
		return response.ValidationError(c, errs)
	}
	k, err := h.service.UpdateKategori(uint(id), req.Nama, req.Urutan)
	if err != nil {
		return h.kategoriError(c, err)
	}
	return response.OK(c, "kategori diperbarui", k)
}

// DeleteKategori — admin: hapus kategori kosong.
func (h *BahanAjarHandler) DeleteKategori(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "id tidak valid")
	}
	if err := h.service.DeleteKategori(uint(id)); err != nil {
		return h.kategoriError(c, err)
	}
	return response.OK(c, "kategori dihapus", nil)
}

// CreateItem — admin: tambah item (multipart: kategoriId, judul, urutan?, filePdf?, filePpt?).
func (h *BahanAjarHandler) CreateItem(c *fiber.Ctx) error {
	kategoriID, _ := strconv.ParseUint(c.FormValue("kategoriId"), 10, 64)
	judul := strings.TrimSpace(c.FormValue("judul"))
	if kategoriID == 0 || judul == "" {
		return response.BadRequest(c, "kategoriId dan judul wajib diisi")
	}
	in := service.BahanAjarItemInput{KategoriID: uint(kategoriID), Judul: judul}
	fillItemInput(c, &in)

	item, err := h.service.CreateItem(in)
	if err != nil {
		return h.itemError(c, err)
	}
	return response.Created(c, "bahan ajar dibuat", dto.ToBahanAjarItemResponse(*item))
}

// UpdateItem — admin: ubah item (multipart; berkas yang dikirim menggantikan lama).
func (h *BahanAjarHandler) UpdateItem(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "id tidak valid")
	}
	judul := strings.TrimSpace(c.FormValue("judul"))
	if judul == "" {
		return response.BadRequest(c, "judul wajib diisi")
	}
	in := service.BahanAjarItemInput{Judul: judul}
	if v, err := strconv.ParseUint(c.FormValue("kategoriId"), 10, 64); err == nil {
		in.KategoriID = uint(v)
	}
	fillItemInput(c, &in)

	item, err := h.service.UpdateItem(uint(id), in)
	if err != nil {
		return h.itemError(c, err)
	}
	return response.OK(c, "bahan ajar diperbarui", dto.ToBahanAjarItemResponse(*item))
}

// DeleteItem — admin: hapus item beserta berkasnya.
func (h *BahanAjarHandler) DeleteItem(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "id tidak valid")
	}
	if err := h.service.DeleteItem(uint(id)); err != nil {
		return h.itemError(c, err)
	}
	return response.OK(c, "bahan ajar dihapus", nil)
}

// fillItemInput mengambil urutan & berkas opsional dari multipart form.
func fillItemInput(c *fiber.Ctx, in *service.BahanAjarItemInput) {
	if v := c.FormValue("urutan"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			in.Urutan = &n
		}
	}
	if f, err := c.FormFile("filePdf"); err == nil && f != nil {
		in.FilePdf = f
	}
	if f, err := c.FormFile("filePpt"); err == nil && f != nil {
		in.FilePpt = f
	}
}

func (h *BahanAjarHandler) kategoriError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, repository.ErrKategoriNotFound):
		return response.NotFound(c, "kategori tidak ditemukan")
	case errors.Is(err, service.ErrNamaKategoriDipakai):
		return response.Conflict(c, "nama kategori sudah dipakai")
	case errors.Is(err, service.ErrKategoriBerisi):
		return response.Conflict(c, "kategori masih berisi bahan ajar — kosongkan dulu")
	default:
		return response.InternalError(c, err)
	}
}

func (h *BahanAjarHandler) itemError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, repository.ErrKategoriNotFound):
		return response.NotFound(c, "kategori tidak ditemukan")
	case errors.Is(err, repository.ErrBahanAjarNotFound):
		return response.NotFound(c, "bahan ajar tidak ditemukan")
	case errors.Is(err, storage.ErrFileType):
		return response.BadRequest(c, "tipe berkas tidak diizinkan (slot PDF: .pdf; slot PPT: .ppt/.pptx)")
	case errors.Is(err, storage.ErrFileTooLarge):
		return response.BadRequest(c, "ukuran berkas melebihi batas")
	default:
		return response.InternalError(c, err)
	}
}
