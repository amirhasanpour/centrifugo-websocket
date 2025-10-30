package service

import (
	"errors"
	"fmt"
	"time"

	"chat-service/internal/models"
	"chat-service/internal/repository"
	"chat-service/pkg/centrifugo"
)

var (
	ErrRoomNotFound     = errors.New("room not found")
	ErrNotRoomMember    = errors.New("user is not a member of this room")
	ErrInvalidMessage   = errors.New("invalid message content")
	ErrCentrifugoFailed = errors.New("failed to publish to centrifugo")
)

type ChatService interface {
	CreateRoom(name, description, userID string) (*models.Room, error)
	GetRooms() ([]*models.Room, error)
	JoinRoom(roomID, userID string) error
	SendMessage(roomID, userID, username, content string) (*models.Message, error)
	GetRoomMessages(roomID string, limit int) ([]*models.Message, error)
}

type chatService struct {
	repo          repository.ChatRepository
	centrifugo    *centrifugo.CentrifugoClient
}

func NewChatService(repo repository.ChatRepository, centrifugo *centrifugo.CentrifugoClient) ChatService {
	return &chatService{
		repo:       repo,
		centrifugo: centrifugo,
	}
}

func (s *chatService) CreateRoom(name, description, userID string) (*models.Room, error) {
	if name == "" {
		return nil, errors.New("room name is required")
	}

	room := &models.Room{
		Name:        name,
		Description: description,
		CreatedBy:   userID,
	}

	err := s.repo.CreateRoom(room)
	if err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}

	// Auto-join the room creator
	err = s.repo.AddRoomMember(room.ID, userID)
	if err != nil && err != repository.ErrAlreadyMember {
		return nil, fmt.Errorf("failed to add room creator as member: %w", err)
	}

	return room, nil
}

func (s *chatService) GetRooms() ([]*models.Room, error) {
	rooms, err := s.repo.GetAllRooms()
	if err != nil {
		return nil, fmt.Errorf("failed to get rooms: %w", err)
	}
	return rooms, nil
}

func (s *chatService) JoinRoom(roomID, userID string) error {
	// Check if room exists
	room, err := s.repo.GetRoomByID(roomID)
	if err != nil {
		if err == repository.ErrRoomNotFound {
			return ErrRoomNotFound
		}
		return fmt.Errorf("failed to get room: %w", err)
	}

	// Add user as room member
	err = s.repo.AddRoomMember(room.ID, userID)
	if err != nil {
		if err == repository.ErrAlreadyMember {
			return nil // Already a member, no error
		}
		return fmt.Errorf("failed to join room: %w", err)
	}

	// Notify room about new member via Centrifugo
	notification := map[string]any{
		"type":      "user_joined",
		"room_id":   roomID,
		"user_id":   userID,
		"timestamp": time.Now().Unix(),
	}

	go func() {
		if err := s.centrifugo.Publish("room:"+roomID, notification); err != nil {
			fmt.Printf("Failed to publish join notification: %v\n", err)
		}
	}()

	return nil
}

func (s *chatService) SendMessage(roomID, userID, username, content string) (*models.Message, error) {
	if content == "" {
		return nil, ErrInvalidMessage
	}

	// Check if room exists
	room, err := s.repo.GetRoomByID(roomID)
	if err != nil {
		if err == repository.ErrRoomNotFound {
			return nil, ErrRoomNotFound
		}
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	// Check if user is room member
	isMember, err := s.repo.IsRoomMember(room.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check room membership: %w", err)
	}
	if !isMember {
		return nil, ErrNotRoomMember
	}

	// Create message
	message := &models.Message{
		RoomID:   roomID,
		UserID:   userID,
		Username: username,
		Content:  content,
	}

	err = s.repo.CreateMessage(message)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Publish message to Centrifugo
	messageData := map[string]any{
		"id":         message.ID,
		"room_id":    message.RoomID,
		"user_id":    message.UserID,
		"username":   message.Username,
		"content":    message.Content,
		"created_at": message.CreatedAt.Format(time.RFC3339),
	}

	go func() {
		if err := s.centrifugo.Publish("room:"+roomID, messageData); err != nil {
			fmt.Printf("Failed to publish message to centrifugo: %v\n", err)
		}
	}()

	return message, nil
}

func (s *chatService) GetRoomMessages(roomID string, limit int) ([]*models.Message, error) {
	// Check if room exists
	_, err := s.repo.GetRoomByID(roomID)
	if err != nil {
		if err == repository.ErrRoomNotFound {
			return nil, ErrRoomNotFound
		}
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	messages, err := s.repo.GetMessagesByRoomID(roomID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return messages, nil
}