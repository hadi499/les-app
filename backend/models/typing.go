package models

import "time"

type LessonProgress struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	UserID       uint `json:"user_id" gorm:"uniqueIndex:idx_user_lesson"`
	User         User `json:"user" gorm:"foreignKey:UserID"`
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
	User   User   `json:"user" gorm:"foreignKey:UserID"`
	Mode   string `json:"mode" gorm:"uniqueIndex:idx_user_mode"` // left, right, both, letters, all
	Score  int    `json:"score"`
}

type GameHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Mode      string    `json:"mode"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

type LessonHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	LessonID  int       `json:"lessonId"`
	WPM       int       `json:"wpm"`
	Accuracy  int       `json:"accuracy"`
	Stars     int       `json:"stars"`
	CreatedAt time.Time `json:"created_at"`
}
