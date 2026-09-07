package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSettings mengembalikan semua pengaturan sistem sebagai map JSON
func GetSettings(c *gin.Context) {
	var settings []models.SystemSetting
	if err := database.DB.Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil pengaturan"})
		return
	}

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}

	c.JSON(http.StatusOK, result)
}

// UpdateSetting memperbarui satu nilai pengaturan berdasarkan key (hanya untuk guru/admin)
func UpdateSetting(c *gin.Context) {
	key := c.Param("key")

	var input struct {
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var setting models.SystemSetting
	if err := database.DB.Where("key = ?", key).First(&setting).Error; err != nil {
		// Jika belum ada, kita bisa membuat baru (upsert)
		setting = models.SystemSetting{Key: key, Value: input.Value}
		database.DB.Create(&setting)
	} else {
		setting.Value = input.Value
		database.DB.Save(&setting)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengaturan berhasil diperbarui", "data": setting})
}
