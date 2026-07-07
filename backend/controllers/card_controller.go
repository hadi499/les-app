package controllers

import (
	"backend/database"
	"backend/models"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func deleteCardImage(card models.Card) {
	if card.Image == "" {
		return
	}
	
	// Format di DB biasanya `/uploads/cards/filename.ext` atau `/uploads/filename.ext`
	// Hapus prefix "/" agar menjadi path relatif ke direktori kerja (backend)
	filePath := strings.TrimPrefix(card.Image, "/")
	
	// Validasi keamanan sederhana
	if filePath == "." || filePath == "" || strings.Contains(filePath, "..") || strings.HasPrefix(filepath.Base(filePath), "?") {
		return
	}
	
	if err := os.Remove(filePath); err != nil {
		log.Printf("Failed to delete image %s: %v", filePath, err)
	}
}

func GetCards(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "6"))
	size := c.Query("size")
	search := c.Query("search")
	cardType := c.Query("cardType")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 10000 {
		limit = 6
	}

	showAll := c.Query("all")

	query := database.DB.Model(&models.Card{}).Preload("CardFolder")

	if showAll != "true" {
		if cardType == "image" {
			query = query.Where("card_type = ?", "image")
		} else {
			query = query.Where("(card_type = ? OR card_type = '' OR card_type IS NULL)", "regular")
		}
	}

	if size != "" {
		query = query.Where("size = ?", size)
	}
	if search != "" {
		query = query.Joins("LEFT JOIN card_folders ON card_folders.id = cards.card_folder_id").
			Where("(cards.title ILIKE ? OR card_folders.name ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var countAll, count6, count4, count2, countGambar int64
	countBaseAll := database.DB.Model(&models.Card{})
	countBase6 := database.DB.Model(&models.Card{}).Where("(card_type = ? OR card_type = '' OR card_type IS NULL)", "regular")
	countBase4 := database.DB.Model(&models.Card{}).Where("(card_type = ? OR card_type = '' OR card_type IS NULL)", "regular")
	countBase2 := database.DB.Model(&models.Card{}).Where("(card_type = ? OR card_type = '' OR card_type IS NULL)", "regular")
	countBaseGambar := database.DB.Model(&models.Card{}).Where("card_type = ?", "image")
	if search != "" {
		countBaseAll = countBaseAll.Joins("LEFT JOIN card_folders ON card_folders.id = cards.card_folder_id").Where("(cards.title ILIKE ? OR card_folders.name ILIKE ?)", "%"+search+"%", "%"+search+"%")
		countBase6 = countBase6.Joins("LEFT JOIN card_folders ON card_folders.id = cards.card_folder_id").Where("(cards.title ILIKE ? OR card_folders.name ILIKE ?)", "%"+search+"%", "%"+search+"%")
		countBase4 = countBase4.Joins("LEFT JOIN card_folders ON card_folders.id = cards.card_folder_id").Where("(cards.title ILIKE ? OR card_folders.name ILIKE ?)", "%"+search+"%", "%"+search+"%")
		countBase2 = countBase2.Joins("LEFT JOIN card_folders ON card_folders.id = cards.card_folder_id").Where("(cards.title ILIKE ? OR card_folders.name ILIKE ?)", "%"+search+"%", "%"+search+"%")
		countBaseGambar = countBaseGambar.Joins("LEFT JOIN card_folders ON card_folders.id = cards.card_folder_id").Where("(cards.title ILIKE ? OR card_folders.name ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}
	countBaseAll.Count(&countAll)
	countBase6.Where("size = ? OR size = '' OR size IS NULL", "6").Count(&count6)
	countBase4.Where("size = ?", "4").Count(&count4)
	countBase2.Where("size = ?", "2").Count(&count2)
	countBaseGambar.Count(&countGambar)

	var cards []models.Card
	offset := (page - 1) * limit
	query.Order("id desc").Offset(offset).Limit(limit).Find(&cards)

	c.JSON(http.StatusOK, gin.H{
		"data":       cards,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int(math.Ceil(float64(total) / float64(limit))),
		"counts": gin.H{
			"all":    countAll,
			"6":      count6,
			"4":      count4,
			"2":      count2,
			"gambar": countGambar,
		},
	})
}

func CreateCard(c *gin.Context) {
	var card models.Card
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&card)
	c.JSON(http.StatusCreated, card)
}

func UpdateCard(c *gin.Context) {
	id := c.Param("id")
	var card models.Card
	if err := database.DB.First(&card, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	var input struct {
		CardFolderID *uint  `json:"card_folder_id"`
		Image        string `json:"image"`
		Title        string `json:"title"`
		Content      string `json:"content"`
		Size         string `json:"size"`
		CardType     string `json:"cardType"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hapus gambar lama jika gambar diubah (URL berbeda) dan tidak kosong
	if input.Image != card.Image && card.Image != "" {
		deleteCardImage(card)
	}

	card.CardFolderID = input.CardFolderID
	card.Image = input.Image
	card.Title = input.Title
	card.Content = input.Content
	card.Size = input.Size
	card.CardType = input.CardType
	database.DB.Save(&card)
	c.JSON(http.StatusOK, card)
}

func DeleteCard(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Card{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Card deleted"})
}

func GetTrashCards(c *gin.Context) {
	var cards []models.Card
	database.DB.Unscoped().Preload("CardFolder").Where("deleted_at IS NOT NULL").Order("deleted_at desc").Find(&cards)
	var total int64
	database.DB.Unscoped().Model(&models.Card{}).Where("deleted_at IS NOT NULL").Count(&total)
	c.JSON(http.StatusOK, gin.H{
		"data":  cards,
		"total": total,
	})
}

func RestoreCard(c *gin.Context) {
	id := c.Param("id")
	var card models.Card
	if err := database.DB.Unscoped().First(&card, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}
	database.DB.Unscoped().Model(&card).Update("deleted_at", nil)
	c.JSON(http.StatusOK, card)
}

func ForceDeleteCard(c *gin.Context) {
	id := c.Param("id")
	var card models.Card
	if err := database.DB.Unscoped().First(&card, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}
	deleteCardImage(card)
	if err := database.DB.Unscoped().Delete(&models.Card{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Card permanently deleted"})
}

func EmptyTrash(c *gin.Context) {
	var cards []models.Card
	database.DB.Unscoped().Where("deleted_at IS NOT NULL").Find(&cards)
	for _, card := range cards {
		deleteCardImage(card)
	}
	database.DB.Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.Card{})
	c.JSON(http.StatusOK, gin.H{"message": "Trash emptied"})
}

func UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	if header.Size > 1*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file maksimal 1MB"})
		return
	}

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Tentukan subfolder berdasarkan query parameter "type"
	uploadType := c.Query("type")
	subfolder := "cards" // default folder untuk image card
	if uploadType == "exam" {
		subfolder = "exams"
	}

	// Buat folder jika belum ada
	dirPath := filepath.Join("uploads", subfolder)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	savePath := filepath.Join(dirPath, filename)
	out, err := os.Create(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}

	// Format URL relatif
	urlPath := "uploads"
	if subfolder != "" {
		urlPath += "/" + subfolder
	}
	url := fmt.Sprintf("/%s/%s", urlPath, filename)
	
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func ListImages(c *gin.Context) {
	entries, err := os.ReadDir("uploads")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read uploads"})
		return
	}

	type ImageEntry struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	var images []ImageEntry
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp" {
			images = append(images, ImageEntry{
				Name: entry.Name(),
				URL:  fmt.Sprintf("/uploads/%s", entry.Name()),
			})
		}
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].Name > images[j].Name
	})

	c.JSON(http.StatusOK, gin.H{"images": images})
}
