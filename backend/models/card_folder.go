package models

import "time"

type CardFolder struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Cards     []Card    `json:"cards" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // If folder is deleted, set CardFolderID to null
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
