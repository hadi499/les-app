package models

type Subject struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name" gorm:"uniqueIndex:idx_name_user;not null"`
	UserID *uint  `json:"user_id" gorm:"uniqueIndex:idx_name_user"`
	User   *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
