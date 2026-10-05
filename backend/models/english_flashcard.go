package models

import "github.com/lib/pq"

type EnglishFlashcardCategory struct {
	ID          uint               `json:"id" gorm:"primaryKey"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Flashcards  []EnglishFlashcard `json:"flashcards,omitempty" gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE;"`
}

type EnglishFlashcard struct {
	ID          uint                     `json:"id" gorm:"primaryKey"`
	CategoryID  uint                     `json:"category_id"`
	Category    EnglishFlashcardCategory `json:"category" gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE;"`
	Question    string                   `json:"question"`
	Options     pq.StringArray           `json:"options" gorm:"type:text[]"`
	Answer      string                   `json:"answer"`
	Explanation string                   `json:"explanation"`
}
