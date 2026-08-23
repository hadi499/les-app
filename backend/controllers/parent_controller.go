package controllers

import (
	"net/http"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
)

// GetMyChildren returns all children linked to the currently logged-in parent
func GetMyChildren(c *gin.Context) {
	parentIDInter, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	parentID := parentIDInter.(uint)

	var children []models.User
	// Fetch users where parent_id = parentID
	if err := database.DB.Where("parent_id = ?", parentID).Find(&children).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data anak"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"children": children,
	})
}
