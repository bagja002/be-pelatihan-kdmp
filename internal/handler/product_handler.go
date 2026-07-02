package handler

import (
	"strconv"

	"knmp-backend/internal/dto"
	"knmp-backend/internal/service"
	"knmp-backend/pkg/response"
	"knmp-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// ProductHandler exposes HTTP handlers for Product.
type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler builds a ProductHandler.
func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateProductRequest
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
	return response.Created(c, "product created", data)
}

func (h *ProductHandler) List(c *fiber.Ctx) error {
	data, err := h.service.GetAll()
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, "product list", data)
}

func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	data, err := h.service.GetByID(uint(id))
	if err != nil {
		return response.NotFound(c, "product not found")
	}
	return response.OK(c, "product detail", data)
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	var req dto.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if errs := validator.Validate(&req); errs != nil {
		return response.ValidationError(c, errs)
	}
	data, err := h.service.Update(uint(id), &req)
	if err != nil {
		return response.NotFound(c, "product not found")
	}
	return response.OK(c, "product updated", data)
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	if err := h.service.Delete(uint(id)); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, "product deleted", nil)
}
