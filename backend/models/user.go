package models

import "time"

const (
	RoleAdmin   = "admin"
	RoleTeacher = "teacher"
	RoleStudent = "student"
	RoleParent  = "parent"
)

type User struct {
	ID               uint             `json:"id" gorm:"primaryKey"`
	CreatedAt        time.Time        `json:"created_at"`
	Username         string           `json:"username" gorm:"unique;not null"`
	Password         string           `json:"-" gorm:"not null"` // Don't expose password in JSON
	Role             string           `json:"role" gorm:"type:varchar(20);default:'student'"`
	Class            string           `json:"class" gorm:"type:varchar(50);default:''"`
	ParentID         *uint            `json:"parent_id,omitempty"`
	Parent           *User            `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children         []User           `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	TeacherID        *uint            `json:"teacher_id,omitempty"`
	Teacher          *User            `json:"teacher,omitempty" gorm:"foreignKey:TeacherID"`
	Students         []User           `json:"students,omitempty" gorm:"foreignKey:TeacherID"`
	LessonProgresses []LessonProgress `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GameHighScores   []GameHighScore  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GameHistories    []GameHistory    `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LessonHistories  []LessonHistory  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Exams            []Exam           `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Notes            []Note           `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Folders          []Folder         `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Absences         []Absence        `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserLogs         []UserLog        `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LastActiveAt     *time.Time       `json:"last_active_at" gorm:"type:timestamp"`
	Points           int              `json:"points" gorm:"default:0"`
	IsSuspended      bool             `json:"is_suspended" gorm:"default:false"`
}
