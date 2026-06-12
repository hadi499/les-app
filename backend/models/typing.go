package models

import "time"

type LessonProgress struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	UserID       uint `json:"user_id" gorm:"uniqueIndex:idx_user_lesson"`
	LessonID     int  `json:"lessonId" gorm:"uniqueIndex:idx_user_lesson"`
	BestWPM      int  `json:"bestWPM"`
	BestAccuracy int  `json:"bestAccuracy"`
	Stars        int  `json:"stars"`
	Completed    bool `json:"completed"`
	Attempts     int  `json:"attempts"`
}

type GameHighScore struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	UserID uint   `json:"user_id" gorm:"uniqueIndex:idx_user_mode"`
	Mode   string `json:"mode" gorm:"uniqueIndex:idx_user_mode"` // left, right, both, letters, all
	Score  int    `json:"score"`
}

type GameHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index"`
	Mode      string    `json:"mode"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

type LessonHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index"`
	LessonID  int       `json:"lessonId"`
	WPM       int       `json:"wpm"`
	Accuracy  int       `json:"accuracy"`
	Stars     int       `json:"stars"`
	CreatedAt time.Time `json:"created_at"`
}
