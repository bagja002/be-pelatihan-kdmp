package handler

import (
	"errors"
	"strconv"

	"knmp-backend/internal/dto"
	"knmp-backend/internal/middleware"
	"knmp-backend/internal/service"
	"knmp-backend/pkg/response"
	"knmp-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type PesertaHandler struct{ service service.PesertaService }

func NewPesertaHandler(s service.PesertaService) *PesertaHandler { return &PesertaHandler{service: s} }

func ctxAuth(c *fiber.Ctx) (uint, string) {
	uid, _ := c.Locals(middleware.LocalUserID).(uint)
	role, _ := c.Locals(middleware.LocalRole).(string)
	return uid, role
}

func (h *PesertaHandler) List(c *fiber.Ctx) error {
	uid, role := ctxAuth(c)
	data, err := h.service.List(role, uid)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, "peserta list", data)
}

func (h *PesertaHandler) GetByID(c *fiber.Ctx) error {
	uid, role := ctxAuth(c)
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	data, err := h.service.Get(role, uid, uint(id))
	if err != nil {
		if errors.Is(err, service.ErrForbiddenScope) {
			return response.Forbidden(c, "di luar cakupan Anda")
		}
		return response.NotFound(c, "peserta not found")
	}
	return response.OK(c, "peserta detail", data)
}

func (h *PesertaHandler) Create(c *fiber.Ctx) error {
	uid, role := ctxAuth(c)
	var req dto.CreatePesertaRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if errs := validator.Validate(&req); errs != nil {
		return response.ValidationError(c, errs)
	}
	data, err := h.service.Create(role, uid, &req)
	if err != nil {
		if errors.Is(err, service.ErrForbiddenScope) {
			return response.Forbidden(c, "di luar cakupan Anda")
		}
		return response.InternalError(c, err)
	}
	return response.Created(c, "peserta created", data)
}

func (h *PesertaHandler) Update(c *fiber.Ctx) error {
	uid, role := ctxAuth(c)
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	var req dto.UpdatePesertaRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if errs := validator.Validate(&req); errs != nil {
		return response.ValidationError(c, errs)
	}
	data, err := h.service.Update(role, uid, uint(id), &req)
	if err != nil {
		if errors.Is(err, service.ErrForbiddenScope) {
			return response.Forbidden(c, "di luar cakupan Anda")
		}
		return response.NotFound(c, "peserta not found")
	}
	return response.OK(c, "peserta updated", data)
}

func (h *PesertaHandler) Delete(c *fiber.Ctx) error {
	uid, role := ctxAuth(c)
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	if err := h.service.Delete(role, uid, uint(id)); err != nil {
		if errors.Is(err, service.ErrForbiddenScope) {
			return response.Forbidden(c, "di luar cakupan Anda")
		}
		return response.InternalError(c, err)
	}
	return response.OK(c, "peserta deleted", nil)
}

// GetSelf & UpdateSelf dipakai route peserta (role=peserta).
func (h *PesertaHandler) GetSelf(c *fiber.Ctx) error {
	uid, _ := ctxAuth(c)
	data, err := h.service.GetSelf(uid)
	if err != nil {
		return response.NotFound(c, "data tidak ditemukan")
	}
	return response.OK(c, "peserta me", data)
}

func (h *PesertaHandler) UpdateSelf(c *fiber.Ctx) error {
	uid, _ := ctxAuth(c)
	var req dto.UpdateSelfRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if errs := validator.Validate(&req); errs != nil {
		return response.ValidationError(c, errs)
	}
	data, err := h.service.UpdateSelf(uid, &req)
	if err != nil {
		return response.NotFound(c, "data tidak ditemukan")
	}
	return response.OK(c, "data tersimpan", data)
}
