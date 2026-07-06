package models

import "time"

// TodoList represents a collection of TodoItems (Title/Category)
type TodoList struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	UserID          uint       `json:"user_id"`
	Title           string     `json:"title"`
	StudentUsername string     `json:"student_username"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Items     []TodoItem `json:"items" gorm:"foreignKey:TodoListID;constraint:OnDelete:CASCADE;"`
}

// TodoItem represents an individual task within a TodoList
type TodoItem struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	TodoListID uint      `json:"todolist_id"`
	Text       string    `json:"text"`
	Completed  bool      `json:"completed" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
