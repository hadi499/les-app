package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"
	"strconv"
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

// CreateExam creates a new exam score record
func CreateExam(c *gin.Context) {
	var input struct {
		UserID   uint   `json:"user_id" binding:"required"`
		ExamDate string `json:"exam_date" binding:"required"`
		ExamName  string `json:"exam_name" binding:"required"`
		SubjectID uint   `json:"subject_id" binding:"required"`
		Score     int    `json:"score"`
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

	exam.UserID = input.UserID
	exam.ExamDate = date
	exam.ExamName = input.ExamName
	exam.SubjectID = input.SubjectID
	exam.Score = input.Score

	if err := database.DB.Save(&exam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exam score"})
		return
	}

	database.DB.Preload("User").Preload("Subject").First(&exam, exam.ID)
	c.JSON(http.StatusOK, exam)
}

// DeleteExam deletes an exam score record
func DeleteExam(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := database.DB.Delete(&models.Exam{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exam score"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exam score deleted successfully"})
}
