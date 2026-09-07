package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// Helper: ekstrak URL gambar dari konten HTML
func extractNoteImagePaths(content string) []string {
	var paths []string
	re := regexp.MustCompile(`src="(/uploads/notes/[^"]+)"`)
	matches := re.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			paths = append(paths, match[1])
		}
	}
	return paths
}

// Helper: hapus file gambar
func deleteNoteImages(paths []string) {
	for _, p := range paths {
		filePath := strings.TrimPrefix(p, "/")
		if filePath == "." || filePath == "" || strings.Contains(filePath, "..") {
			continue
		}
		os.Remove(filePath)
	}
}

// GetNotes mengambil catatan milik pengguna yang sedang login
func GetNotes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var notes []models.Note
	if err := database.DB.Preload("Folder").Where("user_id = ?", userID).Order("created_at desc").Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil catatan"})
		return
	}

	c.JSON(http.StatusOK, notes)
}

// CreateNote
func CreateNote(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Title    string `json:"title" binding:"required"`
		FolderID *uint  `json:"folder_id"`
		Content  string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify if the folder belongs to the user if folder_id is provided
	if input.FolderID != nil {
		var folder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", *input.FolderID, userID).First(&folder).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Folder tidak ditemukan atau Anda tidak memiliki akses ke folder ini"})
			return
		}
	}

	note := models.Note{
		Title:    input.Title,
		FolderID: input.FolderID,
		Content:  input.Content,
		UserID:   userID.(uint),
	}

	if err := database.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat catatan"})
		return
	}
	
	// Preload the folder for the response
	database.DB.Preload("Folder").First(&note, note.ID)

	c.JSON(http.StatusCreated, note)
}

// UpdateNote
func UpdateNote(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	noteID := c.Param("id")

	var note models.Note
	if err := database.DB.Where("id = ? AND user_id = ?", noteID, userID).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Catatan tidak ditemukan"})
		return
	}

	var input struct {
		Title    string `json:"title" binding:"required"`
		FolderID *uint  `json:"folder_id"`
		Content  string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify if the folder belongs to the user if folder_id is provided
	if input.FolderID != nil {
		var folder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", *input.FolderID, userID).First(&folder).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Folder tidak ditemukan atau Anda tidak memiliki akses ke folder ini"})
			return
		}
	}

	note.Title = input.Title
	note.FolderID = input.FolderID
	note.Content = input.Content

	if err := database.DB.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan catatan"})
		return
	}
	
	database.DB.Preload("Folder").First(&note, note.ID)

	c.JSON(http.StatusOK, note)
}

// DeleteNote
func DeleteNote(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	noteID := c.Param("id")

	var note models.Note
	if err := database.DB.Where("id = ? AND user_id = ?", noteID, userID).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Catatan tidak ditemukan"})
		return
	}

	// Ekstrak dan hapus gambar terkait
	paths := extractNoteImagePaths(note.Content)
	deleteNoteImages(paths)

	if err := database.DB.Delete(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus catatan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Catatan berhasil dihapus"})
}
