package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPublicQuotes retrieves all published quotes
func GetPublicQuotes(c *gin.Context) {
	var quotes []models.Quote
	if err := database.DB.Where("is_published = ?", true).Order("created_at asc").Find(&quotes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve quotes"})
		return
	}
	c.JSON(http.StatusOK, quotes)
}

// GetAllQuotes retrieves all quotes for admin
func GetAllQuotes(c *gin.Context) {
	var quotes []models.Quote
	if err := database.DB.Order("is_published desc, created_at asc").Find(&quotes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve quotes"})
		return
	}
	c.JSON(http.StatusOK, quotes)
}

// CreateQuote creates a new quote
func CreateQuote(c *gin.Context) {
	var input struct {
		Quote       string `json:"quote" binding:"required"`
		Arti        string `json:"arti" binding:"required"`
		Author      string `json:"author" binding:"required"`
		IsPublished *bool  `json:"is_published"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isPublished := false
	if input.IsPublished != nil {
		isPublished = *input.IsPublished
	}

	if isPublished {
		var publishedCount int64
		database.DB.Model(&models.Quote{}).Where("is_published = ?", true).Count(&publishedCount)
		if publishedCount >= 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Maksimal 3 quote yang dapat di-publish. Harap ubah quote lain menjadi draft terlebih dahulu."})
			return
		}
	}

	quote := models.Quote{
		Quote:       input.Quote,
		Arti:        input.Arti,
		Author:      input.Author,
		IsPublished: isPublished,
	}

	if err := database.DB.Create(&quote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quote"})
		return
	}

	// Workaround: GORM mengabaikan nilai 'false' saat Create jika ada tag default:true
	// Jadi kita paksa set ke false setelah dibuat jika isPublished adalah false menggunakan raw SQL
	if !isPublished {
		database.DB.Exec("UPDATE quotes SET is_published = false WHERE id = ?", quote.ID)
		quote.IsPublished = false // Pastikan response ke frontend juga false
	}

	c.JSON(http.StatusCreated, quote)
}

// UpdateQuote updates an existing quote
func UpdateQuote(c *gin.Context) {
	id := c.Param("id")
	var quote models.Quote

	if err := database.DB.First(&quote, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		return
	}

	var input struct {
		Quote       string `json:"quote"`
		Arti        string `json:"arti"`
		Author      string `json:"author"`
		IsPublished *bool  `json:"is_published"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Quote != "" {
		quote.Quote = input.Quote
	}
	if input.Arti != "" {
		quote.Arti = input.Arti
	}
	if input.Author != "" {
		quote.Author = input.Author
	}

	if input.IsPublished != nil {
		// Jika kita mengubah status dari draft ke publish
		if *input.IsPublished && !quote.IsPublished {
			var publishedCount int64
			database.DB.Model(&models.Quote{}).Where("is_published = ?", true).Count(&publishedCount)
			if publishedCount >= 3 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Maksimal 3 quote yang dapat di-publish. Harap ubah quote lain menjadi draft terlebih dahulu."})
				return
			}
		}
		quote.IsPublished = *input.IsPublished
	}

	if err := database.DB.Save(&quote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quote"})
		return
	}

	c.JSON(http.StatusOK, quote)
}

// DeleteQuote deletes a quote
func DeleteQuote(c *gin.Context) {
	id := c.Param("id")
	var quote models.Quote

	if err := database.DB.First(&quote, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		return
	}

	if err := database.DB.Delete(&quote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete quote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quote deleted successfully"})
}
