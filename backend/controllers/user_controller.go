package controllers

import (
	"net/http"
	"strconv"
	"time"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers returns all users. Accessible by all authenticated users for the leaderboard.
func GetUsers(c *gin.Context) {

	var users []models.User
	if err := database.DB.Order("id asc").Find(&users).Error; err != nil {
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

	// Hapus ScoreQuiz (berdasarkan username)
	if err := tx.Where("username = ?", user.Username).Delete(&models.ScoreQuiz{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete score quiz"})
		return
	}

	// Hapus data lainnya untuk mencegah foreign key constraint error
	tx.Where("user_id = ?", userId).Delete(&models.TodoList{})
	tx.Where("user_id = ?", userId).Delete(&models.WritingProgress{})
	tx.Where("user_id = ?", userId).Delete(&models.Exam{})
	tx.Where("user_id = ?", userId).Delete(&models.Note{})
	tx.Where("user_id = ?", userId).Delete(&models.Folder{})
	tx.Where("user_id = ?", userId).Delete(&models.Absence{})
	tx.Where("user_id = ?", userId).Delete(&models.UserLog{})
	tx.Where("sender_id = ? OR receiver_id = ?", userId, userId).Delete(&models.ChatMessage{})

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

// UpdateUser updates a user's details (username, role, class, password).
func UpdateUser(c *gin.Context) {
	roleInter, exists := c.Get("role")
	roleStr, _ := roleInter.(string)
	if !exists || (roleStr != "teacher" && roleStr != "admin" && roleStr != "Teacher" && roleStr != "Admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden. Only teachers/admins can perform this action."})
		return
	}

	userIdParam := c.Param("id")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req struct {
		Username    string `json:"username"`
		Role        string `json:"role"`
		Class       string `json:"class"`
		Password    string `json:"password"`
		IsSuspended *bool  `json:"is_suspended"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Class != "" {
		user.Class = req.Class
	}
	if req.Password != "" {
		if len(req.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password minimal 6 karakter"})
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	if req.IsSuspended != nil {
		user.IsSuspended = *req.IsSuspended
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user": gin.H{
			"id":           user.ID,
			"username":     user.Username,
			"role":         user.Role,
			"class":        user.Class,
			"is_suspended": user.IsSuspended,
		},
	})
}

// ResetUserPassword resets a user's password. Only accessible by teachers.
func ResetUserPassword(c *gin.Context) {
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

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input. Password must be at least 6 characters."})
		return
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update password user
	if err := database.DB.Model(&models.User{}).Where("id = ?", userId).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// ResetAllPoints resets all user points to 0 and clears quiz scores. Only accessible by teachers.
func ResetAllPoints(c *gin.Context) {
// Dapatkan role dari context (walaupun sudah ada middleware, sebagai double check)
roleInter, exists := c.Get("role")
roleStr, _ := roleInter.(string)
if !exists || (roleStr != "teacher" && roleStr != "Teacher" && roleStr != "admin" && roleStr != "Admin") {
c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden. Only teachers can perform this action."})
return
}

// Reset semua poin ke 0
if err := database.DB.Model(&models.User{}).Where("1 = 1").Update("points", 0).Error; err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset points"})
return
}

// Reset semua riwayat nilai kuis (menghapus semua skor agar bisa dikerjakan ulang dari nol)
if err := database.DB.Where("1 = 1").Delete(&models.ScoreQuiz{}).Error; err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset quiz scores"})
return
}

// Update bulan terakhir reset
var setting models.SystemSetting
loc, errLoc := time.LoadLocation("Asia/Jakarta")
if errLoc != nil {
loc = time.FixedZone("WIB", 7*3600)
}
currentMonth := time.Now().In(loc).Format("2006-01")

if err := database.DB.Where("key = ?", "last_point_reset").First(&setting).Error; err != nil {
database.DB.Create(&models.SystemSetting{Key: "last_point_reset", Value: currentMonth})
} else {
database.DB.Model(&setting).Update("value", currentMonth)
}

c.JSON(http.StatusOK, gin.H{"message": "All points and quiz scores have been reset successfully"})
}
