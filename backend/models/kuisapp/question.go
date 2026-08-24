package kuisapp

import (
	"time"

	"github.com/lib/pq"
)

type Question struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	QuizID    uint           `json:"quiz_id" gorm:"not null"`
	Question  string         `json:"question"`
	Image     string         `json:"image"`
	Options   pq.StringArray `json:"options" gorm:"type:text[]"`
	Answer    int            `json:"answer"`
}

// TableName overrides the table name used by GORM
func (Question) TableName() string {
	return "kuisapp_questions"
}
