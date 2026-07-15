package models

import "time"

type UserLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Action    string    `json:"action" gorm:"type:varchar(50);not null"` // "login" or "logout"
	IPAddress string    `json:"ip_address" gorm:"type:varchar(45)"`
	UserAgent string    `json:"user_agent" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}
