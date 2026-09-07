package models

import (
	"time"
)

type ChatMessage struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SenderID   uint      `json:"sender_id" gorm:"not null"`
	ReceiverID uint      `json:"receiver_id" gorm:"not null"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	IsRead     bool      `json:"is_read" gorm:"default:false"`
	IsDeleted  bool      `json:"is_deleted" gorm:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Sender   User `json:"sender" gorm:"foreignKey:SenderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Receiver User `json:"receiver" gorm:"foreignKey:ReceiverID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
