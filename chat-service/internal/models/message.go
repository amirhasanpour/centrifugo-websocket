package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	RoomID    string    `gorm:"not null;index" json:"room_id"`
	UserID    string    `gorm:"not null" json:"user_id"`
	Username  string    `gorm:"not null" json:"username"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	m.CreatedAt = time.Now()
	return nil
}