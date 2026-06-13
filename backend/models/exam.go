package models

import "time"

type Exam struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	ExamDate  time.Time `json:"exam_date" gorm:"type:date;not null"`
	ExamName  string    `json:"exam_name" gorm:"not null"`
	SubjectID uint      `json:"subject_id" gorm:"index;not null"`
	Subject   Subject   `json:"subject" gorm:"foreignKey:SubjectID"`
	Score     int       `json:"score" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
