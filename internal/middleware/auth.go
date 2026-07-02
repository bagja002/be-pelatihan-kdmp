// Package middleware holds cross-cutting HTTP middleware.
package middleware

import (
	"strings"

	"knmp-backend/pkg/response"
	"knmp-backend/pkg/token"

	"github.com/gofiber/fiber/v2"
)

// Locals keys used to pass the authenticated identity down the request chain.
const (
	LocalUserID = "user_id"
	LocalRole   = "user_role"
)

// RequireAuth validates the Bearer access token and stores the user id and
// role in c.Locals. Requests without a valid token are rejected with 401.
func RequireAuth(tm *token.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(fiber.HeaderAuthorization)
		const prefix = "Bearer "
		if len(authHeader) <= len(prefix) || !strings.EqualFold(authHeader[:len(prefix)], prefix) {
			return response.Unauthorized(c, "missing or malformed authorization header")
		}

		claims, err := tm.Parse(strings.TrimSpace(authHeader[len(prefix):]), token.Access)
		if err != nil {
			return response.Unauthorized(c, "invalid or expired token")
		}

		c.Locals(LocalUserID, claims.UserID)
		c.Locals(LocalRole, claims.Role)
		return c.Next()
	}
}

// RequireRole authorizes only requests whose authenticated role is in roles.
// It must be chained after RequireAuth.
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals(LocalRole).(string)
		for _, allowed := range roles {
			if role == allowed {
				return c.Next()
			}
		}
		return response.Forbidden(c, "insufficient permissions")
	}
}
