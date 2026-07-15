package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetExams(c *gin.Context) {
	var exams []models.Exam
	
	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")

	query := database.DB.Preload("User").Preload("Subject").Order("exam_date desc")
	if role != "teacher" && role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&exams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exams"})
		return
	}
	c.JSON(http.StatusOK, exams)
}

// GetExamByID fetches a single exam by its ID
func GetExamByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var exam models.Exam
	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")

	query := database.DB.Preload("User").Preload("Subject")
	if role != "teacher" && role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&exam, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam score not found"})
		return
	}

	c.JSON(http.StatusOK, exam)
}

// CreateExam creates a new exam score record
func CreateExam(c *gin.Context) {
	var input struct {
		UserID   uint   `json:"user_id" binding:"required"`
		ExamDate string `json:"exam_date" binding:"required"`
		ExamName  string `json:"exam_name" binding:"required"`
		SubjectID uint   `json:"subject_id" binding:"required"`
		Score     int    `json:"score"`
		Image     string `json:"image"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", input.ExamDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	exam := models.Exam{
		UserID:   input.UserID,
		ExamDate: date,
		ExamName:  input.ExamName,
		SubjectID: input.SubjectID,
		Score:     input.Score,
		Image:     input.Image,
	}

	if err := database.DB.Create(&exam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exam score"})
		return
	}

	// Fetch with user and subject to return
	database.DB.Preload("User").Preload("Subject").First(&exam, exam.ID)
	c.JSON(http.StatusCreated, exam)
}

// UpdateExam updates an existing exam score record
func UpdateExam(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var exam models.Exam
	if err := database.DB.First(&exam, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam score not found"})
		return
	}

	var input struct {
		UserID   uint   `json:"user_id" binding:"required"`
		ExamDate string `json:"exam_date" binding:"required"`
		ExamName  string `json:"exam_name" binding:"required"`
		SubjectID uint   `json:"subject_id" binding:"required"`
		Score     int    `json:"score"`
		Image     string `json:"image"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", input.ExamDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	oldImage := exam.Image

	exam.UserID = input.UserID
	exam.ExamDate = date
	exam.ExamName = input.ExamName
	exam.SubjectID = input.SubjectID
	exam.Score = input.Score
	exam.Image = input.Image

	if err := database.DB.Save(&exam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exam score"})
		return
	}

	if oldImage != "" && oldImage != input.Image {
		filePath := strings.TrimPrefix(oldImage, "/")
		_ = os.Remove(filePath)
	}

	database.DB.Preload("User").Preload("Subject").First(&exam, exam.ID)
	c.JSON(http.StatusOK, exam)
}

func DeleteExam(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var exam models.Exam
	if err := database.DB.First(&exam, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam score not found"})
		return
	}

	if err := database.DB.Delete(&exam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exam score"})
		return
	}

	if exam.Image != "" {
		filePath := strings.TrimPrefix(exam.Image, "/")
		_ = os.Remove(filePath)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exam score deleted successfully"})
}

// BulkDeleteExams deletes multiple exam score records
func BulkDeleteExams(c *gin.Context) {
	var input struct {
		IDs []int `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if len(input.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No IDs provided"})
		return
	}

	var exams []models.Exam
	if err := database.DB.Where("id IN ?", input.IDs).Find(&exams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exams to delete"})
		return
	}

	if err := database.DB.Where("id IN ?", input.IDs).Delete(&models.Exam{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exam scores"})
		return
	}

	for _, exam := range exams {
		if exam.Image != "" {
			filePath := strings.TrimPrefix(exam.Image, "/")
			_ = os.Remove(filePath)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exam scores deleted successfully"})
}
