package controllers

import (
	"net/http"
	"strconv"
	"time"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetUsers returns all users. Accessible by all authenticated users for the leaderboard.
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))
	offset := (page - 1) * limit

	roleFilter := c.Query("role")

	query := database.DB.Model(&models.User{})

	// Role-based filtering logic
	roleInter, _ := c.Get("role")
	roleStr, _ := roleInter.(string)
	userIDInter, _ := c.Get("user_id")
	userID, _ := userIDInter.(uint)

	if roleStr == "teacher" {
		var parentIDs []uint
		database.DB.Model(&models.User{}).Where("teacher_id = ?", userID).Where("parent_id IS NOT NULL").Pluck("parent_id", &parentIDs)
		if len(parentIDs) > 0 {
			query = query.Where("teacher_id = ? OR id = ? OR id IN ?", userID, userID, parentIDs)
		} else {
			query = query.Where("teacher_id = ? OR id = ?", userID, userID)
		}
	} else if roleStr == "student" {
		var currentStudent models.User
		database.DB.First(&currentStudent, userID)
		if currentStudent.TeacherID != nil {
			query = query.Where("teacher_id = ? OR id = ?", *currentStudent.TeacherID, *currentStudent.TeacherID)
		} else {
			query = query.Where("id = ?", userID)
		}
	} else if roleStr == "parent" {
		var children []models.User
		database.DB.Where("parent_id = ?", userID).Find(&children)
		
		var teacherIDs []uint
		for _, child := range children {
			if child.TeacherID != nil {
				teacherIDs = append(teacherIDs, *child.TeacherID)
			}
		}
		
		if len(teacherIDs) > 0 {
			query = query.Where("id IN ? OR id = ? OR parent_id = ?", teacherIDs, userID, userID)
		} else {
			query = query.Where("id = ? OR parent_id = ?", userID, userID)
		}
	}

	if roleFilter != "" {
		query = query.Where("role = ?", roleFilter)
	}

	var users []models.User
	if err := query.Preload("Parent").Preload("Children").Preload("Teacher").Preload("Students").Order("id asc").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var totalUsers, totalTeachers, totalStudents, totalParents int64
	
	// Create base queries for counting without Limit and Offset
	countQuery := database.DB.Model(&models.User{})
	if roleStr == "teacher" {
		var parentIDs []uint
		database.DB.Model(&models.User{}).Where("teacher_id = ?", userID).Where("parent_id IS NOT NULL").Pluck("parent_id", &parentIDs)
		if len(parentIDs) > 0 {
			countQuery = countQuery.Where("teacher_id = ? OR id = ? OR id IN ?", userID, userID, parentIDs)
		} else {
			countQuery = countQuery.Where("teacher_id = ? OR id = ?", userID, userID)
		}
	} else if roleStr == "student" {
		var currentStudent models.User
		database.DB.First(&currentStudent, userID)
		if currentStudent.TeacherID != nil {
			countQuery = countQuery.Where("teacher_id = ? OR id = ?", *currentStudent.TeacherID, *currentStudent.TeacherID)
		} else {
			countQuery = countQuery.Where("id = ?", userID)
		}
	} else if roleStr == "parent" {
		var children []models.User
		database.DB.Where("parent_id = ?", userID).Find(&children)
		var teacherIDs []uint
		for _, child := range children {
			if child.TeacherID != nil {
				teacherIDs = append(teacherIDs, *child.TeacherID)
			}
		}
		if len(teacherIDs) > 0 {
			countQuery = countQuery.Where("id IN ? OR id = ? OR parent_id = ?", teacherIDs, userID, userID)
		} else {
			countQuery = countQuery.Where("id = ? OR parent_id = ?", userID, userID)
		}
	}

	if roleFilter != "" {
		countQuery = countQuery.Where("role = ?", roleFilter)
	}
	countQuery.Count(&totalUsers)

	// Since Count() doesn't mutate the base Where conditions of countQuery, we can clone a new session for subsequent role counts
	teacherCountQuery := countQuery.Session(&gorm.Session{})
	teacherCountQuery.Where("role = ?", "teacher").Count(&totalTeachers)
	
	studentCountQuery := countQuery.Session(&gorm.Session{})
	studentCountQuery.Where("role = ?", "student").Count(&totalStudents)

	parentCountQuery := countQuery.Session(&gorm.Session{})
	parentCountQuery.Where("role = ?", "parent").Count(&totalParents)

	totalPages := 1
	if limit > 0 {
		totalPages = int((totalUsers + int64(limit) - 1) / int64(limit))
	}

	c.JSON(http.StatusOK, gin.H{
		"users":          users,
		"total_users":    totalUsers,
		"total_teachers": totalTeachers,
		"total_students": totalStudents,
		"total_parents":  totalParents,
		"current_page":   page,
		"total_pages":    totalPages,
	})
}

