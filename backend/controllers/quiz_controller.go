package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetQuizzes - Mendapatkan daftar semua kuis
func GetQuizzes(c *gin.Context) {
	var quizzes []models.Quiz
	if err := database.DB.Find(&quizzes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kuis"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": quizzes})
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
		Title     string            `json:"title" binding:"required"`
		Category  string            `json:"category" binding:"required"`
		TimeLimit int               `json:"timeLimit" binding:"required"`
		Questions []models.Question `json:"questions"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quiz := models.Quiz{
		Title:     input.Title,
		Category:  input.Category,
		TimeLimit: input.TimeLimit,
		Questions: input.Questions,
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
		Title     string            `json:"title" binding:"required"`
		Category  string            `json:"category" binding:"required"`
		TimeLimit int               `json:"timeLimit" binding:"required"`
		Questions []models.Question `json:"questions"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Memulai transaction untuk menghindari data korup
	tx := database.DB.Begin()

	// Update detail kuis
	if err := tx.Model(&quiz).Updates(models.Quiz{
		Title:     input.Title,
		Category:  input.Category,
		TimeLimit: input.TimeLimit,
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

	if err := database.DB.Delete(&quiz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kuis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kuis berhasil dihapus"})
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

	scoreQuiz := models.ScoreQuiz{
		Username:  usernameStr,
		QuizID:    input.QuizID,
		Score:     input.Score,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&scoreQuiz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan skor"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Skor berhasil disimpan", "data": scoreQuiz})
}

// GetQuizScores - Mengambil daftar skor (hanya untuk guru/admin)
func GetQuizScores(c *gin.Context) {
	var scores []models.ScoreQuiz
	if err := database.DB.Preload("Quiz").Order("created_at desc").Find(&scores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar skor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": scores})
}

// GetMyQuizScores - Mengambil riwayat skor kuis pengguna
func GetMyQuizScores(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi"})
		return
	}
	usernameStr := username.(string)

	var scores []models.ScoreQuiz
	if err := database.DB.Preload("Quiz").Where("username = ?", usernameStr).Order("created_at desc").Find(&scores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar skor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": scores})
}
