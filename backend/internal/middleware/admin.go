package middleware

import (
	"github.com/gofiber/fiber/v2"
	"bobastream/internal/models"
	"bobastream/internal/utils"
)

// AdminOnly middleware checks if user is admin
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user role from context (set by AuthRequired middleware)
		role, ok := c.Locals("user_role").(models.UserRole)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required")
		}

		if role != models.RoleAdmin {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Admin access required")
		}

		return c.Next()
	}
}