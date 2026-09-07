package models

import "time"

// SystemSetting digunakan untuk menyimpan konfigurasi global (seperti status kelas buka/libur)
type SystemSetting struct {
	Key       string    `gorm:"primaryKey;type:varchar(100)" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}
