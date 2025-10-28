package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"api-gateway/pkg/clients"
)

func AuthMiddleware(authClient *clients.AuthClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		token := parts[1]

		// Validate token with auth service
		resp, err := authClient.ValidateToken(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		if !resp.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Set user info in context
		c.Locals("userID", resp.UserId)
		c.Locals("username", resp.Username)

		return c.Next()
	}
}

// OptionalAuthMiddleware is similar to AuthMiddleware but doesn't require authentication
func OptionalAuthMiddleware(authClient *clients.AuthClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				resp, err := authClient.ValidateToken(c.Context(), token)
				if err == nil && resp.Valid {
					c.Locals("userID", resp.UserId)
					c.Locals("username", resp.Username)
				}
			}
		}
		return c.Next()
	}
}