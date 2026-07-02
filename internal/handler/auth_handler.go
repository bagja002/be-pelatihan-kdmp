package handler

import (
	"errors"

	"knmp-backend/internal/dto"
	"knmp-backend/internal/entity"
	"knmp-backend/internal/middleware"
	"knmp-backend/internal/service"
	"knmp-backend/pkg/response"
	"knmp-backend/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler exposes authentication endpoints.
type AuthHandler struct {
	service service.AuthService
}

// NewAuthHandler builds an AuthHandler.
func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if errs := validator.Validate(&req); errs != nil {
		return response.ValidationError(c, errs)
	}

	u, err := h.service.Register(&req)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			return response.Conflict(c, "email already registered")
		}
		return response.InternalError(c, err)
	}
	return response.Created(c, "registered", toUserResponse(u))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if errs := validator.Validate(&req); errs != nil {
		return response.ValidationError(c, errs)
	}

	tokens, err := h.service.Login(&req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return response.Unauthorized(c, "invalid email or password")
		}
		return response.InternalError(c, err)
	}
	return response.OK(c, "login successful", tokens)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req dto.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if errs := validator.Validate(&req); errs != nil {
		return response.ValidationError(c, errs)
	}

	tokens, err := h.service.Refresh(req.RefreshToken)
	if err != nil {
		return response.Unauthorized(c, "invalid refresh token")
	}
	return response.OK(c, "token refreshed", tokens)
}

// Me returns the profile of the authenticated user (requires RequireAuth).
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, _ := c.Locals(middleware.LocalUserID).(uint)
	u, err := h.service.Profile(userID)
	if err != nil {
		return response.NotFound(c, "user not found")
	}
	return response.OK(c, "profile", toUserResponse(u))
}

func toUserResponse(u *entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    u.ID,
		Email: u.Email,
		Role:  u.Role,
		Phone: u.Phone,
	}
}
