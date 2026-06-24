package models

import "time"

type Folder struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Notes     []Note    `json:"notes" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Cascade delete
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
