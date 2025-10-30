package models

import (
	"time"

	"gorm.io/gorm"
)

type RoomMember struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	RoomID    string    `gorm:"not null;index" json:"room_id"`
	UserID    string    `gorm:"not null" json:"user_id"`
	JoinedAt  time.Time `json:"joined_at"`
}

func (rm *RoomMember) BeforeCreate(tx *gorm.DB) error {
	rm.JoinedAt = time.Now()
	return nil
}