package models

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"` // Don't expose password in JSON
	Role     string `json:"role" gorm:"type:varchar(20);default:'student'"`
}
