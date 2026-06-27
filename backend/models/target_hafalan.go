package models

import "time"

// TargetHafalan represents a memorization target assigned to a student
type TargetHafalan struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	StudentID uint      `json:"student_id" gorm:"not null"`
	Text      string    `json:"text" gorm:"not null"`
	Completed bool      `json:"completed" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Student User `json:"student" gorm:"foreignKey:StudentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
