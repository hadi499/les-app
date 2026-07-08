package models

import "time"

type WritingProgress struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Date      time.Time `json:"date" gorm:"type:date;not null"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
