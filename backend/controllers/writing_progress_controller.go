package controllers

import (
	"backend/database"
	"backend/models"
	"backend/services"
	"fmt"
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

func BackupToDrive(c *gin.Context) {
	// Import os and services above manually if not exists (handled by goimports usually, but let's assume we need to import)
	// We will rely on Go's LSP / our manual import fix if needed. Wait, we should just import it.
	
	// Cari semua writing progress yang memiliki image tetapi belum dibackup (drive_file_id kosong atau NULL)
	var progresses []models.WritingProgress
	if err := database.DB.Preload("User").Where("image != ? AND (drive_file_id = ? OR drive_file_id IS NULL)", "", "").Find(&progresses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data untuk backup"})
		return
	}

	if len(progresses) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Semua data sudah dibackup ke Google Drive."})
		return
	}

	folderId := os.Getenv("DRIVE_FOLDER_ID")
	
	// TODO: Actually we need to import "backend/services" and "fmt"
	// For now we assume we'll fix imports right after
	
	successCount := 0
	var errorDetails []string
	
	for _, wp := range progresses {
		localPath := strings.TrimPrefix(wp.Image, "/")
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			errorDetails = append(errorDetails, fmt.Sprintf("File lokal tidak ditemukan: %s", localPath))
			continue // File tidak ditemukan, lewati
		}

		// Format nama: ProgressMenulis_[Username]_[Tanggal].jpg
		ext := ".jpg" // default
		if strings.HasSuffix(localPath, ".png") {
			ext = ".png"
		}
		
		dateStr := wp.Date.Format("2006-01-02")
		customName := fmt.Sprintf("ProgressMenulis_%s_%s%s", wp.User.Username, dateStr, ext)

		driveFileId, err := services.UploadFileToDrive(localPath, customName, folderId)
		if err == nil && driveFileId != "" {
			// Update database
			database.DB.Model(&wp).Update("drive_file_id", driveFileId)
			successCount++
		} else {
			errorDetails = append(errorDetails, fmt.Sprintf("Gagal upload %s: %v", customName, err))
		}
	}

	msg := fmt.Sprintf("Berhasil membackup %d dari %d gambar.", successCount, len(progresses))
	if len(errorDetails) > 0 {
		c.JSON(http.StatusOK, gin.H{"message": msg, "error_details": errorDetails})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": msg})
	}
}
