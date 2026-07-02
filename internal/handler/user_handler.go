package handler

import (
	"errors"
	"strconv"

	"knmp-backend/internal/dto"
	"knmp-backend/internal/entity"
	"knmp-backend/internal/service"
	"knmp-backend/pkg/response"
	"knmp-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct{ service service.UserService }

func NewUserHandler(s service.UserService) *UserHandler { return &UserHandler{service: s} }

func userToResponse(u *entity.User) dto.UserResponse {
	return dto.UserResponse{ID: u.ID, Nama: u.Nama, Username: u.Username, Type: u.Type, IDSatdik: u.IDSatdik}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if errs := validator.Validate(&req); errs != nil {
		return response.ValidationError(c, errs)
	}
	u, err := h.service.Create(&req)
	if err != nil {
		if errors.Is(err, service.ErrUsernameTaken) {
			return response.Conflict(c, "username sudah dipakai")
		}
		return response.InternalError(c, err)
	}
	return response.Created(c, "user created", userToResponse(u))
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	items, err := h.service.GetAll()
	if err != nil {
		return response.InternalError(c, err)
	}
	out := make([]dto.UserResponse, 0, len(items))
	for i := range items {
		out = append(out, userToResponse(&items[i]))
	}
	return response.OK(c, "user list", out)
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	u, err := h.service.GetByID(uint(id))
	if err != nil {
		return response.NotFound(c, "user not found")
	}
	return response.OK(c, "user detail", userToResponse(u))
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if errs := validator.Validate(&req); errs != nil {
		return response.ValidationError(c, errs)
	}
	u, err := h.service.Update(uint(id), &req)
	if err != nil {
		return response.NotFound(c, "user not found")
	}
	return response.OK(c, "user updated", userToResponse(u))
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}
	if err := h.service.Delete(uint(id)); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, "user deleted", nil)
}
