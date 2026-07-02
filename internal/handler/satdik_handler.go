package handler

import (
	"strconv"

	"knmp-backend/internal/dto"
	"knmp-backend/internal/service"
	"knmp-backend/pkg/response"
	"knmp-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type SatdikHandler struct{ service service.SatdikService }

func NewSatdikHandler(s service.SatdikService) *SatdikHandler { return &SatdikHandler{service: s} }

func (h *SatdikHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateSatdikRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if errs := validator.Validate(&req); errs != nil {
		return response.ValidationError(c, errs)
	}
	data, err := h.service.Create(&req)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.Created(c, "satdik created", data)
}

func (h *SatdikHandler) List(c *fiber.Ctx) error {
	data, err := h.service.GetAll()
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, "satdik list", data)
}

func (h *SatdikHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	data, err := h.service.GetByID(uint(id))
	if err != nil {
		return response.NotFound(c, "satdik not found")
	}
	return response.OK(c, "satdik detail", data)
}

func (h *SatdikHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	var req dto.UpdateSatdikRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if errs := validator.Validate(&req); errs != nil {
		return response.ValidationError(c, errs)
	}
	data, err := h.service.Update(uint(id), &req)
	if err != nil {
		return response.NotFound(c, "satdik not found")
	}
	return response.OK(c, "satdik updated", data)
}

func (h *SatdikHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	if err := h.service.Delete(uint(id)); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, "satdik deleted", nil)
}
