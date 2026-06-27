package models

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	CardFolderID *uint       `json:"card_folder_id"` // nullable
	CardFolder   CardFolder  `json:"card_folder" gorm:"foreignKey:CardFolderID"`
	Image        string      `json:"image"`
	Title        string      `json:"title"`
	Content      string      `json:"content"`
	Size         string      `json:"size" gorm:"default:6"`
	CardType     string      `json:"cardType" gorm:"default:regular"`
}
