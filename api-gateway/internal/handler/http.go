package handler

import (
	"api-gateway/internal/models"
	"api-gateway/pkg/clients"
	"time"

	"github.com/gofiber/fiber/v2"
)

type HTTPHandler struct {
	authClient *clients.AuthClient
	chatClient *clients.ChatClient
}

func NewHTTPHandler(authClient *clients.AuthClient, chatClient *clients.ChatClient) *HTTPHandler {
	return &HTTPHandler{
		authClient: authClient,
		chatClient: chatClient,
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

func (h *HTTPHandler) CreateRoom(c *fiber.Ctx) error {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	userID := c.Locals("userID").(string)

	resp, err := h.chatClient.CreateRoom(c.Context(), req.Name, req.Description, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to create room",
			Message: err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"room": resp.Room,
	})
}

func (h *HTTPHandler) GetRooms(c *fiber.Ctx) error {
	resp, err := h.chatClient.GetRooms(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to get rooms",
			Message: err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"rooms": resp.Rooms,
	})
}

func (h *HTTPHandler) GetRoomMessages(c *fiber.Ctx) error {
	roomID := c.Params("roomId")
	limit := c.QueryInt("limit", 50)

	resp, err := h.chatClient.GetRoomMessages(c.Context(), roomID, int32(limit))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to get room messages",
			Message: err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"messages": resp.Messages,
	})
}

func (h *HTTPHandler) SendMessage(c *fiber.Ctx) error {
	req := models.SendMessageRequest{}
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if req.RoomID == "" || req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Room ID and content are required",
		})
	}

	if len(req.Content) > 1000 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Message content too long (max 1000 characters)",
		})
	}

	userID := c.Locals("userID").(string)

	resp, err := h.chatClient.SendMessage(c.Context(), req.RoomID, userID, req.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to send message",
			Message: err.Error(),
		})
	}

	return c.JSON(models.SendMessageResponse{
		Message: models.MessageResponse{
			ID:        resp.Message.Id,
			RoomID:    resp.Message.RoomId,
			UserID:    resp.Message.UserId,
			Username:  resp.Message.Username,
			Content:   resp.Message.Content,
			CreatedAt: resp.Message.CreatedAt,
		},
	})
}

func (h *HTTPHandler) JoinRoom(c *fiber.Ctx) error {
	req := models.JoinRoomRequest{}
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if req.RoomID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Room ID is required",
		})
	}

	userID := c.Locals("userID").(string)

	resp, err := h.chatClient.JoinRoom(c.Context(), req.RoomID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to join room",
			Message: err.Error(),
		})
	}

	return c.JSON(models.JoinRoomResponse{
		Success: resp.Success,
		RoomID:  resp.RoomId,
	})
}