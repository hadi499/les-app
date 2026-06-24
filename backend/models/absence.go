package models

import "time"

type Absence struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	User       User      `json:"user"`
	Date       time.Time `json:"date" gorm:"type:date;not null"`
	ReasonType string    `json:"reason_type" gorm:"type:varchar(20);not null"`
	Note       string    `json:"note" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
