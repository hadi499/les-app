package kuisapp

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Role        string    `json:"role" gorm:"type:varchar(20);default:'user'"` // "user" atau "admin"
	Status      string    `json:"status" gorm:"type:varchar(50);default:'umum'"` // "pelajar" atau "umum"
	Points      int       `json:"points" gorm:"default:0"`
	IsSuspended bool      `json:"is_suspended" gorm:"default:false"`
}

// TableName overrides the table name used by GORM
func (User) TableName() string {
	return "kuisapp_users"
}
