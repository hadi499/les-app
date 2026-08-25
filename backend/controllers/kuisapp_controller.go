package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"backend/database"
	"backend/models/kuisapp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ======== KATEGORI ========

func GetKuisAppCategories(c *gin.Context) {
	var categories []kuisapp.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func CreateKuisAppCategory(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userObj, _ := c.Get("kuisapp_user")
	user := userObj.(kuisapp.User)

	category := kuisapp.Category{
		Name:        input.Name,
		Description: input.Description,
		CreatedByID: user.ID,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kategori"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Kategori berhasil dibuat", "data": category})
}

func UpdateKuisAppCategory(c *gin.Context) {
	id := c.Param("id")
	var category kuisapp.Category

	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		category.Name = input.Name
	}
	if input.Description != "" {
		category.Description = input.Description
	}

	if err := database.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kategori diupdate", "data": category})
}

func DeleteKuisAppCategory(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&kuisapp.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kategori dihapus"})
}

// ======== KUIS ========

func GetKuisAppQuizzes(c *gin.Context) {
	var quizzes []kuisapp.Quiz
	if err := database.DB.Preload("Questions").Find(&quizzes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil kuis"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": quizzes})
}

func GetKuisAppQuiz(c *gin.Context) {
	id := c.Param("id")
	var quiz kuisapp.Quiz
	// Sertakan soal (tapi nanti jawaban benar akan di-filter di frontend saat ngerjain)
	if err := database.DB.Preload("Questions").First(&quiz, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kuis tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": quiz})
}

func CreateKuisAppQuiz(c *gin.Context) {
	var input struct {
		Title       string             `json:"title" binding:"required"`
		TimeLimit   int                `json:"timeLimit" binding:"required"`
		CategoryID  *uint              `json:"category_id"`
		IsPublished *bool              `json:"is_published"`
		Questions   []kuisapp.Question `json:"questions"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isPublished := false
	if input.IsPublished != nil {
		isPublished = *input.IsPublished
	}

	userObj, _ := c.Get("kuisapp_user")
	user := userObj.(kuisapp.User)

	quiz := kuisapp.Quiz{
		Title:       input.Title,
		TimeLimit:   input.TimeLimit,
		CategoryID:  input.CategoryID,
		IsPublished: isPublished,
		Questions:   input.Questions,
		UserID:      &user.ID,
	}

	if err := database.DB.Create(&quiz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kuis"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Kuis berhasil dibuat", "data": quiz})
}

func UpdateKuisAppQuiz(c *gin.Context) {
	id := c.Param("id")
	var quiz kuisapp.Quiz
	if err := database.DB.Preload("Questions").First(&quiz, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kuis tidak ditemukan"})
		return
	}

	var input struct {
		Title       string             `json:"title" binding:"required"`
		TimeLimit   int                `json:"timeLimit" binding:"required"`
		CategoryID  *uint              `json:"category_id"`
		IsPublished *bool              `json:"is_published"`
		Questions   []kuisapp.Question `json:"questions"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := database.DB.Begin()

	isPublished := true
	if input.IsPublished != nil {
		isPublished = *input.IsPublished
	}

	if err := tx.Model(&quiz).Updates(map[string]interface{}{
		"title":        input.Title,
		"time_limit":   input.TimeLimit,
		"category_id":  input.CategoryID,
		"is_published": isPublished,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui kuis"})
		return
	}

	// Hapus soal lama dan buat yang baru
	if err := tx.Where("quiz_id = ?", quiz.ID).Delete(&kuisapp.Question{}).Error; err != nil {
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

	database.DB.Preload("Questions").First(&quiz, id)
	c.JSON(http.StatusOK, gin.H{"message": "Kuis berhasil diperbarui", "data": quiz})
}

func DeleteKuisAppQuiz(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&kuisapp.Quiz{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kuis"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kuis dihapus"})
}

func DuplicateKuisAppQuiz(c *gin.Context) {
	id := c.Param("id")
	var quiz kuisapp.Quiz
	if err := database.DB.Preload("Questions").First(&quiz, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kuis tidak ditemukan"})
		return
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi"})
		return
	}

	newQuiz := kuisapp.Quiz{
		Title:       quiz.Title + " (Salinan)",
		TimeLimit:   quiz.TimeLimit,
		CategoryID:  quiz.CategoryID,
		IsPublished: false, // Default to false
		UserID:      quiz.UserID,
	}

	if err := tx.Create(&newQuiz).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menduplikasi kuis"})
		return
	}

	for _, q := range quiz.Questions {
		newQ := kuisapp.Question{
			QuizID:   newQuiz.ID,
			Question: q.Question,
			Image:    q.Image,
			Options:  q.Options,
			Answer:   q.Answer,
		}
		if err := tx.Create(&newQ).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menduplikasi soal"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"message": "Kuis berhasil diduplikasi", "data": newQuiz})
}

// ======== QUESTION ========

func CreateKuisAppQuestion(c *gin.Context) {
	id := c.Param("id") // QuizID
	quizID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Quiz ID"})
		return
	}

	var input struct {
		Question string   `json:"question" binding:"required"`
		Image    string   `json:"image"`
		Options  []string `json:"options" binding:"required"` 
		Answer   int      `json:"answer"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	q := kuisapp.Question{
		QuizID:   uint(quizID),
		Question: input.Question,
		Image:    input.Image,
		Options:  input.Options,
		Answer:   input.Answer,
	}

	if err := database.DB.Create(&q).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pertanyaan"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Pertanyaan ditambahkan", "data": q})
}

func DeleteKuisAppQuestion(c *gin.Context) {
	questionID := c.Param("question_id")
	if err := database.DB.Delete(&kuisapp.Question{}, questionID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus soal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Soal dihapus"})
}

func UpdateKuisAppQuestion(c *gin.Context) {
	questionID := c.Param("question_id")
	var q kuisapp.Question
	if err := database.DB.First(&q, questionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Soal tidak ditemukan"})
		return
	}

	var input struct {
		Question string   `json:"question"`
		Image    string   `json:"image"`
		Options  []string `json:"options"`
		Answer   *int     `json:"answer"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Question != "" {
		q.Question = input.Question
	}
	q.Image = input.Image // can be cleared
	if len(input.Options) > 0 {
		q.Options = input.Options
	}
	if input.Answer != nil {
		q.Answer = *input.Answer
	}

	if err := database.DB.Save(&q).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update soal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Soal diupdate", "data": q})
}

func BulkSaveKuisAppQuestions(c *gin.Context) {
	quizIDStr := c.Param("id")
	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kuis tidak valid"})
		return
	}

	var input struct {
		Questions []struct {
			ID       uint     `json:"id"`
			Question string   `json:"question"`
			Image    string   `json:"image"`
			Options  []string `json:"options"`
			Answer   int      `json:"answer"`
		} `json:"questions"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi"})
		return
	}

	incomingIDs := make(map[uint]bool)
	for _, q := range input.Questions {
		if q.ID != 0 {
			incomingIDs[q.ID] = true
		}
	}

	var existingQuestions []kuisapp.Question
	if err := tx.Where("quiz_id = ?", quizID).Find(&existingQuestions).Error; err == nil {
		for _, eq := range existingQuestions {
			if !incomingIDs[eq.ID] {
				tx.Delete(&eq)
			}
		}
	}

	for _, q := range input.Questions {
		if q.ID == 0 {
			newQ := kuisapp.Question{
				QuizID:   uint(quizID),
				Question: q.Question,
				Image:    q.Image,
				Options:  q.Options,
				Answer:   q.Answer,
			}
			if err := tx.Create(&newQ).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan soal baru"})
				return
			}
		} else {
			var existingQ kuisapp.Question
			if err := tx.First(&existingQ, q.ID).Error; err == nil && existingQ.QuizID == uint(quizID) {
				existingQ.Question = q.Question
				existingQ.Image = q.Image
				existingQ.Options = q.Options
				existingQ.Answer = q.Answer
				if err := tx.Save(&existingQ).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate soal"})
					return
				}
			}
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Semua soal berhasil disimpan"})
}

// ======== UPLOADS ========

func UploadKuisAppImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak ada file yang diunggah"})
		return
	}

	dirPath := filepath.Join("uploads", "kuisapp")
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat folder upload"})
		return
	}

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(dirPath, newFileName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
		return
	}

	urlPath := fmt.Sprintf("/uploads/kuisapp/%s", newFileName)
	c.JSON(http.StatusOK, gin.H{"url": urlPath})
}

// ======== RESULTS ========

func SubmitKuisAppQuiz(c *gin.Context) {
	quizID := c.Param("id")
	
	var input struct {
		Answers map[string]int `json:"answers"` // "questionID" -> answer index
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hitung nilai
	var questions []kuisapp.Question
	if err := database.DB.Where("quiz_id = ?", quizID).Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil soal"})
		return
	}

	correctCount := 0
	totalQuestions := len(questions)

	if totalQuestions > 0 {
		for _, q := range questions {
			qIDStr := strconv.FormatUint(uint64(q.ID), 10)
			userAns, exists := input.Answers[qIDStr]
			if exists && userAns == q.Answer {
				correctCount++
			}
		}
	}

	score := float64(0)
	if totalQuestions > 0 {
		score = float64(correctCount) / float64(totalQuestions) * 100
	}

	userObj, _ := c.Get("kuisapp_user")
	user := userObj.(kuisapp.User)

	qID, _ := strconv.ParseUint(quizID, 10, 32)

	var previousAttempts int64
	database.DB.Model(&kuisapp.Result{}).Where("user_id = ? AND quiz_id = ?", user.ID, uint(qID)).Count(&previousAttempts)

	pointsToAdd := 0
	if previousAttempts == 0 {
		if score >= 100 {
			pointsToAdd = 5
		} else if score >= 90 {
			pointsToAdd = 4
		} else if score >= 80 {
			pointsToAdd = 3
		}
	}

	result := kuisapp.Result{
		QuizID:       uint(qID),
		UserID:       user.ID,
		Score:        score,
		PointsEarned: pointsToAdd,
		FinishedAt:   time.Now(),
	}

	if err := database.DB.Create(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan hasil kuis"})
		return
	}

	if pointsToAdd > 0 {
		database.DB.Model(&kuisapp.User{}).Where("id = ?", user.ID).UpdateColumn("points", gorm.Expr("points + ?", pointsToAdd))
	}

	pointsAlreadyClaimed := previousAttempts > 0

	c.JSON(http.StatusOK, gin.H{
		"message":                "Berhasil",
		"score":                  score,
		"points_earned":          pointsToAdd,
		"points_already_claimed": pointsAlreadyClaimed,
	})
}

func GetKuisAppMyResults(c *gin.Context) {
	userObj, _ := c.Get("kuisapp_user")
	user := userObj.(kuisapp.User)

	var results []kuisapp.Result
	if err := database.DB.Preload("Quiz").Where("user_id = ?", user.ID).Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil hasil"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results})
}

func GetKuisAppAllResults(c *gin.Context) {
	var results []kuisapp.Result
	var total int64

	if err := database.DB.Model(&kuisapp.Result{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung data peserta"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "30"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 30
	}
	offset := (page - 1) * limit

	if err := database.DB.Preload("User").Preload("Quiz").Order("finished_at desc").Limit(limit).Offset(offset).Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peserta"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results, "total": total})
}

// ======== USERS ========

func GetKuisAppUsers(c *gin.Context) {
	var users []kuisapp.User
	query := database.DB

	search := c.Query("search")
	if search != "" {
		query = query.Where("username LIKE ?", "%"+search+"%")
	}

	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func CreateKuisAppUser(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
		Status   string `json:"status" binding:"required"`
		Points   int    `json:"points"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username exists
	var existingUser kuisapp.User
	if err := database.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	user := kuisapp.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     input.Role,
		Status:   input.Status,
		Points:   input.Points,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User berhasil dibuat", "data": user})
}

func UpdateKuisAppUser(c *gin.Context) {
	id := c.Param("id")
	var user kuisapp.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	var input struct {
		Username string `json:"username" binding:"required"`
		Role     string `json:"role" binding:"required"`
		Status   string `json:"status" binding:"required"`
		Points   int    `json:"points"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username is being changed to one that already exists
	if user.Username != input.Username {
		var existingUser kuisapp.User
		if err := database.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan"})
			return
		}
	}

	user.Username = input.Username
	user.Role = input.Role
	user.Status = input.Status
	user.Points = input.Points

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User diupdate", "data": user})
}

func DeleteKuisAppUser(c *gin.Context) {
	id := c.Param("id")
	
	// Delete results first
	database.DB.Where("user_id = ?", id).Delete(&kuisapp.Result{})

	if err := database.DB.Delete(&kuisapp.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User dihapus"})
}

func ResetKuisAppUserPassword(c *gin.Context) {
	id := c.Param("id")
	var user kuisapp.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	var input struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	user.Password = string(hashedPassword)

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password direset"})
}

func ToggleKuisAppUserSuspend(c *gin.Context) {
	id := c.Param("id")
	var user kuisapp.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	user.IsSuspended = !user.IsSuspended

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah status suspend user"})
		return
	}

	status := "di-suspend"
	if !user.IsSuspended {
		status = "diaktifkan kembali"
	}

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil " + status, "is_suspended": user.IsSuspended})
}

func ResetKuisAppUserPoints(c *gin.Context) {
id := c.Param("id")
var user kuisapp.User

if err := database.DB.First(&user, id).Error; err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
return
}

user.Points = 0
if err := database.DB.Save(&user).Error; err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal reset poin"})
return
}

c.JSON(http.StatusOK, gin.H{"message": "Poin berhasil direset", "data": user})
}

func ResetAllKuisAppPointsAndHistory(c *gin.Context) {
// Reset all users' points to 0
if err := database.DB.Model(&kuisapp.User{}).Where("role = ?", "user").Update("points", 0).Error; err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mereset poin peserta"})
return
}

// Delete all results (no user_answers table exists)
if err := database.DB.Exec("DELETE FROM kuisapp_results").Error; err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mereset riwayat kuis"})
return
}

c.JSON(http.StatusOK, gin.H{"message": "Seluruh poin dan riwayat kuis berhasil direset"})
}
