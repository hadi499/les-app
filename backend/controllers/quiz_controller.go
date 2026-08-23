package controllers

import (
	"backend/database"
	"backend/models"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func deleteQuizImage(imageURL string) {
	if imageURL == "" {
		return
	}
	filePath := strings.TrimPrefix(imageURL, "/")
	if filePath == "." || filePath == "" || strings.Contains(filePath, "..") || strings.HasPrefix(filepath.Base(filePath), "?") {
		return
	}
	if err := os.Remove(filePath); err != nil {
		log.Printf("Failed to delete quiz image %s: %v", filePath, err)
	}
}

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

	roleInter, _ := c.Get("role")
	roleStr, _ := roleInter.(string)
	userIDInter, _ := c.Get("user_id")
	userID, _ := userIDInter.(uint)

	if roleStr == "teacher" {
		query = query.Where("user_id = ?", userID)
	} else if roleStr == "student" {
		query = query.Where("is_published = ?", true)
		
		var currentUser models.User
		database.DB.First(&currentUser, userID)
		if currentUser.TeacherID != nil {
			query = query.Where("user_id = ? OR user_id IS NULL", *currentUser.TeacherID)
		} else {
			query = query.Where("user_id IS NULL")
		}
	} else if roleStr == "parent" {
		query = query.Where("is_published = ?", true)
		
		var children []models.User
		database.DB.Where("parent_id = ?", userID).Find(&children)
		
		var teacherIDs []uint
		for _, child := range children {
			if child.TeacherID != nil {
				teacherIDs = append(teacherIDs, *child.TeacherID)
			}
		}
		if len(teacherIDs) > 0 {
			query = query.Where("user_id IN ? OR user_id IS NULL", teacherIDs)
		} else {
			query = query.Where("user_id IS NULL")
		}
	}
	// Admin sees all, no additional filter

	if err := query.Count(&totalItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung data kuis"})
		return
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	if err := query.Order("is_published desc, id desc").Limit(limit).Offset(offset).Find(&quizzes).Error; err != nil {
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

	roleInter, _ := c.Get("role")
	roleStr, _ := roleInter.(string)
	userIDInter, _ := c.Get("user_id")
	userID, _ := userIDInter.(uint)

	if !quiz.IsPublished && roleStr != "teacher" && roleStr != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Kuis ini belum dipublikasikan atau sudah ditutup"})
		return
	}

	if roleStr == "teacher" {
		if quiz.UserID != nil && *quiz.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke kuis ini"})
			return
		}
	} else if roleStr == "student" {
		var currentStudent models.User
		database.DB.First(&currentStudent, userID)
		if quiz.UserID != nil && (currentStudent.TeacherID == nil || *quiz.UserID != *currentStudent.TeacherID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke kuis ini"})
			return
		}
	} else if roleStr == "parent" {
		if quiz.UserID != nil {
			var children []models.User
			database.DB.Where("parent_id = ?", userID).Find(&children)
			hasAccess := false
			for _, child := range children {
				if child.TeacherID != nil && *child.TeacherID == *quiz.UserID {
					hasAccess = true
					break
				}
			}
			if !hasAccess {
				c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke kuis ini"})
				return
			}
		}
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

	isPublished := false
	if input.IsPublished != nil {
		isPublished = *input.IsPublished
	}

	userIDInter, _ := c.Get("user_id")
	userID, _ := userIDInter.(uint)

	quiz := models.Quiz{
		Title:       input.Title,
		Category:    input.Category,
		TimeLimit:   input.TimeLimit,
		IsPublished: isPublished,
		Questions:   input.Questions,
		UserID:      &userID,
	}

	if err := database.DB.Create(&quiz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kuis"})
		return
	}

	// Workaround: GORM mengabaikan nilai 'false' saat Create jika ada tag default:true
	if !isPublished {
		database.DB.Exec("UPDATE quizzes SET is_published = false WHERE id = ?", quiz.ID)
		quiz.IsPublished = false
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Kuis berhasil dibuat", "data": quiz})
}

// DuplicateQuiz - Menduplikat kuis yang sudah ada
func DuplicateQuiz(c *gin.Context) {
	id := c.Param("id")
	var originalQuiz models.Quiz

	if err := database.DB.Preload("Questions").First(&originalQuiz, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kuis tidak ditemukan"})
		return
	}

	newQuestions := make([]models.Question, len(originalQuiz.Questions))
	for i, q := range originalQuiz.Questions {
		newQuestions[i] = models.Question{
			Question: q.Question,
			Image:    q.Image,
			Options:  q.Options,
			Answer:   q.Answer,
		}
	}

	userIDInter, _ := c.Get("user_id")
	userID, _ := userIDInter.(uint)

	duplicatedQuiz := models.Quiz{
		Title:       originalQuiz.Title + " (Salinan)",
		Category:    originalQuiz.Category,
		TimeLimit:   originalQuiz.TimeLimit,
		IsPublished: false,
		Questions:   newQuestions,
		UserID:      &userID,
	}

	if err := database.DB.Create(&duplicatedQuiz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menduplikat kuis"})
		return
	}

	// Workaround: GORM mengabaikan nilai 'false' saat Create jika ada tag default:true
	database.DB.Exec("UPDATE quizzes SET is_published = false WHERE id = ?", duplicatedQuiz.ID)
	duplicatedQuiz.IsPublished = false

	c.JSON(http.StatusCreated, gin.H{"message": "Kuis berhasil diduplikat", "data": duplicatedQuiz})
}

// UpdateQuiz - Mengedit kuis dan soal-soalnya
func UpdateQuiz(c *gin.Context) {
	id := c.Param("id")
	var quiz models.Quiz

	if err := database.DB.Preload("Questions").First(&quiz, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kuis tidak ditemukan"})
		return
	}

	roleInter, _ := c.Get("role")
	roleStr, _ := roleInter.(string)
	userIDInter, _ := c.Get("user_id")
	userID, _ := userIDInter.(uint)
	
	if roleStr != "admin" {
		if quiz.UserID == nil || *quiz.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak berhak mengubah kuis ini"})
			return
		}
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

	// Hapus gambar lama yang tidak dipakai lagi
	var oldQuestions []models.Question
	database.DB.Where("quiz_id = ?", quiz.ID).Find(&oldQuestions)

	newImages := make(map[string]bool)
	for _, q := range input.Questions {
		if q.Image != "" {
			newImages[q.Image] = true
		}
	}

	for _, oldQ := range oldQuestions {
		if oldQ.Image != "" && !newImages[oldQ.Image] {
			deleteQuizImage(oldQ.Image)
		}
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

	roleInter, _ := c.Get("role")
	roleStr, _ := roleInter.(string)
	userIDInter, _ := c.Get("user_id")
	userID, _ := userIDInter.(uint)

	if roleStr != "admin" {
		if quiz.UserID == nil || *quiz.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak berhak menghapus kuis ini"})
			return
		}
	}

	tx := database.DB.Begin()

	// Hapus gambar dari soal yang terkait
	var oldQuestions []models.Question
	database.DB.Where("quiz_id = ?", quiz.ID).Find(&oldQuestions)
	for _, oldQ := range oldQuestions {
		if oldQ.Image != "" {
			deleteQuizImage(oldQ.Image)
		}
	}

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



	tx := database.DB.Begin()

	// Ambil semua skor untuk kuis ini sebelum dihapus
	var scoresToRevoke []models.ScoreQuiz
	if err := tx.Where("quiz_id = ? AND points_earned > 0", quiz.ID).Find(&scoresToRevoke).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data poin yang akan ditarik"})
		return
	}

	// Tarik poin dari setiap user yang mendapatkannya dari kuis ini
	for _, s := range scoresToRevoke {
		if err := tx.Model(&models.User{}).Where("username = ?", s.Username).UpdateColumn("points", gorm.Expr("CASE WHEN points - ? < 0 THEN 0 ELSE points - ? END", s.PointsEarned, s.PointsEarned)).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik poin dari pengguna"})
			return
		}
	}

	// Hapus semua skor untuk kuis ini
	if err := tx.Where("quiz_id = ?", quiz.ID).Delete(&models.ScoreQuiz{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus riwayat nilai kuis"})
		return
	}

	// Update LastResetAt
	if err := tx.Model(&quiz).Update("last_reset_at", time.Now()).Error; err != nil {
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

	var quiz models.Quiz
	if err := database.DB.First(&quiz, input.QuizID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kuis tidak ditemukan"})
		return
	}

	role, _ := c.Get("role")
	if role == "teacher" || role == "admin" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Uji coba kuis berhasil. Nilai guru/admin tidak disimpan ke database.",
			"points_earned": 0,
			"points_already_claimed": false,
		})
		return
	}

	if !quiz.IsPublished {
		c.JSON(http.StatusForbidden, gin.H{"error": "Kuis ini belum dipublikasikan atau sudah ditutup"})
		return
	}

	var previousAttempts int64
	database.DB.Model(&models.ScoreQuiz{}).Where("username = ? AND quiz_id = ?", usernameStr, input.QuizID).Count(&previousAttempts)

	pointsToAdd := 0
	if previousAttempts == 0 {
		if input.Score >= 100 {
			pointsToAdd = 5
		} else if input.Score >= 90 {
			pointsToAdd = 4
		} else if input.Score >= 80 {
			pointsToAdd = 3
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

	pointsAlreadyClaimed := previousAttempts > 0

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

	roleInter, _ := c.Get("role")
	roleStr, _ := roleInter.(string)
	userIDInter, _ := c.Get("user_id")
	userID, _ := userIDInter.(uint)

	if strings.ToLower(strings.TrimSpace(roleStr)) == "teacher" {
		query = query.Joins("JOIN users ON users.username = score_quizzes.username").
			Joins("JOIN quizzes ON quizzes.id = score_quizzes.quiz_id").
			Where("(users.teacher_id = ? OR users.id = ?) AND (quizzes.user_id = ? OR quizzes.user_id IS NULL)", userID, userID, userID)
	}

	searchUsername := c.Query("username")
	if searchUsername != "" {
		query = query.Where("score_quizzes.username LIKE ?", "%"+searchUsername+"%")
	}

	if err := query.Count(&totalItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung skor"})
		return
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	if err := query.Preload("Quiz").Preload("User").Order("score_quizzes.created_at desc").Limit(limit).Offset(offset).Find(&scores).Error; err != nil {
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
