package controllers

import (
	"net/http"
	"strconv"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
)

// GetUsers returns all users. Only accessible by teachers.
func GetUsers(c *gin.Context) {
	// Dapatkan role dari context yang di-set oleh AuthMiddleware
	roleInter, exists := c.Get("role")
	roleStr, _ := roleInter.(string)
	if !exists || (roleStr != "teacher" && roleStr != "Teacher" && roleStr != "admin" && roleStr != "Admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden. Only teachers can access this."})
		return
	}

	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// DeleteUser deletes a user and all their associated data (cascade manually).
func DeleteUser(c *gin.Context) {
	// Dapatkan role dari context
	roleInter, exists := c.Get("role")
	roleStr, _ := roleInter.(string)
	if !exists || (roleStr != "teacher" && roleStr != "Teacher" && roleStr != "admin" && roleStr != "Admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden. Only teachers can perform this action."})
		return
	}

	userIdParam := c.Param("id")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Jangan izinkan user menghapus dirinya sendiri
	currentUserId, exists := c.Get("user_id")
	if exists && uint(userId) == currentUserId.(uint) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot delete yourself."})
		return
	}

	// Cek apakah user ada
	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Mulai transaksi untuk memastikan konsistensi data
	tx := database.DB.Begin()

	// 1. Hapus LessonProgress
	if err := tx.Where("user_id = ?", userId).Delete(&models.LessonProgress{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete lesson progress"})
		return
	}

	// 2. Hapus GameHighScore
	if err := tx.Where("user_id = ?", userId).Delete(&models.GameHighScore{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete game high scores"})
		return
	}

	// 3. Hapus GameHistory
	if err := tx.Where("user_id = ?", userId).Delete(&models.GameHistory{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete game history"})
		return
	}

	// 4. Hapus LessonHistory
	if err := tx.Where("user_id = ?", userId).Delete(&models.LessonHistory{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete lesson history"})
		return
	}

	// 5. Hapus User
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Commit transaksi jika semua berhasil
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "User and associated data deleted successfully"})
}
