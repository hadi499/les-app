package models

import "time"

type Note struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content"`
	FolderID  *uint     `json:"folder_id"` // nullable
	Folder    Folder    `json:"folder" gorm:"foreignKey:FolderID"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
