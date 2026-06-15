package models

import "github.com/lib/pq"

type Question struct {
	ID       uint           `json:"id" gorm:"primaryKey"`
	QuizID   uint           `json:"quiz_id"` // Relasi ke Quiz
	Question string         `json:"question"`
	Options  pq.StringArray `json:"options" gorm:"type:text[]"`
	Answer   int            `json:"answer"`
}
