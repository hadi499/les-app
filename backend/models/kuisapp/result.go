package kuisapp

import "time"

type Result struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time `json:"created_at"`
	QuizID       uint      `json:"quiz_id" gorm:"not null"`
	UserID       uint      `json:"user_id" gorm:"not null"` // Refers to kuisapp.User
	Score        float64   `json:"score"`
	PointsEarned int       `json:"points_earned" gorm:"default:0"`
	FinishedAt   time.Time `json:"finished_at"`

	// Relationships
	Quiz Quiz `json:"quiz" gorm:"foreignKey:QuizID"`
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// TableName overrides the table name used by GORM
func (Result) TableName() string {
	return "kuisapp_results"
}
