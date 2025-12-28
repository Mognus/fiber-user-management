package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	apperrors "template/modules/core/pkg/errors"
)

// RequireAuth middleware ensures user is authenticated
func (m *Module) RequireAuth(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		appErr := apperrors.Unauthorized("Authentication required")
		return c.Status(appErr.Code).JSON(appErr)
	}
	return c.Next()
}

// RequireAdmin middleware ensures user is authenticated and has admin role
func (m *Module) RequireAdmin(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		appErr := apperrors.Unauthorized("Authentication required")
		return c.Status(appErr.Code).JSON(appErr)
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	role := claims["role"].(string)

	if role != string(RoleAdmin) {
		appErr := apperrors.Forbidden("Admin access required")
		return c.Status(appErr.Code).JSON(appErr)
	}

	return c.Next()
}

// GetUserIDFromContext extracts user ID from JWT token in context
func GetUserIDFromContext(c *fiber.Ctx) (uint, error) {
	user := c.Locals("user")
	if user == nil {
		return 0, apperrors.Unauthorized("Authentication required")
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	return userID, nil
}

// GetUserRoleFromContext extracts user role from JWT token in context
func GetUserRoleFromContext(c *fiber.Ctx) (UserRole, error) {
	user := c.Locals("user")
	if user == nil {
		return "", apperrors.Unauthorized("Authentication required")
	}

	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	role := UserRole(claims["role"].(string))

	return role, nil
}
