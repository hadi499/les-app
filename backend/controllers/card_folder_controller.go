package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCardFolders(c *gin.Context) {
	var folders []models.CardFolder
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	userIDVal, _ := c.Get("user_id")
	userID, _ := userIDVal.(uint)

	query := database.DB.Order("name asc")
	if role == "teacher" {
		query = query.Where("user_id = ?", userID)
	}

	query.Find(&folders)
	c.JSON(http.StatusOK, gin.H{"data": folders})
}

func CreateCardFolder(c *gin.Context) {
	var folder models.CardFolder
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDVal, _ := c.Get("user_id")
	userID, _ := userIDVal.(uint)
	folder.UserID = userID

	if err := database.DB.Create(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, folder)
}

func UpdateCardFolder(c *gin.Context) {
	id := c.Param("id")
	var folder models.CardFolder
	if err := database.DB.First(&folder, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card folder not found"})
		return
	}

	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	userIDVal, _ := c.Get("user_id")
	userID, _ := userIDVal.(uint)

	if role == "teacher" && folder.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var input struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder.Name = input.Name
	database.DB.Save(&folder)
	c.JSON(http.StatusOK, folder)
}

func DeleteCardFolder(c *gin.Context) {
	id := c.Param("id")
	var folder models.CardFolder
	if err := database.DB.First(&folder, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card folder not found"})
		return
	}

	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	userIDVal, _ := c.Get("user_id")
	userID, _ := userIDVal.(uint)

	if role == "teacher" && folder.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := database.DB.Delete(&models.CardFolder{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Card folder deleted"})
}