// DeleteUser deletes a user and all their associated data (cascade manually).
func DeleteUser(c *gin.Context) {
	// Dapatkan role dari context
	roleInter, exists := c.Get("role")
	roleStr, _ := roleInter.(string)
	if !exists || (roleStr != "admin" && roleStr != "Admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden. Only admins can perform this action."})
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
	if !exists || (roleStr != "admin" && roleStr != "Admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden. Only admins can perform this action."})
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
		ParentID      *uint  `json:"parent_id"`
		RemoveParent  bool   `json:"remove_parent"`
		TeacherID     *uint  `json:"teacher_id"`
		RemoveTeacher bool   `json:"remove_teacher"`
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
		if req.Role != models.RoleAdmin && req.Role != models.RoleTeacher && req.Role != models.RoleStudent && req.Role != models.RoleParent {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid"})
			return
		}
		user.Role = req.Role
	}
	if req.Class != "" {
		user.Class = req.Class
	}
	
	if req.RemoveParent {
		user.ParentID = nil
	} else if req.ParentID != nil {
		user.ParentID = req.ParentID
	}
	
	if req.RemoveTeacher {
		user.TeacherID = nil
	} else if req.TeacherID != nil {
		user.TeacherID = req.TeacherID
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
			"parent_id":    user.ParentID,
		},
	})
}

// ResetUserPassword resets a user's password. Only accessible by teachers.
func ResetUserPassword(c *gin.Context) {
	// Dapatkan role dari context
	roleInter, exists := c.Get("role")
	roleStr, _ := roleInter.(string)
	if !exists || (roleStr != "admin" && roleStr != "Admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden. Only admins can perform this action."})
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

func ResetAllPoints(c *gin.Context) {
	// Dapatkan role dari context (walaupun sudah ada middleware, sebagai double check)
	roleInter, exists := c.Get("role")
	roleStr, _ := roleInter.(string)
	if !exists || (roleStr != "admin" && roleStr != "Admin" && roleStr != "teacher" && roleStr != "Teacher") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden. Only admins and teachers can perform this action."})
		return
	}

	if roleStr == "admin" || roleStr == "Admin" {
		// Reset semua poin ke 0
		if err := database.DB.Model(&models.User{}).Where("1 = 1").Update("points", 0).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset points"})
			return
		}

		// Reset semua riwayat nilai kuis
		if err := database.DB.Where("1 = 1").Delete(&models.ScoreQuiz{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset quiz scores"})
			return
		}
	} else {
		// Teacher resets only their students
		userIDInter, _ := c.Get("user_id")
		userID, _ := userIDInter.(uint)

		var students []models.User
		if err := database.DB.Where("teacher_id = ?", userID).Find(&students).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
			return
		}

		if len(students) > 0 {
			var studentUsernames []string
			for _, s := range students {
				studentUsernames = append(studentUsernames, s.Username)
			}

			if err := database.DB.Model(&models.User{}).Where("teacher_id = ?", userID).Update("points", 0).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset points for students"})
				return
			}

			if err := database.DB.Where("username IN ?", studentUsernames).Delete(&models.ScoreQuiz{}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset quiz scores for students"})
				return
			}
		}
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

	c.JSON(http.StatusOK, gin.H{"message": "Points and quiz scores have been reset successfully"})
}

// GetUserByID returns details of a single user by ID.
func GetUserByID(c *gin.Context) {
	userIdParam := c.Param("id")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var user models.User
	if err := database.DB.Preload("Parent").Preload("Children").Preload("Teacher").Preload("Students").First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
