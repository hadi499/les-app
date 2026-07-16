package controllers

import (
	"backend/database"
	"backend/models"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Helper: ekstrak URL gambar dari konten HTML
func extractMateriImagePaths(content string) []string {
	var paths []string
	re := regexp.MustCompile(`src="(/uploads/materis/[^"]+)"`)
	matches := re.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			paths = append(paths, match[1])
		}
	}
	return paths
}

// Helper: hapus file gambar
func deleteMateriImages(paths []string) {
	for _, p := range paths {
		filePath := strings.TrimPrefix(p, "/")
		if filePath == "." || filePath == "" || strings.Contains(filePath, "..") {
			continue
		}
		os.Remove(filePath)
	}
}

// GetMateris mengambil materi dengan paginasi
func GetMateris(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var materis []models.Materi
	query := database.DB.Model(&models.Materi{}).Preload("Subject").Preload("Users")

	if role != "teacher" {
		query = query.Joins("JOIN materi_users ON materi_users.materi_id = materis.id").Where("materi_users.user_id = ?", userID)
	}

	search := c.Query("search")
	if search != "" {
		query = query.Where("LOWER(materis.title) LIKE LOWER(?) OR LOWER(materis.content) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%")
	}

	subjectID := c.Query("subject_id")
	if subjectID != "" && subjectID != "all" {
		query = query.Where("materis.subject_id = ?", subjectID)
	}

	var total int64
	query.Count(&total)

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "12")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 12
	}
	offset := (page - 1) * limit

	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&materis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil materi"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data":        materis,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// GetMateriByID
func GetMateriByID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	materiID := c.Param("id")

	var materi models.Materi
	query := database.DB.Preload("Subject").Preload("Users").Where("materis.id = ?", materiID)

	if role != "teacher" {
		query = query.Joins("JOIN materi_users ON materi_users.materi_id = materis.id").Where("materi_users.user_id = ?", userID)
	}

	if err := query.First(&materi).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Materi tidak ditemukan atau Anda tidak memiliki akses"})
		return
	}

	c.JSON(http.StatusOK, materi)
}

// CreateMateri
func CreateMateri(c *gin.Context) {
	var input struct {
		Title     string `json:"title" binding:"required"`
		SubjectID uint   `json:"subject_id" binding:"required"`
		UserIDs   []uint `json:"user_ids"` // Bisa kosong jika ditujukan untuk semua atau opsional
		Content   string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	materi := models.Materi{
		Title:     input.Title,
		SubjectID: input.SubjectID,
		Content:   input.Content,
	}

	if err := database.DB.Create(&materi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat materi"})
		return
	}

	// Tambahkan relasi ke users jika ada
	if len(input.UserIDs) > 0 {
		var users []models.User
		database.DB.Where("id IN ?", input.UserIDs).Find(&users)
		database.DB.Model(&materi).Association("Users").Append(users)
	}

	database.DB.Preload("Subject").Preload("Users").First(&materi, materi.ID)
	c.JSON(http.StatusCreated, materi)
}

// UpdateMateri
func UpdateMateri(c *gin.Context) {
	materiID := c.Param("id")

	var materi models.Materi
	if err := database.DB.Where("id = ?", materiID).First(&materi).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Materi tidak ditemukan"})
		return
	}

	var input struct {
		Title     string `json:"title" binding:"required"`
		SubjectID uint   `json:"subject_id" binding:"required"`
		UserIDs   []uint `json:"user_ids"`
		Content   string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oldContent := materi.Content
	materi.Title = input.Title
	materi.SubjectID = input.SubjectID
	materi.Content = input.Content

	if err := database.DB.Save(&materi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate materi"})
		return
	}

	// Hapus gambar lama yang tidak ada di konten baru
	oldPaths := extractMateriImagePaths(oldContent)
	newPaths := extractMateriImagePaths(input.Content)
	newPathsMap := make(map[string]bool)
	for _, p := range newPaths {
		newPathsMap[p] = true
	}

	var pathsToDelete []string
	for _, p := range oldPaths {
		if !newPathsMap[p] {
			pathsToDelete = append(pathsToDelete, p)
		}
	}
	deleteMateriImages(pathsToDelete)

	// Update relasi users
	var users []models.User
	if len(input.UserIDs) > 0 {
		database.DB.Where("id IN ?", input.UserIDs).Find(&users)
	}
	database.DB.Model(&materi).Association("Users").Replace(users)

	database.DB.Preload("Subject").Preload("Users").First(&materi, materi.ID)
	c.JSON(http.StatusOK, materi)
}

// DeleteMateri
func DeleteMateri(c *gin.Context) {
	materiID := c.Param("id")

	var materi models.Materi
	if err := database.DB.Where("id = ?", materiID).First(&materi).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Materi tidak ditemukan"})
		return
	}

	// Ambil path gambar sebelum materi dihapus dari DB
	imagePaths := extractMateriImagePaths(materi.Content)

	if err := database.DB.Delete(&materi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus materi"})
		return
	}

	// Hapus file gambar dari server
	deleteMateriImages(imagePaths)

	c.JSON(http.StatusOK, gin.H{"message": "Materi berhasil dihapus"})
}
