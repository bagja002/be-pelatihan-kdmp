package handler

import (
	"errors"

	"knmp-backend/internal/dto"
	"knmp-backend/internal/service"
	"knmp-backend/pkg/response"
	"knmp-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type PesertaAuthHandler struct{ service service.PesertaAuthService }

func NewPesertaAuthHandler(s service.PesertaAuthService) *PesertaAuthHandler {
	return &PesertaAuthHandler{service: s}
}

func (h *PesertaAuthHandler) Verify(c *fiber.Ctx) error {
	var req dto.VerifyNIKRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if errs := validator.Validate(&req); errs != nil {
		return response.ValidationError(c, errs)
	}
	tokens, err := h.service.Verify(&req)
	if err != nil {
		if errors.Is(err, service.ErrPesertaVerification) {
			return response.NotFound(c, "NIK tidak ditemukan pada Satdik ini")
		}
		return response.InternalError(c, err)
	}
	return response.OK(c, "verifikasi berhasil", tokens)
}
