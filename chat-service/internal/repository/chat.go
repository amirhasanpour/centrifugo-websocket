package repository

import (
	"errors"

	"chat-service/internal/models"
	"gorm.io/gorm"
)

var (
	ErrRoomNotFound    = errors.New("room not found")
	ErrMessageNotFound = errors.New("message not found")
	ErrAlreadyMember   = errors.New("user already a member of this room")
)

type ChatRepository interface {
	CreateRoom(room *models.Room) error
	GetRoomByID(roomID string) (*models.Room, error)
	GetAllRooms() ([]*models.Room, error)
	CreateMessage(message *models.Message) error
	GetMessagesByRoomID(roomID string, limit int) ([]*models.Message, error)
	AddRoomMember(roomID, userID string) error
	IsRoomMember(roomID, userID string) (bool, error)
	GetRoomMembers(roomID string) ([]*models.RoomMember, error)
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) CreateRoom(room *models.Room) error {
	result := r.db.Create(room)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *chatRepository) GetRoomByID(roomID string) (*models.Room, error) {
	var room models.Room
	result := r.db.First(&room, "id = ?", roomID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}
		return nil, result.Error
	}
	return &room, nil
}

func (r *chatRepository) GetAllRooms() ([]*models.Room, error) {
	var rooms []*models.Room
	result := r.db.Find(&rooms)
	if result.Error != nil {
		return nil, result.Error
	}
	return rooms, nil
}

func (r *chatRepository) CreateMessage(message *models.Message) error {
	result := r.db.Create(message)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *chatRepository) GetMessagesByRoomID(roomID string, limit int) ([]*models.Message, error) {
	var messages []*models.Message
	result := r.db.Where("room_id = ?", roomID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *chatRepository) AddRoomMember(roomID, userID string) error {
	// Check if already a member
	var existingMember models.RoomMember
	result := r.db.Where("room_id = ? AND user_id = ?", roomID, userID).
		First(&existingMember)

	if result.Error == nil {
		return ErrAlreadyMember
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	member := &models.RoomMember{
		RoomID: roomID,
		UserID: userID,
	}

	result = r.db.Create(member)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *chatRepository) IsRoomMember(roomID, userID string) (bool, error) {
	var member models.RoomMember
	result := r.db.Where("room_id = ? AND user_id = ?", roomID, userID).
		First(&member)
		
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func (r *chatRepository) GetRoomMembers(roomID string) ([]*models.RoomMember, error) {
	var members []*models.RoomMember
	result := r.db.Where("room_id = ?", roomID).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}