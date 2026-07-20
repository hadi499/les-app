package models

import "time"

type User struct {
	ID               uint             `json:"id" gorm:"primaryKey"`
	Username         string           `json:"username" gorm:"unique;not null"`
	Password         string           `json:"-" gorm:"not null"` // Don't expose password in JSON
	Role             string           `json:"role" gorm:"type:varchar(20);default:'student'"`
	LessonProgresses []LessonProgress `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GameHighScores   []GameHighScore  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GameHistories    []GameHistory    `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LessonHistories  []LessonHistory  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Exams            []Exam           `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Notes            []Note           `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Folders          []Folder         `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Absences         []Absence        `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserLogs         []UserLog        `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LastActiveAt     *time.Time       `json:"last_active_at" gorm:"type:timestamp"`
	Points           int              `json:"points" gorm:"default:0"`
}
