package kuisapp

import "time"

type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedByID uint      `json:"created_by_id"`
	Quizzes     []Quiz    `json:"quizzes,omitempty" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// TableName overrides the table name used by GORM
func (Category) TableName() string {
	return "kuisapp_categories"
}
