package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"math"

	"github.com/gin-gonic/gin"
)

func GetWritingProgresses(c *gin.Context) {
	var progresses []models.WritingProgress
	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit < 1 {
		limit = 50
	}
	offset := (page - 1) * limit

	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")

	query := database.DB.Model(&models.WritingProgress{})
	if role != "teacher" && role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung data perkembangan menulis"})
		return
	}

	query = query.Preload("User").Order("created_at desc").Limit(limit).Offset(offset)

	if err := query.Find(&progresses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data perkembangan menulis"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        progresses,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

func GetWritingProgressByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID tidak valid"})
		return
	}

	var progress models.WritingProgress
	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")

	query := database.DB.Preload("User")
	if role != "teacher" && role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&progress, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data perkembangan menulis tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

func CreateWritingProgress(c *gin.Context) {
	var input struct {
		UserID uint   `json:"user_id" binding:"required"`
		Date   string `json:"date" binding:"required"`
		Image  string `json:"image" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal salah. Gunakan YYYY-MM-DD"})
		return
	}

	progress := models.WritingProgress{
		UserID: input.UserID,
		Date:   date,
		Image:  input.Image,
	}

	if err := database.DB.Create(&progress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data perkembangan menulis"})
		return
	}

	database.DB.Preload("User").First(&progress, progress.ID)
	c.JSON(http.StatusCreated, progress)
}

func UpdateWritingProgress(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var progress models.WritingProgress
	if err := database.DB.First(&progress, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data perkembangan menulis tidak ditemukan"})
		return
	}

	var input struct {
		UserID uint   `json:"user_id" binding:"required"`
		Date   string `json:"date" binding:"required"`
		Image  string `json:"image" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal salah. Gunakan YYYY-MM-DD"})
		return
	}

	oldImage := progress.Image

	progress.UserID = input.UserID
	progress.Date = date
	progress.Image = input.Image

	if err := database.DB.Save(&progress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate data perkembangan menulis"})
		return
	}

	if oldImage != "" && oldImage != input.Image {
		filePath := strings.TrimPrefix(oldImage, "/")
		_ = os.Remove(filePath)
	}

	database.DB.Preload("User").First(&progress, progress.ID)
	c.JSON(http.StatusOK, progress)
}

func DeleteWritingProgress(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var progress models.WritingProgress
	if err := database.DB.First(&progress, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data perkembangan menulis tidak ditemukan"})
		return
	}

	if err := database.DB.Delete(&progress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data perkembangan menulis"})
		return
	}

	if progress.Image != "" {
		filePath := strings.TrimPrefix(progress.Image, "/")
		_ = os.Remove(filePath)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data perkembangan menulis berhasil dihapus"})
}
