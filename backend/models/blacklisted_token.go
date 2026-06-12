package models

import "time"

type BlacklistedToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TokenHash string    `json:"-" gorm:"unique;not null"`
	ExpiresAt time.Time `json:"expires_at"`
}
