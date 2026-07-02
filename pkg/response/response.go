package response

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// Response is the standard JSON envelope for every API reply.
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

// OK returns a 200 response with data.
func OK(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(Response{Success: true, Message: message, Data: data})
}

// Created returns a 201 response with data.
func Created(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Response{Success: true, Message: message, Data: data})
}

// BadRequest returns a 400 response.
func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{Success: false, Message: message})
}

// Unauthorized returns a 401 response.
func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{Success: false, Message: message})
}

// Forbidden returns a 403 response.
func Forbidden(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(Response{Success: false, Message: message})
}

// NotFound returns a 404 response.
func NotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{Success: false, Message: message})
}

// Conflict returns a 409 response.
func Conflict(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusConflict).JSON(Response{Success: false, Message: message})
}

// TooManyRequests returns a 429 response.
func TooManyRequests(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(Response{Success: false, Message: message})
}

// ValidationError returns a 422 response carrying field-level errors.
func ValidationError(c *fiber.Ctx, errs any) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{Success: false, Message: "validation failed", Errors: errs})
}

// InternalError logs the real error (with request id) for auditing and returns
// a generic 500 so internal details never leak to the client.
func InternalError(c *fiber.Ctx, err error) error {
	log.Printf("[ERROR] rid=%v method=%s path=%s error=%v",
		c.Locals("requestid"), c.Method(), c.OriginalURL(), err)
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Success: false,
		Message: "internal server error",
	})
}

// ErrorHandler is Fiber's global error handler. It maps *fiber.Error to the
// standard envelope and hides 5xx details in production while still logging
// them for auditing.
func ErrorHandler(production bool) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "internal server error"

		var fe *fiber.Error
		if e, ok := err.(*fiber.Error); ok {
			fe = e
			code = e.Code
			message = e.Message
		}

		if code >= 500 {
			log.Printf("[ERROR] rid=%v method=%s path=%s error=%v",
				c.Locals("requestid"), c.Method(), c.OriginalURL(), err)
			if production {
				message = "internal server error"
			} else if fe == nil {
				message = err.Error()
			}
		}

		return c.Status(code).JSON(Response{Success: false, Message: message})
	}
}
