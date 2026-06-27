package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTargets mengambil semua target hafalan, bisa difilter dengan query `student_id`
func GetTargets(c *gin.Context) {
	// API ini hanya untuk teacher, diatur via RoleMiddleware di routes
	
	studentID := c.Query("student_id")
	
	var targets []models.TargetHafalan
	query := database.DB.Model(&models.TargetHafalan{})
	
	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	
	// Preload Student to get student details if needed
	if err := query.Preload("Student").Order("created_at desc").Find(&targets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil target hafalan"})
		return
	}
	
	c.JSON(http.StatusOK, targets)
}

// CreateTarget membuat target hafalan baru
func CreateTarget(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		StudentID uint   `json:"student_id" binding:"required"`
		Text      string `json:"text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Pastikan student ada
	var student models.User
	if err := database.DB.Where("id = ? AND role = ?", input.StudentID, "student").First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student tidak ditemukan"})
		return
	}

	target := models.TargetHafalan{
		StudentID: input.StudentID,
		Text:      input.Text,
		Completed: false,
	}

	if err := database.DB.Create(&target).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat target hafalan"})
		return
	}

	database.DB.Preload("Student").First(&target, target.ID)
	c.JSON(http.StatusCreated, target)
}

// UpdateTarget mengubah status target hafalan (completed/uncompleted)
func UpdateTarget(c *gin.Context) {
	targetID := c.Param("id")

	var target models.TargetHafalan
	if err := database.DB.First(&target, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target hafalan tidak ditemukan"})
		return
	}

	var input struct {
		Completed *bool  `json:"completed"`
		Text      string `json:"text"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Completed != nil {
		target.Completed = *input.Completed
	}
	if input.Text != "" {
		target.Text = input.Text
	}

	if err := database.DB.Save(&target).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate target hafalan"})
		return
	}

	c.JSON(http.StatusOK, target)
}

// DeleteTarget menghapus target hafalan
func DeleteTarget(c *gin.Context) {
	targetID := c.Param("id")

	if err := database.DB.Delete(&models.TargetHafalan{}, targetID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus target hafalan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Target hafalan berhasil dihapus"})
}
