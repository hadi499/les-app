package models

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
}
