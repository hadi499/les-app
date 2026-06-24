package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetFolders mengambil semua folder milik pengguna
func GetFolders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var folders []models.Folder
	if err := database.DB.Where("user_id = ?", userID).Order("created_at asc").Find(&folders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil folder"})
		return
	}

	c.JSON(http.StatusOK, folders)
}

// CreateFolder membuat folder baru
func CreateFolder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder := models.Folder{
		Name:   input.Name,
		UserID: userID.(uint),
	}

	if err := database.DB.Create(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat folder"})
		return
	}

	c.JSON(http.StatusCreated, folder)
}

// UpdateFolder mengubah nama folder
func UpdateFolder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	folderID := c.Param("id")

	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder tidak ditemukan"})
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder.Name = input.Name

	if err := database.DB.Save(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan folder"})
		return
	}

	c.JSON(http.StatusOK, folder)
}

// DeleteFolder menghapus folder (beserta isinya karena constraint CASCADE)
func DeleteFolder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	folderID := c.Param("id")

	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).Delete(&models.Folder{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus folder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder berhasil dihapus"})
}
