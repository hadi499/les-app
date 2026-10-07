package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- CATEGORY ENDPOINTS ---

func GetEnglishFlashcardCategories(c *gin.Context) {
	var categories []models.EnglishFlashcardCategory
	if err := database.DB.Order("id DESC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil kategori"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func CreateEnglishFlashcardCategory(c *gin.Context) {
	var input models.EnglishFlashcardCategory
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kategori"})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func UpdateEnglishFlashcardCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.EnglishFlashcardCategory
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := database.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update kategori"})
		return
	}
	c.JSON(http.StatusOK, category)
}

func DeleteEnglishFlashcardCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.EnglishFlashcardCategory
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}
	if err := database.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
}


// --- FLASHCARD ENDPOINTS ---

func GetEnglishFlashcards(c *gin.Context) {
	categoryID := c.Query("category_id")
	search := c.Query("search")
	var flashcards []models.EnglishFlashcard
	
	query := database.DB.Preload("Category").Order("english_flashcards.id DESC").Model(&models.EnglishFlashcard{})
	
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Joins("LEFT JOIN english_flashcard_categories ON english_flashcard_categories.id = english_flashcards.category_id").
			Where("english_flashcards.question ILIKE ? OR english_flashcards.answer ILIKE ? OR english_flashcard_categories.name ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if err := query.Find(&flashcards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data flashcard"})
		return
	}
	c.JSON(http.StatusOK, flashcards)
}

func GetEnglishFlashcardByID(c *gin.Context) {
	id := c.Param("id")
	var flashcard models.EnglishFlashcard
	if err := database.DB.Preload("Category").First(&flashcard, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flashcard tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, flashcard)
}

func CreateEnglishFlashcard(c *gin.Context) {
	var input models.EnglishFlashcard
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat flashcard"})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func UpdateEnglishFlashcard(c *gin.Context) {
	id := c.Param("id")
	var flashcard models.EnglishFlashcard
	if err := database.DB.First(&flashcard, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flashcard tidak ditemukan"})
		return
	}

	if err := c.ShouldBindJSON(&flashcard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := database.DB.Save(&flashcard).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update flashcard"})
		return
	}
	c.JSON(http.StatusOK, flashcard)
}

func DeleteEnglishFlashcard(c *gin.Context) {
	id := c.Param("id")
	var flashcard models.EnglishFlashcard
	if err := database.DB.First(&flashcard, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flashcard tidak ditemukan"})
		return
	}
	if err := database.DB.Delete(&flashcard).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus flashcard"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Flashcard berhasil dihapus"})
}
