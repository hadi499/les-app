package controllers

import (
	"backend/database"
	"backend/models"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetQuizzes - Mendapatkan daftar semua kuis
func GetQuizzes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "24"))
	if limit < 1 {
		limit = 24
	}
	offset := (page - 1) * limit

	var quizzes []models.Quiz
	var totalItems int64

	query := database.DB.Model(&models.Quiz{})

	role, _ := c.Get("role")
	if role != "teacher" && role != "admin" {
		query = query.Where("is_published = ?", true)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung data kuis"})
		return
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	if err := query.Order("id desc").Limit(limit).Offset(offset).Find(&quizzes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kuis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":         quizzes,
		"total_pages":  totalPages,
		"current_page": page,
		"total_items":  totalItems,
	})
}

// GetQuizByID - Mendapatkan detail kuis beserta soal-soalnya
func GetQuizByID(c *gin.Context) {
	id := c.Param("id")
	var quiz models.Quiz

	if err := database.DB.Preload("Questions").First(&quiz, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kuis tidak ditemukan"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan server"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": quiz})
}

// CreateQuiz - Membuat kuis baru beserta soal
func CreateQuiz(c *gin.Context) {
	var input struct {
		Title       string            `json:"title" binding:"required"`
		Category    string            `json:"category" binding:"required"`
		TimeLimit   int               `json:"timeLimit" binding:"required"`
		IsPublished *bool             `json:"is_published"` // Use pointer to distinguish true/false vs null
		Questions   []models.Question `json:"questions"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isPublished := true
	if input.IsPublished != nil {
		isPublished = *input.IsPublished
	}

	quiz := models.Quiz{
		Title:       input.Title,
		Category:    input.Category,
		TimeLimit:   input.TimeLimit,
		IsPublished: isPublished,
		Questions:   input.Questions,
	}

	if err := database.DB.Create(&quiz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kuis"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Kuis berhasil dibuat", "data": quiz})
}

// UpdateQuiz - Mengedit kuis dan soal-soalnya
func UpdateQuiz(c *gin.Context) {
	id := c.Param("id")
	var quiz models.Quiz

	if err := database.DB.Preload("Questions").First(&quiz, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kuis tidak ditemukan"})
		return
	}

	var input struct {
		Title       string            `json:"title" binding:"required"`
		Category    string            `json:"category" binding:"required"`
		TimeLimit   int               `json:"timeLimit" binding:"required"`
		IsPublished *bool             `json:"is_published"`
		Questions   []models.Question `json:"questions"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Memulai transaction untuk menghindari data korup
	tx := database.DB.Begin()

	// Menentukan nilai IsPublished
	isPublished := true
	if input.IsPublished != nil {
		isPublished = *input.IsPublished
	}

	// Update detail kuis
	if err := tx.Model(&quiz).Updates(map[string]interface{}{
		"title":        input.Title,
		"category":     input.Category,
		"time_limit":   input.TimeLimit,
		"is_published": isPublished,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui kuis"})
		return
	}

	// Hapus soal lama dan buat yang baru
	if err := tx.Where("quiz_id = ?", quiz.ID).Delete(&models.Question{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus soal lama"})
		return
	}

	for i := range input.Questions {
		input.Questions[i].QuizID = quiz.ID
	}

	if len(input.Questions) > 0 {
		if err := tx.Create(&input.Questions).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan soal baru"})
			return
		}
	}

	tx.Commit()

	// Muat ulang kuis
	database.DB.Preload("Questions").First(&quiz, id)
	c.JSON(http.StatusOK, gin.H{"message": "Kuis berhasil diperbarui", "data": quiz})
}

// DeleteQuiz - Menghapus kuis
func DeleteQuiz(c *gin.Context) {
	id := c.Param("id")
	var quiz models.Quiz

	if err := database.DB.First(&quiz, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kuis tidak ditemukan"})
		return
	}

	tx := database.DB.Begin()

	// Hapus soal yang terkait dengan kuis ini
	if err := tx.Where("quiz_id = ?", quiz.ID).Delete(&models.Question{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus soal kuis"})
		return
	}

	// Hapus skor yang terkait dengan kuis ini
	if err := tx.Where("quiz_id = ?", quiz.ID).Delete(&models.ScoreQuiz{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus skor kuis"})
		return
	}

	if err := tx.Delete(&quiz).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kuis"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Kuis berhasil dihapus"})
}

// ResetQuizScores - Menghapus semua riwayat nilai pada kuis tertentu
func ResetQuizScores(c *gin.Context) {
	id := c.Param("id")
	var quiz models.Quiz

	// Verifikasi kuis ada
	if err := database.DB.First(&quiz, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kuis tidak ditemukan"})
		return
	}

	// Cek apakah sudah direset bulan ini
	now := time.Now()
	if quiz.LastResetAt != nil && quiz.LastResetAt.Year() == now.Year() && quiz.LastResetAt.Month() == now.Month() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kuis sudah direset pada bulan ini"})
		return
	}

	tx := database.DB.Begin()

	// Hapus semua skor untuk kuis ini
	if err := tx.Where("quiz_id = ?", quiz.ID).Delete(&models.ScoreQuiz{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus riwayat nilai kuis"})
		return
	}

	// Update LastResetAt
	if err := tx.Model(&quiz).Update("last_reset_at", now).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate waktu reset"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Riwayat nilai kuis berhasil direset"})
}

// SubmitQuizScore - Menyimpan hasil kuis pengguna
func SubmitQuizScore(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi"})
		return
	}
	usernameStr := username.(string)

	var input struct {
		QuizID uint `json:"quiz_id" binding:"required"`
		Score  int  `json:"score"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var previousPerfectScore int64
	database.DB.Model(&models.ScoreQuiz{}).Where("username = ? AND quiz_id = ? AND score = 100", usernameStr, input.QuizID).Count(&previousPerfectScore)

	var previousAttempts int64
	database.DB.Model(&models.ScoreQuiz{}).Where("username = ? AND quiz_id = ?", usernameStr, input.QuizID).Count(&previousAttempts)

	pointsToAdd := 0
	if input.Score == 100 && previousPerfectScore == 0 {
		if previousAttempts == 0 {
			pointsToAdd = 5
		} else if previousAttempts == 1 {
			pointsToAdd = 3
		} else if previousAttempts == 2 {
			pointsToAdd = 1
		}
	}

	scoreQuiz := models.ScoreQuiz{
		Username:     usernameStr,
		QuizID:       input.QuizID,
		Score:        input.Score,
		PointsEarned: pointsToAdd,
		CreatedAt:    time.Now(),
	}

	if err := database.DB.Create(&scoreQuiz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan skor"})
		return
	}

	if pointsToAdd > 0 {
		database.DB.Model(&models.User{}).Where("username = ?", usernameStr).UpdateColumn("points", gorm.Expr("points + ?", pointsToAdd))
	}

	pointsAlreadyClaimed := previousPerfectScore > 0

	c.JSON(http.StatusCreated, gin.H{
		"message": "Skor berhasil disimpan", 
		"data": scoreQuiz,
		"points_earned": pointsToAdd,
		"points_already_claimed": pointsAlreadyClaimed,
	})
}

// GetQuizScores - Mengambil daftar skor (hanya untuk guru/admin)
func GetQuizScores(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "24"))
	if limit < 1 {
		limit = 24
	}
	offset := (page - 1) * limit

	var scores []models.ScoreQuiz
	var totalItems int64

	query := database.DB.Model(&models.ScoreQuiz{})

	searchUsername := c.Query("username")
	if searchUsername != "" {
		query = query.Where("username LIKE ?", "%"+searchUsername+"%")
	}

	if err := query.Count(&totalItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung skor"})
		return
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	if err := query.Preload("Quiz").Order("created_at desc").Limit(limit).Offset(offset).Find(&scores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar skor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":         scores,
		"total_pages":  totalPages,
		"current_page": page,
		"total_items":  totalItems,
	})
}

// GetMyQuizScores - Mengambil riwayat skor kuis pengguna
func GetMyQuizScores(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi"})
		return
	}
	usernameStr := username.(string)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "24"))
	if limit < 1 {
		limit = 24
	}
	offset := (page - 1) * limit

	var scores []models.ScoreQuiz
	var totalItems int64

	query := database.DB.Model(&models.ScoreQuiz{}).Where("username = ?", usernameStr)

	if err := query.Count(&totalItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung skor"})
		return
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	if err := query.Preload("Quiz").Order("created_at desc").Limit(limit).Offset(offset).Find(&scores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar skor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":         scores,
		"total_pages":  totalPages,
		"current_page": page,
		"total_items":  totalItems,
	})
}
