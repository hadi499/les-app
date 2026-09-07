package controllers

import (
	"backend/database"
	"backend/models"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetProgress fetches all lesson progress for the authenticated user
func GetProgress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var progress []models.LessonProgress
	if err := database.DB.Where("user_id = ?", userID).Find(&progress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// SaveProgress updates or creates progress for a lesson
func SaveProgress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input models.LessonProgress
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.UserID = userID.(uint)

	var existing models.LessonProgress
	err := database.DB.Where("user_id = ? AND lesson_id = ?", input.UserID, input.LessonID).First(&existing).Error
	if err == nil {
		// Update logic
		if input.BestWPM > existing.BestWPM {
			existing.BestWPM = input.BestWPM
		}
		if input.BestAccuracy > existing.BestAccuracy {
			existing.BestAccuracy = input.BestAccuracy
		}
		if input.Stars > existing.Stars {
			existing.Stars = input.Stars
		}
		existing.Completed = existing.Completed || input.Completed
		existing.Attempts += 1

		database.DB.Save(&existing)

		// Save history
		history := models.LessonHistory{
			UserID:    input.UserID,
			LessonID:  input.LessonID,
			WPM:       input.BestWPM,
			Accuracy:  input.BestAccuracy,
			Stars:     input.Stars,
			CreatedAt: time.Now(),
		}
		database.DB.Create(&history)

		c.JSON(http.StatusOK, existing)
		return
	}

	// Create new
	input.Attempts = 1
	database.DB.Create(&input)

	history := models.LessonHistory{
		UserID:    input.UserID,
		LessonID:  input.LessonID,
		WPM:       input.BestWPM,
		Accuracy:  input.BestAccuracy,
		Stars:     input.Stars,
		CreatedAt: time.Now(),
	}
	database.DB.Create(&history)

	c.JSON(http.StatusOK, input)
}

// GetGameScores fetches all game scores for the authenticated user
func GetGameScores(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var scores []models.GameHighScore
	if err := database.DB.Where("user_id = ?", userID).Find(&scores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scores"})
		return
	}

	result := map[string]int{
		"left":    0,
		"right":   0,
		"both":    0,
		"letters": 0,
		"all":     0,
	}

	for _, s := range scores {
		result[s.Mode] = s.Score
	}

	c.JSON(http.StatusOK, result)
}

// SaveGameScore updates a game score if it's higher
func SaveGameScore(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Mode  string `json:"mode" binding:"required"`
		Score int    `json:"score" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.GameHighScore
	err := database.DB.Where("user_id = ? AND mode = ?", userID, input.Mode).First(&existing).Error
	if err == nil {
		if input.Score > existing.Score {
			existing.Score = input.Score
			database.DB.Save(&existing)
		}

		// Save history
		history := models.GameHistory{
			UserID:    userID.(uint),
			Mode:      input.Mode,
			Score:     input.Score,
			CreatedAt: time.Now(),
		}
		database.DB.Create(&history)

		c.JSON(http.StatusOK, existing)
		return
	}

	newScore := models.GameHighScore{
		UserID: userID.(uint),
		Mode:   input.Mode,
		Score:  input.Score,
	}
	database.DB.Create(&newScore)

	history := models.GameHistory{
		UserID:    userID.(uint),
		Mode:      input.Mode,
		Score:     input.Score,
		CreatedAt: time.Now(),
	}
	database.DB.Create(&history)

	c.JSON(http.StatusOK, newScore)
}

func GetGameHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var history []models.GameHistory
	if err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(20).Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}

	c.JSON(http.StatusOK, history)
}

func GetLessonHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var history []models.LessonHistory
	if err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(20).Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}

	c.JSON(http.StatusOK, history)
}

// --- Admin / Teacher Endpoints ---

func GetAllLessonProgress(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "24"))
	offset := (page - 1) * limit

	var progress []models.LessonProgress
	var total int64

	database.DB.Model(&models.LessonProgress{}).Count(&total)

	if err := database.DB.Preload("User").Order("lesson_id asc").Offset(offset).Limit(limit).Find(&progress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch progress"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data":         progress,
		"total_pages":  totalPages,
		"current_page": page,
		"total_items":  total,
	})
}

func GetAllGameScores(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "24"))
	offset := (page - 1) * limit

	var scores []models.GameHighScore
	var total int64

	database.DB.Model(&models.GameHighScore{}).Count(&total)

	if err := database.DB.Preload("User").Order("score desc").Offset(offset).Limit(limit).Find(&scores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scores"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data":         scores,
		"total_pages":  totalPages,
		"current_page": page,
		"total_items":  total,
	})
}

func GetAllGameHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "24"))
	offset := (page - 1) * limit

	var history []models.GameHistory
	var total int64

	database.DB.Model(&models.GameHistory{}).Count(&total)

	if err := database.DB.Preload("User").Order("created_at desc").Offset(offset).Limit(limit).Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data":         history,
		"total_pages":  totalPages,
		"current_page": page,
		"total_items":  total,
	})
}

func GetAllLessonHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "24"))
	offset := (page - 1) * limit

	var history []models.LessonHistory
	var total int64

	database.DB.Model(&models.LessonHistory{}).Count(&total)

	if err := database.DB.Preload("User").Order("created_at desc").Offset(offset).Limit(limit).Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data":         history,
		"total_pages":  totalPages,
		"current_page": page,
		"total_items":  total,
	})
}
