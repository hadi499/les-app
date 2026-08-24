package kuisapp

import "time"

type Quiz struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time  `json:"created_at"`
	Title       string     `json:"title" gorm:"not null"`
	CategoryID  *uint      `json:"category_id"` // Reference to KuisApp Category/Folder
	TimeLimit   int        `json:"timeLimit"`
	IsPublished bool       `json:"is_published" gorm:"default:true"`
	LastResetAt *time.Time `json:"last_reset_at"`
	UserID      *uint      `json:"user_id,omitempty"` // Formerly CreatedByID
	User        *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Questions   []Question `json:"questions,omitempty" gorm:"foreignKey:QuizID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName overrides the table name used by GORM
func (Quiz) TableName() string {
	return "kuisapp_quizzes"
}
