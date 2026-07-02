package handler

import (
	"knmp-backend/internal/importer"
	"knmp-backend/internal/service"
	"knmp-backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type PesertaImportHandler struct{ service service.PesertaImportService }

func NewPesertaImportHandler(s service.PesertaImportService) *PesertaImportHandler {
	return &PesertaImportHandler{service: s}
}

func (h *PesertaImportHandler) Import(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return response.BadRequest(c, "file .xlsx wajib diunggah pada field 'file'")
	}
	f, err := fileHeader.Open()
	if err != nil {
		return response.InternalError(c, err)
	}
	defer f.Close()

	rows, err := importer.Parse(f)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	summary, err := h.service.Import(rows)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, "import selesai", summary)
}
