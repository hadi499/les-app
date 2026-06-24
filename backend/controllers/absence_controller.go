package controllers

import (
	"net/http"
	"time"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
)

type AbsenceInput struct {
	UserID     uint   `json:"user_id" binding:"required"`
	Date       string `json:"date" binding:"required"`
	ReasonType string `json:"reason_type" binding:"required"`
	Note       string `json:"note"`
}

func CreateAbsence(c *gin.Context) {
	var input AbsenceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	absence := models.Absence{
		UserID:     input.UserID,
		Date:       date,
		ReasonType: input.ReasonType,
		Note:       input.Note,
	}

	if err := database.DB.Create(&absence).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record absence"})
		return
	}

	// Fetch user details for the response if needed
	database.DB.Preload("User").First(&absence, absence.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Absence recorded successfully", "absence": absence})
}

type RecapData struct {
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"`
	Total    int    `json:"total"`
	Sakit    int    `json:"sakit"`
	Izin     int    `json:"izin"`
	Alpa     int    `json:"alpa"`
	TotalAll int    `json:"total_all"`
}

func GetAbsenceRecap(c *gin.Context) {
	// Filter for current month and year
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.Local)

	var absences []models.Absence
	// Fetch all absences for total_all calculation
	if err := database.DB.Preload("User").Find(&absences).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch absence recap"})
		return
	}

	// Group by UserID
	recapMap := make(map[uint]*RecapData)
	
	// Also fetch all students to ensure students with 0 absences are included
	var students []models.User
	if err := database.DB.Where("role = ?", "student").Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	for _, student := range students {
		recapMap[student.ID] = &RecapData{
			UserID:   student.ID,
			Name:     student.Username,
			Total:    0,
			Sakit:    0,
			Izin:     0,
			Alpa:     0,
			TotalAll: 0,
		}
	}

	for _, absence := range absences {
		if recap, exists := recapMap[absence.UserID]; exists {
			recap.TotalAll++

			// Only increment month specific counts if it falls within current month
			if (absence.Date.Equal(startOfMonth) || absence.Date.After(startOfMonth)) && absence.Date.Before(nextMonth) {
				recap.Total++
				if absence.ReasonType == "Sakit" {
					recap.Sakit++
				} else if absence.ReasonType == "Izin" {
					recap.Izin++
				} else if absence.ReasonType == "Alpa" {
					recap.Alpa++
				}
			}
		}
	}

	var results []RecapData
	for _, recap := range recapMap {
		results = append(results, *recap)
	}

	c.JSON(http.StatusOK, gin.H{"recap": results})
}
func ResetAbsences(c *gin.Context) {
	if err := database.DB.Exec("DELETE FROM absences").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset absences"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Absences reset successfully"})
}

func GetAbsenceHistory(c *gin.Context) {
	userId := c.Param("id")

	var absences []models.Absence
	if err := database.DB.Preload("User").Where("user_id = ?", userId).Order("date desc").Find(&absences).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch absence history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": absences})
}

type UpdateAbsenceInput struct {
	Date       string `json:"date"`
	ReasonType string `json:"reason_type"`
	Note       string `json:"note"`
}

func UpdateAbsence(c *gin.Context) {
	absenceId := c.Param("id")

	var absence models.Absence
	if err := database.DB.First(&absence, absenceId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Absence record not found"})
		return
	}

	var input UpdateAbsenceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.Date != "" {
		date, err := time.Parse("2006-01-02", input.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
			return
		}
		absence.Date = date
	}
	
	if input.ReasonType != "" {
		absence.ReasonType = input.ReasonType
	}
	
	absence.Note = input.Note

	if err := database.DB.Save(&absence).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update absence record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Absence updated successfully", "absence": absence})
}

func DeleteAbsence(c *gin.Context) {
	absenceId := c.Param("id")

	var absence models.Absence
	if err := database.DB.First(&absence, absenceId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Absence record not found"})
		return
	}

	if err := database.DB.Delete(&absence).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete absence record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Absence record deleted successfully"})
}
