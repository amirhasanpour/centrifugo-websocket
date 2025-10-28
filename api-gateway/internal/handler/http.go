package handler

import (
	"api-gateway/internal/models"
	"api-gateway/pkg/clients"
	"time"

	"github.com/gofiber/fiber/v2"
)

type HTTPHandler struct {
	authClient *clients.AuthClient
}

func NewHTTPHandler(authClient *clients.AuthClient) *HTTPHandler {
	return &HTTPHandler{
		authClient: authClient,
	}
}

func (h *HTTPHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	// Basic validation
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Username, email, and password are required",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Password must be at least 6 characters",
		})
	}

	// Call auth service
	resp, err := h.authClient.Register(c.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Registration failed",
			Message: err.Error(),
		})
	}

	return c.JSON(models.AuthResponse{
		User: models.UserResponse{
			ID:        resp.User.Id,
			Username:  resp.User.Username,
			Email:     resp.User.Email,
			CreatedAt: resp.User.CreatedAt,
			UpdatedAt: resp.User.UpdatedAt,
		},
		Token: resp.Token,
	})
}

func (h *HTTPHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Email and password are required",
		})
	}

	// Call auth service
	resp, err := h.authClient.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Login failed",
			Message: "Invalid email or password",
		})
	}

	return c.JSON(models.AuthResponse{
		User: models.UserResponse{
			ID:        resp.User.Id,
			Username:  resp.User.Username,
			Email:     resp.User.Email,
			CreatedAt: resp.User.CreatedAt,
			UpdatedAt: resp.User.UpdatedAt,
		},
		Token: resp.Token,
	})
}

func (h *HTTPHandler) ValidateToken(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if req.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Token is required",
		})
	}

	resp, err := h.authClient.ValidateToken(c.Context(), req.Token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Token validation failed",
		})
	}

	return c.JSON(models.ValidateTokenResponse{
		Valid:    resp.Valid,
		UserID:   resp.UserId,
		Username: resp.Username,
	})
}

func (h *HTTPHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	resp, err := h.authClient.GetUser(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "Failed to get user profile",
		})
	}

	return c.JSON(models.UserResponse{
		ID:        resp.User.Id,
		Username:  resp.User.Username,
		Email:     resp.User.Email,
		CreatedAt: resp.User.CreatedAt,
		UpdatedAt: resp.User.UpdatedAt,
	})
}

func (h *HTTPHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"service":   "api-gateway",
		"timestamp": time.Now(),
	})
}