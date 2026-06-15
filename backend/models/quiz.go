package models

// Quiz merepresentasikan sebuah kuis dengan judul dan kategori tertentu
type Quiz struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Title     string     `json:"title"`
	Category  string     `json:"category"`
	TimeLimit int        `json:"timeLimit"`
	Questions []Question `json:"questions" gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE;"` // Relasi one-to-many
}
