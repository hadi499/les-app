package models

import "time"

type ScoreQuiz struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username"`
	QuizID    uint      `json:"quiz_id"`
	Quiz      *Quiz     `json:"quiz" gorm:"foreignKey:QuizID"`
	Score        int       `json:"score"`
	PointsEarned int       `json:"points_earned"`
	CreatedAt    time.Time `json:"created_at"`
}
