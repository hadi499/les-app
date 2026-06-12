package models

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	Category string `json:"category"`
	Image    string `json:"image"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Size     string `json:"size" gorm:"default:6"`
	CardType string `json:"cardType" gorm:"default:regular"`
}
