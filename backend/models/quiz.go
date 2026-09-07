package models

import "time"

// Quiz merepresentasikan sebuah kuis dengan judul dan kategori tertentu
type Quiz struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title"`
	Category    string     `json:"category"`
	TimeLimit   int        `json:"timeLimit"`
	IsPublished bool       `json:"is_published" gorm:"default:true"`
	LastResetAt *time.Time `json:"last_reset_at"`
	Questions   []Question `json:"questions" gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE;"` // Relasi one-to-many
}
