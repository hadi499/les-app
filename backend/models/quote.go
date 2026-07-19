package models

import (
	"time"
)

type Quote struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Quote       string    `gorm:"type:text;not null" json:"quote"`
	Arti        string    `gorm:"type:text;not null" json:"arti"`
	Author      string    `gorm:"type:varchar(100);not null" json:"author"`
	IsPublished bool      `gorm:"default:false" json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
